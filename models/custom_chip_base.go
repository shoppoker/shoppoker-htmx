package models

import (
	"github.com/w1png/go-htmx-ecommerce-template/file_storage"
	"gorm.io/gorm"
)

const CUSTOM_CHIP_BASES_PER_PAGE = 1

type CustomChipBase struct {
	gorm.Model

	ID uint

	Title        string
	Slug         string `gorm:"unique"`
	Description  string
	StickerScale float64

	Price         int
	DiscountPrice int

	PriceFoil         int
	DiscountPriceFoil int

	VectorImageString string
	VectorImage       file_storage.ObjectStorageId `gorm:"-"`

	IsEnabled bool
}

func (c *CustomChipBase) AfterFind(tx *gorm.DB) error {
	c.VectorImage = file_storage.ObjectStorageId(c.VectorImageString)
	return nil
}

func NewCustomChipBase(
	title string,
	slug string,
	description string,
	sticker_scale float64,
	vector_image file_storage.ObjectStorageId,
	is_enabled bool,
) *CustomChipBase {
	return &CustomChipBase{
		Title:             title,
		Slug:              slug,
		Description:       description,
		StickerScale:      sticker_scale,
		VectorImageString: string(vector_image),
		VectorImage:       vector_image,
		IsEnabled:         is_enabled,
	}
}
