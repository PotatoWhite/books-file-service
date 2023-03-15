package util

import (
	"github.com/potatowhite/books/file-service/graph/model"
	"github.com/potatowhite/books/file-service/pkg/repository/entity"
)

func ToFolderDto(folder *entity.Folder) *model.Folder {
	parentID := UItoAOrNil(folder.ParentId)

	return &model.Folder{
		ID:       *UItoAOrNil(&folder.ID),
		Name:     folder.Name,
		ParentID: parentID,
		UserID:   *UItoAOrNil(folder.UserId),
		Path:     &folder.Path,
	}
}

func ToFileDto(file *entity.File) *model.File {
	return &model.File{
		ID:        *UItoAOrNil(&file.ID),
		Name:      file.Name,
		FolderID:  *UItoAOrNil(&file.FolderId),
		Type:      &file.Type,
		Extension: &file.Extension,
		Size:      &file.Size,
		Modified:  &file.Modified,
		UserID:    *UItoAOrNil(file.UserId),
		Path:      &file.Path,
	}
}
