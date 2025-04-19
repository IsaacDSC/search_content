package serverhttp

import (
	"github.com/IsaacDSC/search_content/internal/content/infra/api/handler"
	"github.com/IsaacDSC/search_content/internal/content/reader"
	"github.com/IsaacDSC/search_content/internal/content/writer"
	"net/http"
)

func GetRouters(wh writer.Handler, rh reader.Handler) *http.ServeMux {
	m := http.NewServeMux()

	videoHandler := handler.NewHandler(wh, rh)

	for path, fn := range videoHandler.GetRoutes() {
		m.HandleFunc(path, func(responseWriter http.ResponseWriter, request *http.Request) {
			if err := fn(responseWriter, request); err != nil {
				http.Error(responseWriter, err.Error(), http.StatusInternalServerError)
				return
			}
		})
	}

	return m
}
