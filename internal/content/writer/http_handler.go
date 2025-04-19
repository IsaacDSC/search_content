package writer

import "net/http"

type Handler interface {
	SaveContent(w http.ResponseWriter, r *http.Request) error
}

type HttpHandler struct {
	service Service
}

var _ Handler = (*HttpHandler)(nil)

func NewHandler(service Service) *HttpHandler {
	return &HttpHandler{service: service}
}

func (h *HttpHandler) SaveContent(w http.ResponseWriter, r *http.Request) error {
	_, err := w.Write([]byte("OK"))

	return err
}
