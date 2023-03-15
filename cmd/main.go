package main

import (
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/potatowhite/books/file-service/cmd/config"
	"github.com/potatowhite/books/file-service/db"
	"github.com/potatowhite/books/file-service/graph"
	"github.com/potatowhite/books/file-service/pkg/repository"
	"github.com/potatowhite/books/file-service/pkg/resolver"
	"github.com/potatowhite/books/file-service/pkg/service"
	"log"
	"net/http"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	db, err := db.InitDB(cfg)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	folderRepo := repository.NewFolderRepository(db)
	folderSvc := service.NewFolderService(folderRepo)
	server := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: resolver.NewResolver(folderSvc)}))

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", server)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", cfg.Server.Port)
	log.Fatal(http.ListenAndServe(":"+cfg.Server.Port, nil))
}
