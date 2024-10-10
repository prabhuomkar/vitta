package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"vitta/config"

	"github.com/jackc/pgx/v5"
)

// Handler configuration.
type Handler struct {
	cfg *config.Config
	db  *pgx.Conn
}

func New(cfg *config.Config, db *pgx.Conn) http.Handler {
	h := &Handler{
		cfg: cfg,
		db:  db,
	}

	mux := http.NewServeMux()

	// accounts
	mux.HandleFunc("POST /v1/accounts", h.CreateAccount)
	mux.HandleFunc("PATCH /v1/accounts/{id}", h.UpdateAccount)
	mux.HandleFunc("DELETE /v1/accounts/{id}", h.DeleteAccount)
	mux.HandleFunc("GET /v1/accounts", h.GetAccounts)
	// payees
	mux.HandleFunc("POST /v1/payees", h.CreatePayee)
	mux.HandleFunc("PATCH /v1/payees/{id}", h.UpdatePayee)
	mux.HandleFunc("DELETE /v1/payees/{id}", h.DeletePayee)
	mux.HandleFunc("GET /v1/payees", h.GetPayees)
	// transactions
	mux.HandleFunc("POST /v1/transactions", h.CreateTransaction)
	mux.HandleFunc("PUT /v1/transactions/{id}", h.ImportTransactions)
	mux.HandleFunc("PATCH /v1/transactions/{id}", h.UpdateTransaction)
	mux.HandleFunc("DELETE /v1/transactions/{id}", h.DeleteTransaction)
	mux.HandleFunc("GET /v1/transactions", h.GetTransactions)
	// budgets
	mux.HandleFunc("POST /v1/groups", h.CreateGroup)
	mux.HandleFunc("PATCH /v1/groups/{id}", h.UpdateGroup)
	mux.HandleFunc("DELETE /v1/groups/{id}", h.DeleteGroup)
	mux.HandleFunc("POST /v1/categories", h.CreateCategory)
	mux.HandleFunc("PATCH /v1/categories/{id}", h.UpdateCategory)
	mux.HandleFunc("DELETE /v1/categories/{id}", h.DeleteCategory)
	mux.HandleFunc("GET /v1/budgets", h.GetBudget)
	mux.HandleFunc("PUT /v1/budgets", h.SetBudget)

	return h.basicAuthMiddleware(mux)
}

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
