package handler

import "net/http"

type Adapter interface {
	GetRoutes() map[string]func(w http.ResponseWriter, r *http.Request) error
}
