package main

import (
	"github.com/IsaacDSC/search_content/internal/content/infra/container"
	"github.com/IsaacDSC/search_content/pkg/serverhttp"
	"log"
)

func main() {

	repositories := container.NewRepositoryContainer()
	services := container.NewServicesContainer(repositories)
	handlers := container.GetHandlers(services)

	if err := serverhttp.StartServer(handlers); err != nil {
		log.Fatal(err)
	}

}
