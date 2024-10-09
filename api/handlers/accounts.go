package handlers

import (
	"fmt"
	"log/slog"
	"net/http"
)

func (h *Handler) CreateAccount(w http.ResponseWriter, r *http.Request) {
	slog.Info("request url", "url", r.URL.Path)
	fmt.Fprintln(w, "CreateAccount")
}

func (h *Handler) UpdateAccount(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	slog.Info("UpdateAccount", "id", id)
	fmt.Fprintln(w, "UpdateAccount")
}

func (h *Handler) DeleteAccount(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	slog.Info("DeleteAccount", "id", id)
	fmt.Fprintln(w, "DeleteAccount")
}

func (h *Handler) GetAccounts(w http.ResponseWriter, r *http.Request) {
	slog.Info("request url", "url", r.URL.Path)
	fmt.Fprintln(w, "GetAccounts")
}
