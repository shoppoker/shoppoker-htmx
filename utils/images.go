package utils

import (
	"io"

	"github.com/h2non/bimg"
)

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

	thumbnail, err := bimg.NewImage(processed).Thumbnail(100)
	if err != nil {
		return []byte{}, []byte{}, err
	}

	return processed, thumbnail, nil
}
