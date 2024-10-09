package handlers

import (
	"fmt"
	"log/slog"
	"net/http"
)

func (h *Handler) CreateGroup(w http.ResponseWriter, r *http.Request) {
	slog.Info("request url", "url", r.URL.Path)
	fmt.Fprintln(w, "CreateGroup")
}

func (h *Handler) UpdateGroup(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	slog.Info("UpdateGroup", "id", id)
	fmt.Fprintln(w, "UpdateGroup")
}

func (h *Handler) DeleteGroup(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	slog.Info("DeleteGroup", "id", id)
	fmt.Fprintln(w, "DeleteGroup")
}

func (h *Handler) CreateCategory(w http.ResponseWriter, r *http.Request) {
	slog.Info("request url", "url", r.URL.Path)
	fmt.Fprintln(w, "CreateCategory")
}

func (h *Handler) UpdateCategory(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	slog.Info("UpdateCategory", "id", id)
	fmt.Fprintln(w, "UpdateCategory")
}

func (h *Handler) DeleteCategory(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	slog.Info("DeleteCategory", "id", id)
	fmt.Fprintln(w, "DeleteCategory")
}

func (h *Handler) SetBudget(w http.ResponseWriter, r *http.Request) {
	slog.Info("request url", "url", r.URL.Path)
	fmt.Fprintln(w, "SetBudget")
}

func (h *Handler) GetBudget(w http.ResponseWriter, r *http.Request) {
	slog.Info("request url", "url", r.URL.Path)
	fmt.Fprintln(w, "GetBudget")
}
