package resizer

import (
	"image"

	"github.com/disintegration/imaging"
)

func ResizeImg(img image.Image, w int, p int) (image.Image, error) {
	resized := imaging.Fill(img, w, p, imaging.Center, imaging.Lanczos)

	return resized, nil
}
