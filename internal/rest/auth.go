package rest

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/nei7/gls/internal"
	"github.com/nei7/gls/internal/db"
	"github.com/nei7/gls/internal/dto"
	"golang.org/x/crypto/bcrypt"
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

func (h *UserHandler) Register(r *chi.Mux) {
	r.Post("/signup", h.signUp)
	r.Post("/login", h.logIn)
}

func (h *UserHandler) signUp(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateUserParams
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		renderErrorResponse(w, r, "Invalid request", internal.WrapErrorf(err, internal.ErrorCodeInvalidArgument, "json decoder"))
		return
	}
	defer r.Body.Close()

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		renderErrorResponse(w, r, "Internal server error", err)
		return
	}
	req.Password = string(hashedPassword)

	user, err := h.userService.Create(r.Context(), req)
	if err != nil {
		renderErrorResponse(w, r, "Failed to create user", err)
		return
	}

	renderResponse(w, r, map[string]string{
		"email": user.Email,
		"id":    user.ID.String(),
	}, http.StatusCreated)
}

func (h *UserHandler) logIn(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateUserParams
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		renderErrorResponse(w, r, "Invalid request", internal.WrapErrorf(err, internal.ErrorCodeInvalidArgument, "json decoder"))
		return
	}
	defer r.Body.Close()

	user, err := h.userService.Find(r.Context(), req.Email)
	if err != nil {
		renderErrorResponse(w, r, "user not found", err)
		return
	}

	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		renderErrorResponse(w, r, "invalid password", err)
		return
	}

	renderResponse(w, r, map[string]string{
		"email": user.Email,
		"id":    user.ID.String(),
	}, http.StatusCreated)
}
