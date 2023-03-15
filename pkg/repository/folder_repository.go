package repository

import (
	"fmt"
	"github.com/potatowhite/books/file-service/pkg/repository/entity"
	"gorm.io/gorm"
	"log"
	"os"
)

var (
	logger = log.New(os.Stdout, "file-service: ", log.LstdFlags|log.Lshortfile)
)

func NewFolderRepository(db *gorm.DB) FolderRepository {
	return &folderRepository{
		db: db,
	}
}

type FolderRepository interface {
	CreateRootFolder(userId uint) (*entity.Folder, error)
	CreateFolder(userId uint, name string, parentId uint) (*entity.Folder, error)
	UpdateFolder(userId uint, folder *entity.Folder) error
	DeleteFolder(userId uint, id uint) (bool, error)

	GetRootFolder(userId uint) (*entity.Folder, error)
	GetFolder(userId uint, id uint) (*entity.Folder, error)
	GetChildren(userId uint, id uint) ([]*entity.Folder, error)
	GetFolderByNameAndParentId(userId uint, name string, parentId uint) (*entity.Folder, error)
}

type folderRepository struct {
	db *gorm.DB
}

func (f *folderRepository) GetFolderByNameAndParentId(userId uint, name string, parentId uint) (*entity.Folder, error) {
	var folder entity.Folder
	err := f.db.Where("user_id = ? AND name = ? AND parent_id = ?", userId, name, parentId).First(&folder).Error
	if err != nil {
		return nil, err
	}

	return &folder, nil
}

func (f *folderRepository) DeleteFolder(userId uint, id uint) (bool, error) {
	result := f.db.Where("id = ?", id).Delete(&entity.Folder{})
	if result.Error != nil {
		return false, result.Error
	}

	return result.RowsAffected > 0, nil
}

func (f *folderRepository) UpdateFolder(userId uint, folder *entity.Folder) error {
	return f.db.Model(&entity.Folder{}).Where("id = ?", folder.ID).Updates(folder).Error
}

func (f *folderRepository) GetChildren(userId uint, id uint) ([]*entity.Folder, error) {
	var children []*entity.Folder
	err := f.db.Where("user_id = ? AND parent_id = ?", userId, id).Find(&children).Error
	if err != nil {
		return nil, err
	}

	parent, err := f.GetFolder(userId, id)
	if err != nil {
		return nil, err
	}

	parentPath := f.GetPathCTE(parent)
	for _, child := range children {
		child.Path = parentPath + "/" + child.Name
	}

	return children, nil
}

func (f *folderRepository) GetFolder(userId uint, id uint) (*entity.Folder, error) {
	var folder entity.Folder
	err := f.db.Where("id = ?", id).First(&folder).Error
	if err != nil {
		return nil, err
	}

	folder.Path = f.GetPathCTE(&folder)

	return &folder, nil
}

func (f *folderRepository) GetRootFolder(userId uint) (*entity.Folder, error) {
	var rootFolder entity.Folder
	err := f.db.Where("user_id = ? AND parent_id IS NULL", userId).First(&rootFolder).Error
	if err != nil {
		return nil, err
	}

	rootFolder.Path = f.GetPathCTE(&rootFolder)

	return &rootFolder, nil
}

// get the path of a folder(cte version) by traversing the parent folders
func (f *folderRepository) GetPathCTE(folder *entity.Folder) string {
	var path string
	err := f.db.Raw("WITH RECURSIVE cte AS ( SELECT id, name, parent_id, name AS full_path FROM public.folders WHERE id = ? and user_id = ? UNION ALL SELECT f.id, f.name, f.parent_id,  f.name || '/' || cte.full_path FROM public.folders f JOIN cte ON cte.parent_id = f.id ) SELECT full_path FROM cte WHERE parent_id is null", folder.ID, folder.UserId).Scan(&path).Error

	if err != nil {
		logger.Println(fmt.Sprintf("failed to get path of folder %d: %v", folder.ID, err))
		return ""
	}
	return path
}

func (f *folderRepository) CreateRootFolder(userId uint) (*entity.Folder, error) {
	rootFolder := entity.Folder{
		Name:   "",
		UserId: &userId,
	}

	err := f.db.Create(&rootFolder).Error
	if err != nil {
		return nil, err
	}

	rootFolder.Path = f.GetPathCTE(&rootFolder)

	return &rootFolder, nil
}

func (f *folderRepository) CreateFolder(userId uint, name string, parentId uint) (*entity.Folder, error) {
	folder := entity.Folder{
		Name:     name,
		ParentId: &parentId,
		UserId:   &userId,
	}

	err := f.db.Create(&folder).Error
	if err != nil {
		return nil, err
	}

	folder.Path = f.GetPathCTE(&folder)

	return &folder, nil
}
