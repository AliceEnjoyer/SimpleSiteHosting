package getPage

import (
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/AliceEnjoyer/SimpleSiteHosting/lib/api/response"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	"golang.org/x/exp/slog"
)

type Response struct {
	response.Response
}

type contextKey struct {
	name string
}

func (k *contextKey) String() string {
	return "chi render context value " + k.name
}

func New(log *slog.Logger, pagesPath string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const fn = "handlers.getPage"

		log = log.With(
			slog.String("fn", fn),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		pageName := r.URL.String()[1:]
		page, err := readPage(pagesPath, pageName)
		if err != nil {
			log.Error("cannot get page: ", fn, "\t", err)

			file, err := readPage(pagesPath, "NotFound.html")
			if err != nil {
				log.Error("cannot get page 404: ", fn, "\t", err)

				render.JSON(w, r, response.Error("cannot get page 404: "+err.Error()))

				return
			}
			fileBytes, err := io.ReadAll(file)
			if err != nil {
				log.Error("cannot read page 404: ", fn, "\t", err)

				render.JSON(w, r, response.Error("cannot read page 404: "+err.Error()))

				return
			}
			render.HTML(w, r, string(fileBytes))

			return
		}
		pageBytes, err := io.ReadAll(page)
		if err != nil {
			log.Error("cannot read page: ", fn, "\t", err)

			render.JSON(w, r, response.Error("cannot read page 404: "+err.Error()))

			return
		}
		switch filepath.Ext(pageName) {
		case ".html":
			render.HTML(w, r, string(pageBytes))
		case ".css":
			w.Header().Set("Content-Type", "text/css;")
			w.Write(pageBytes)
		case ".js":
			w.Header().Set("Content-Type", "text/plain;")
			w.Write(pageBytes)
		case ".jpg":
			w.Header().Set("Content-Type", "image/jpeg;")
			w.Write(pageBytes)
		default:

		}
	}
}

func readPage(pagesPath, pageName string) (*os.File, error) {
	return os.OpenFile(pagesPath+pageName, os.O_RDONLY, 0644)
}
