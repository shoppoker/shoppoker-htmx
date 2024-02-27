package models

import (
	"gorm.io/gorm"
)

type CartProduct struct {
	gorm.Model

	ID        uint
	ProductId uint
	CartID    uint

	Slug          string
	Title         string
	Price         int
	Thumbnail     string
	DiscountPrice int
	Quantity      int
}

func NewCartProduct(
	product_id uint,
	cart_id uint,
	slug string,
	name string,
	price int,
	discount_price int,
	thumbnail string,
	quantity int,
) *CartProduct {
	return &CartProduct{
		ProductId:     product_id,
		CartID:        cart_id,
		Slug:          slug,
		Title:         name,
		Price:         price,
		DiscountPrice: discount_price,
		Quantity:      quantity,
		Thumbnail:     thumbnail,
	}
}

func (c *CartProduct) AfterFind(tx *gorm.DB) error {
	var product Product
	err := tx.Model(&Product{}).Where("id = ?", c.ProductId).First(&product).Error
	if err == nil {
		if err != gorm.ErrRecordNotFound {
			return err
		}
		err := tx.Model(c).Where("id = ?", c.ID).Unscoped().Delete(&CartProduct{}).Error
		if err != nil {
			return err
		}
		return nil
	}

	if product.StockType == StockTypeOutOfStock {
		c.Quantity = 0
	}

	c.Slug = product.Slug
	c.Title = product.Title
	c.Price = product.Price
	c.DiscountPrice = product.DiscountPrice

	if err = tx.Save(c).Error; err != nil {
		return nil
	}

	return nil
}

func (c *CartProduct) TotalPrice() int {
	if c.DiscountPrice != 0 {
		return c.DiscountPrice * c.Quantity
	}

	return c.Price * c.Quantity
}
