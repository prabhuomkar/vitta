package handlers

import (
	"net/http"
	"strconv"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/pashagolub/pgxmock/v4"
)

var (
	testAccountName = "ICICI Bank"
	testOffBudget   = true
	testCategory    = "CASH/CHECK IN"
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
			"error due to bad request", http.MethodPatch, "/v1/accounts/invalid-uuid", true, "invalid-body",
			nil,
			http.StatusBadRequest, "invalid UUID",
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
		// TODO(omkar): Add more unit test cases
	}
	for _, tc := range tests {
		tc.Run(t)
	}
}
