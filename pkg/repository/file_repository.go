package repository

import (
	"github.com/potatowhite/books/file-service/pkg/repository/entity"
	"gorm.io/gorm"
	"log"
	"os"
)

var (
	logger = log.New(os.Stdout, "", log.LstdFlags|log.Lshortfile)
)

func NewFileRepository(db *gorm.DB) FileRepository {
	return &fileRepository{db: db}
}

type FileRepository interface {
	CreateFile(userId uint, name string, folderId uint) (*entity.File, error)
	UpdateFile(userId uint, file *entity.File) error
	DeleteFile(userId uint, id uint) (bool, error)
	GetFile(userId uint, id uint) (*entity.File, error)
	GetFileByNameAndFolderId(userId uint, name string, folderId uint) (*entity.File, error)
	GetFilesByFolderId(userId uint, folderId uint) ([]*entity.File, error)
	DeleteAllFiles(userId uint) (int64, error)
}
type fileRepository struct {
	db *gorm.DB
}

func (f *fileRepository) DeleteAllFiles(userId uint) (int64, error) {
	tx := f.db.Where("user_id = ?", userId).Delete(&entity.File{})

	if tx.Error != nil {
		return -1, tx.Error
	}

	if tx.RowsAffected == 0 {
		return -1, nil
	}

	return tx.RowsAffected, nil
}

func (f *fileRepository) CreateFile(userId uint, name string, folderId uint) (*entity.File, error) {
	create := &entity.File{
		Name:     name,
		FolderId: folderId,
		UserId:   userId,
	}

	if err := f.db.Create(create).Error; err != nil {
		return nil, err
	}

	return create, nil
}

func (f *fileRepository) UpdateFile(userId uint, file *entity.File) error {
	err := f.db.Save(file).Error
	if err != nil {
		return err
	}

	return nil
}

func (f *fileRepository) DeleteFile(userId uint, id uint) (bool, error) {
	tx := f.db.Where("user_id = ? AND id = ?", userId, id).Delete(&entity.File{})

	if tx.Error != nil {
		return false, tx.Error
	}

	if tx.RowsAffected == 0 {
		return false, nil
	}

	return true, nil
}

func (f *fileRepository) GetFile(userId uint, id uint) (*entity.File, error) {
	var file entity.File

	tx := f.db.Where("user_id = ? AND id = ?", userId, id).First(&file)

	if tx.Error != nil {
		if tx.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, tx.Error
	}

	return &file, nil
}

func (f *fileRepository) GetFileByNameAndFolderId(userId uint, name string, folderId uint) (*entity.File, error) {
	var file entity.File

	tx := f.db.Where("user_id = ? AND name = ? AND folder_id = ?", userId, name, folderId).First(&file)

	if tx.Error != nil {
		if tx.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, tx.Error
	}

	return &file, nil
}

func (f *fileRepository) GetFilesByFolderId(userId uint, folderId uint) ([]*entity.File, error) {
	var files []*entity.File
	f.db.Where("user_id = ? AND folder_id = ?", userId, folderId).Find(&files)
	return files, nil
}
