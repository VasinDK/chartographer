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

	id := "I.Utils.GetRandLen(lenIdImg)" // временно

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
	I.Create(400, 400) // временно
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
