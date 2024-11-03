package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	uuid "github.com/google/uuid"
)

// Payee model.
type Payee struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

const (
	queryCreatePayee = `INSERT INTO payees (id, name, created_at, updated_at) VALUES ($1, $2, $3, $4)`
	queryUpdatePayee = `UPDATE payees SET name=$1, updated_at=$2 WHERE id=$3`
	queryDeletePayee = `DELETE FROM payees WHERE id=$1`
	queryGetPayees   = `SELECT * FROM payees WHERE (name ILIKE '%' || COALESCE(NULLIF($1, ''), '')` +
		` || '%') ORDER BY created_at DESC`
)

func (h *Handler) CreatePayee(w http.ResponseWriter, r *http.Request) {
	var payee Payee

	err := json.NewDecoder(r.Body).Decode(&payee)
	if err != nil {
		slog.Error("error decoding create payee request", "error", err)
		buildErrorResponse(w, err.Error(), http.StatusBadRequest)

		return
	}

	payee.ID, err = uuid.NewV7()
	if err != nil {
		slog.Error("error creating new payee id", "error", err)
		buildErrorResponse(w, err.Error(), http.StatusInternalServerError)

		return
	}

	payee.CreatedAt = time.Now()
	payee.UpdatedAt = payee.CreatedAt

	_, err = h.db.Exec(r.Context(), queryCreatePayee,
		payee.ID, payee.Name, payee.CreatedAt, payee.UpdatedAt)
	if err != nil {
		slog.Error("error creating payee in database", "error", err)
		buildErrorResponse(w, err.Error(), http.StatusInternalServerError)

		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	err = json.NewEncoder(w).Encode(payee)
	if err != nil {
		slog.Error("error encoding payee response", "error", err)
	}
}

func (h *Handler) UpdatePayee(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	payeeID, err := uuid.Parse(id)
	if err != nil {
		slog.Error("error parsing payee id", "error", err)
		buildErrorResponse(w, err.Error(), http.StatusBadRequest)

		return
	}

	var payee Payee

	err = json.NewDecoder(r.Body).Decode(&payee)
	if err != nil {
		slog.Error("error decoding update payee request", "error", err)
		buildErrorResponse(w, err.Error(), http.StatusBadRequest)

		return
	}

	payee.UpdatedAt = time.Now()

	_, err = h.db.Exec(r.Context(), queryUpdatePayee,
		payee.Name, payee.UpdatedAt, payeeID)
	if err != nil {
		slog.Error("error updating payee in database", "error", err)
		buildErrorResponse(w, err.Error(), http.StatusInternalServerError)

		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) DeletePayee(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	payeeID, err := uuid.Parse(id)
	if err != nil {
		slog.Error("error parsing payee id", "error", err)
		buildErrorResponse(w, err.Error(), http.StatusBadRequest)

		return
	}

	_, err = h.db.Exec(r.Context(), queryDeletePayee, payeeID)
	if err != nil {
		slog.Error("error deleting payee in database", "error", err)
		buildErrorResponse(w, err.Error(), http.StatusInternalServerError)

		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) GetPayees(w http.ResponseWriter, r *http.Request) {
	searchQuery := r.URL.Query().Get("q")

	rows, err := h.db.Query(r.Context(), queryGetPayees, searchQuery)
	if err != nil {
		slog.Error("error getting payees from database", "error", err)
		buildErrorResponse(w, err.Error(), http.StatusInternalServerError)

		return
	}
	defer rows.Close()

	payees := []Payee{}

	for rows.Next() {
		var payee Payee

		err := rows.Scan(&payee.ID, &payee.Name, &payee.CreatedAt, &payee.UpdatedAt)
		if err != nil {
			slog.Error("error scanning payees row from database", "error", err)
			buildErrorResponse(w, err.Error(), http.StatusInternalServerError)

			return
		}

		payees = append(payees, payee)
	}

	if err := rows.Err(); err != nil {
		slog.Error("error reading payees rows from database", "error", err)
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
