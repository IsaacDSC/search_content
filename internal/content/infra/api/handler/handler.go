package handler

import (
	"github.com/IsaacDSC/search_content/internal/content/reader"
	"github.com/IsaacDSC/search_content/internal/content/writer"
	"net/http"
)

type Handler struct {
	wh writer.Handler
	rh reader.Handler
}

func NewHandler(wh writer.Handler, rh reader.Handler) *Handler {
	return &Handler{
		wh: wh,
		rh: rh,
	}
}

func (h Handler) GetRoutes() map[string]func(w http.ResponseWriter, r *http.Request) error {
	return map[string]func(w http.ResponseWriter, r *http.Request) error{
		"GET /health":             h.health,
		"POST /content":           h.wh.SaveContent,
		"GET /content/{endpoint}": h.rh.GetContent,
	}
}

func (h Handler) health(w http.ResponseWriter, r *http.Request) error {
	w.WriteHeader(http.StatusOK)

	return nil
}
