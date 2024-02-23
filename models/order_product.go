package models

import "gorm.io/gorm"

type OrderProduct struct {
	gorm.Model

	ID            uint
	OrderID       uint
	ProductId     uint
	Slug          string
	Name          string
	Price         int
	DiscountPrice int
	Quantity      int
}

func NewOrderProduct(
	product_id uint,
	order_id uint,
	slug string,
	name string,
	price int,
	discount_price int,
	quantity int,
) *OrderProduct {
	return &OrderProduct{
		ProductId:     product_id,
		OrderID:       order_id,
		Slug:          slug,
		Name:          name,
		Price:         price,
		DiscountPrice: discount_price,
		Quantity:      quantity,
	}
}

func (op *OrderProduct) GetTotalPrice() int {
	p := op.Price
	if op.DiscountPrice != -1 {
		p = op.DiscountPrice
	}
	return p * op.Quantity
}
