package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"
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
		Budgeted      float64    `json:"budgeted"`
		Spent         float64    `json:"spent"`
		Year          uint16     `json:"year"`
		Month         uint8      `json:"month"`
		CategoryID    *uuid.UUID `json:"categoryId"`
		CategoryName  *string    `json:"categoryName"`
		CategoryNotes *string    `json:"categoryNotes"`
		GroupID       uuid.UUID  `json:"groupId"`
		GroupName     string     `json:"groupName"`
		GroupNotes    string     `json:"groupNotes"`
	}
)

const (
	queryCreateGroup = `INSERT INTO groups (id, name, notes, created_at, updated_at)` +
		` VALUES ($1, $2, $3, $4, $5)`
	queryUpdateGroup    = `UPDATE groups SET name=$1, notes=$2, updated_at=$3 WHERE id=$4`
	queryDeleteGroup    = `DELETE FROM groups WHERE id=$1`
	queryGetTotalGroups = `SELECT COUNT(*) as total FROM groups WHERE (name ILIKE '%' ||` +
		` COALESCE(NULLIF($1, ''), '') || '%')`
	queryGetGroups = `SELECT * FROM groups WHERE (name ILIKE '%' || COALESCE(NULLIF($1, ''), '')` +
		` || '%') ORDER BY created_at ASC`
	queryCreateCategory = `INSERT INTO categories (id, group_id, name, notes, created_at, updated_at)` +
		` VALUES ($1, $2, $3, $4, $5, $6)`
	queryUpdateCategory     = `UPDATE categories SET name=$1, notes=$2, group_id=$3, updated_at=$4 WHERE id=$5`
	queryDeleteCategory     = `DELETE FROM categories WHERE id=$1`
	queryGetTotalCategories = `SELECT COUNT(*) as total FROM categories WHERE (name ILIKE '%' ||` +
		` COALESCE(NULLIF($1, ''), '') || '%')`
	queryGetCategories = `SELECT * FROM categories WHERE (name ILIKE '%' || COALESCE(NULLIF($1, ''), '')` +
		` || '%') ORDER BY created_at ASC`
	querySetBudget = `INSERT INTO budgets (id, category_id, year, month, budgeted, created_at,` +
		` updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7) ON CONFLICT (year, month, category_id) DO UPDATE SET budgeted=$5`
	queryGetBudget = `SELECT COALESCE(budgets.budgeted, 0) AS budgeted, COALESCE(t.spent, 0) AS spent,` +
		` COALESCE(budgets.year, $1) AS year, COALESCE(budgets.month, $2) AS month, cg.id AS category_id,` +
		` cg.name AS category_name, cg.notes as category_notes, g.id AS group_id, g.name AS group_name,` +
		` g.notes as group_notes FROM groups AS g LEFT JOIN categories AS cg` +
		` ON cg.group_id = g.id LEFT JOIN budgets ON budgets.category_id = cg.id AND budgets.year = $1 AND` +
		` budgets.month = $2 LEFT JOIN(SELECT category_id, SUM(credit) AS total_credit, SUM(debit) AS total_debit,` +
		` SUM(credit - debit) AS spent FROM transactions WHERE EXTRACT(YEAR FROM cleared_at) = $1 AND` +
		` EXTRACT(MONTH FROM cleared_at) = $2 GROUP BY category_id) AS t ON cg.id = t.category_id` +
		` ORDER BY g.created_at ASC, cg.created_at ASC`
)

func (h *Handler) CreateGroup(w http.ResponseWriter, r *http.Request) {
	var group Group

	err := json.NewDecoder(r.Body).Decode(&group)
	if err != nil {
		slog.Error("error decoding create group request", "error", err)
		buildErrorResponse(w, err.Error(), http.StatusBadRequest)

		return
	}

	group.ID, err = uuid.NewV7()
	if err != nil {
		slog.Error("error creating group id", "error", err)
		buildErrorResponse(w, err.Error(), http.StatusInternalServerError)

		return
	}

	group.CreatedAt = time.Now()
	group.UpdatedAt = group.CreatedAt

	_, err = h.db.Exec(r.Context(), queryCreateGroup,
		group.ID, group.Name, group.Notes, group.CreatedAt, group.UpdatedAt)
	if err != nil {
		slog.Error("error creating group in database", "error", err)
		buildErrorResponse(w, err.Error(), http.StatusInternalServerError)

		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	err = json.NewEncoder(w).Encode(group)
	if err != nil {
		slog.Error("error encoding group response", "error", err)
	}
}

func (h *Handler) UpdateGroup(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	groupID, err := uuid.Parse(id)
	if err != nil {
		slog.Error("error parsing group id", "error", err)
		buildErrorResponse(w, err.Error(), http.StatusBadRequest)

		return
	}

	var group Group

	err = json.NewDecoder(r.Body).Decode(&group)
	if err != nil {
		slog.Error("error decoding update group request", "error", err)
		buildErrorResponse(w, err.Error(), http.StatusBadRequest)

		return
	}

	group.UpdatedAt = time.Now()

	_, err = h.db.Exec(r.Context(), queryUpdateGroup,
		group.Name, group.Notes, group.UpdatedAt, groupID)
	if err != nil {
		slog.Error("error updating group in database", "error", err)
		buildErrorResponse(w, err.Error(), http.StatusInternalServerError)

		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) DeleteGroup(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	groupID, err := uuid.Parse(id)
	if err != nil {
		slog.Error("error parsing group id", "error", err)
		buildErrorResponse(w, err.Error(), http.StatusBadRequest)

		return
	}

	_, err = h.db.Exec(r.Context(), queryDeleteGroup, groupID)
	if err != nil {
		slog.Error("error deleting group in database", "error", err)
		buildErrorResponse(w, err.Error(), http.StatusInternalServerError)

		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) GetGroups(w http.ResponseWriter, r *http.Request) { //nolint: funlen
	searchQuery := r.URL.Query().Get("q")

	rows, err := h.db.Query(r.Context(), queryGetTotalGroups, searchQuery)
	if err != nil {
		slog.Error("error getting total groups from database", "error", err)
		buildErrorResponse(w, err.Error(), http.StatusInternalServerError)

		return
	}
	defer rows.Close()

	var total int

	for rows.Next() {
		err = rows.Scan(&total)
		if err != nil {
			slog.Error("error scanning total groups row from database", "error", err)
			buildErrorResponse(w, err.Error(), http.StatusInternalServerError)

			return
		}
	}

	rows, err = h.db.Query(r.Context(), queryGetGroups, searchQuery)
	if err != nil {
		slog.Error("error getting groups from database", "error", err)
		buildErrorResponse(w, err.Error(), http.StatusInternalServerError)

		return
	}
	defer rows.Close()

	groups := []Group{}

	for rows.Next() {
		var group Group

		err := rows.Scan(&group.ID, &group.Name, &group.Notes, &group.CreatedAt, &group.UpdatedAt)
		if err != nil {
			slog.Error("error scanning groups row from database", "error", err)
			buildErrorResponse(w, err.Error(), http.StatusInternalServerError)

			return
		}

		groups = append(groups, group)
	}

	if err := rows.Err(); err != nil {
		slog.Error("error reading groups rows from database", "error", err)
		buildErrorResponse(w, err.Error(), http.StatusInternalServerError)

		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	err = json.NewEncoder(w).Encode(map[string]interface{}{"total": total, "groups": groups})
	if err != nil {
		slog.Error("error encoding groups response", "error", err)
	}
}

func (h *Handler) CreateCategory(w http.ResponseWriter, r *http.Request) {
	var category Category

	err := json.NewDecoder(r.Body).Decode(&category)
	if err != nil {
		slog.Error("error decoding create category request", "error", err)
		buildErrorResponse(w, err.Error(), http.StatusBadRequest)

		return
	}

	category.ID, err = uuid.NewV7()
	if err != nil {
		slog.Error("error creating category id", "error", err)
		buildErrorResponse(w, err.Error(), http.StatusInternalServerError)

		return
	}

	category.CreatedAt = time.Now()
	category.UpdatedAt = category.CreatedAt

	_, err = h.db.Exec(r.Context(), queryCreateCategory,
		category.ID, category.GroupID, category.Name, category.Notes, category.CreatedAt, category.UpdatedAt)
	if err != nil {
		slog.Error("error creating category in database", "error", err)
		buildErrorResponse(w, err.Error(), http.StatusInternalServerError)

		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	err = json.NewEncoder(w).Encode(category)
	if err != nil {
		slog.Error("error encoding category response", "error", err)
	}
}

func (h *Handler) UpdateCategory(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	categoryID, err := uuid.Parse(id)
	if err != nil {
		slog.Error("error parsing category id", "error", err)
		buildErrorResponse(w, err.Error(), http.StatusBadRequest)

		return
	}

	var category Category

	err = json.NewDecoder(r.Body).Decode(&category)
	if err != nil {
		slog.Error("error decoding update category request", "error", err)
		buildErrorResponse(w, err.Error(), http.StatusBadRequest)

		return
	}

	category.UpdatedAt = time.Now()

	_, err = h.db.Exec(r.Context(), queryUpdateCategory,
		category.Name, category.Notes, category.GroupID, category.UpdatedAt, categoryID)
	if err != nil {
		slog.Error("error udpating category in database", "error", err)
		buildErrorResponse(w, err.Error(), http.StatusInternalServerError)

		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) DeleteCategory(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	categoryID, err := uuid.Parse(id)
	if err != nil {
		slog.Error("error parsing category id", "error", err)
		buildErrorResponse(w, err.Error(), http.StatusBadRequest)

		return
	}

	_, err = h.db.Exec(r.Context(), queryDeleteCategory, categoryID)
	if err != nil {
		slog.Error("error deleting category in database", "error", err)
		buildErrorResponse(w, err.Error(), http.StatusInternalServerError)

		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) GetCategories(w http.ResponseWriter, r *http.Request) { //nolint: funlen
	searchQuery := r.URL.Query().Get("q")

	rows, err := h.db.Query(r.Context(), queryGetTotalCategories, searchQuery)
	if err != nil {
		slog.Error("error getting total categories from database", "error", err)
		buildErrorResponse(w, err.Error(), http.StatusInternalServerError)

		return
	}
	defer rows.Close()

	var total int

	for rows.Next() {
		err = rows.Scan(&total)
		if err != nil {
			slog.Error("error scanning total categories row from database", "error", err)
			buildErrorResponse(w, err.Error(), http.StatusInternalServerError)

			return
		}
	}

	rows, err = h.db.Query(r.Context(), queryGetCategories, searchQuery)
	if err != nil {
		slog.Error("error getting categories from database", "error", err)
		buildErrorResponse(w, err.Error(), http.StatusInternalServerError)

		return
	}
	defer rows.Close()

	categories := []Category{}

	for rows.Next() {
		var category Category

		err := rows.Scan(&category.ID, &category.GroupID, &category.Name, &category.Notes,
			&category.CreatedAt, &category.UpdatedAt)
		if err != nil {
			slog.Error("error scanning categories row from database", "error", err)
			buildErrorResponse(w, err.Error(), http.StatusInternalServerError)

			return
		}

		categories = append(categories, category)
	}

	if err := rows.Err(); err != nil {
		slog.Error("error reading categories rows from database", "error", err)
		buildErrorResponse(w, err.Error(), http.StatusInternalServerError)

		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	err = json.NewEncoder(w).Encode(map[string]interface{}{"total": total, "categories": categories})
	if err != nil {
		slog.Error("error encoding categories response", "error", err)
	}
}

func (h *Handler) SetBudget(w http.ResponseWriter, r *http.Request) {
	var budget Budget

	err := json.NewDecoder(r.Body).Decode(&budget)
	if err != nil {
		slog.Error("error decoding set budget request", "error", err)
		buildErrorResponse(w, err.Error(), http.StatusBadRequest)

		return
	}

	budget.ID, err = uuid.NewV7()
	if err != nil {
		slog.Error("error creating budget id", "error", err)
		buildErrorResponse(w, err.Error(), http.StatusInternalServerError)

		return
	}

	budget.CreatedAt = time.Now()
	budget.UpdatedAt = budget.CreatedAt

	_, err = h.db.Exec(r.Context(), querySetBudget,
		budget.ID, budget.CategoryID, budget.Year, budget.Month, budget.Budgeted, budget.CreatedAt, budget.UpdatedAt)
	if err != nil {
		slog.Error("error setting budget in database", "error", err)
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
	yearQ := r.URL.Query().Get("year")
	monthQ := r.URL.Query().Get("month")

	year, err := strconv.Atoi(yearQ)
	if err != nil {
		slog.Error("error converting year to int", "error", err)
		buildErrorResponse(w, err.Error(), http.StatusBadRequest)

		return
	}

	month, err := strconv.Atoi(monthQ)
	if err != nil {
		slog.Error("error converting month to int", "error", err)
		buildErrorResponse(w, err.Error(), http.StatusBadRequest)

		return
	}

	rows, err := h.db.Query(r.Context(), queryGetBudget, year, month)
	if err != nil {
		slog.Error("error getting budgets from database", "error", err)
		buildErrorResponse(w, err.Error(), http.StatusInternalServerError)

		return
	}
	defer rows.Close()

	budgets := []BudgetResult{}

	for rows.Next() {
		var budget BudgetResult

		err := rows.Scan(&budget.Budgeted, &budget.Spent, &budget.Year, &budget.Month, &budget.CategoryID,
			&budget.CategoryName, &budget.CategoryNotes, &budget.GroupID, &budget.GroupName, &budget.GroupNotes)
		if err != nil {
			slog.Error("error scanning budgets row from database", "error", err)
			buildErrorResponse(w, err.Error(), http.StatusInternalServerError)

			return
		}

		budgets = append(budgets, budget)
	}

	if err := rows.Err(); err != nil {
		slog.Error("error reading budgets rows from database", "error", err)
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
