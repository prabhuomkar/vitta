package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"time"

	uuid "github.com/google/uuid"
)

type (
	// Rules model.
	Rules struct {
		Includes   []string `json:"includes,omitempty"`
		Excludes   []string `json:"excludes,omitempty"`
		StartsWith []string `json:"startsWith,omitempty"`
		EndsWith   []string `json:"endsWith,omitempty"`
	}

	// Payee model.
	Payee struct {
		ID             uuid.UUID  `json:"id"`
		Name           string     `json:"name"`
		Rules          *Rules     `json:"rules,omitempty"`
		AutoCategoryID *uuid.UUID `json:"autoCategoryId,omitempty"`
		CreatedAt      time.Time  `json:"createdAt"`
		UpdatedAt      time.Time  `json:"updatedAt"`
	}
)

const (
	queryCreatePayee = `INSERT INTO payees (id, name, rules, auto_category_id, created_at, updated_at)` +
		` VALUES ($1, $2, $3, $4, $5, $6)`
	queryUpdatePayee    = `UPDATE payees SET name=$1, rules=$2, auto_category_id=$3, updated_at=$4 WHERE id=$5`
	queryDeletePayee    = `DELETE FROM payees WHERE id=$1`
	queryGetTotalPayees = `SELECT COUNT(*) as total FROM payees WHERE (name ILIKE '%' ||` +
		` COALESCE(NULLIF($1, ''), '') || '%')`
	queryGetPayees = `SELECT * FROM payees WHERE (name ILIKE '%' || COALESCE(NULLIF($1, ''), '')` +
		` || '%') ORDER BY created_at DESC`
)

func (h *Handler) CreatePayee(w http.ResponseWriter, r *http.Request) { //nolint: cyclop
	updateTransactionsQuery := r.URL.Query().Get("updateTransactions")

	updateTransactions, err := strconv.ParseBool(updateTransactionsQuery)
	if err != nil {
		updateTransactions = false
	}

	var payee Payee

	err = json.NewDecoder(r.Body).Decode(&payee)
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
		payee.ID, payee.Name, payee.Rules, payee.AutoCategoryID, payee.CreatedAt, payee.UpdatedAt)
	if err != nil {
		slog.Error("error creating payee in database", "error", err)
		buildErrorResponse(w, err.Error(), http.StatusInternalServerError)

		return
	}

	if updateTransactions && (payee.Rules != nil || payee.AutoCategoryID != nil) {
		getPayeeCategory, err := h.assignPayeeAndCategory(r.Context(), []Payee{payee})
		if err != nil {
			slog.Error("error creating payee category assigner", "error", err)
			buildErrorResponse(w, err.Error(), http.StatusInternalServerError)

			return
		}

		err = h.updateTransactions(r.Context(), getPayeeCategory)
		if err != nil {
			slog.Error("error updating transactions", "error", err)
			buildErrorResponse(w, err.Error(), http.StatusInternalServerError)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	err = json.NewEncoder(w).Encode(payee)
	if err != nil {
		slog.Error("error encoding payee response", "error", err)
	}
}

func (h *Handler) UpdatePayee(w http.ResponseWriter, r *http.Request) {
	updateTransactionsQuery := r.URL.Query().Get("updateTransactions")

	updateTransactions, err := strconv.ParseBool(updateTransactionsQuery)
	if err != nil {
		updateTransactions = false
	}

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
		payee.Name, payee.Rules, payee.AutoCategoryID, payee.UpdatedAt, payeeID)
	if err != nil {
		slog.Error("error updating payee in database", "error", err)
		buildErrorResponse(w, err.Error(), http.StatusInternalServerError)

		return
	}

	if updateTransactions && (payee.Rules != nil || payee.AutoCategoryID != nil) {
		getPayeeCategory, err := h.assignPayeeAndCategory(r.Context(), []Payee{payee})
		if err != nil {
			slog.Error("error creating payee category assigner", "error", err)
			buildErrorResponse(w, err.Error(), http.StatusInternalServerError)

			return
		}

		err = h.updateTransactions(r.Context(), getPayeeCategory)
		if err != nil {
			slog.Error("error updating transactions", "error", err)
			buildErrorResponse(w, err.Error(), http.StatusInternalServerError)
		}
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

func (h *Handler) GetPayees(w http.ResponseWriter, r *http.Request) { //nolint: funlen
	searchQuery := r.URL.Query().Get("q")

	rows, err := h.db.Query(r.Context(), queryGetTotalPayees, searchQuery)
	if err != nil {
		slog.Error("error getting total payees from database", "error", err)
		buildErrorResponse(w, err.Error(), http.StatusInternalServerError)

		return
	}
	defer rows.Close()

	var total int

	for rows.Next() {
		err = rows.Scan(&total)
		if err != nil {
			slog.Error("error scanning total payees row from database", "error", err)
			buildErrorResponse(w, err.Error(), http.StatusInternalServerError)

			return
		}
	}

	rows, err = h.db.Query(r.Context(), queryGetPayees, searchQuery)
	if err != nil {
		slog.Error("error getting payees from database", "error", err)
		buildErrorResponse(w, err.Error(), http.StatusInternalServerError)

		return
	}
	defer rows.Close()

	payees := []Payee{}

	for rows.Next() {
		var payee Payee

		err := rows.Scan(&payee.ID, &payee.Name, &payee.Rules, &payee.AutoCategoryID, &payee.CreatedAt, &payee.UpdatedAt)
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

	err = json.NewEncoder(w).Encode(map[string]interface{}{"total": total, "payees": payees})
	if err != nil {
		slog.Error("error encoding payees response", "error", err)
	}
}

func (h *Handler) assignPayeeAndCategory(ctx context.Context, payees []Payee) ( //nolint: funlen,gocognit,cyclop
	func(string) (*uuid.UUID, *uuid.UUID), error,
) {
	if len(payees) == 0 {
		rows, err := h.db.Query(ctx, queryGetPayees, "")
		if err != nil {
			slog.Error("error getting payees from database", "error", err)

			return nil, fmt.Errorf("error getting payees: %w", err)
		}
		defer rows.Close()

		for rows.Next() {
			var payee Payee

			err := rows.Scan(&payee.ID, &payee.Name, &payee.Rules, &payee.AutoCategoryID, &payee.CreatedAt, &payee.UpdatedAt)
			if err != nil {
				slog.Error("error scanning payees row from database", "error", err)

				return nil, fmt.Errorf("error scanning payees row: %w", err)
			}

			payees = append(payees, payee)
		}

		if err := rows.Err(); err != nil {
			slog.Error("error reading payees rows from database", "error", err)

			return nil, fmt.Errorf("error reading payees rows: %w", err)
		}
	}

	return func(input string) (*uuid.UUID, *uuid.UUID) {
		for _, payee := range payees {
			if payee.Rules != nil { //nolint: nestif
				ri, re, rsw, rew := false, false, false, false

				input = strings.ToLower(input)

				for _, includes := range payee.Rules.Includes {
					if strings.Contains(input, strings.ToLower(includes)) {
						ri = true

						break
					}
				}

				if len(payee.Rules.Excludes) > 0 {
					for _, excludes := range payee.Rules.Excludes {
						if !strings.Contains(input, strings.ToLower(excludes)) {
							re = true

							break
						}
					}
				} else {
					re = true
				}

				for _, startsWith := range payee.Rules.StartsWith {
					if strings.HasPrefix(input, strings.ToLower(startsWith)) {
						rsw = true

						break
					}
				}

				for _, endsWith := range payee.Rules.EndsWith {
					if strings.HasSuffix(input, strings.ToLower(endsWith)) {
						rew = true

						break
					}
				}

				if (ri && re) || rsw || rew {
					return &payee.ID, payee.AutoCategoryID
				}
			}
		}

		return nil, nil
	}, nil
}
