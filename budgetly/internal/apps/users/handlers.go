package users

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/ZiadMansourM/budgetly/utils"
)

// userHandler is an HTTP handler for user-related operations
// (e.g., registration, fetching by ID, etc.)
type userHandler struct {
	userService *userService
	logger      *slog.Logger
	router      *http.ServeMux
}

// newUserHandler creates a new user handler with the provided user service and logger
func newUserHandler(userService *userService, logger *slog.Logger, router *http.ServeMux) *userHandler {
	userHandler := &userHandler{
		userService: userService,
		logger:      logger,
		router:      router,
	}
	userHandler.registerRoutes()
	userHandler.registerSSRRoutes()
	return userHandler
}

// Register routes for user-related actions
func (h *userHandler) registerRoutes() {
	h.router.HandleFunc("POST /users/register", h.register)
	h.router.HandleFunc("GET /users/{id}", h.getByID)
}

// registerSSRRoutes registers SSR routes for user-related actions
func (h *userHandler) registerSSRRoutes() {
	h.router.HandleFunc("/users", h.renderUserListPage) // SSR route for rendering user list
}

// renderUserListPage is a dummy SSR handler for rendering the user list page
func (h *userHandler) renderUserListPage(w http.ResponseWriter, r *http.Request) {
	// Render the user list page
	page := `
		<h1>User List</h1>
		<ul>
			<li>User 1</li>
			<li>User 2</li>
		</ul>
	`
	w.Write([]byte(page))
}

// Register is an HTTP handler for registering a new user
func (h *userHandler) register(w http.ResponseWriter, r *http.Request) {
	var req UserRequest

	// Decode the JSON request body
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Warn("Invalid request body", "error", err)
		utils.WriteJson(
			w,
			http.StatusBadRequest,
			map[string]string{"error": "Invalid request body"},
		)
		return
	}

	// Call the service layer to register the user
	user, err := h.userService.register(req)
	if err != nil {
		h.logger.Warn("Error registering user", "error", err)
		utils.WriteJson(
			w,
			http.StatusBadRequest,
			map[string]string{"error": err.Error()},
		)
		return
	}

	// Return the newly registered user as a JSON response
	utils.WriteJson(w, http.StatusCreated, user)
}

// GetByID is an HTTP handler for fetching a user by ID
func (h *userHandler) getByID(w http.ResponseWriter, r *http.Request) {
	// Parse the user ID from the URL
	userID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || userID <= 0 {
		utils.WriteJson(
			w,
			http.StatusBadRequest,
			map[string]string{"error": "Invalid user ID"},
		)
		return
	}

	// Call the service layer to retrieve the user
	user, err := h.userService.getByID(userID)
	if err != nil {
		utils.WriteJson(
			w,
			http.StatusNotFound,
			map[string]string{"error": "User not found"},
		)
		return
	}

	// Return the user as a JSON response
	utils.WriteJson(w, http.StatusOK, user)
}
