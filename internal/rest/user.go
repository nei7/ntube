package rest

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/nei7/gls/internal"
	"github.com/nei7/gls/internal/db"
	"github.com/nei7/gls/internal/dto"
)

type UserService interface {
	Create(ctx context.Context, params dto.CreateUserParams) (db.User, error)
	Find(ctx context.Context, email string) (db.User, error)
}

type UserHandler struct {
	userService UserService
}

func NewUserHandler(userService UserService) *UserHandler {
	return &UserHandler{
		userService,
	}
}

func (h *UserHandler) Register(r *chi.Mux) {}

func (h *UserHandler) create(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateUserParams
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		renderErrorResponse(w, r, "invalid request", internal.WrapErrorf(err, internal.ErrorCodeInvalidArgument, "json decoder"))
		return
	}

	defer r.Body.Close()

	user, err := h.userService.Create(r.Context(), req)
	if err != nil {
		renderErrorResponse(w, r, "failed to create", err)
		return
	}

	renderResponse(w, r,,)

}
