package main

import (
	"fmt"
	"net/url"
)

func main() {
	parsedURL, err := url.Parse("https://site.com.br/camisa/*")
	if err != nil {
		panic(err)
	}

	origin := parsedURL.Scheme + "://" + parsedURL.Host

	// Obter o path
	path := parsedURL.Path

	fmt.Printf("Origin: %s\n", origin)
	fmt.Printf("Path: %s\n", path)

	// Se quiser separar os componentes do path
	pathSegments := parsedURL.Path[1:] // Remove a primeira barra
	fmt.Printf("Segmento do path: %s\n", pathSegments)

}
