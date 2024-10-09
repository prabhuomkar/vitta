package handlers

import (
	"net/http"
	"vitta/config"

	"github.com/jackc/pgx/v5"
)

// Handler configuration.
type Handler struct {
	cfg *config.Config
	db  *pgx.Conn
}

func New(cfg *config.Config, db *pgx.Conn) *http.ServeMux {
	handler := &Handler{
		cfg: cfg,
		db:  db,
	}

	mux := http.NewServeMux()

	// accounts
	mux.HandleFunc("POST /v1/accounts", handler.CreateAccount)
	mux.HandleFunc("PATCH /v1/accounts/{id}", handler.UpdateAccount)
	mux.HandleFunc("DELETE /v1/accounts/{id}", handler.DeleteAccount)
	mux.HandleFunc("GET /v1/accounts", handler.GetAccounts)
	// payees
	mux.HandleFunc("POST /v1/payees", handler.CreatePayee)
	mux.HandleFunc("PATCH /v1/payees/{id}", handler.UpdatePayee)
	mux.HandleFunc("DELETE /v1/payees/{id}", handler.DeletePayee)
	mux.HandleFunc("GET /v1/payees", handler.GetPayees)
	// transactions
	mux.HandleFunc("POST /v1/transactions", handler.CreateTransaction)
	mux.HandleFunc("PUT /v1/transactions/{id}", handler.ImportTransactions)
	mux.HandleFunc("PATCH /v1/transactions/{id}", handler.UpdateTransaction)
	mux.HandleFunc("DELETE /v1/transactions/{id}", handler.DeleteTransaction)
	mux.HandleFunc("GET /v1/transactions", handler.GetTransactions)
	// budgets
	mux.HandleFunc("POST /v1/groups", handler.CreateGroup)
	mux.HandleFunc("PATCH /v1/groups/{id}", handler.UpdateGroup)
	mux.HandleFunc("DELETE /v1/groups/{id}", handler.DeleteGroup)
	mux.HandleFunc("POST /v1/categories", handler.CreateCategory)
	mux.HandleFunc("PATCH /v1/categories/{id}", handler.CreateCategory)
	mux.HandleFunc("DELETE /v1/categories/{id}", handler.CreateCategory)
	mux.HandleFunc("GET /v1/budgets", handler.GetBudget)
	mux.HandleFunc("PATCH /v1/budgets", handler.SetBudget)

	return mux
}
