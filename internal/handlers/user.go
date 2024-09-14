package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/ZiadMansourM/budgetly/internal/services"
)

// UserHandler is an HTTP handler for user-related operations
// (e.g., registration, fetching by ID, etc.)
type UserHandler struct {
	UserService *services.UserService
}

// Register routes for user-related actions
func (h *UserHandler) RegisterRoutes(router *http.ServeMux) {
	router.HandleFunc("POST /users/register", h.Register)
	router.HandleFunc("GET /users/{id}", h.GetByID)
}

// Register is an HTTP handler for registering a new user
func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	// Decode the JSON request body
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	// Call the service layer to register the user
	user, err := h.UserService.Register(req.Username, req.Email, req.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return the newly registered user as a JSON response
	json.NewEncoder(w).Encode(user)
}

// GetByID is an HTTP handler for fetching a user by ID
func (h *UserHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	// Parse the user ID from the URL
	userID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || userID <= 0 {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	// Call the service layer to retrieve the user
	user, err := h.UserService.GetByID(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return the user as a JSON response
	json.NewEncoder(w).Encode(user)
}
