package models

import (
	"github.com/google/uuid"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

const EXAMPLES_PER_PAGE = 1

type Example struct {
	gorm.Model
	UUID uuid.UUID `gorm:"unique;type:uuid;default:gen_random_uuid()"`

	Title string
	Tags  pq.StringArray

	Images     pq.StringArray
	Thumbnails pq.StringArray
}

func NewExample(
	uuid uuid.UUID,
	title string,
	tags pq.StringArray,
	images pq.StringArray,
	thumbnails pq.StringArray,
) *Example {
	return &Example{
		UUID:       uuid,
		Title:      title,
		Tags:       tags,
		Images:     images,
		Thumbnails: thumbnails,
	}
}
