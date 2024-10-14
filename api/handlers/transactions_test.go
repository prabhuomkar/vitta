package handlers

import (
	"errors"
	"net/http"
	"testing"

	uuid "github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/pashagolub/pgxmock/v4"
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
			"error due to auth", http.MethodPost, "/v1/accounts/" + testAccountID.String() + "/transactions", false, "invalid-body",
			nil,
			http.StatusUnauthorized, "Unauthorized",
		},
		{
			"error due to invalid account id", http.MethodPost, "/v1/accounts/invalid-account-id/transactions", true, "invalid-body",
			nil,
			http.StatusBadRequest, "invalid UUID",
		},
		{
			"error due to bad request", http.MethodPost, "/v1/accounts/" + testAccountID.String() + "/transactions", true, "invalid-body",
			nil,
			http.StatusBadRequest, "invalid character",
		},
		{
			"error inserting transaction to database", http.MethodPost, "/v1/accounts/" + testAccountID.String() + "/transactions", true,
			`{"name":"` + testTransactionName + `","accountId":"` + testAccountID.String() + `","payeeId":"` +
				testPayeeID.String() + `","categoryId":"` + testCategoryID.String() + `","credit":4.20}`,
			func(mock pgxmock.PgxPoolIface) {
				mock.ExpectExec("INSERT INTO transactions").WithArgs(pgxmock.AnyArg(),
					testAccountID, &testCategoryID, &testPayeeID, 4.20, 0.0, testTransactionName, "",
					pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg()).WillReturnError(pgx.ErrTxClosed)
			},
			http.StatusInternalServerError, "tx is closed",
		},
		{
			"success creating transaction", http.MethodPost, "/v1/accounts/" + testAccountID.String() + "/transactions", true,
			`{"name":"` + testTransactionName + `","accountId":"` + testAccountID.String() + `","payeeId":"` +
				testPayeeID.String() + `","categoryId":"` + testCategoryID.String() + `","credit":4.20}`,
			func(mock pgxmock.PgxPoolIface) {
				mock.ExpectExec("INSERT INTO transactions").WithArgs(pgxmock.AnyArg(),
					testAccountID, &testCategoryID, &testPayeeID, 4.20, 0.0, testTransactionName, "",
					pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg()).WillReturnResult(pgxmock.NewResult("INSERT", 1))
			},
			http.StatusCreated, testTransactionName,
		},
	}
	for _, tc := range tests {
		tc.Run(t)
	}
}

func TestUpdateTransaction(t *testing.T) {
	tests := []testCase{
		{
			"error due to auth", http.MethodPatch, "/v1/accounts/" + testAccountID.String() + "/transactions/invalid-uuid", false, "invalid-body",
			nil,
			http.StatusUnauthorized, "Unauthorized",
		},
		{
			"error due to invalid account id", http.MethodPatch, "/v1/accounts/invalid-account-id/transactions/invalid-uuid", true, "invalid-body",
			nil,
			http.StatusBadRequest, "invalid UUID",
		},
		{
			"error due to bad transaction id", http.MethodPatch, "/v1/accounts/" + testAccountID.String() + "/transactions/invalid-uuid", true, "invalid-body",
			nil,
			http.StatusBadRequest, "invalid UUID",
		},
		{
			"error due to bad request", http.MethodPatch, "/v1/accounts/" + testAccountID.String() + "/transactions/" + testTransactionID.String(), true, "invalid-body",
			nil,
			http.StatusBadRequest, "invalid character",
		},
		{
			"error updating transaction in database", http.MethodPatch, "/v1/accounts/" + testAccountID.String() + "/transactions/" + testTransactionID.String(), true,
			`{"accountId":"` + testAccountID.String() + `","categoryId":"` + testCategoryID.String() +
				`","payeeId":"` + testPayeeID.String() + `","name":"Some name","notes":"Some notes",` +
				`"credit":4.20,"debit":4.20}`,
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
			`{"accountId":"` + testAccountID.String() + `","categoryId":"` + testCategoryID.String() +
				`","payeeId":"` + testPayeeID.String() + `","name":"Some name","notes":"Some notes",` +
				`"credit":4.20,"debit":4.20}`,
			func(mock pgxmock.PgxPoolIface) {
				mock.ExpectExec("UPDATE transactions").WithArgs(
					testAccountID, &testCategoryID, &testPayeeID, 4.20, 4.20, "Some name",
					"Some notes", pgxmock.AnyArg(), pgxmock.AnyArg(), testTransactionID,
				).WillReturnResult(pgxmock.NewResult("UPDATE", 1))
			},
			http.StatusNoContent, "",
		},
	}
	for _, tc := range tests {
		tc.Run(t)
	}
}

func TestDeleteTransaction(t *testing.T) {
	tests := []testCase{
		{
			"error due to auth", http.MethodDelete, "/v1/accounts/" + testTransactionID.String() + "/transactions/invalid-uuid", false, "",
			nil,
			http.StatusUnauthorized, "Unauthorized",
		},
		{
			"error due to invalid account id", http.MethodDelete, "/v1/accounts/invalid-account-id/transactions/invalid-uuid", true, "",
			nil,
			http.StatusBadRequest, "invalid UUID",
		},
		{
			"error due to bad transaction id", http.MethodDelete, "/v1/accounts/" + testAccountID.String() + "/transactions/invalid-uuid", true, "",
			nil,
			http.StatusBadRequest, "invalid UUID",
		},
		{
			"error deleting transaction in database", http.MethodDelete, "/v1/accounts/" + testAccountID.String() + "/transactions/" + testTransactionID.String(), true,
			``,
			func(mock pgxmock.PgxPoolIface) {
				mock.ExpectExec("DELETE FROM transactions").WithArgs(testAccountID, testTransactionID).WillReturnError(pgx.ErrTxClosed)
			},
			http.StatusInternalServerError, "tx is closed",
		},
		{
			"success deleting transaction", http.MethodDelete, "/v1/accounts/" + testAccountID.String() + "/transactions/" + testTransactionID.String(), true,
			``,
			func(mock pgxmock.PgxPoolIface) {
				mock.ExpectExec("DELETE FROM transactions").WithArgs(testAccountID, testTransactionID).WillReturnResult(pgxmock.NewResult("DELETE", 1))
			},
			http.StatusNoContent, "",
		},
	}
	for _, tc := range tests {
		tc.Run(t)
	}
}

func TestGetTransactions(t *testing.T) {
	tests := []testCase{
		{
			"error due to auth", http.MethodGet, "/v1/accounts/" + testAccountID.String() + "/transactions", false, "",
			nil,
			http.StatusUnauthorized, "Unauthorized",
		},
		{
			"error getting transactions from db", http.MethodGet, "/v1/accounts/" + testAccountID.String() + "/transactions", true, "",
			func(mock pgxmock.PgxPoolIface) {
				mock.ExpectQuery("SELECT *").WithArgs(testAccountID).WillReturnError(pgx.ErrNoRows)
			},
			http.StatusInternalServerError, "no rows",
		},
		{
			"error scanning transactions rows from db", http.MethodGet, "/v1/accounts/" + testAccountID.String() + "/transactions", true, "",
			func(mock pgxmock.PgxPoolIface) {
				mock.ExpectQuery("SELECT *").WithArgs(testAccountID).WillReturnRows(pgxmock.NewRows(transactionRowCols).AddRow("invalid", "invalid", "invalid",
					"invalid", "invalid", "invalid", "invalid", "invalid", "invalid", "invalid", "invalid", "invalid", "invalid"))
			},
			http.StatusInternalServerError, "Scanning value error",
		},
		{
			"error reading transactions rows from db", http.MethodGet, "/v1/accounts/" + testAccountID.String() + "/transactions", true, "",
			func(mock pgxmock.PgxPoolIface) {
				mock.ExpectQuery("SELECT *").WithArgs(testAccountID).WillReturnRows(pgxmock.NewRows(transactionRowCols).RowError(0, errors.New("some error in db")))
			},
			http.StatusInternalServerError, "some error in db",
		},
		{
			"success", http.MethodGet, "/v1/accounts/" + testAccountID.String() + "/transactions", true, "",
			func(mock pgxmock.PgxPoolIface) {
				mock.ExpectQuery("SELECT *").WithArgs(testAccountID).WillReturnRows(pgxmock.NewRows(transactionRowCols).AddRow(testTransactionID,
					testAccountID, &testCategoryID, &testPayeeID, "Some transaction", 4.20, 4.20, "Some notes",
					&testAccountTime, testAccountTime, testAccountTime, &testCategoryName, &testPayeeName))
			},
			http.StatusOK, testAccountID.String(),
		},
	}
	for _, tc := range tests {
		tc.Run(t)
	}
}
