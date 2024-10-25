package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	uuid "github.com/google/uuid"
)

// Account model.
type Account struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	OffBudget *bool     `json:"offBudget"`
	Category  string    `json:"category"`
	Adapter   string    `json:"adapter"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

const (
	queryCreateAccount = `INSERT INTO accounts (id, name, off_budget, category, adapter, created_at, updated_at)` +
		` VALUES ($1, $2, $3, $4, $5, $6, $7)`
	queryUpdateAccount = `UPDATE accounts SET name=$1, off_budget=$2, category=$3, adapter=$4, updated_at=$5 WHERE id=$6`
	queryDeleteAccount = `DELETE FROM accounts WHERE id=$1`
	queryGetAccount    = "SELECT * FROM accounts WHERE id=$1"
	queryGetAccounts   = `SELECT * FROM accounts ORDER BY created_at DESC`
)

func (h *Handler) CreateAccount(w http.ResponseWriter, r *http.Request) {
	var account Account

	err := json.NewDecoder(r.Body).Decode(&account)
	if err != nil {
		slog.Error("error decoding create account request", "error", err)
		buildErrorResponse(w, err.Error(), http.StatusBadRequest)

		return
	}

	account.ID, err = uuid.NewV7()
	if err != nil {
		slog.Error("error creating account id", "error", err)
		buildErrorResponse(w, err.Error(), http.StatusInternalServerError)

		return
	}

	account.CreatedAt = time.Now()
	account.UpdatedAt = account.CreatedAt

	_, err = h.db.Exec(r.Context(), queryCreateAccount,
		account.ID, account.Name, account.OffBudget, account.Category, account.Adapter,
		account.CreatedAt, account.UpdatedAt)
	if err != nil {
		slog.Error("error creating account in database", "error", err)
		buildErrorResponse(w, err.Error(), http.StatusInternalServerError)

		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	err = json.NewEncoder(w).Encode(account)
	if err != nil {
		slog.Error("error encoding account response", "error", err)
	}
}

func (h *Handler) UpdateAccount(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	accountID, err := uuid.Parse(id)
	if err != nil {
		slog.Error("error parsing account id", "error", err)
		buildErrorResponse(w, err.Error(), http.StatusBadRequest)

		return
	}

	var account Account

	err = json.NewDecoder(r.Body).Decode(&account)
	if err != nil {
		slog.Error("error decoding update account request", "error", err)
		buildErrorResponse(w, err.Error(), http.StatusBadRequest)

		return
	}

	account.UpdatedAt = time.Now()

	_, err = h.db.Exec(r.Context(), queryUpdateAccount,
		account.Name, account.OffBudget, account.Category, account.Adapter, account.UpdatedAt, accountID)
	if err != nil {
		slog.Error("error updating account in database", "error", err)
		buildErrorResponse(w, err.Error(), http.StatusInternalServerError)

		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) DeleteAccount(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	accountID, err := uuid.Parse(id)
	if err != nil {
		slog.Error("error parsing account id", "error", err)
		buildErrorResponse(w, err.Error(), http.StatusBadRequest)

		return
	}

	_, err = h.db.Exec(r.Context(), queryDeleteAccount, accountID)
	if err != nil {
		slog.Error("error deleting account in database", "error", err)
		buildErrorResponse(w, err.Error(), http.StatusInternalServerError)

		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) GetAccounts(w http.ResponseWriter, r *http.Request) {
	rows, err := h.db.Query(r.Context(), queryGetAccounts)
	if err != nil {
		slog.Error("error getting accounts from database", "error", err)
		buildErrorResponse(w, err.Error(), http.StatusInternalServerError)

		return
	}
	defer rows.Close()

	accounts := []Account{}

	for rows.Next() {
		var acc Account

		err := rows.Scan(&acc.ID, &acc.Name, &acc.OffBudget, &acc.Category, &acc.Adapter, &acc.CreatedAt, &acc.UpdatedAt)
		if err != nil {
			slog.Error("error scanning accounts row from database", "error", err)
			buildErrorResponse(w, err.Error(), http.StatusInternalServerError)

			return
		}

		accounts = append(accounts, acc)
	}

	if err := rows.Err(); err != nil {
		slog.Error("error reading accounts rows from database", "error", err)
		buildErrorResponse(w, err.Error(), http.StatusInternalServerError)

		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	err = json.NewEncoder(w).Encode(accounts)
	if err != nil {
		slog.Error("error encoding accounts response", "error", err)
	}
}
