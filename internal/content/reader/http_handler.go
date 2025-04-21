package reader

import (
	"encoding/base64"
	"encoding/json"
	"net/http"
)

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
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
	w.Header().Set("Cache-Control", "no-cache")

	defer r.Body.Close()

	endpoint := r.PathValue("endpoint")

	decodedBytes, err := base64.StdEncoding.DecodeString(endpoint)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid base64 encoding"))
		return err
	}

	decodedEndpoint := NewEndpointDto(string(decodedBytes))

	content, err := h.service.GetContent(r.Context(), decodedEndpoint)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Failed to get content"))
		return err
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(content); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Failed to encode content"))
		return err
	}

	return nil
}
