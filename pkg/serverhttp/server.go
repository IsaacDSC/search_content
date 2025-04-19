package serverhttp

import (
	"github.com/IsaacDSC/search_content/internal/content/infra/container"
	"log"
	"net/http"
)

func StartServer(handlers container.Handlers) error {
	routers := GetRouters(handlers.WriterHandler, handlers.ReaderHandler)

	server := &http.Server{
		Addr:    ":8080",
		Handler: routers,
	}

	log.Printf("Server listening on %s", server.Addr)

	if err := server.ListenAndServe(); err != nil {
		return err
	}

	return nil
}
