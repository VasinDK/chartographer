package usecase

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	_ "image/jpeg"
	"io"
	"net/http"
	"os"
	"this_module/internal/entity"
)

type Utiler interface {
	GetRandLen(int) string
	EncodeBMP(io.Writer, image.Image) error
	DecodeBMP(io.Reader) (image.Image, error)
	NewChart(*http.Request) (*entity.Chart, error)
}

type ImageUC struct {
	Utils Utiler
}

const (
	lenIdImg      = 15
	addressImages = "./docs/img/"
	fileExtension = ".bmp"
	heightMax     = 50000
	widthMax      = 20000
)

// New создает экземпляр usecase
func New(Utils Utiler) *ImageUC {
	return &ImageUC{
		Utils,
	}
}

// Create создает изображение
func (I *ImageUC) Create(width, height int) (string, error) {
	if width <= 0 || width > widthMax || height <= 0 || height > heightMax {
		return "", fmt.Errorf("The image dimensions exceed the allowed ones")
	}

	img := image.NewRGBA(image.Rect(0, 0, width, height))
	black := color.RGBA{0, 0, 0, 1}
	draw.Draw(img, img.Bounds(), &image.Uniform{black}, image.Point{0, 0}, draw.Src)

	id := "I.Utils.GetRandLen(lenIdImg)" // временно. Разкоментить

	file, err := os.Create(addressImages + id + fileExtension)
	if err != nil {
		return "", err
	}

	defer file.Close()

	err = I.Utils.EncodeBMP(file, img)
	if err != nil {
		return "", err
	}

	return id, nil
}

// Add добавляет изображение на имеющееся
func (I *ImageUC) Add(chart *entity.Chart) error {
	I.Create(400, 400) // временно. Удалить
	path := addressImages + chart.IdParent + fileExtension

	dstFile, err := os.Open(path)
	if err != nil {
		return err
	}

	dst, err := I.Utils.DecodeBMP(dstFile)
	if err != nil {
		return err
	}
	dstFile.Close()

	newRGBA := dst.(draw.Image)

	src, err := I.Utils.DecodeBMP(chart.File)
	if err != nil {
		return err
	}

	rectWidth := newRGBA.Bounds().Dx() - chart.X
	rectHeight := newRGBA.Bounds().Dy() - chart.Y

	if rectWidth >= chart.Width {
		rectWidth = chart.Width
	}

	rectWidth += chart.X

	if rectHeight >= chart.Height {
		rectHeight = chart.Height
	}

	rectHeight += chart.Y

	r := image.Rectangle{
		image.Point{
			chart.X,
			chart.Y,
		},

		image.Point{
			rectWidth,
			rectHeight,
		},
	}

	draw.Draw(newRGBA, r, src, src.Bounds().Min, draw.Src)

	tempPath := addressImages + "tempFile" + fileExtension

	tempFile, err := os.Create(tempPath)
	if err != nil {
		return err
	}

	err = I.Utils.EncodeBMP(tempFile, newRGBA)
	if err != nil {
		return err
	}
	tempFile.Close()

	err = os.Remove(path)
	if err != nil {
		return err
	}

	err = os.Rename(tempPath, path)

	return nil
}

// Part возвращает часть изображения
func (I *ImageUC) Part(chart *entity.Chart) (*image.RGBA, error) {
	dstFile, err := os.Open(addressImages + chart.IdParent + fileExtension)
	if err != nil {
		return nil, err
	}

	dst, err := I.Utils.DecodeBMP(dstFile)
	if err != nil {
		return nil, err
	}
	dstFile.Close()

	r := image.Rect(chart.X, chart.Y, chart.Width, chart.Height)
	newRGBA := image.NewRGBA(r)

	draw.Draw(newRGBA, r, dst, image.Point{}, draw.Src)

	return newRGBA, nil
}

// Delete удаляет изображение с идентификатором
func (I *ImageUC) Delete(id string) error {
	err := os.Remove(addressImages + id + fileExtension)
	if err != nil {
		return err
	}

	return nil
}
