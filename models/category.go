package models

import (
	"github.com/w1png/go-htmx-ecommerce-template/file_storage"
	"gorm.io/gorm"
)

type Category struct {
	gorm.Model

	ID   uint
	Name string
	Slug string `gorm:"unique"`
	Tags string

	Image     file_storage.ObjectStorageId
	Thumbnail file_storage.ObjectStorageId

	ParentId uint
	Parent   *Category `gorm:"-"`

	Children []*Category `gorm:"-"`
	Products []*Product  `gorm:"-"`

	IsEnabled bool
}

const CATEGORIES_PER_PAGE = 10

func (c *Category) BeforeDelete(tx *gorm.DB) error {
	return tx.Model(&Product{}).Where("category_id = ?", c.ID).Update("category_id", 0).Error
}

func (c *Category) AfterFind(tx *gorm.DB) error {
	if err := tx.Model(&Category{}).Where("parent_id = ? and is_enabled = ?", c.ID, true).Find(&c.Children).Error; err != nil {
		return err
	}

	return nil
}

func NewCategory(name, slug, tags string, parent_id uint) *Category {
	return &Category{
		Name:      name,
		Slug:      slug,
		Tags:      tags,
		ParentId:  parent_id,
		IsEnabled: true,
	}
}
