package handlers

import (
	"bytes"
	"encoding/csv"
	"errors"
	"io"
	"mime/multipart"
	"net/http"
	"strings"
	"testing"

	uuid "github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/pashagolub/pgxmock/v4"
	"github.com/stretchr/testify/require"
)

var (
	testTransactionName = "House Party Food"
	testCategoryID      = uuid.MustParse("01927f3e-11b6-7a79-bbc7-affae59272ae")
	testPayeeID         = uuid.MustParse("01927f3e-5609-703b-b067-f9b9dd9d8ee2")
	testAccountID       = uuid.MustParse("01927f3e-6ecf-7091-987f-8aa23adcda09")
	testTransactionID   = uuid.MustParse("01927f3e-6ecf-7091-987f-8aa23addda09")
	transactionRowCols  = []string{"id", "account_id", "category_id", "payee_id", "name", "credit",
		"debit", "notes", "cleared_at", "created_at", "updated_at", "category_name", "payee_name"}
)

func TestCreateTransaction(t *testing.T) {
	tests := []testCase{
		{
			"error due to auth", http.MethodPost, "/v1/accounts/" + testAccountID.String() + "/transactions", false, strings.NewReader("invalid-body"),
			nil, nil,
			http.StatusUnauthorized, "Unauthorized",
		},
		{
			"error due to invalid account id", http.MethodPost, "/v1/accounts/invalid-account-id/transactions", true, strings.NewReader("invalid-body"),
			nil, nil,
			http.StatusBadRequest, "invalid UUID",
		},
		{
			"error due to bad request", http.MethodPost, "/v1/accounts/" + testAccountID.String() + "/transactions", true, strings.NewReader("invalid-body"),
			nil, nil,
			http.StatusBadRequest, "invalid character",
		},
		{
			"error inserting transaction to database", http.MethodPost, "/v1/accounts/" + testAccountID.String() + "/transactions", true,
			strings.NewReader(`{"name":"` + testTransactionName + `","accountId":"` + testAccountID.String() + `","payeeId":"` +
				testPayeeID.String() + `","categoryId":"` + testCategoryID.String() + `","credit":4.20}`),
			nil,
			func(mock pgxmock.PgxPoolIface) {
				mock.ExpectExec("INSERT INTO transactions").WithArgs(pgxmock.AnyArg(),
					testAccountID, &testCategoryID, &testPayeeID, 4.20, 0.0, testTransactionName, "",
					pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg()).WillReturnError(pgx.ErrTxClosed)
			},
			http.StatusInternalServerError, "tx is closed",
		},
		{
			"success creating transaction", http.MethodPost, "/v1/accounts/" + testAccountID.String() + "/transactions", true,
			strings.NewReader(`{"name":"` + testTransactionName + `","accountId":"` + testAccountID.String() + `","payeeId":"` +
				testPayeeID.String() + `","categoryId":"` + testCategoryID.String() + `","credit":4.20}`),
			nil,
			func(mock pgxmock.PgxPoolIface) {
				mock.ExpectExec("INSERT INTO transactions").WithArgs(pgxmock.AnyArg(),
					testAccountID, &testCategoryID, &testPayeeID, 4.20, 0.0, testTransactionName, "",
					pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg()).WillReturnResult(pgxmock.NewResult("INSERT", 1))
			},
			http.StatusCreated, testTransactionName,
		},
	}
	executeTests(t, tests)
}

func TestUpdateTransaction(t *testing.T) {
	tests := []testCase{
		{
			"error due to auth", http.MethodPatch, "/v1/accounts/" + testAccountID.String() + "/transactions/invalid-uuid", false, strings.NewReader("invalid-body"),
			nil, nil,
			http.StatusUnauthorized, "Unauthorized",
		},
		{
			"error due to invalid account id", http.MethodPatch, "/v1/accounts/invalid-account-id/transactions/invalid-uuid", true, strings.NewReader("invalid-body"),
			nil, nil,
			http.StatusBadRequest, "invalid UUID",
		},
		{
			"error due to bad transaction id", http.MethodPatch, "/v1/accounts/" + testAccountID.String() + "/transactions/invalid-uuid", true, strings.NewReader("invalid-body"),
			nil, nil,
			http.StatusBadRequest, "invalid UUID",
		},
		{
			"error due to bad request", http.MethodPatch, "/v1/accounts/" + testAccountID.String() + "/transactions/" + testTransactionID.String(), true, strings.NewReader("invalid-body"),
			nil, nil,
			http.StatusBadRequest, "invalid character",
		},
		{
			"error updating transaction in database", http.MethodPatch, "/v1/accounts/" + testAccountID.String() + "/transactions/" + testTransactionID.String(), true,
			strings.NewReader(`{"accountId":"` + testAccountID.String() + `","categoryId":"` + testCategoryID.String() +
				`","payeeId":"` + testPayeeID.String() + `","name":"Some name","notes":"Some notes",` +
				`"credit":4.20,"debit":4.20}`),
			nil,
			func(mock pgxmock.PgxPoolIface) {
				mock.ExpectExec("UPDATE transactions").WithArgs(
					testAccountID, &testCategoryID, &testPayeeID, 4.20, 4.20, "Some name",
					"Some notes", pgxmock.AnyArg(), pgxmock.AnyArg(), testTransactionID,
				).WillReturnError(pgx.ErrTxClosed)
			},
			http.StatusInternalServerError, "tx is closed",
		},
		{
			"success updating transaction", http.MethodPatch, "/v1/accounts/" + testAccountID.String() + "/transactions/" + testTransactionID.String(), true,
			strings.NewReader(`{"accountId":"` + testAccountID.String() + `","categoryId":"` + testCategoryID.String() +
				`","payeeId":"` + testPayeeID.String() + `","name":"Some name","notes":"Some notes",` +
				`"credit":4.20,"debit":4.20}`),
			nil,
			func(mock pgxmock.PgxPoolIface) {
				mock.ExpectExec("UPDATE transactions").WithArgs(
					testAccountID, &testCategoryID, &testPayeeID, 4.20, 4.20, "Some name",
					"Some notes", pgxmock.AnyArg(), pgxmock.AnyArg(), testTransactionID,
				).WillReturnResult(pgxmock.NewResult("UPDATE", 1))
			},
			http.StatusNoContent, "",
		},
	}
	executeTests(t, tests)
}

func TestDeleteTransaction(t *testing.T) {
	tests := []testCase{
		{
			"error due to auth", http.MethodDelete, "/v1/accounts/" + testTransactionID.String() + "/transactions/invalid-uuid", false, nil,
			nil, nil,
			http.StatusUnauthorized, "Unauthorized",
		},
		{
			"error due to invalid account id", http.MethodDelete, "/v1/accounts/invalid-account-id/transactions/invalid-uuid", true, nil,
			nil, nil,
			http.StatusBadRequest, "invalid UUID",
		},
		{
			"error due to bad transaction id", http.MethodDelete, "/v1/accounts/" + testAccountID.String() + "/transactions/invalid-uuid", true, nil,
			nil, nil,
			http.StatusBadRequest, "invalid UUID",
		},
		{
			"error deleting transaction in database", http.MethodDelete, "/v1/accounts/" + testAccountID.String() + "/transactions/" + testTransactionID.String(), true,
			nil, nil,
			func(mock pgxmock.PgxPoolIface) {
				mock.ExpectExec("DELETE FROM transactions").WithArgs(testAccountID, testTransactionID).WillReturnError(pgx.ErrTxClosed)
			},
			http.StatusInternalServerError, "tx is closed",
		},
		{
			"success deleting transaction", http.MethodDelete, "/v1/accounts/" + testAccountID.String() + "/transactions/" + testTransactionID.String(), true,
			nil, nil,
			func(mock pgxmock.PgxPoolIface) {
				mock.ExpectExec("DELETE FROM transactions").WithArgs(testAccountID, testTransactionID).WillReturnResult(pgxmock.NewResult("DELETE", 1))
			},
			http.StatusNoContent, "",
		},
	}
	executeTests(t, tests)
}

func TestGetTransactions(t *testing.T) {
	tests := []testCase{
		{
			"error due to auth", http.MethodGet, "/v1/accounts/" + testAccountID.String() + "/transactions", false, nil,
			nil, nil,
			http.StatusUnauthorized, "Unauthorized",
		},
		{
			"error getting transactions from db", http.MethodGet, "/v1/accounts/" + testAccountID.String() + "/transactions", true, nil,
			nil,
			func(mock pgxmock.PgxPoolIface) {
				mock.ExpectQuery("SELECT *").WithArgs(testAccountID).WillReturnError(pgx.ErrNoRows)
			},
			http.StatusInternalServerError, "no rows",
		},
		{
			"error scanning transactions rows from db", http.MethodGet, "/v1/accounts/" + testAccountID.String() + "/transactions", true, nil,
			nil,
			func(mock pgxmock.PgxPoolIface) {
				mock.ExpectQuery("SELECT *").WithArgs(testAccountID).WillReturnRows(pgxmock.NewRows(transactionRowCols).AddRow("invalid", "invalid", "invalid",
					"invalid", "invalid", "invalid", "invalid", "invalid", "invalid", "invalid", "invalid", "invalid", "invalid"))
			},
			http.StatusInternalServerError, "Scanning value error",
		},
		{
			"error reading transactions rows from db", http.MethodGet, "/v1/accounts/" + testAccountID.String() + "/transactions", true, nil,
			nil,
			func(mock pgxmock.PgxPoolIface) {
				mock.ExpectQuery("SELECT *").WithArgs(testAccountID).WillReturnRows(pgxmock.NewRows(transactionRowCols).RowError(0, errors.New("some error in db")))
			},
			http.StatusInternalServerError, "some error in db",
		},
		{
			"success", http.MethodGet, "/v1/accounts/" + testAccountID.String() + "/transactions", true, nil,
			nil,
			func(mock pgxmock.PgxPoolIface) {
				mock.ExpectQuery("SELECT *").WithArgs(testAccountID).WillReturnRows(pgxmock.NewRows(transactionRowCols).AddRow(testTransactionID,
					testAccountID, &testCategoryID, &testPayeeID, "Some transaction", 4.20, 4.20, "Some notes",
					&testAccountTime, testAccountTime, testAccountTime, &testCategoryName, &testPayeeName))
			},
			http.StatusOK, testAccountID.String(),
		},
	}
	executeTests(t, tests)
}

func TestImportTransactions(t *testing.T) {
	sampleBytes1, ctype1 := getMockCSV(t, true)
	sampleBytes2, ctype2 := getMockCSV(t, false)
	sampleBytes3, ctype3 := getMockCSV(t, false)
	tests := []testCase{
		{
			"error due to auth", http.MethodPut, "/v1/accounts/" + testAccountID.String() + "/transactions", false, nil,
			nil, nil,
			http.StatusUnauthorized, "Unauthorized",
		},
		{
			"error getting account from db", http.MethodPut, "/v1/accounts/" + testAccountID.String() + "/transactions", true, nil,
			nil,
			func(mock pgxmock.PgxPoolIface) {
				mock.ExpectQuery("SELECT *").WithArgs(testAccountID).WillReturnError(pgx.ErrNoRows)
			},
			http.StatusInternalServerError, "no rows",
		},
		{
			"error parsing multi part form", http.MethodPut, "/v1/accounts/" + testAccountID.String() + "/transactions", true, nil,
			nil,
			func(mock pgxmock.PgxPoolIface) {
				mock.ExpectQuery("SELECT *").WithArgs(testAccountID).WillReturnRows(pgxmock.NewRows(accountRowCols).AddRow(testAccountID.String(), testAccountName, &testOffBudget, testCategory,
					testAdapter, testAccountTime, testAccountTime))
			},
			http.StatusBadRequest, "request Content-Type isn't multipart/form-data",
		},
		{
			"error getting data rows", http.MethodPut, "/v1/accounts/" + testAccountID.String() + "/transactions", true,
			sampleBytes1, ctype1,
			func(mock pgxmock.PgxPoolIface) {
				mock.ExpectQuery("SELECT *").WithArgs(testAccountID).WillReturnRows(pgxmock.NewRows(accountRowCols).AddRow(testAccountID.String(), testAccountName, &testOffBudget, testCategory,
					testAdapter, testAccountTime, testAccountTime))
			},
			http.StatusInternalServerError, "error opening file",
		},
		{
			"error inserting transaction to database", http.MethodPut, "/v1/accounts/" + testAccountID.String() + "/transactions", true,
			sampleBytes2, ctype2,
			func(mock pgxmock.PgxPoolIface) {
				mock.ExpectQuery("SELECT *").WithArgs(testAccountID).WillReturnRows(pgxmock.NewRows(accountRowCols).AddRow(testAccountID.String(), testAccountName, &testOffBudget, testCategory,
					testAdapter, testAccountTime, testAccountTime))
				mock.ExpectBeginTx(pgx.TxOptions{})
				mock.ExpectExec("SAVEPOINT sp1").WillReturnResult(pgxmock.NewResult("SAVEPOINT", 1))
				mock.ExpectExec("INSERT INTO transactions").WithArgs(pgxmock.AnyArg(),
					testAccountID, nil, nil, 0.0, 4.20, "imported transaction", "John Doe",
					pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg()).WillReturnError(pgx.ErrTxClosed)
				mock.ExpectExec("ROLLBACK TO SAVEPOINT sp1").WillReturnResult(pgxmock.NewResult("SAVEPOINT", 1))
			},
			http.StatusInternalServerError, "tx is closed",
		},
		{
			"success", http.MethodPut, "/v1/accounts/" + testAccountID.String() + "/transactions", true,
			sampleBytes3, ctype3,
			func(mock pgxmock.PgxPoolIface) {
				mock.ExpectQuery("SELECT *").WithArgs(testAccountID).WillReturnRows(pgxmock.NewRows(accountRowCols).AddRow(testAccountID.String(), testAccountName, &testOffBudget, testCategory,
					testAdapter, testAccountTime, testAccountTime))
				mock.ExpectBeginTx(pgx.TxOptions{})
				mock.ExpectExec("SAVEPOINT sp1").WillReturnResult(pgxmock.NewResult("SAVEPOINT", 1))
				mock.ExpectExec("INSERT INTO transactions").WithArgs(pgxmock.AnyArg(),
					testAccountID, nil, nil, 0.0, 4.20, "imported transaction", "John Doe",
					pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg()).WillReturnResult(pgxmock.NewResult("INSERT", 1))
				mock.ExpectCommit()
			},
			http.StatusOK, `"total":1,"imported":1`,
		},
	}
	executeTests(t, tests)
}

func getMockCSV(t *testing.T, fail bool) (io.Reader, http.Header) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	ext := ".csv"
	rows := [][]string{
		{"Transaction Date", "Details", "Amount (INR)"},
		{"18/10/2024", "John Doe", "4.20 Dr."},
	}

	if fail {
		ext = ".pdf"
		rows = [][]string{
			{"Transaction Date", "Details", "Amount (INR)"},
			{"invalid-data", "no"},
		}
	}

	part, err := writer.CreateFormFile("file", "SomeFile"+ext)
	require.NoError(t, err)

	csvWriter := csv.NewWriter(part)

	err = csvWriter.WriteAll(rows)
	require.NoError(t, err)
	csvWriter.Flush()

	err = writer.Close()
	require.NoError(t, err)

	return body, http.Header{"Content-Type": []string{writer.FormDataContentType()}}
}
