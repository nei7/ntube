package rest

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/lib/pq"
	"github.com/nei7/gls/internal/dto"
	"github.com/nei7/gls/internal/service"
	"golang.org/x/crypto/bcrypt"
)

type UserHandler struct {
	userService  service.UserService
	tokenManager service.TokenManager
}

func NewUserHandler(userService service.UserService, tokenManager service.TokenManager) *UserHandler {
	return &UserHandler{
		userService,
		tokenManager,
	}
}

func (h *UserHandler) Register(r *chi.Mux) {
	r.Post("/signup", h.signUp)
	r.Post("/login", h.logIn)
}

func (h *UserHandler) signUp(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateUserParams
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		renderErrorResponse(w, r, errors.New("Invalid request"), http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		renderErrorResponse(w, r, errors.New("Invalid request"), http.StatusInternalServerError)
		return
	}
	req.Password = string(hashedPassword)

	user, err := h.userService.Create(r.Context(), req)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "unique_violation":
				renderErrorResponse(w, r, errors.New("User already exists"), http.StatusConflict)
				return
			}
		}
		renderErrorResponse(w, r, errors.New("Failed to create user"), http.StatusInternalServerError)
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
		renderErrorResponse(w, r, errors.New("Invalid request"), http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	user, err := h.userService.Find(r.Context(), req.Email)
	if err != nil {
		renderErrorResponse(w, r, errors.New("User not found"), http.StatusNotFound)
		return
	}

	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		renderErrorResponse(w, r, errors.New("Invalid password"), http.StatusConflict)
		return
	}

	accessToken, err := h.tokenManager.NewJWT(user.ID.String(), time.Now().Add(15*time.Minute).Unix())
	if err != nil {
		renderErrorResponse(w, r, errors.New("Failed to generate access token"), http.StatusInternalServerError)
		return
	}

	refreshToken, err := h.tokenManager.NewJWT(user.ID.String(), time.Now().Add(7*24*time.Hour).Unix())
	if err != nil {
		renderErrorResponse(w, r, errors.New("Failed to generate refresh token"), http.StatusInternalServerError)
		return
	}

	renderResponse(w, r, map[string]string{
		"email":        user.Email,
		"id":           user.ID.String(),
		"refreshToken": refreshToken,
		"accessToken":  accessToken,
	}, http.StatusCreated)
}

type renewAccessTokenRequest struct {
	RefreshToken string `json:"refresh_token"`
}

func (h *UserHandler) renewToken(w http.ResponseWriter, r *http.Request) {
	var req renewAccessTokenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		renderErrorResponse(w, r, errors.New("Invalid request"), http.StatusInternalServerError)
		return
	}

	id, err := h.tokenManager.Parse(req.RefreshToken)
	if err != nil {
		renderErrorResponse(w, r, errors.New("Invalid token"), http.StatusUnauthorized)
		return
	}

	accessToken, err := h.tokenManager.NewJWT(id, time.Now().Add(time.Minute*15).Unix())
	if err != nil {
		renderErrorResponse(w, r, errors.New("Invalid token"), http.StatusInternalServerError)
		return
	}

	renderResponse(w, r, map[string]string{"accessToken": accessToken}, http.StatusOK)
}
