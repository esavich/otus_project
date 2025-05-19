package resizer

import (
	"image"

	"github.com/disintegration/imaging"
)

type Resizer struct{}

func NewResizer() *Resizer {
	return &Resizer{}
}

func (*Resizer) ResizeImg(img image.Image, w int, p int) image.Image {
	resized := imaging.Fill(img, w, p, imaging.Center, imaging.Lanczos)

	return resized
}
