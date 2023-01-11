package rest

import (
	"net/http"

	"github.com/go-chi/chi"
)

type VideoHandler struct{}

func NewVideoHandler() *VideoHandler {
	return &VideoHandler{}
}

func (h *VideoHandler) Register(r *chi.Mux) {}

func (h *VideoHandler) serve(w http.ResponseWriter, r *http.Request) {
	chi.URLParam()

	http.ServeContent(w, r)
}
