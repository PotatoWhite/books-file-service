package resolver

import "github.com/potatowhite/books/file-service/pkg/service"

type Resolver struct {
	FolderSvc service.FolderService
	FileSvc   service.FileService
}

func NewResolver(folderSvc service.FolderService, fileSvc service.FileService) *Resolver {
	return &Resolver{FolderSvc: folderSvc, FileSvc: fileSvc}
}
