package utils

import (
	"image"
	"io"
	"net/http"
	"strconv"
	"this_module/internal/entity"
	"this_module/pkg/random"

	"github.com/go-chi/chi"
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

func (u *Utils) NewChart(r *http.Request) (*entity.Chart, error) {
	chart := &entity.Chart{}
	chart.IdParent = chi.URLParam(r, "idParent")
	chart.Width, _ = strconv.Atoi(r.URL.Query().Get("width"))
	chart.Height, _ = strconv.Atoi(r.URL.Query().Get("height"))
	chart.X, _ = strconv.Atoi(r.URL.Query().Get("x"))
	chart.Y, _ = strconv.Atoi(r.URL.Query().Get("y"))

	return chart, nil
}

func New() *Utils {
	return &Utils{}
}
