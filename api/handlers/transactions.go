package handlers

import (
	"fmt"
	"log/slog"
	"net/http"
)

func (h *Handler) CreateTransaction(w http.ResponseWriter, r *http.Request) {
	slog.Info("request url", "url", r.URL.Path)
	fmt.Fprintln(w, "CreateTransaction")
}

func (h *Handler) UpdateTransaction(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	slog.Info("UpdateTransaction", "id", id)
	fmt.Fprintln(w, "UpdateTransaction")
}

func (h *Handler) DeleteTransaction(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	slog.Info("DeleteTransaction", "id", id)
	fmt.Fprintln(w, "DeleteTransaction")
}

func (h *Handler) GetTransactions(w http.ResponseWriter, r *http.Request) {
	slog.Info("request url", "url", r.URL.Path)
	fmt.Fprintln(w, "GetTransactions")
}

func (h *Handler) ImportTransactions(w http.ResponseWriter, r *http.Request) {
	slog.Info("request url", "url", r.URL.Path)
	fmt.Fprintln(w, "ImportTransactions")
}

func (h *Handler) CreatePayee(w http.ResponseWriter, r *http.Request) {
	slog.Info("request url", "url", r.URL.Path)
	fmt.Fprintln(w, "CreatePayee")
}

func (h *Handler) UpdatePayee(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	slog.Info("UpdatePayee", "id", id)
	fmt.Fprintln(w, "UpdatePayee")
}

func (h *Handler) DeletePayee(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	slog.Info("DeletePayee", "id", id)
	fmt.Fprintln(w, "DeletePayee")
}

func (h *Handler) GetPayees(w http.ResponseWriter, r *http.Request) {
	slog.Info("request url", "url", r.URL.Path)
	fmt.Fprintln(w, "GetPayees")
}
