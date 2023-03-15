package resolver

import "github.com/potatowhite/books/file-service/pkg/service"

type Resolver struct {
	FolderSvc service.FolderService
}

func NewResolver(folderSvc service.FolderService) *Resolver {
	return &Resolver{FolderSvc: folderSvc}
}
