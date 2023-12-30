package v1

import (
	"errors"
	"log/slog"
	"net/http"
	"strconv"
	"this_module/internal/entity"
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
		r.Get("/read", func(w http.ResponseWriter, r *http.Request) {

			w.Write([]byte("Hello"))

			return
		})

		r.Post("/chartas/", func(w http.ResponseWriter, r *http.Request) {
			width, _ := strconv.Atoi(r.URL.Query().Get("width"))
			height, _ := strconv.Atoi(r.URL.Query().Get("height"))

			id, err := uc.Create(width, height)
			if err != nil {
				l.Error(err.Error())
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			w.WriteHeader(201)
			w.Write([]byte(id))

			return
		})

		r.Post("/chartas/{IdParent}/", func(w http.ResponseWriter, r *http.Request) {
			chart := &entity.Chart{}
			err := errors.New("")
			// сделать проверку параметров
			chart.IdParent = chi.URLParam(r, "IdParent")
			chart.Width, _ = strconv.Atoi(r.URL.Query().Get("width"))
			chart.Height, _ = strconv.Atoi(r.URL.Query().Get("height"))
			chart.X, _ = strconv.Atoi(r.URL.Query().Get("x"))
			chart.Y, _ = strconv.Atoi(r.URL.Query().Get("y"))

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

			w.WriteHeader(200)

			return
		})
	})

	return r, nil
}

// http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
//
// w.Header().Set(strconv.Itoa(http.StatusCreated), http.StatusText(http.StatusCreated))
/* type Hello struct {
	Greeting string `json:"greeting"`
}

h := Hello{
	Greeting: "Hello world " + fmt.Sprint(id),
}

ec := json.NewEncoder(w)
err = ec.Encode(h) */

// w.Write([]byte(h))
