package models

import (
	"gorm.io/gorm"
)

type CartProduct struct {
	gorm.Model

	ID        uint
	ProductId uint
	Product   *Product `gorm:"foreignKey:ProductId"`
	CartID    uint

	Slug          string
	Title         string
	Price         int
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
	}
}

func (c *CartProduct) AfterFind(tx *gorm.DB) error {
	err := tx.Model(&Product{}).Where("id = ?", c.ProductId).First(&c.Product).Error
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

	c.Slug = c.Product.Slug
	c.Title = c.Product.Title
	c.Price = c.Product.Price
	c.DiscountPrice = c.Product.DiscountPrice

	if err = tx.Save(c).Error; err != nil {
		return nil
	}

	return nil
}

func (c *CartProduct) TotalPrice() int {
	if c.DiscountPrice != -1 {
		return c.DiscountPrice * c.Quantity
	}

	return c.Price * c.Quantity
}
