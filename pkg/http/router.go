package http

import (
	"github.com/IsaacDSC/search_content/internal/content/infra/api/handler"
	"net/http"
)

func GetRouters() *http.ServeMux {
	m := http.NewServeMux()

	videoHandler := handler.NewHandler()

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
