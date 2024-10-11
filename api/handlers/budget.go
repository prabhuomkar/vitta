package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	uuid "github.com/google/uuid"
)

type (
	// Group model.
	Group struct {
		ID        uuid.UUID `json:"id"`
		Name      string    `json:"name"`
		Notes     string    `json:"notes,omitempty"`
		CreatedAt time.Time `json:"createdAt"`
		UpdatedAt time.Time `json:"updatedAt"`
	}

	// Category model.
	Category struct {
		ID        uuid.UUID `json:"id"`
		GroupID   uuid.UUID `json:"groupId"`
		Name      string    `json:"name"`
		Notes     string    `json:"notes,omitempty"`
		CreatedAt time.Time `json:"createdAt"`
		UpdatedAt time.Time `json:"updatedAt"`
	}

	// Budget model.
	Budget struct {
		ID         uuid.UUID `json:"id"`
		CategoryID uuid.UUID `json:"categoryId"`
		Year       uint16    `json:"year"`
		Month      uint8     `json:"month"`
		Budgeted   float64   `json:"budgeted"`
		CreatedAt  time.Time `json:"createdAt"`
		UpdatedAt  time.Time `json:"updatedAt"`
	}

	// BudgetResult model.
	BudgetResult struct {
		ID           uuid.UUID `json:"id"`
		Budgeted     float64   `json:"budgeted"`
		Spent        float64   `json:"spent"`
		Year         uint16    `json:"year"`
		Month        uint8     `json:"month"`
		CategoryID   uuid.UUID `json:"categoryId"`
		CategoryName string    `json:"categoryName"`
		GroupID      uuid.UUID `json:"groupId"`
		GroupName    string    `json:"groupName"`
	}
)

const (
	queryCreateGroup = `INSERT into groups (id, name, notes, created_at, updated_at)` +
		` VALUES ($1, $2, $3, $4, $5)`
	queryUpdateGroup    = `UPDATE groups SET name=$1, notes=$2, updated_at=$3 WHERE id=$4`
	queryDeleteGroup    = `DELETE FROM groups WHERE id=$1`
	queryCreateCategory = `INSERT into categories (id, group_id, name, notes, created_at, updated_at)` +
		` VALUES ($1, $2, $3, $4, $5, $6)`
	queryUpdateCategory = `UPDATE categories SET name=$1, notes=$2, updated_at=$3 WHERE id=$4`
	queryDeleteCategory = `DELETE FROM categories WHERE id=$1`
	querySetBudget      = `INSERT into budgets (id, category_id, year, month, budgeted, created_at,` +
		` updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7) ON CONFLICT (year, month, category_id) DO UPDATE SET budgeted=$5`
	queryGetBudget = `SELECT budgets.id AS id, COALESCE(budgets.budgeted, 0) AS budgeted, COALESCE(t.spent, 0) AS spent,` +
		` budgets.year AS year, budgets.month AS month, cg.id AS category_id, cg.name AS category_name,` +
		` cg.group_id AS group_id, cg.group_name AS group_name FROM budgets RIGHT JOIN` +
		` (SELECT c.*, g.name as group_name FROM categories AS c LEFT JOIN groups as g ON c.group_id = g.id` +
		`) AS cg ON budgets.category_id = cg.id LEFT JOIN (SELECT category_id, SUM(credit) AS total_credit,` +
		` SUM(debit) AS total_debit, SUM(credit - debit) AS spent FROM transactions GROUP BY category_id` +
		`) AS t ON cg.id = t.category_id;`
)

func (h *Handler) CreateGroup(w http.ResponseWriter, r *http.Request) {
	var group Group

	err := json.NewDecoder(r.Body).Decode(&group)
	if err != nil {
		buildErrorResponse(w, err.Error(), http.StatusBadRequest)

		return
	}

	group.ID, err = uuid.NewV7()
	if err != nil {
		buildErrorResponse(w, err.Error(), http.StatusInternalServerError)

		return
	}

	group.CreatedAt = time.Now()
	group.UpdatedAt = group.CreatedAt

	_, err = h.db.Exec(r.Context(), queryCreateGroup,
		group.ID, group.Name, group.Notes, group.CreatedAt, group.UpdatedAt)
	if err != nil {
		buildErrorResponse(w, err.Error(), http.StatusInternalServerError)

		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	err = json.NewEncoder(w).Encode(group)
	if err != nil {
		slog.Error("error encoding group response", "error", err)
	}
}

func (h *Handler) UpdateGroup(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	var group Group

	err := json.NewDecoder(r.Body).Decode(&group)
	if err != nil {
		buildErrorResponse(w, err.Error(), http.StatusBadRequest)

		return
	}

	group.ID = uuid.MustParse(id)
	group.UpdatedAt = time.Now()

	_, err = h.db.Exec(r.Context(), queryUpdateGroup,
		group.Name, group.Notes, group.UpdatedAt, group.ID)
	if err != nil {
		buildErrorResponse(w, err.Error(), http.StatusInternalServerError)

		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) DeleteGroup(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	groupID := uuid.MustParse(id)

	_, err := h.db.Exec(r.Context(), queryDeleteGroup, groupID)
	if err != nil {
		buildErrorResponse(w, err.Error(), http.StatusInternalServerError)

		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) CreateCategory(w http.ResponseWriter, r *http.Request) {
	var category Category

	err := json.NewDecoder(r.Body).Decode(&category)
	if err != nil {
		buildErrorResponse(w, err.Error(), http.StatusBadRequest)

		return
	}

	category.ID, err = uuid.NewV7()
	if err != nil {
		buildErrorResponse(w, err.Error(), http.StatusInternalServerError)

		return
	}

	category.CreatedAt = time.Now()
	category.UpdatedAt = category.CreatedAt

	_, err = h.db.Exec(r.Context(), queryCreateCategory,
		category.ID, category.GroupID, category.Name, category.Notes, category.CreatedAt, category.UpdatedAt)
	if err != nil {
		buildErrorResponse(w, err.Error(), http.StatusInternalServerError)

		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	err = json.NewEncoder(w).Encode(category)
	if err != nil {
		slog.Error("error encoding category response", "error", err)
	}
}

func (h *Handler) UpdateCategory(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	var category Category

	err := json.NewDecoder(r.Body).Decode(&category)
	if err != nil {
		buildErrorResponse(w, err.Error(), http.StatusBadRequest)

		return
	}

	category.ID = uuid.MustParse(id)
	category.UpdatedAt = time.Now()

	_, err = h.db.Exec(r.Context(), queryUpdateCategory,
		category.Name, category.Notes, category.UpdatedAt, category.ID)
	if err != nil {
		buildErrorResponse(w, err.Error(), http.StatusInternalServerError)

		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) DeleteCategory(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	categoryID := uuid.MustParse(id)

	_, err := h.db.Exec(r.Context(), queryDeleteCategory, categoryID)
	if err != nil {
		buildErrorResponse(w, err.Error(), http.StatusInternalServerError)

		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) SetBudget(w http.ResponseWriter, r *http.Request) {
	var budget Budget

	err := json.NewDecoder(r.Body).Decode(&budget)
	if err != nil {
		buildErrorResponse(w, err.Error(), http.StatusBadRequest)

		return
	}

	budget.ID, err = uuid.NewV7()
	if err != nil {
		buildErrorResponse(w, err.Error(), http.StatusInternalServerError)

		return
	}

	budget.CreatedAt = time.Now()
	budget.UpdatedAt = budget.CreatedAt

	_, err = h.db.Exec(r.Context(), querySetBudget,
		budget.ID, budget.CategoryID, budget.Year, budget.Month, budget.Budgeted, budget.CreatedAt, budget.UpdatedAt)
	if err != nil {
		buildErrorResponse(w, err.Error(), http.StatusInternalServerError)

		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	err = json.NewEncoder(w).Encode(budget)
	if err != nil {
		slog.Error("error encoding budget response", "error", err)
	}
}

func (h *Handler) GetBudget(w http.ResponseWriter, r *http.Request) {
	rows, err := h.db.Query(r.Context(), queryGetBudget)
	if err != nil {
		buildErrorResponse(w, err.Error(), http.StatusInternalServerError)

		return
	}
	defer rows.Close()

	budgets := []BudgetResult{}

	for rows.Next() {
		var budget BudgetResult

		err := rows.Scan(&budget.ID, &budget.Budgeted, &budget.Spent, &budget.Year, &budget.Month,
			&budget.CategoryID, &budget.CategoryName, &budget.GroupID, &budget.GroupName)
		if err != nil {
			buildErrorResponse(w, err.Error(), http.StatusInternalServerError)

			return
		}

		budgets = append(budgets, budget)
	}

	if err := rows.Err(); err != nil {
		buildErrorResponse(w, err.Error(), http.StatusInternalServerError)

		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	err = json.NewEncoder(w).Encode(budgets)
	if err != nil {
		slog.Error("error encoding budgets response", "error", err)
	}
}
