package http

import "net/http"

func StartServer() error {
	routers := GetRouters()

	// Initialize the server
	server := &http.Server{
		Addr:    ":8080",
		Handler: routers, // Set your handler here
	}

	// Start the server
	if err := server.ListenAndServe(); err != nil {
		return err
	}

	return nil
}
