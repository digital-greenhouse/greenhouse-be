package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"digital-greenhouse/greenhouse-be/internal/domain"
	"digital-greenhouse/greenhouse-be/internal/http/dto"

	"github.com/go-chi/chi/v5"
)

type UserHandler struct {
	userService domain.UserService
}

func NewUserHandler(userService domain.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

func jsonResponse(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if data != nil {
		json.NewEncoder(w).Encode(data)
	}
}

func errResponse(w http.ResponseWriter, status int, message string) {
	jsonResponse(w, status, map[string]string{"error": message})
}

func toUserResponse(u *domain.User) dto.UserResponse {
	return dto.UserResponse{
		ID:        u.ID,
		Name:      u.Name,
		Email:     u.Email,
		Role:      string(u.Role),
		Phone:     u.Phone,
		IsActive:  u.IsActive,
		CreatedAt: u.CreatedAt.Format(time.RFC3339),
	}
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		errResponse(w, http.StatusBadRequest, "payload inválido")
		return
	}

	user := &domain.User{
		Name:         req.Name,
		Email:        req.Email,
		PasswordHash: req.Password,
		Role:         domain.Role(req.Role),
	}
	if req.Phone != "" {
		user.Phone = &req.Phone
	}

	if err := h.userService.CreateUser(r.Context(), user); err != nil {
		errResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	jsonResponse(w, http.StatusCreated, toUserResponse(user))
}

func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req dto.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		errResponse(w, http.StatusBadRequest, "payload inválido")
		return
	}

	token, user, err := h.userService.Login(r.Context(), req.Email, req.Password)
	if err != nil {
		errResponse(w, http.StatusUnauthorized, err.Error())
		return
	}

	resp := dto.LoginResponse{
		Token: token,
		User:  toUserResponse(user),
	}

	jsonResponse(w, http.StatusOK, resp)
}

func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 32)
	if err != nil {
		errResponse(w, http.StatusBadRequest, "ID de usuario inválido")
		return
	}

	user, err := h.userService.GetUserByID(r.Context(), uint(id))
	if err != nil {
		errResponse(w, http.StatusNotFound, "usuario no encontrado")
		return
	}

	jsonResponse(w, http.StatusOK, toUserResponse(user))
}

func (h *UserHandler) GetUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.userService.GetAllUsers(r.Context())
	if err != nil {
		errResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	resp := make([]dto.UserResponse, len(users))
	for i := range users {
		resp[i] = toUserResponse(&users[i])
	}

	jsonResponse(w, http.StatusOK, resp)
}

func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 32)
	if err != nil {
		errResponse(w, http.StatusBadRequest, "ID de usuario inválido")
		return
	}

	var req dto.UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		errResponse(w, http.StatusBadRequest, "payload inválido")
		return
	}

	user := &domain.User{
		ID:   uint(id),
		Name: req.Name,
		Role: domain.Role(req.Role),
	}
	if req.Phone != "" {
		user.Phone = &req.Phone
	}

	if err := h.userService.UpdateUser(r.Context(), user); err != nil {
		errResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	jsonResponse(w, http.StatusOK, toUserResponse(user))
}

func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 32)
	if err != nil {
		errResponse(w, http.StatusBadRequest, "ID de usuario inválido")
		return
	}

	if err := h.userService.DeleteUser(r.Context(), uint(id)); err != nil {
		errResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	jsonResponse(w, http.StatusNoContent, nil)
}
