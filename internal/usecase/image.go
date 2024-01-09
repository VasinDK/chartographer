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

type Repository interface {
	GetAllFile(string) ([]string, error)
	AddFileInList(string, string)
	DelFileInList(string)
	LockFile(string)
	UnlockFile(string)
	RLockFile(string)
	RUnlockFile(string)
}

type ImageUC struct {
	Utils Utiler
	Repo  Repository
}

const (
	LenIdImg      = 15
	AddressImages = "./docs/img/"
	FileExtension = ".bmp"
	HeightMax     = 50000
	WidthMax      = 20000
)

// New создает экземпляр usecase
func New(Utils Utiler, Repo Repository) *ImageUC {
	return &ImageUC{
		Utils,
		Repo,
	}
}

// Create создает изображение
func (I *ImageUC) Create(width, height int) (string, error) {
	if width <= 0 || width > WidthMax || height <= 0 || height > HeightMax {
		return "", fmt.Errorf("The image dimensions exceed the allowed ones")
	}

	img := image.NewRGBA(image.Rect(0, 0, width, height))
	black := color.RGBA{0, 0, 0, 1}
	draw.Draw(img, img.Bounds(), &image.Uniform{black}, image.Point{0, 0}, draw.Src)

	id := I.Utils.GetRandLen(LenIdImg)

	file, err := os.Create(AddressImages + id + FileExtension)
	if err != nil {
		return "", err
	}

	defer file.Close()

	err = I.Utils.EncodeBMP(file, img)
	if err != nil {
		return "", err
	}

	I.Repo.AddFileInList(id, FileExtension)

	return id, nil
}

// Add добавляет изображение на имеющееся
func (I *ImageUC) Add(chart *entity.Chart) error {
	path := AddressImages + chart.IdParent + FileExtension

	I.Repo.LockFile(chart.IdParent)
	defer I.Repo.UnlockFile(chart.IdParent)

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

	tempPath := AddressImages + "tempFile" + FileExtension

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
	I.Repo.RLockFile(chart.IdParent)
	defer I.Repo.RUnlockFile(chart.IdParent)

	dstFile, err := os.Open(AddressImages + chart.IdParent + FileExtension)
	if err != nil {
		return nil, err
	}

	dst, err := I.Utils.DecodeBMP(dstFile)
	if err != nil {
		return nil, err
	}
	dstFile.Close()

	// I.Repo.RUnlockFile(chart.IdParent)

	r := image.Rect(chart.X, chart.Y, chart.Width, chart.Height)
	newRGBA := image.NewRGBA(r)

	draw.Draw(newRGBA, r, dst, image.Point{}, draw.Src)

	return newRGBA, nil
}

// Delete удаляет изображение с идентификатором
func (I *ImageUC) Delete(id string) error {
	I.Repo.LockFile(id)
	defer func() {
		I.Repo.UnlockFile(id)
		I.Repo.DelFileInList(id)
	}()

	err := os.Remove(AddressImages + id + FileExtension)
	if err != nil {
		return err
	}

	return nil
}
