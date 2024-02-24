package utils

import (
	"io"

	"github.com/h2non/bimg"
)

const MAX_SIZE_W = 1920
const MAX_SIZE_H = 1080

const THUMBNAIL_SIZE = 300

// returns image, thumbnail, error
func ProcessImage(file io.Reader) ([]byte, []byte, error) {
	buffer, err := io.ReadAll(file)
	if err != nil {
		return []byte{}, []byte{}, err
	}

	webp, err := bimg.NewImage(buffer).Convert(bimg.WEBP)
	if err != nil {
		return []byte{}, []byte{}, err
	}

	processed, err := bimg.NewImage(webp).Process(bimg.Options{Quality: bimg.Quality})
	if err != nil {
		return []byte{}, []byte{}, err
	}

	b := bimg.NewImage(processed)
	if size, err := b.Size(); err == nil && size.Width > MAX_SIZE_W {
		aspect_ratio := float64(size.Height) / float64(size.Width)
		new_height := int(float64(MAX_SIZE_W) * aspect_ratio)

		processed, err = bimg.NewImage(processed).Resize(MAX_SIZE_W, new_height)
		if err != nil {
			return []byte{}, []byte{}, err
		}
	}

	thumbnail, err := bimg.NewImage(processed).Thumbnail(THUMBNAIL_SIZE)
	if err != nil {
		return []byte{}, []byte{}, err
	}

	return processed, thumbnail, nil
}
