package rest

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/nei7/ntube/internal/db"
	"github.com/nei7/ntube/internal/dto"
	"github.com/nei7/ntube/internal/service"
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
	r.Post("/renew", h.renewToken)
}

type userResponse struct {
	email         string
	username      string
	id            string
	followers     int32
	description   string
	avatar        string
	created_at    string
	access_token  string
	refresh_token string
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
		if errors.Is(err, db.DuplicateKeyValueError) {
			renderErrorResponse(w, r, errors.New("user already exists"), http.StatusConflict)
			return
		}

		renderErrorResponse(w, r, errors.New("Failed to create user"), http.StatusInternalServerError)
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

	renderResponse(w, r, userResponse{
		email:         user.Email,
		username:      user.Username,
		id:            user.ID.String(),
		avatar:        user.Avatar.String,
		created_at:    user.CreatedAt.Time.String(),
		followers:     user.Followers,
		access_token:  accessToken,
		refresh_token: refreshToken,
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

	renderResponse(w, r, userResponse{
		email:         user.Email,
		username:      user.Username,
		id:            user.ID.String(),
		avatar:        user.Avatar.String,
		created_at:    user.CreatedAt.Time.String(),
		followers:     user.Followers,
		access_token:  accessToken,
		refresh_token: refreshToken,
	}, http.StatusOK)
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
