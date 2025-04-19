package reader

import "net/http"

type Handler interface {
	GetContent(w http.ResponseWriter, r *http.Request) error
}

type HttpHandler struct {
	service Service
}

func NewHandler(service Service) *HttpHandler {
	return &HttpHandler{service: service}
}

func (h *HttpHandler) GetContent(w http.ResponseWriter, r *http.Request) error {
	_, err := w.Write([]byte("OK"))

	return err
}
