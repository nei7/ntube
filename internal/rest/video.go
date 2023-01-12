package rest

import (
	"errors"
	"net/http"
	"os"
	"path"
	"time"

	"github.com/go-chi/chi"
)

type VideoHandler struct {
	videoStoragePath string
}

func NewVideoHandler(videoStoragePath string) *VideoHandler {
	return &VideoHandler{videoStoragePath}
}

func (h *VideoHandler) Register(r *chi.Mux) {
	r.Get("/videos/{video}", h.serve)
}

func (h *VideoHandler) serve(w http.ResponseWriter, r *http.Request) {
	videoName := chi.URLParam(r, "video")

	file, err := os.Open(path.Join(h.videoStoragePath, "mp4", videoName+".mp4"))
	if err != nil {
		renderErrorResponse(w, r, errors.New("not found"), http.StatusNotFound)
		return
	}

	http.ServeContent(w, r, videoName, time.Unix(0, 0), file)
}
