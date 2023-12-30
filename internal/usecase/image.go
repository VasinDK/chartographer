package usecase

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	_ "image/jpeg"
	"io"
	"os"
	"this_module/internal/entity"
)

type Utiler interface {
	GetRandLen(int) string
	EncodeBMP(io.Writer, image.Image) error
	DecodeBMP(io.Reader) (image.Image, error)
}

type ImageUC struct {
	Utils Utiler
}

const lenIdImg = 15
const addressImages = "./docs/img/"
const fileExtension = ".bmp"

func New(Utils Utiler) *ImageUC {
	return &ImageUC{
		Utils,
	}
}

func (I *ImageUC) Create(width, height int) (string, error) {
	if width < 1 || width > 20000 || height < 1 || height > 50000 {
		return "", fmt.Errorf("The image dimensions exceed the allowed ones")
	}

	/* reader, err := os.Open("./docs/img/porfolio.jpg")
	if err != nil {
		fmt.Println(err.Error())
	}

	defer reader.Close() */

	// conf, _, err := image.DecodeConfig(reader)
	// fmt.Println(conf.Height, conf.Width)

	// m, _, err := image.Decode(reader)
	// fmt.Println(m.Bounds().Dx(), m.Bounds().Dy(), image.Pt(0, 0).In(m.Bounds()))

	// if err != nil {
	// fmt.Println(err.Error())
	// }

	// r := image.Rect(1, 1, 10, 10)
	// fmt.Println(r.Dx(), r.Dy(), image.Pt(0, 0).In(r))

	// bounds := m.Bounds()
	// fmt.Println("width, height", width, height)

	img := image.NewRGBA(image.Rect(0, 0, width, height))
	black := color.RGBA{0, 0, 0, 1}
	draw.Draw(img, img.Bounds(), &image.Uniform{black}, image.Point{0, 0}, draw.Src)

	id := "I.Utils.GetRandLen(lenIdImg)"

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

func (I *ImageUC) Add(chart *entity.Chart) error {
	nameSrcFile := I.Utils.GetRandLen(lenIdImg)

	dstFile, err := os.Open(addressImages + chart.IdParent + fileExtension)
	if err != nil {
		fmt.Println(err.Error())
	}

	dst, err := I.Utils.DecodeBMP(dstFile)
	if err != nil {
		fmt.Println(err.Error())
	}
	dstFile.Close()

	newRGBA := image.NewRGBA(dst.Bounds())

	draw.Draw(newRGBA, dst.Bounds(), dst, image.Point{}, draw.Src)

	srcFile, err := os.Create(addressImages + nameSrcFile + fileExtension)
	if err != nil {
		return err
	}

	_, err = io.Copy(srcFile, chart.File)
	if err != nil {
		return err
	}
	srcFile.Close()

	srcFile, err = os.Open(addressImages + nameSrcFile + fileExtension)
	if err != nil {
		return err
	}

	src, err := I.Utils.DecodeBMP(srcFile) // chart.File chatGPT
	if err != nil {
		return err
	}

	r := image.Rectangle{
		image.Point{
			0, //chart.X,
			0, //chart.Y,
		},
		image.Point{
			50, //chart.Width,
			50, //chart.Height,
		},
	}

	draw.Draw(newRGBA, r, src, src.Bounds().Min, draw.Src)

	srcFile.Close()

	dstFile, err = os.Create(addressImages + chart.IdParent + fileExtension)
	if err != nil {
		fmt.Println(err.Error())
	}

	err = I.Utils.EncodeBMP(dstFile, newRGBA)
	if err != nil {
		fmt.Println(err.Error())
	}

	dstFile.Close()

	return nil
}
