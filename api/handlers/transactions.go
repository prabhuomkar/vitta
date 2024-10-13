package handlers

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log/slog"
	"mime/multipart"
	"net/http"
	"strings"
	"time"
	"vitta/adapters"

	"github.com/extrame/xls"
	uuid "github.com/google/uuid"
	"github.com/xuri/excelize/v2"
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
		ClearedAt    *time.Time `json:"clearedAt,omitempty"`
		CreatedAt    time.Time  `json:"createdAt"`
		UpdatedAt    time.Time  `json:"updatedAt"`
	}
)

const (
	queryCreateTransaction = `INSERT INTO transactions (id, account_id, category_id, payee_id, credit,` +
		` debit, name, notes, cleared_at, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`
	queryUpdateTransaction = `UPDATE transactions SET account_id=$1, category_id=$2, payee_id=$3,` +
		` credit=$4, debit=$5, name=$6, notes=$7, cleared_at=$8, updated_at=$9 WHERE id=$10`
	queryDeleteTransaction = `DELETE FROM transactions WHERE id=$1`
	queryGetTransactions   = `SELECT t.*, c.name as category_name, p.name as payee_name FROM transactions AS t` +
		` LEFT JOIN categories AS c ON t.category_id = c.id LEFT JOIN payees AS p ON t.payee_id = p.id`

	queryCreatePayee = `INSERT INTO payees (id, name, created_at, updated_at) VALUES ($1, $2, $3, $4)`
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
		nil, transaction.CreatedAt, transaction.UpdatedAt)
	if err != nil {
		buildErrorResponse(w, err.Error(), http.StatusInternalServerError)

		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	err = json.NewEncoder(w).Encode(transaction)
	if err != nil {
		slog.Error("error encoding transaction response", "error", err)
	}
}

func (h *Handler) UpdateTransaction(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	transactionID, err := uuid.Parse(id)
	if err != nil {
		buildErrorResponse(w, err.Error(), http.StatusBadRequest)

		return
	}

	var transaction Transaction

	err = json.NewDecoder(r.Body).Decode(&transaction)
	if err != nil {
		buildErrorResponse(w, err.Error(), http.StatusBadRequest)

		return
	}

	transaction.UpdatedAt = time.Now()

	_, err = h.db.Exec(r.Context(), queryUpdateTransaction,
		transaction.AccountID, transaction.CategoryID, transaction.PayeeID,
		transaction.Credit, transaction.Debit, transaction.Name, transaction.Notes,
		transaction.ClearedAt, transaction.UpdatedAt, transactionID)
	if err != nil {
		buildErrorResponse(w, err.Error(), http.StatusInternalServerError)

		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) DeleteTransaction(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	transactionID, err := uuid.Parse(id)
	if err != nil {
		buildErrorResponse(w, err.Error(), http.StatusBadRequest)

		return
	}

	_, err = h.db.Exec(r.Context(), queryDeleteTransaction, transactionID)
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

func (h *Handler) ImportTransactions(w http.ResponseWriter, r *http.Request) { //nolint: funlen,cyclop
	accountIDQuery := r.URL.Query().Get("accountId")
	accountCategory := r.URL.Query().Get("accountCategory")
	adapter := r.URL.Query().Get("adapter")

	accountID, err := uuid.Parse(accountIDQuery)
	if err != nil {
		buildErrorResponse(w, err.Error(), http.StatusBadRequest)

		return
	}

	err = r.ParseMultipartForm(h.cfg.UploadMemoryLimit)
	if err != nil {
		slog.Error("error parsing multipart form", "error", err)
		buildErrorResponse(w, err.Error(), http.StatusBadRequest)

		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		slog.Error("error getting file", "error", err)
		buildErrorResponse(w, err.Error(), http.StatusBadRequest)

		return
	}
	defer file.Close()

	slog.Info("uploaded file", "name", header.Filename, "size", header.Size)

	rows, err := h.getDataRows(header.Filename, file)
	if err != nil {
		slog.Error("error getting data rows", "error", err)
		buildErrorResponse(w, err.Error(), http.StatusInternalServerError)

		return
	}

	bankAdapter := adapters.New(adapter, accountCategory, nil)
	adapterTransactions := bankAdapter.GetTransactions(rows)

	tx, err := h.db.Begin(r.Context())
	if err != nil {
		slog.Error("error creating database txn", "error", err)
		buildErrorResponse(w, err.Error(), http.StatusInternalServerError)

		return
	}

	transactionTime := time.Now()

	for _, adapterTransaction := range adapterTransactions {
		transactionID, err := uuid.NewV7()
		if err != nil {
			slog.Error("error creating transaction id", "error", err)
			buildErrorResponse(w, err.Error(), http.StatusInternalServerError)

			return
		}

		_, err = tx.Exec(r.Context(), queryCreateTransaction, transactionID, accountID, nil, nil,
			adapterTransaction.Credit, adapterTransaction.Debit, "imported transaction", adapterTransaction.Remarks,
			adapterTransaction.Date, transactionTime, transactionTime)
		if err != nil {
			slog.Error("error inserting transaction", "error", err)

			err = tx.Rollback(r.Context())
			if err != nil {
				slog.Error("error rolling back database txn", "error", err)
			}

			buildErrorResponse(w, err.Error(), http.StatusInternalServerError)

			return
		}
	}

	if err := tx.Commit(r.Context()); err != nil {
		slog.Error("error committing database txn", "error", err)
		buildErrorResponse(w, err.Error(), http.StatusInternalServerError)

		return
	}
}

func (h *Handler) getDataRows(fileName string, file multipart.File) ([][]string, error) {
	var (
		rows [][]string
		err  error
	)

	switch {
	case strings.HasSuffix(strings.ToLower(fileName), ".csv"):
		reader := csv.NewReader(file)

		rows, err = reader.ReadAll()
		if err != nil {
			slog.Error("error reading rows", "error", err)

			return nil, fmt.Errorf("error reading rows: %w", err)
		}

	case strings.HasSuffix(strings.ToLower(fileName), ".xls"):
		xlsFile, err := xls.OpenReader(file, "utf-8")
		if err != nil {
			slog.Error("error opening file", "error", err)

			return nil, fmt.Errorf("error opening file: %w", err)
		}

		sheet := xlsFile.GetSheet(0)
		for rowIndex := 0; rowIndex <= int(sheet.MaxRow); rowIndex++ {
			row := sheet.Row(rowIndex)

			rowData := []string{}

			for colIndex := range row.LastCol() {
				rowData = append(rowData, row.Col(colIndex))
			}

			rows = append(rows, rowData)
		}
	default:
		xFile, err := excelize.OpenReader(file)
		if err != nil {
			slog.Error("error opening file", "error", err)

			return nil, fmt.Errorf("error opening file: %w", err)
		}

		rows, err = xFile.GetRows(xFile.GetSheetName(0))
		if err != nil {
			slog.Error("error reading rows", "error", err)

			return nil, fmt.Errorf("error reading rows: %w", err)
		}
	}

	return rows, nil
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
		buildErrorResponse(w, err.Error(), http.StatusBadRequest)

		return
	}

	var payee Payee

	err = json.NewDecoder(r.Body).Decode(&payee)
	if err != nil {
		buildErrorResponse(w, err.Error(), http.StatusBadRequest)

		return
	}

	payee.UpdatedAt = time.Now()

	_, err = h.db.Exec(r.Context(), queryUpdatePayee,
		payee.Name, payee.UpdatedAt, payeeID)
	if err != nil {
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
		buildErrorResponse(w, err.Error(), http.StatusBadRequest)

		return
	}

	_, err = h.db.Exec(r.Context(), queryDeletePayee, payeeID)
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
