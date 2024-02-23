package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Cart struct {
	gorm.Model

	ID       uint
	UUID     uuid.UUID `gorm:"unique;type:uuid;default:gen_random_uuid()"`
	Products []*CartProduct
}

func (c *Cart) AfterFind(tx *gorm.DB) error {
	if err := tx.Model(&CartProduct{}).Where("cart_id = ?", c.ID).Find(&c.Products).Error; err != nil {
		return err
	}
	return nil
}

func NewCart() *Cart {
	return &Cart{}
}

func (c *Cart) GetTotalPrice() int {
	total := 0
	for _, product := range c.Products {
		p := product.Price
		if product.DiscountPrice != -1 {
			p = product.DiscountPrice
		}
		total += p * product.Quantity
	}
	return total
}
