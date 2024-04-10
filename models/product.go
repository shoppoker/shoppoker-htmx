package models

import (
	"strings"

	"github.com/lib/pq"
	"gorm.io/gorm"
)

type StockType int

var STOCK_TYPES_ARRAY = []StockType{
	StockTypeInStock,
	StockTypeOutOfStock,
	StockTypeOrder,
}

const PRODUCTS_PER_PAGE = 20000

const (
	StockTypeInStock StockType = iota
	StockTypeOutOfStock
	StockTypeOrder
)

func (s StockType) ToString() string {
	switch s {
	case StockTypeInStock:
		return "В наличии"
	case StockTypeOutOfStock:
		return "Нет в наличии"
	case StockTypeOrder:
		return "Под заказ"
	default:
		return ""
	}
}

func (s StockType) Color() string {
	switch s {
	case StockTypeInStock:
		return "text-in-stock"
	case StockTypeOutOfStock:
		return "text-out-of-stock"
	case StockTypeOrder:
		return "text-order"
	default:
		return ""
	}
}

type Product struct {
	gorm.Model

	ID          uint
	Slug        string `gorm:"unique"`
	Title       string
	Description string

	Price         int
	DiscountPrice int
	StockType     StockType
	Tags          string
	Priority      int

	CategoryId uint
	Category   Category `gorm:"-"`

	Images     pq.StringArray `gorm:"type:text[]"`
	Thumbnails pq.StringArray `gorm:"type:text[]"`

	IsEnabled  bool
	IsFeatured bool
}

func NewProduct(
	slug string,
	title string,
	description string,
	price int,
	discount_price int,
	stock_type StockType,
	tags string,
	priority int,
	category_id uint,
	images []string,
	thumbnails []string,
	is_enabled bool,
	is_featured bool,
) *Product {
	dashed_slug := strings.ReplaceAll(slug, " ", "-")

	return &Product{
		Slug:          dashed_slug,
		Title:         title,
		Description:   description,
		Price:         price,
		DiscountPrice: discount_price,
		StockType:     stock_type,
		Tags:          tags,
		Priority:      priority,
		CategoryId:    category_id,
		Images:        images,
		Thumbnails:    thumbnails,
		IsEnabled:     is_enabled,
		IsFeatured:    is_featured,
	}
}

func (product *Product) AfterFind(tx *gorm.DB) error {
	if err := tx.Where("id = ?", product.CategoryId).First(&product.Category).Error; err != nil && err != gorm.ErrRecordNotFound {
		return err
	}

	return nil
}
