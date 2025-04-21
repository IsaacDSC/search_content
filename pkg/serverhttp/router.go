package serverhttp

import (
	"github.com/IsaacDSC/search_content/internal/content/infra/api/handler"
	"github.com/IsaacDSC/search_content/internal/content/infra/container"
	"github.com/IsaacDSC/search_content/internal/content/reader"
	"github.com/IsaacDSC/search_content/internal/content/writer"
	"net/http"
)

func GetRouters(strategies container.CacheStrategies, wh writer.Handler, rh reader.Handler) *http.ServeMux {
	m := http.NewServeMux()

	cacheMw := reader.NewCacheMiddleware(strategies.LRUCache)
	videoHandler := handler.NewHandler(wh, rh)

	for path, fn := range videoHandler.GetRoutes() {
		m.HandleFunc(path, func(responseWriter http.ResponseWriter, request *http.Request) {
			if request.Method == http.MethodGet {
				cacheMw.WithCache(func(w http.ResponseWriter, r *http.Request) {
					if err := fn(w, r); err != nil {
						http.Error(w, err.Error(), http.StatusInternalServerError)
					}
				})(responseWriter, request)
			} else {
				// No caching for non-GET requests
				if err := fn(responseWriter, request); err != nil {
					http.Error(responseWriter, err.Error(), http.StatusInternalServerError)
				}
			}
		})
	}

	return m
}
