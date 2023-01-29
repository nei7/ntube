package rest

import (
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"path"
	"time"

	"github.com/go-chi/chi"
	"github.com/nei7/ntube/internal/dto"
	"github.com/nei7/ntube/internal/service"
)

type VideoHandler struct {
	videoStoragePath string
	svc              service.VideoService
}

func NewVideoHandler(videoStoragePath string, svc service.VideoService) *VideoHandler {
	return &VideoHandler{videoStoragePath, svc}
}

func (h *VideoHandler) Register(r *chi.Mux) {
	r.Get("/videos/{video}", h.serveVideo)
	r.Get("/thumbnails/{thumbnail}", h.serveThumbnail)
	r.Post("/videos/search", h.search)
}

func (h *VideoHandler) serveVideo(w http.ResponseWriter, r *http.Request) {
	folderId := chi.URLParam(r, "video")

	quality := r.URL.Query().Get("quality")
	if quality == "" {
		quality = "720p"
	}

	file, err := os.Open(path.Join(h.videoStoragePath, folderId, quality+".mp4"))
	if err != nil {
		renderErrorResponse(w, r, errors.New("not found"), http.StatusNotFound)
		return
	}

	http.ServeContent(w, r, folderId, time.Unix(0, 0), file)
}

func (h *VideoHandler) serveThumbnail(w http.ResponseWriter, r *http.Request) {
	folderId := chi.URLParam(r, "thumbnail")

	file, err := os.Open(path.Join(h.videoStoragePath, folderId, "thumbnail1.jpg"))
	if err != nil {
		renderErrorResponse(w, r, errors.New("not found"), http.StatusNotFound)
		return
	}

	http.ServeContent(w, r, folderId, time.Unix(0, 0), file)
}

func (h *VideoHandler) search(w http.ResponseWriter, r *http.Request) {
	var req dto.VideoSearchParams
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		renderErrorResponse(w, r, err, http.StatusBadRequest)
		return
	}

	defer r.Body.Close()

	res, err := h.svc.Search(r.Context(), req)
	if err != nil {
		renderErrorResponse(w, r, err, http.StatusInternalServerError)
		return
	}

	renderResponse(w, r, res, http.StatusOK)
}
