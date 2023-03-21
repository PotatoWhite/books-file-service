package main

import (
	"context"
	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/potatowhite/books/file-service/cmd/config"
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
	"os/signal"
	"regexp"
	"syscall"
)

var (
	logger = log.New(os.Stdout, "", log.LstdFlags|log.Lshortfile)
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	database, err := db.InitDB(cfg)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	defer db.CloseDB(database)

	folderRepo, fileRepo := initRepository(database)
	folderSvc, fileSvc := initService(folderRepo, fileRepo)
	server := initGraphqlServer(folderSvc, fileSvc)
	userConsumer, err := initConsumer(cfg, fileSvc, folderSvc)
	if err != nil {
		log.Fatalf("failed to create users: %v", err)
	}

	userConsumer.Run()
	defer userConsumer.Close()

	startServer(server, cfg.Server.Port)

}

func initConsumer(cfg *config.Config, fileSvc service.FileService, folderSvc service.FolderService) (consumer.Consumer, error) {
	userConsumer, err := consumer.NewConsumer(cfg.Policy.Users.BootStrapServers, cfg.Policy.Users.GroupId, &users.UserEventHandler{
		FileSvc:   fileSvc,
		FolderSvc: folderSvc,
	})

	go func() {
		if err := userConsumer.Run(); err != nil {
			log.Fatalf("failed to run Kafka consumer: %v", err)
		}
	}()

	// Listen for a SIGINT or SIGTERM signal and close the consumer when received
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
	<-signals

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
