package handler

import "net/http"

type Handler struct {
}

func NewHandler() *Handler {
	return &Handler{}
}

func (h Handler) GetRoutes() map[string]func(w http.ResponseWriter, r *http.Request) error {
	return map[string]func(w http.ResponseWriter, r *http.Request) error{
		"GET /health":    h.health,
		"POST /endpoint": h.saveNewRule,
	}
}

func (h Handler) health(w http.ResponseWriter, r *http.Request) error {
	w.WriteHeader(http.StatusOK)

	return nil
}

func (h Handler) saveNewRule(w http.ResponseWriter, r *http.Request) error {
	_, err := w.Write([]byte("OK"))

	return err
}
