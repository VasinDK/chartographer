package utils

import (
	"image"
	"io"
	"this_module/pkg/random"

	"golang.org/x/image/bmp"
)

type Utils struct{}

func (u *Utils) GetRandLen(aliasLangth int) string {
	return random.GetRandomByLength(aliasLangth)
}

func (u *Utils) EncodeBMP(file io.Writer, img image.Image) error {
	return bmp.Encode(file, img)
}

func (u *Utils) DecodeBMP(file io.Reader) (image.Image, error) {
	return bmp.Decode(file)
}

func New() *Utils {
	return &Utils{}
}
