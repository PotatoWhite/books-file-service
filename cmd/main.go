package main

import (
	"context"
	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/potatowhite/books/file-service/config"
	"github.com/potatowhite/books/file-service/consumer"
	"github.com/potatowhite/books/file-service/db"
	"github.com/potatowhite/books/file-service/graph"
	"github.com/potatowhite/books/file-service/handler/users"
	"github.com/potatowhite/books/file-service/pkg/repository"
	"github.com/potatowhite/books/file-service/pkg/resolver"
	"github.com/potatowhite/books/file-service/pkg/service"
	"gorm.io/gorm"
	"log"
	"net/http"
	"os"
	"regexp"
)

var (
	logger = log.New(os.Stdout, "", log.LstdFlags|log.Lshortfile)
)

func main() {

	// port from args
	port := os.Args[1]

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	if port != "" {
		cfg.Server.Port = port
	}

	database, err := db.InitDB(cfg)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	defer db.CloseDB(database)

	folderRepo, fileRepo := initRepository(database)
	folderSvc, fileSvc := initService(folderRepo, fileRepo)

	userConsumer, err := initUserConsumer(cfg, fileSvc, folderSvc)
	defer userConsumer.Close()

	server := initGraphqlServer(folderSvc, fileSvc)
	startServer(server, cfg.Server.Port)

}

func initUserConsumer(cfg *config.Config, fileSvc service.FileService, folderSvc service.FolderService) (consumer.Consumer, error) {
	userConsumer, err := consumer.NewConsumer(cfg.Policy.Users.BootStrapServers, cfg.Policy.Users.GroupId, &users.UserEventHandler{
		FileSvc:   fileSvc,
		FolderSvc: folderSvc,
	})

	go func() {
		if err := userConsumer.Run(); err != nil {
			log.Fatalf("failed to run Kafka consumer: %v", err)
		}
	}()

	return userConsumer, err
}

func initRepository(db *gorm.DB) (folderRepo repository.FolderRepository, fileRepo repository.FileRepository) {
	folderRepo = repository.NewFolderRepository(db)
	fileRepo = repository.NewFileRepository(db)
	return
}

func initService(folderRepo repository.FolderRepository, fileRepo repository.FileRepository) (folderSvc service.FolderService, fileSvc service.FileService) {
	folderSvc = service.NewFolderService(folderRepo)
	fileSvc = service.NewFileService(fileRepo)
	return
}

func initGraphqlServer(folderSvc service.FolderService, fileSvc service.FileService) *handler.Server {
	resolver := resolver.NewResolver(folderSvc, fileSvc)
	schema := graph.NewExecutableSchema(graph.Config{Resolvers: resolver})
	server := handler.NewDefaultServer(schema)
	server.AroundResponses(func(ctx context.Context, next graphql.ResponseHandler) *graphql.Response {
		res := next(ctx)
		if len(res.Errors) > 0 {
			op := graphql.GetOperationContext(ctx)
			rawQuery := normalizeQuery(op.RawQuery)
			logger.Printf("Failed operation %s with query: %s", op.OperationName, rawQuery)
		}
		return res
	})

	server.SetQueryCache(lru.New(1000))

	return server
}

func startServer(server *handler.Server, port string) {
	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", server)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func normalizeQuery(query string) string {
	return regexp.MustCompile(`\s+`).ReplaceAllString(query, " ")
}
