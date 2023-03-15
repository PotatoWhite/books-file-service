package service

import (
	"fmt"
	"github.com/potatowhite/books/file-service/pkg/repository"
	"github.com/potatowhite/books/file-service/pkg/repository/entity"
	"gorm.io/gorm"
	"log"
	"os"
)

var (
	// log with method name and line number
	logger = log.New(os.Stdout, "", log.LstdFlags|log.Lshortfile)
)

type FolderService interface {
	CreateFolder(userId uint, name string, parentId uint) (*entity.Folder, error)
	RenameFolder(userId uint, id uint, newName string) (*entity.Folder, error)
	DeleteFolder(userId uint, id uint) (bool, error)
	GetFolder(userId uint, id uint) (*entity.Folder, error)
	GetChildren(userId uint, parentID uint) ([]*entity.Folder, error)
	CreateRootFolder(userId uint) (*entity.Folder, error)
	GetRootFolder(userId uint) (*entity.Folder, error)
}

type folderService struct {
	repo repository.FolderRepository
}

func (f *folderService) GetRootFolder(userId uint) (*entity.Folder, error) {
	return f.repo.GetRootFolder(userId)
}

func (f *folderService) RenameFolder(userId uint, id uint, newName string) (*entity.Folder, error) {
	folder, err := f.repo.GetFolder(userId, id)
	if err != nil {
		return nil, err
	}

	folder.Name = newName
	err = f.repo.UpdateFolder(userId, folder)
	if err != nil {
		return nil, err
	}

	return folder, nil
}

func (f *folderService) CreateRootFolder(userId uint) (*entity.Folder, error) {
	// must not exist a root folder for the user
	folder, err := f.repo.GetRootFolder(userId)
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	if folder != nil {
		logger.Printf("Root folder already exists for user %v", userId)
		return nil, fmt.Errorf("root folder already exists for user %v", userId)
	}

	return f.repo.CreateRootFolder(userId)
}

func (f *folderService) GetFolder(userId uint, id uint) (*entity.Folder, error) {
	return f.repo.GetFolder(userId, id)
}

func (f *folderService) GetChildren(userId uint, parentID uint) ([]*entity.Folder, error) {
	return f.repo.GetChildren(userId, parentID)
}

func (f *folderService) CreateFolder(userId uint, name string, parentId uint) (*entity.Folder, error) {
	// must not exist a folder with the same name for the user
	folder, err := f.repo.GetFolder(userId, parentId)
	if err != nil {
		return nil, err
	}

	// check if the folder already exists
	_, err = f.repo.GetFolderByNameAndParentId(userId, name, folder.ID)
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	} else if err == gorm.ErrRecordNotFound {
		return f.repo.CreateFolder(userId, name, parentId)
	}

	return nil, fmt.Errorf("folder with name %v already exists", name)
}

func (f folderService) DeleteFolder(userId uint, id uint) (bool, error) {
	return f.repo.DeleteFolder(userId, id)
}

func NewFolderService(folderRepo repository.FolderRepository) FolderService {
	return &folderService{
		repo: folderRepo,
	}
}
