package users

import (
	"encoding/json"
	"net/http"
	"strconv"
)

// userHandler is an HTTP handler for user-related operations
// (e.g., registration, fetching by ID, etc.)
type userHandler struct {
	userService *userService
}

// Register routes for user-related actions
func (h *userHandler) RegisterRoutes(router *http.ServeMux) {
	router.HandleFunc("POST /users/register", h.register)
	router.HandleFunc("GET /users/{id}", h.getByID)
}

// Register is an HTTP handler for registering a new user
func (h *userHandler) register(w http.ResponseWriter, r *http.Request) {
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
	user, err := h.userService.register(req.Username, req.Email, req.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return the newly registered user as a JSON response
	json.NewEncoder(w).Encode(user)
}

// GetByID is an HTTP handler for fetching a user by ID
func (h *userHandler) getByID(w http.ResponseWriter, r *http.Request) {
	// Parse the user ID from the URL
	userID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || userID <= 0 {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	// Call the service layer to retrieve the user
	user, err := h.userService.getByID(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return the user as a JSON response
	json.NewEncoder(w).Encode(user)
}
