package main

import (
	"context"
	"github.com/IsaacDSC/search_content/internal/content/infra/container"
	"github.com/IsaacDSC/search_content/pkg/serverhttp"
	"github.com/redis/go-redis/v9"
	"log"
)

var client *redis.Client

func init() {
	client = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379", //TODO: levar para variavel de ambiente
		Password: "",               // No password set
		DB:       0,                // Use default DB
		Protocol: 2,                // Connection protocol
	})

	if err := client.Ping(context.Background()).Err(); err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
}

func main() {
	cacheStrategies := container.NewCacheStrategies(client)
	repositories := container.NewRepositoryContainer()
	services := container.NewServicesContainer(repositories)
	handlers := container.GetHandlers(services)

	if err := serverhttp.StartServer(handlers, cacheStrategies); err != nil {
		log.Fatal(err)
	}

}
