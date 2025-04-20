package writer

import (
	"encoding/json"
	"log"
	"net/http"
)

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
	w.Header().Set("Content-Type", "application/json")
	defer r.Body.Close()

	var body VideoInputDto
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return nil
	}

	if err := h.service.Register(r.Context(), body); err != nil {
		log.Printf("failed to register video: %v", err)
		http.Error(w, "Failed to register content", http.StatusInternalServerError)
		return nil
	}

	w.WriteHeader(http.StatusCreated)
	return nil
}
