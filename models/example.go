package models

import (
	"github.com/google/uuid"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

const EXAMPLES_PER_PAGE = 10

type Example struct {
	gorm.Model
	UUID uuid.UUID `gorm:"unique;type:uuid;default:gen_random_uuid()"`

	Title string
	Tags  pq.StringArray `gorm:"type:text[]"`

	CustomChipBaseID uint
	CustomChipBase   CustomChipBase `gorm:"foreignKey:CustomChipBaseID"`

	Images     pq.StringArray `gorm:"type:text[]"`
	Thumbnails pq.StringArray `gorm:"type:text[]"`
}

func NewExample(
	title string,
	tags pq.StringArray,
	images pq.StringArray,
	thumbnails pq.StringArray,
	custom_chip_base_id uint,
) *Example {
	return &Example{
		Title:            title,
		Tags:             tags,
		Images:           images,
		Thumbnails:       thumbnails,
		CustomChipBaseID: custom_chip_base_id,
	}
}
