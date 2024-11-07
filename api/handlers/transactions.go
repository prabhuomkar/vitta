package handlers

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log/slog"
	"mime/multipart"
	"net/http"
	"strconv"
	"strings"
	"time"
	"vitta/adapters"

	"github.com/extrame/xls"
	uuid "github.com/google/uuid"
	"github.com/xuri/excelize/v2"
)

type (
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

	// TransactionsResult model.
	TransactionsResult struct {
		Total    int `json:"total"`
		Imported int `json:"imported"`
	}
)

const (
	queryCreateTransaction = `INSERT INTO transactions (id, account_id, category_id, payee_id, credit,` +
		` debit, name, notes, cleared_at, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`
	queryUpdateTransaction = `UPDATE transactions SET category_id=$2, payee_id=$3,` +
		` credit=$4, debit=$5, name=$6, notes=$7, cleared_at=$8, updated_at=$9 WHERE account_id=$1 AND id=$10`
	queryDeleteTransaction = `DELETE FROM transactions WHERE account_id=$1 AND id=$2`
	queryGetTransactions   = `SELECT t.*, c.name as category_name, p.name as payee_name FROM transactions AS t` +
		` LEFT JOIN categories AS c ON t.category_id = c.id LEFT JOIN payees AS p ON t.payee_id = p.id` +
		` WHERE t.account_id=$1 AND (name ILIKE '%' || COALESCE(NULLIF($2, ''), '') || '%') ORDER BY t.created_at DESC` +
		` LIMIT $3 OFFSET $4`
)

func (h *Handler) CreateTransaction(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	accountID, err := uuid.Parse(id)
	if err != nil {
		slog.Error("error parsing account id", "error", err)
		buildErrorResponse(w, err.Error(), http.StatusBadRequest)

		return
	}

	var transaction Transaction

	err = json.NewDecoder(r.Body).Decode(&transaction)
	if err != nil {
		slog.Error("error decoding create transaction request", "error", err)
		buildErrorResponse(w, err.Error(), http.StatusBadRequest)

		return
	}

	transaction.ID, err = uuid.NewV7()
	if err != nil {
		slog.Error("error creating transaction id", "error", err)
		buildErrorResponse(w, err.Error(), http.StatusInternalServerError)

		return
	}

	transaction.AccountID = accountID
	transaction.CreatedAt = time.Now()
	transaction.UpdatedAt = transaction.CreatedAt

	_, err = h.db.Exec(r.Context(), queryCreateTransaction,
		transaction.ID, accountID, transaction.CategoryID, transaction.PayeeID,
		transaction.Credit, transaction.Debit, transaction.Name, transaction.Notes,
		nil, transaction.CreatedAt, transaction.UpdatedAt)
	if err != nil {
		slog.Error("error creating transaction in database", "error", err)
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
	tID := r.PathValue("tId")

	accountID, err := uuid.Parse(id)
	if err != nil {
		slog.Error("error parsing account id", "error", err)
		buildErrorResponse(w, err.Error(), http.StatusBadRequest)

		return
	}

	transactionID, err := uuid.Parse(tID)
	if err != nil {
		slog.Error("error parsing transaction id", "error", err)
		buildErrorResponse(w, err.Error(), http.StatusBadRequest)

		return
	}

	var transaction Transaction

	err = json.NewDecoder(r.Body).Decode(&transaction)
	if err != nil {
		slog.Error("error decoding update transaction request", "error", err)
		buildErrorResponse(w, err.Error(), http.StatusBadRequest)

		return
	}

	transaction.UpdatedAt = time.Now()

	_, err = h.db.Exec(r.Context(), queryUpdateTransaction,
		accountID, transaction.CategoryID, transaction.PayeeID,
		transaction.Credit, transaction.Debit, transaction.Name, transaction.Notes,
		transaction.ClearedAt, transaction.UpdatedAt, transactionID)
	if err != nil {
		slog.Error("error updating transaction in database", "error", err)
		buildErrorResponse(w, err.Error(), http.StatusInternalServerError)

		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) DeleteTransaction(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	tID := r.PathValue("tId")

	accountID, err := uuid.Parse(id)
	if err != nil {
		slog.Error("error parsing account id", "error", err)
		buildErrorResponse(w, err.Error(), http.StatusBadRequest)

		return
	}

	transactionID, err := uuid.Parse(tID)
	if err != nil {
		slog.Error("error parsing transaction id", "error", err)
		buildErrorResponse(w, err.Error(), http.StatusBadRequest)

		return
	}

	_, err = h.db.Exec(r.Context(), queryDeleteTransaction, accountID, transactionID)
	if err != nil {
		slog.Error("error deleting transaction in database", "error", err)
		buildErrorResponse(w, err.Error(), http.StatusInternalServerError)

		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) GetTransactions(w http.ResponseWriter, r *http.Request) { //nolint: funlen
	id := r.PathValue("id")
	searchQuery := r.URL.Query().Get("q")

	page := 0
	limit := 50

	if p, err := strconv.Atoi(r.URL.Query().Get("page")); err == nil {
		page = p - 1
	}

	if l, err := strconv.Atoi(r.URL.Query().Get("limit")); err == nil {
		limit = l
	}

	accountID, err := uuid.Parse(id)
	if err != nil {
		slog.Error("error parsing account id", "error", err)
		buildErrorResponse(w, err.Error(), http.StatusBadRequest)

		return
	}

	rows, err := h.db.Query(r.Context(), queryGetTransactions, accountID, searchQuery, page, limit)
	if err != nil {
		slog.Error("error getting transactions from database", "error", err)
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
			slog.Error("error scanning transactions row from database", "error", err)
			buildErrorResponse(w, err.Error(), http.StatusInternalServerError)

			return
		}

		transactions = append(transactions, transaction)
	}

	if err := rows.Err(); err != nil {
		slog.Error("error reading transactions rows from database", "error", err)
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
	id := r.PathValue("id")

	accountID, err := uuid.Parse(id)
	if err != nil {
		slog.Error("error parsing account id", "error", err)
		buildErrorResponse(w, err.Error(), http.StatusBadRequest)

		return
	}

	var account Account

	err = h.db.QueryRow(r.Context(), queryGetAccount, accountID).Scan(&account.ID,
		&account.Name, &account.OffBudget, &account.Category, &account.Adapter, &account.CreatedAt, &account.UpdatedAt)
	if err != nil {
		slog.Error("error getting account from database", "error", err)
		buildErrorResponse(w, err.Error(), http.StatusInternalServerError)

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

	slog.Info("adapter", "config", h.adapters[account.Adapter+"-"+account.Category])
	adapterTransactions := adapters.GetTransactions(h.adapters[account.Adapter+"-"+account.Category], rows)
	importedTransactions := 0

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

		_, err = tx.Exec(r.Context(), "SAVEPOINT sp1")
		if err != nil {
			slog.Error("error storing savepoint", "error", err)
			buildErrorResponse(w, err.Error(), http.StatusInternalServerError)

			return
		}

		_, err = tx.Exec(r.Context(), queryCreateTransaction, transactionID, accountID, nil, nil,
			adapterTransaction.Credit, adapterTransaction.Debit, "imported transaction", adapterTransaction.Remarks,
			adapterTransaction.Date, transactionTime, transactionTime)

		if err != nil {
			slog.Error("error inserting transaction", "error", err, "adapterTransaction", adapterTransaction)

			_, rollbackErr := tx.Exec(r.Context(), "ROLLBACK TO SAVEPOINT sp1")
			if rollbackErr != nil {
				slog.Error("error rolling back savepoint", "error", rollbackErr)
			}

			if !strings.Contains(err.Error(), "duplicate key value") {
				buildErrorResponse(w, err.Error(), http.StatusInternalServerError)

				return
			}
		} else {
			importedTransactions++
		}
	}

	defer func(ctx context.Context) {
		if err != nil {
			rollBackErr := tx.Rollback(ctx)
			if rollBackErr != nil {
				slog.Error("error rolling back database txn", "error", rollBackErr)
			}
		}
	}(r.Context())

	if err := tx.Commit(r.Context()); err != nil {
		slog.Error("error committing database txn", "error", err)
		buildErrorResponse(w, err.Error(), http.StatusInternalServerError)

		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	err = json.NewEncoder(w).Encode(TransactionsResult{Total: len(adapterTransactions), Imported: importedTransactions})
	if err != nil {
		slog.Error("error encoding transactions result response", "error", err)
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
