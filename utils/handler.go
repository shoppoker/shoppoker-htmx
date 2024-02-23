package utils

import (
	"math"
)

func GetOffsetAndLimit(page, limit int) (int, int) {
	return limit * (page - 1), limit
}

func GetNextPage(page int, f func() (int, error), limit int) (int, error) {
	count, err := f()
	if err != nil {
		return -1, err
	}

	total_pages := int(math.Ceil(float64(count) / float64(limit)))
	if total_pages <= page {
		return -1, nil
	}

	return page + 1, nil
}
