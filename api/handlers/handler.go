package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"vitta/adapters"
	"vitta/config"
	"vitta/database"
)

// Handler configuration.
type Handler struct {
	cfg      *config.Config
	db       database.DBIface
	adapters map[string]adapters.Config
}

func New(cfg *config.Config, db database.DBIface, adapters map[string]adapters.Config) http.Handler {
	h := &Handler{
		cfg:      cfg,
		db:       db,
		adapters: adapters,
	}

	mux := http.NewServeMux()

	// payees
	mux.HandleFunc("POST /v1/payees", h.CreatePayee)
	mux.HandleFunc("PATCH /v1/payees/{id}", h.UpdatePayee)
	mux.HandleFunc("DELETE /v1/payees/{id}", h.DeletePayee)
	mux.HandleFunc("GET /v1/payees", h.GetPayees)
	// accounts
	mux.HandleFunc("POST /v1/accounts", h.CreateAccount)
	mux.HandleFunc("PATCH /v1/accounts/{id}", h.UpdateAccount)
	mux.HandleFunc("DELETE /v1/accounts/{id}", h.DeleteAccount)
	mux.HandleFunc("GET /v1/accounts", h.GetAccounts)
	// transactions
	mux.HandleFunc("POST /v1/accounts/{id}/transactions", h.CreateTransaction)
	mux.HandleFunc("PUT /v1/accounts/{id}/transactions", h.ImportTransactions)
	mux.HandleFunc("PATCH /v1/accounts/{id}/transactions/{tId}", h.UpdateTransaction)
	mux.HandleFunc("DELETE /v1/accounts/{id}/transactions/{tId}", h.DeleteTransaction)
	mux.HandleFunc("GET /v1/accounts/{id}/transactions", h.GetTransactions)
	// budgets
	mux.HandleFunc("POST /v1/groups", h.CreateGroup)
	mux.HandleFunc("PATCH /v1/groups/{id}", h.UpdateGroup)
	mux.HandleFunc("DELETE /v1/groups/{id}", h.DeleteGroup)
	mux.HandleFunc("GET /v1/groups", h.GetGroups)
	mux.HandleFunc("POST /v1/categories", h.CreateCategory)
	mux.HandleFunc("PATCH /v1/categories/{id}", h.UpdateCategory)
	mux.HandleFunc("DELETE /v1/categories/{id}", h.DeleteCategory)
	mux.HandleFunc("GET /v1/categories", h.GetCategories)
	mux.HandleFunc("GET /v1/budgets", h.GetBudget)
	mux.HandleFunc("PUT /v1/budgets", h.SetBudget)

	return h.corsMiddleware(h.basicAuthMiddleware(mux))
}

// ErrorResponse model.
type ErrorResponse struct {
	Error string `json:"error"`
}

func buildErrorResponse(w http.ResponseWriter, err string, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	errorResponse := ErrorResponse{Error: err}

	encodeErr := json.NewEncoder(w).Encode(errorResponse)
	if encodeErr != nil {
		slog.Error("error encoding error response", "error", encodeErr)
	}
}

// Middleware for checking basic auth credentials.
func (h *Handler) basicAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		username, password, ok := r.BasicAuth()

		if !ok || username != h.cfg.AdminUsername || password != h.cfg.AdminPassword {
			buildErrorResponse(w, "Unauthorized", http.StatusUnauthorized)

			return
		}

		next.ServeHTTP(w, r)
	})
}

// Middleware for CORS related operations.
func (h *Handler) corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PATCH, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)

			return
		}

		next.ServeHTTP(w, r)
	})
}
