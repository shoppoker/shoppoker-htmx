package models

import (
	"time"

	"github.com/a-h/templ"
	"github.com/w1png/go-htmx-ecommerce-template/templates/components/icons"
	"gorm.io/gorm"
)

const ORDERS_PER_PAGE = 20

type DeliveryType int

const (
	DeliveryTypePickup DeliveryType = iota
	DeliveryTypeDelivery
)

var DELIVERY_TYPES_ARRAY = []DeliveryType{
	DeliveryTypePickup,
	DeliveryTypeDelivery,
}

func (d DeliveryType) ToString() string {
	switch d {
	case DeliveryTypePickup:
		return "Самовывоз"
	case DeliveryTypeDelivery:
		return "Доставка"
	default:
		return "Неизвестно"
	}
}

func (d DeliveryType) GetColorClass() string {
	switch d {
	case DeliveryTypePickup:
		return "bg-blue-200"
	case DeliveryTypeDelivery:
		return "bg-green-200"
	default:
		return "bg-gray-200"
	}
}

func (d DeliveryType) GetIconImg(class string) templ.Component {
	switch d {
	case DeliveryTypePickup:
		return icons.Pickup(class)
	case DeliveryTypeDelivery:
		return icons.Truck(class)
	default:
		return icons.Question(class)
	}
}

type OrderStatus int

const (
	OrderStatusAny OrderStatus = -1
	OrderStatusNew OrderStatus = iota
	OrderStatusProcessing
	OrderStatusDelivery
	OrderStatusCompleted
	OrderStatusCancelled
)

var ORDER_STATUSES_ARRAY = []OrderStatus{
	OrderStatusNew,
	OrderStatusProcessing,
	OrderStatusDelivery,
	OrderStatusCompleted,
	OrderStatusCancelled,
}

func QueryParamToOrderStatus(query_param string) OrderStatus {
	switch query_param {
	case "all":
		return OrderStatusAny
	case "new":
		return OrderStatusNew
	case "processing":
		return OrderStatusProcessing
	case "delivery":
		return OrderStatusDelivery
	case "completed":
		return OrderStatusCompleted
	case "cancelled":
		return OrderStatusCancelled
	default:
		return OrderStatusAny
	}
}

func (o OrderStatus) ToString() string {
	switch o {
	case OrderStatusAny:
		return "Все заказы"
	case OrderStatusNew:
		return "Зарегистрирован"
	case OrderStatusProcessing:
		return "В обработке"
	case OrderStatusDelivery:
		return "В доставке"
	case OrderStatusCompleted:
		return "Завершен"
	case OrderStatusCancelled:
		return "Отменен"
	default:
		return "Неизвестно"
	}
}

func (o OrderStatus) GetColorClass() string {
	switch o {
	case OrderStatusAny:
		return "bg-gray-200"
	case OrderStatusNew:
		return "bg-yellow-200"
	case OrderStatusProcessing:
		return "bg-orange-200"
	case OrderStatusDelivery:
		return "bg-blue-200"
	case OrderStatusCompleted:
		return "bg-green-200"
	case OrderStatusCancelled:
		return "bg-red-200"
	default:
		return "bg-gray-200"
	}
}

type Order struct {
	gorm.Model

	ID           uint
	Name         string
	PhoneNumber  string
	Email        string
	City         string
	Message      string
	Adress       string
	DeliveryType DeliveryType
	Products     []*OrderProduct
	Status       OrderStatus
}

func NewOrder(
	name string,
	phone_number string,
	email string,
	message string,
	delivery_type DeliveryType,
	adress string,
	city string,
) *Order {
	return &Order{
		Name:         name,
		PhoneNumber:  phone_number,
		Email:        email,
		Message:      message,
		DeliveryType: delivery_type,
		Adress:       adress,
		City:         city,
		Status:       OrderStatusNew,
	}
}

func (o *Order) AfterFind(tx *gorm.DB) error {
	if err := tx.Model(&OrderProduct{}).Where("order_id = ?", o.ID).Find(&o.Products).Error; err != nil {
		return err
	}
	return nil
}

func (o *Order) GetTotalPrice() int {
	total := 0
	for _, product := range o.Products {
		p := product.Price
		if product.DiscountPrice != -1 {
			p = product.DiscountPrice
		}
		total += p * product.Quantity
	}

	return total
}

func (o *Order) FormatTime() string {
	return o.CreatedAt.In(time.Local).Format("15:04 02-01-2006")
}
