package getPage

import (
	"io"
	"net/http"
	"os"

	"github.com/AliceEnjoyer/SimpleSiteHosting/lib/api/response"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	"golang.org/x/exp/slog"
)

type Response struct {
	response.Response
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
		render.HTML(w, r, string(pageBytes))

		style, _ := readPage(pagesPath, "style.css")
		styleBytes, _ := io.ReadAll(style)
		render.Data(w, r, []byte("<style>"+string(styleBytes)+"</style>"))
	}
}

func readPage(pagesPath, pageName string) (*os.File, error) {
	return os.OpenFile(pagesPath+pageName, os.O_RDONLY, 0644)
}
