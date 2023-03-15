package service

import (
	"fmt"
	"github.com/potatowhite/books/file-service/pkg/repository"
	"github.com/potatowhite/books/file-service/pkg/repository/entity"
)

func NewFileService(repo repository.FileRepository) FileService {
	return &fileService{repo: repo}
}

type FileService interface {
	CreateFile(userId uint, name string, folderId uint) (*entity.File, error)
	PatchFile(userId uint, id uint, name *string, fileType *string, fileExtension *string, size *uint64) (*entity.File, error)

	GetFile(userId uint, id uint) (*entity.File, error)
	GetChildren(userId uint, folderId uint) ([]*entity.File, error)
	DeleteFile(userId uint, id uint) (bool, error)
}

type fileService struct {
	repo repository.FileRepository
}

func (f *fileService) DeleteFile(userId uint, id uint) (bool, error) {
	return f.repo.DeleteFile(userId, id)
}

func (f *fileService) PatchFile(userId uint, id uint, name *string, fileType *string, fileExtension *string, size *uint64) (*entity.File, error) {
	file, err := f.repo.GetFile(userId, id)
	if err != nil {
		return nil, err
	} else if file == nil {
		return nil, fmt.Errorf("file with id %v not found", id)
	}

	updateField(&file.Name, name)
	updateField(&file.Type, fileType)
	updateField(&file.Extension, fileExtension)
	updateSize(&file.Size, size)

	if err = f.repo.UpdateFile(userId, file); err != nil {
		return nil, err
	}

	return file, nil
}

func (f *fileService) GetChildren(userId uint, folderId uint) ([]*entity.File, error) {
	return f.repo.GetFilesByFolderId(userId, folderId)
}

func (f *fileService) GetFile(userId uint, id uint) (*entity.File, error) {
	return f.repo.GetFile(userId, id)
}

func (f *fileService) CreateFile(userId uint, name string, folderId uint) (*entity.File, error) {
	// unique name in folder
	file, err := f.repo.GetFileByNameAndFolderId(userId, name, folderId)
	if err != nil {
		return nil, err
	} else if file != nil {
		return nil, fmt.Errorf("file with name %v already exists in folder %v", name, folderId)
	}

	return f.repo.CreateFile(userId, name, folderId)
}

func updateField(field *string, value *string) {
	if value != nil {
		*field = *value
	}
}

func updateSize(field *uint64, value *uint64) {
	if value != nil {
		*field = *value
	}
}
