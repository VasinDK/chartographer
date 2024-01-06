package v1

import (
	"log/slog"
	"net/http"
	"strconv"
	"this_module/internal/usecase"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

const formName = "aaa"

func NewRouter(l *slog.Logger, uc *usecase.ImageUC) (*chi.Mux, error) {
	r := chi.NewRouter()

	// r.Use(middleware.RequestID)
	// r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Route("/v1", func(r chi.Router) {
		r.Post("/chartas/", func(w http.ResponseWriter, r *http.Request) {
			width, _ := strconv.Atoi(r.URL.Query().Get("width"))
			height, _ := strconv.Atoi(r.URL.Query().Get("height"))

			id, err := uc.Create(width, height)
			if err != nil {
				l.Error(err.Error())
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			w.WriteHeader(http.StatusCreated)
			w.Write([]byte(id))

			return
		})

		r.Post("/chartas/{idParent}/", func(w http.ResponseWriter, r *http.Request) {
			chart, err := uc.Utils.NewChart(r)
			// сделать проверку параметров
			chart.File, chart.Handle, err = r.FormFile(formName)
			defer chart.File.Close()
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}

			err = uc.Add(chart)
			if err != nil {
				l.Error(err.Error())
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			w.WriteHeader(http.StatusOK)

			return
		})

		r.Get("/chartas/{idParent}/", func(w http.ResponseWriter, r *http.Request) {
			chart, err := uc.Utils.NewChart(r)

			filePart, err := uc.Part(chart)
			if err != nil {
				l.Error(err.Error())
				http.Error(w, err.Error(), http.StatusBadRequest)
			}

			err = uc.Utils.EncodeBMP(w, filePart)
			if err != nil {
				l.Error(err.Error())
				http.Error(w, err.Error(), http.StatusBadRequest)
			}

			return
		})

		r.Delete("/chartas/{id}/", func(w http.ResponseWriter, r *http.Request) {
			err := uc.Delete(chi.URLParam(r, "id"))

			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte(err.Error()))
			}

			return
		})
	})

	return r, nil
}
