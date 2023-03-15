package entity

import (
	"gorm.io/gorm"
)

type Folder struct {
	gorm.Model

	Name     string  `json:"name" gorm:"not null"`
	ParentId *uint   `json:"parentId" gorm:"index"`
	Parent   *Folder `json:"parent,omitempty"`
	UserId   *uint   `json:"userId" gorm:"not null;index"`
	Path     string  `json:"path" gorm:"-"`
}

type File struct {
	gorm.Model

	Name      string  `json:"name" gorm:"not null""`
	FolderId  uint    `json:"folderId" gorm:"not null;index"`
	Folder    *Folder `json:"folder"`
	Type      string  `json:"type"`
	Extension string  `json:"extension"`
	Size      int     `json:"size"`
	Modified  string  `json:"modified"`
	UserId    *uint   `json:"userId" gorm:"not null;index"`
	Path      string  `json:"path" gorm:"-"`
}
