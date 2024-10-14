package handlers

import (
	"errors"
	"net/http"
	"strconv"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/pashagolub/pgxmock/v4"
)

var (
	testAccountName = "ICICI Bank"
	testOffBudget   = true
	testCategory    = "CASH/CHECK IN"
	testAccountTime = time.Now()
	accountRowCols  = []string{"id", "name", "off_budget", "category", "created_at", "updated_at"}
)

func TestCreateAccount(t *testing.T) {
	tests := []testCase{
		{
			"error due to auth", http.MethodPost, "/v1/accounts", false, "",
			nil,
			http.StatusUnauthorized, "Unauthorized",
		},
		{
			"error due to bad request", http.MethodPost, "/v1/accounts", true, "invalid-body",
			nil,
			http.StatusBadRequest, "invalid character",
		},
		{
			"error inserting account to database", http.MethodPost, "/v1/accounts", true,
			`{"name":"` + testAccountName + `","offBudget":` + strconv.FormatBool(testOffBudget) +
				`,"category":"` + testCategory + `"}`,
			func(mock pgxmock.PgxPoolIface) {
				mock.ExpectExec("INSERT INTO accounts").WithArgs(pgxmock.AnyArg(),
					testAccountName, &testOffBudget, testCategory, pgxmock.AnyArg(),
					pgxmock.AnyArg()).WillReturnError(pgx.ErrTxClosed)
			},
			http.StatusInternalServerError, "tx is closed",
		},
		{
			"success creating account", http.MethodPost, "/v1/accounts", true,
			`{"name":"` + testAccountName + `","offBudget":` + strconv.FormatBool(testOffBudget) +
				`,"category":"` + testCategory + `"}`,
			func(mock pgxmock.PgxPoolIface) {
				mock.ExpectExec("INSERT INTO accounts").WithArgs(pgxmock.AnyArg(),
					testAccountName, &testOffBudget, testCategory, pgxmock.AnyArg(),
					pgxmock.AnyArg()).WillReturnResult(pgxmock.NewResult("INSERT", 1))
			},
			http.StatusCreated, testAccountName,
		},
	}
	for _, tc := range tests {
		tc.Run(t)
	}
}

func TestUpdateAccount(t *testing.T) {
	tests := []testCase{
		{
			"error due to auth", http.MethodPatch, "/v1/accounts/invalid-uuid", false, "invalid-body",
			nil,
			http.StatusUnauthorized, "Unauthorized",
		},
		{
			"error due to bad account id", http.MethodPatch, "/v1/accounts/invalid-uuid", true, "invalid-body",
			nil,
			http.StatusBadRequest, "invalid UUID",
		},
		{
			"error due to bad request", http.MethodPatch, "/v1/accounts/" + testAccountID.String(), true, "invalid-body",
			nil,
			http.StatusBadRequest, "invalid character",
		},
		{
			"error updating account in database", http.MethodPatch, "/v1/accounts/" + testAccountID.String(), true,
			`{"name":"` + testAccountName + `","offBudget":` + strconv.FormatBool(testOffBudget) +
				`,"category":"` + testCategory + `"}`,
			func(mock pgxmock.PgxPoolIface) {
				mock.ExpectExec("UPDATE accounts").WithArgs(
					testAccountName, &testOffBudget, testCategory, pgxmock.AnyArg(), testAccountID,
				).WillReturnError(pgx.ErrTxClosed)
			},
			http.StatusInternalServerError, "tx is closed",
		},
		{
			"success updating account", http.MethodPatch, "/v1/accounts/" + testAccountID.String(), true,
			`{"name":"` + testAccountName + `","offBudget":` + strconv.FormatBool(testOffBudget) +
				`,"category":"` + testCategory + `"}`,
			func(mock pgxmock.PgxPoolIface) {
				mock.ExpectExec("UPDATE accounts").WithArgs(
					testAccountName, &testOffBudget, testCategory, pgxmock.AnyArg(), testAccountID,
				).WillReturnResult(pgxmock.NewResult("UPDATE", 1))
			},
			http.StatusNoContent, "",
		},
	}
	for _, tc := range tests {
		tc.Run(t)
	}
}

func TestDeleteAccount(t *testing.T) {
	tests := []testCase{
		{
			"error due to auth", http.MethodDelete, "/v1/accounts/invalid-uuid", false, "invalid-body",
			nil,
			http.StatusUnauthorized, "Unauthorized",
		},
		{
			"error due to bad request", http.MethodDelete, "/v1/accounts/invalid-uuid", true, "invalid-body",
			nil,
			http.StatusBadRequest, "invalid UUID",
		},
		{
			"error deleting account in database", http.MethodDelete, "/v1/accounts/" + testAccountID.String(), true,
			``,
			func(mock pgxmock.PgxPoolIface) {
				mock.ExpectExec("DELETE FROM accounts").WithArgs(testAccountID).WillReturnError(pgx.ErrTxClosed)
			},
			http.StatusInternalServerError, "tx is closed",
		},
		{
			"success deleting account", http.MethodDelete, "/v1/accounts/" + testAccountID.String(), true,
			``,
			func(mock pgxmock.PgxPoolIface) {
				mock.ExpectExec("DELETE FROM accounts").WithArgs(testAccountID).WillReturnResult(pgxmock.NewResult("DELETE", 1))
			},
			http.StatusNoContent, "",
		},
	}
	for _, tc := range tests {
		tc.Run(t)
	}
}

func TestGetAccounts(t *testing.T) {
	tests := []testCase{
		{
			"error due to auth", http.MethodGet, "/v1/accounts", false, "",
			nil,
			http.StatusUnauthorized, "Unauthorized",
		},
		{
			"error getting accounts from db", http.MethodGet, "/v1/accounts", true, "",
			func(mock pgxmock.PgxPoolIface) {
				mock.ExpectQuery("SELECT *").WillReturnError(pgx.ErrNoRows)
			},
			http.StatusInternalServerError, "no rows",
		},
		{
			"error scanning accounts rows from db", http.MethodGet, "/v1/accounts", true, "",
			func(mock pgxmock.PgxPoolIface) {
				mock.ExpectQuery("SELECT *").WillReturnRows(pgxmock.NewRows(accountRowCols).AddRow("invalid", "ok", "false", "category", "bad-time", "bad-time"))
			},
			http.StatusInternalServerError, "Scanning value error",
		},
		{
			"error reading accounts rows from db", http.MethodGet, "/v1/accounts", true, "",
			func(mock pgxmock.PgxPoolIface) {
				mock.ExpectQuery("SELECT *").WillReturnRows(pgxmock.NewRows(accountRowCols).RowError(0, errors.New("some error in db")))
			},
			http.StatusInternalServerError, "some error in db",
		},
		{
			"success", http.MethodGet, "/v1/accounts", true, "",
			func(mock pgxmock.PgxPoolIface) {
				mock.ExpectQuery("SELECT *").WillReturnRows(pgxmock.NewRows(accountRowCols).AddRow(testAccountID.String(), testAccountName, &testOffBudget, testCategory,
					testAccountTime, testAccountTime))
			},
			http.StatusOK, testAccountID.String(),
		},
	}
	for _, tc := range tests {
		tc.Run(t)
	}
}
