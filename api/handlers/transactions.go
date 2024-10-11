package handlers

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	uuid "github.com/google/uuid"
)

type (
	// Payee model.
	Payee struct {
		ID        uuid.UUID `json:"id"`
		Name      string    `json:"name"`
		CreatedAt time.Time `json:"createdAt"`
		UpdatedAt time.Time `json:"updatedAt"`
	}

	// Transaction model.
	Transaction struct {
		ID           uuid.UUID  `json:"id"`
		AccountID    uuid.UUID  `json:"accountId"`
		CategoryID   *uuid.UUID `json:"categoryId,omitempty"`
		CategoryName *string    `json:"categoryName,omitempty"`
		PayeeID      *uuid.UUID `json:"payeeId,omitempty"`
		PayeeName    *string    `json:"payeeName,omitempty"`
		Credit       float64    `json:"credit,omitempty"`
		Debit        float64    `json:"debit,omitempty"`
		Name         string     `json:"name"`
		Notes        string     `json:"notes,omitempty"`
		ClearedAt    time.Time  `json:"clearedAt"`
		CreatedAt    time.Time  `json:"createdAt"`
		UpdatedAt    time.Time  `json:"updatedAt"`
	}
)

const (
	queryCreateTransaction = `INSERT into transactions (id, account_id, category_id, payee_id, credit,` +
		` debit, name, notes, cleared_at, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`
	queryUpdateTransaction = `UPDATE transactions SET account_id=$1, category_id=$2, payee_id=$3,` +
		` credit=$4, debit=$5, name=$6, notes=$7, cleared_at=$8, updated_at=$9 WHERE id=$10`
	queryDeleteTransaction = `DELETE FROM transactions WHERE id=$1`
	queryGetTransactions   = `SELECT t.*, c.name as category_name, p.name as payee_name FROM transactions AS t` +
		` LEFT JOIN categories AS c ON t.category_id = c.id LEFT JOIN payees AS p ON t.payee_id = p.id`

	queryCreatePayee = `INSERT into payees (id, name, created_at, updated_at) VALUES ($1, $2, $3, $4)`
	queryUpdatePayee = `UPDATE payees SET name=$1, updated_at=$2 WHERE id=$3`
	queryDeletePayee = `DELETE FROM payees WHERE id=$1`
	queryGetPayees   = `SELECT * FROM payees`
)

func (h *Handler) CreateTransaction(w http.ResponseWriter, r *http.Request) {
	var transaction Transaction

	err := json.NewDecoder(r.Body).Decode(&transaction)
	if err != nil {
		buildErrorResponse(w, err.Error(), http.StatusBadRequest)

		return
	}

	transaction.ID, err = uuid.NewV7()
	if err != nil {
		buildErrorResponse(w, err.Error(), http.StatusInternalServerError)

		return
	}

	transaction.CreatedAt = time.Now()
	transaction.UpdatedAt = transaction.CreatedAt

	_, err = h.db.Exec(r.Context(), queryCreateTransaction,
		transaction.ID, transaction.AccountID, transaction.CategoryID, transaction.PayeeID,
		transaction.Credit, transaction.Debit, transaction.Name, transaction.Notes,
		transaction.ClearedAt, transaction.CreatedAt, transaction.UpdatedAt)
	if err != nil {
		buildErrorResponse(w, err.Error(), http.StatusInternalServerError)

		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	err = json.NewEncoder(w).Encode(transaction)
	if err != nil {
		slog.Error("error encoding transaction response", "error", err)
	}
}

func (h *Handler) UpdateTransaction(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	var transaction Transaction

	err := json.NewDecoder(r.Body).Decode(&transaction)
	if err != nil {
		buildErrorResponse(w, err.Error(), http.StatusBadRequest)

		return
	}

	transaction.ID = uuid.MustParse(id)
	transaction.UpdatedAt = time.Now()

	_, err = h.db.Exec(r.Context(), queryUpdateTransaction,
		transaction.AccountID, transaction.CategoryID, transaction.PayeeID,
		transaction.Credit, transaction.Debit, transaction.Name, transaction.Notes,
		transaction.ClearedAt, transaction.UpdatedAt, transaction.ID)
	if err != nil {
		buildErrorResponse(w, err.Error(), http.StatusInternalServerError)

		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) DeleteTransaction(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	transactionID := uuid.MustParse(id)

	_, err := h.db.Exec(r.Context(), queryDeleteTransaction, transactionID)
	if err != nil {
		buildErrorResponse(w, err.Error(), http.StatusInternalServerError)

		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) GetTransactions(w http.ResponseWriter, r *http.Request) {
	rows, err := h.db.Query(r.Context(), queryGetTransactions)
	if err != nil {
		buildErrorResponse(w, err.Error(), http.StatusInternalServerError)

		return
	}
	defer rows.Close()

	transactions := []Transaction{}

	for rows.Next() {
		var transaction Transaction

		err := rows.Scan(&transaction.ID, &transaction.AccountID, &transaction.CategoryID, &transaction.PayeeID,
			&transaction.Name, &transaction.Credit, &transaction.Debit, &transaction.Notes, &transaction.ClearedAt,
			&transaction.CreatedAt, &transaction.UpdatedAt, &transaction.CategoryName, &transaction.PayeeName)
		if err != nil {
			buildErrorResponse(w, err.Error(), http.StatusInternalServerError)

			return
		}

		transactions = append(transactions, transaction)
	}

	if err := rows.Err(); err != nil {
		buildErrorResponse(w, err.Error(), http.StatusInternalServerError)

		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	err = json.NewEncoder(w).Encode(transactions)
	if err != nil {
		slog.Error("error encoding transactions response", "error", err)
	}
}

func (h *Handler) ImportTransactions(w http.ResponseWriter, r *http.Request) {
	slog.Info("request url", "url", r.URL.Path)
	fmt.Fprintln(w, "ImportTransactions")
}

func (h *Handler) CreatePayee(w http.ResponseWriter, r *http.Request) {
	var payee Payee

	err := json.NewDecoder(r.Body).Decode(&payee)
	if err != nil {
		buildErrorResponse(w, err.Error(), http.StatusBadRequest)

		return
	}

	payee.ID, err = uuid.NewV7()
	if err != nil {
		buildErrorResponse(w, err.Error(), http.StatusInternalServerError)

		return
	}

	payee.CreatedAt = time.Now()
	payee.UpdatedAt = payee.CreatedAt

	_, err = h.db.Exec(r.Context(), queryCreatePayee,
		payee.ID, payee.Name, payee.CreatedAt, payee.UpdatedAt)
	if err != nil {
		buildErrorResponse(w, err.Error(), http.StatusInternalServerError)

		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	err = json.NewEncoder(w).Encode(payee)
	if err != nil {
		slog.Error("error encoding payee response", "error", err)
	}
}

func (h *Handler) UpdatePayee(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	var payee Payee

	err := json.NewDecoder(r.Body).Decode(&payee)
	if err != nil {
		buildErrorResponse(w, err.Error(), http.StatusBadRequest)

		return
	}

	payee.ID = uuid.MustParse(id)
	payee.UpdatedAt = time.Now()

	_, err = h.db.Exec(r.Context(), queryUpdatePayee,
		payee.Name, payee.UpdatedAt, payee.ID)
	if err != nil {
		buildErrorResponse(w, err.Error(), http.StatusInternalServerError)

		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) DeletePayee(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	payeeID := uuid.MustParse(id)

	_, err := h.db.Exec(r.Context(), queryDeletePayee, payeeID)
	if err != nil {
		buildErrorResponse(w, err.Error(), http.StatusInternalServerError)

		return
	}

	w.Header().Set("Content-Type", "application/json")
}

func (h *Handler) GetPayees(w http.ResponseWriter, r *http.Request) {
	rows, err := h.db.Query(r.Context(), queryGetPayees)
	if err != nil {
		buildErrorResponse(w, err.Error(), http.StatusInternalServerError)

		return
	}
	defer rows.Close()

	payees := []Payee{}

	for rows.Next() {
		var payee Payee

		err := rows.Scan(&payee.ID, &payee.Name, &payee.CreatedAt, &payee.UpdatedAt)
		if err != nil {
			buildErrorResponse(w, err.Error(), http.StatusInternalServerError)

			return
		}

		payees = append(payees, payee)
	}

	if err := rows.Err(); err != nil {
		buildErrorResponse(w, err.Error(), http.StatusInternalServerError)

		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	err = json.NewEncoder(w).Encode(payees)
	if err != nil {
		slog.Error("error encoding payees response", "error", err)
	}
}
