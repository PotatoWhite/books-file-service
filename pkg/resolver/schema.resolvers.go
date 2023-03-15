package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.26

import (
	"context"
	"github.com/potatowhite/books/file-service/graph"
	"github.com/potatowhite/books/file-service/graph/model"
	"github.com/potatowhite/books/file-service/pkg/util"
	"log"
	"os"
)

var (
	logger = log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lshortfile)
)

// Mutation returns MutationResolver implementation.
func (r *Resolver) Mutation() graph.MutationResolver { return &mutationResolver{r} }

// Query returns QueryResolver implementation.
func (r *Resolver) Query() graph.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }

type queryResolver struct{ *Resolver }

// CreateRootFolder is the resolver for the createRootFolder field.
func (r *mutationResolver) CreateRootFolder(ctx context.Context, userID string) (*model.Folder, error) {
	userIDInt := *util.AtoUIOrNil(&userID)
	rootFolder, err := r.FolderSvc.CreateRootFolder(userIDInt)
	if err != nil {
		return nil, err
	}

	return util.ToFolderDto(rootFolder), nil
}

// CreateFolder is the resolver for the createFolder field.
func (r *mutationResolver) CreateFolder(ctx context.Context, userID string, name string, parentID string) (*model.Folder, error) {
	userIDInt := *util.AtoUIOrNil(&userID)
	parentIDInt := *util.AtoUIOrNil(&parentID)

	subFolder, err := r.FolderSvc.CreateFolder(userIDInt, name, parentIDInt)
	if err != nil {
		return nil, err
	}

	return util.ToFolderDto(subFolder), nil
}

// RenameFolder is the resolver for the renameFolder field.
func (r *mutationResolver) RenameFolder(ctx context.Context, userID string, id string, name string) (*model.Folder, error) {
	userIDInt := *util.AtoUIOrNil(&userID)
	idInt := *util.AtoUIOrNil(&id)

	folder, err := r.FolderSvc.RenameFolder(userIDInt, idInt, name)
	if err != nil {
		return nil, err
	}

	return util.ToFolderDto(folder), nil
}

// DeleteFolder is the resolver for the deleteFolder field.
func (r *mutationResolver) DeleteFolder(ctx context.Context, userID string, id string) (bool, error) {
	userIDInt := *util.AtoUIOrNil(&userID)
	idInt := *util.AtoUIOrNil(&id)

	_, err := r.FolderSvc.DeleteFolder(userIDInt, idInt)
	if err != nil {
		return false, err
	}

	return true, nil
}

// CreateFile is the resolver for the createFile field.
func (r *mutationResolver) CreateFile(ctx context.Context, userID string, name string, folderID string) (*model.File, error) {
	userIDInt := *util.AtoUIOrNil(&userID)
	folderIDInt := *util.AtoUIOrNil(&folderID)

	file, err := r.FileSvc.CreateFile(userIDInt, name, folderIDInt)
	if err != nil {
		return nil, err
	}

	return util.ToFileDto(file), nil
}

// UpdateFile is the resolver for the updateFile field.
func (r *mutationResolver) UpdateFile(ctx context.Context, userID string, id string, name *string, typeArg *string, extension *string, size *int) (*model.File, error) {
	userIDInt := *util.AtoUIOrNil(&userID)
	idInt := *util.AtoUIOrNil(&id)

	sizeUInt := uint64(*size)
	file, err := r.FileSvc.PatchFile(userIDInt, idInt, name, typeArg, extension, &sizeUInt)
	if err != nil {
		return nil, err
	}

	return util.ToFileDto(file), nil
}

// DeleteFile is the resolver for the deleteFile field.
func (r *mutationResolver) DeleteFile(ctx context.Context, userID string, id string) (bool, error) {
	userIDInt := *util.AtoUIOrNil(&userID)
	idInt := *util.AtoUIOrNil(&id)

	_, err := r.FileSvc.DeleteFile(userIDInt, idInt)
	if err != nil {
		return false, err
	}

	return true, nil
}

// RootFolder is the resolver for the rootFolder field.
func (r *queryResolver) RootFolder(ctx context.Context, userID string) (*model.Folder, error) {
	rootFolder, err := r.FolderSvc.GetRootFolder(*util.AtoUIOrNil(&userID))
	if err != nil {
		return nil, err
	}

	return util.ToFolderDto(rootFolder), nil
}

// ChildrenFolders is the resolver for the childrenFolders field.
func (r *queryResolver) ChildrenFolders(ctx context.Context, userID string, id string) ([]*model.Folder, error) {
	userIDInt := *util.AtoUIOrNil(&userID)
	folderIDInt := *util.AtoUIOrNil(&id)

	folders, err := r.FolderSvc.GetChildren(userIDInt, folderIDInt)
	if err != nil {
		return nil, err
	}

	// Convert to DTO
	foldersDto := make([]*model.Folder, len(folders))
	for i, folder := range folders {
		foldersDto[i] = util.ToFolderDto(folder)
	}

	return foldersDto, nil
}

// ChildrenFiles is the resolver for the childrenFiles field.
func (r *queryResolver) ChildrenFiles(ctx context.Context, userID string, id string) ([]*model.File, error) {
	userIDInt := *util.AtoUIOrNil(&userID)
	folderIDInt := *util.AtoUIOrNil(&id)

	files, err := r.FileSvc.GetChildren(userIDInt, folderIDInt)
	if err != nil {
		return nil, err
	}

	// Convert to DTO
	filesDto := make([]*model.File, len(files))
	for i, file := range files {
		filesDto[i] = util.ToFileDto(file)
	}

	return filesDto, nil
}

// Folder is the resolver for the folder field.
func (r *queryResolver) Folder(ctx context.Context, userID string, id string) (*model.Folder, error) {
	folder, err := r.FolderSvc.GetFolder(*util.AtoUIOrNil(&userID), *util.AtoUIOrNil(&id))
	if err != nil {
		return nil, err
	}

	return util.ToFolderDto(folder), nil
}

// File is the resolver for the file field.
func (r *queryResolver) File(ctx context.Context, userID string, id string) (*model.File, error) {
	file, err := r.FileSvc.GetFile(*util.AtoUIOrNil(&userID), *util.AtoUIOrNil(&id))
	if err != nil {
		return nil, err
	}

	return util.ToFileDto(file), nil
}
