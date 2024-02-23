package models

import (
	"gorm.io/gorm"
)

type Category struct {
	gorm.Model

	ID   uint
	Name string
	Slug string `gorm:"unique"`
	Tags string

	ParentId uint
	Parent   *Category `gorm:"-"`

	Children []*Category `gorm:"-"`
	Products []*Product  `gorm:"-"`

	IsEnabled bool
}

const CATEGORIES_PER_PAGE = 1

func (c *Category) BeforeDelete(tx *gorm.DB) error {
	return tx.Model(&Product{}).Where("category_id = ?", c.ID).Update("category_id", 0).Error
}

func (c *Category) AfterFind(tx *gorm.DB) error {
	if err := tx.Model(&Category{}).Where("parent_id = ?", c.ID).Find(&c.Children).Error; err != nil {
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
