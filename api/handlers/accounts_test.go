package handlers

import (
	"errors"
	"net/http"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/pashagolub/pgxmock/v4"
)

var (
	testAccountName = "ICICI Bank"
	testOffBudget   = true
	testCategory    = "CC"
	testAdapter     = "icici"
	testAccountTime = time.Now()
	accountRowCols  = []string{"id", "name", "off_budget", "category", "adapter", "created_at", "updated_at"}
	accountsRowCols = []string{"id", "name", "off_budget", "category", "adapter", "created_at", "updated_at", "balance"}
)

func TestCreateAccount(t *testing.T) {
	tests := []testCase{
		{
			"error due to auth", http.MethodPost, "/v1/accounts", false, nil,
			nil, nil,
			http.StatusUnauthorized, "Unauthorized",
		},
		{
			"error due to bad request", http.MethodPost, "/v1/accounts", true, strings.NewReader("invalid-body"),
			nil, nil,
			http.StatusBadRequest, "invalid character",
		},
		{
			"error inserting account to database", http.MethodPost, "/v1/accounts", true,
			strings.NewReader(`{"name":"` + testAccountName + `","offBudget":` + strconv.FormatBool(testOffBudget) +
				`,"category":"` + testCategory + `","adapter":"` + testAdapter + `"}`),
			nil,
			func(mock pgxmock.PgxPoolIface) {
				mock.ExpectExec("INSERT INTO accounts").WithArgs(pgxmock.AnyArg(),
					testAccountName, &testOffBudget, testCategory, testAdapter, pgxmock.AnyArg(),
					pgxmock.AnyArg()).WillReturnError(pgx.ErrTxClosed)
			},
			http.StatusInternalServerError, "tx is closed",
		},
		{
			"success creating account", http.MethodPost, "/v1/accounts", true,
			strings.NewReader(`{"name":"` + testAccountName + `","offBudget":` + strconv.FormatBool(testOffBudget) +
				`,"category":"` + testCategory + `","adapter":"` + testAdapter + `"}`),
			nil,
			func(mock pgxmock.PgxPoolIface) {
				mock.ExpectExec("INSERT INTO accounts").WithArgs(pgxmock.AnyArg(),
					testAccountName, &testOffBudget, testCategory, testAdapter, pgxmock.AnyArg(),
					pgxmock.AnyArg()).WillReturnResult(pgxmock.NewResult("INSERT", 1))
			},
			http.StatusCreated, testAccountName,
		},
	}
	executeTests(t, tests)
}

func TestUpdateAccount(t *testing.T) {
	tests := []testCase{
		{
			"error due to auth", http.MethodPatch, "/v1/accounts/invalid-uuid", false, strings.NewReader("invalid-body"),
			nil, nil,
			http.StatusUnauthorized, "Unauthorized",
		},
		{
			"error due to bad account id", http.MethodPatch, "/v1/accounts/invalid-uuid", true, strings.NewReader("invalid-body"),
			nil, nil,
			http.StatusBadRequest, "invalid UUID",
		},
		{
			"error due to bad request", http.MethodPatch, "/v1/accounts/" + testAccountID.String(), true, strings.NewReader("invalid-body"),
			nil, nil,
			http.StatusBadRequest, "invalid character",
		},
		{
			"error updating account in database", http.MethodPatch, "/v1/accounts/" + testAccountID.String(), true,
			strings.NewReader(`{"name":"` + testAccountName + `","offBudget":` + strconv.FormatBool(testOffBudget) +
				`,"category":"` + testCategory + `","adapter":"` + testAdapter + `"}`),
			nil,
			func(mock pgxmock.PgxPoolIface) {
				mock.ExpectExec("UPDATE accounts").WithArgs(
					testAccountName, &testOffBudget, testCategory, testAdapter, pgxmock.AnyArg(), testAccountID,
				).WillReturnError(pgx.ErrTxClosed)
			},
			http.StatusInternalServerError, "tx is closed",
		},
		{
			"success updating account", http.MethodPatch, "/v1/accounts/" + testAccountID.String(), true,
			strings.NewReader(`{"name":"` + testAccountName + `","offBudget":` + strconv.FormatBool(testOffBudget) +
				`,"category":"` + testCategory + `","adapter":"` + testAdapter + `"}`),
			nil,
			func(mock pgxmock.PgxPoolIface) {
				mock.ExpectExec("UPDATE accounts").WithArgs(
					testAccountName, &testOffBudget, testCategory, testAdapter, pgxmock.AnyArg(), testAccountID,
				).WillReturnResult(pgxmock.NewResult("UPDATE", 1))
			},
			http.StatusNoContent, "",
		},
	}
	executeTests(t, tests)
}

func TestDeleteAccount(t *testing.T) {
	tests := []testCase{
		{
			"error due to auth", http.MethodDelete, "/v1/accounts/invalid-uuid", false, nil,
			nil, nil,
			http.StatusUnauthorized, "Unauthorized",
		},
		{
			"error due to bad account id", http.MethodDelete, "/v1/accounts/invalid-uuid", true, nil,
			nil, nil,
			http.StatusBadRequest, "invalid UUID",
		},
		{
			"error deleting account in database", http.MethodDelete, "/v1/accounts/" + testAccountID.String(), true,
			nil,
			nil,
			func(mock pgxmock.PgxPoolIface) {
				mock.ExpectExec("DELETE FROM accounts").WithArgs(testAccountID).WillReturnError(pgx.ErrTxClosed)
			},
			http.StatusInternalServerError, "tx is closed",
		},
		{
			"success deleting account", http.MethodDelete, "/v1/accounts/" + testAccountID.String(), true,
			nil,
			nil,
			func(mock pgxmock.PgxPoolIface) {
				mock.ExpectExec("DELETE FROM accounts").WithArgs(testAccountID).WillReturnResult(pgxmock.NewResult("DELETE", 1))
			},
			http.StatusNoContent, "",
		},
	}
	executeTests(t, tests)
}

func TestGetAccounts(t *testing.T) {
	tests := []testCase{
		{
			"error due to auth", http.MethodGet, "/v1/accounts", false, nil,
			nil, nil,
			http.StatusUnauthorized, "Unauthorized",
		},
		{
			"error getting accounts from db", http.MethodGet, "/v1/accounts", true, nil,
			nil,
			func(mock pgxmock.PgxPoolIface) {
				mock.ExpectQuery("SELECT *").WillReturnError(pgx.ErrNoRows)
			},
			http.StatusInternalServerError, "no rows",
		},
		{
			"error scanning accounts rows from db", http.MethodGet, "/v1/accounts", true, nil,
			nil,
			func(mock pgxmock.PgxPoolIface) {
				mock.ExpectQuery("SELECT *").WillReturnRows(pgxmock.NewRows(accountsRowCols).AddRow("invalid", "ok", "false", "category", "adapter", "bad-time", "bad-time", "bad-value"))
			},
			http.StatusInternalServerError, "Scanning value error",
		},
		{
			"error reading accounts rows from db", http.MethodGet, "/v1/accounts", true, nil,
			nil,
			func(mock pgxmock.PgxPoolIface) {
				mock.ExpectQuery("SELECT *").WillReturnRows(pgxmock.NewRows(accountsRowCols).RowError(0, errors.New("some error in db")))
			},
			http.StatusInternalServerError, "some error in db",
		},
		{
			"success", http.MethodGet, "/v1/accounts", true, nil,
			nil,
			func(mock pgxmock.PgxPoolIface) {
				mock.ExpectQuery("SELECT *").WillReturnRows(pgxmock.NewRows(accountsRowCols).AddRow(testAccountID.String(), testAccountName, &testOffBudget, testCategory,
					testAdapter, testAccountTime, testAccountTime, testAmount))
			},
			http.StatusOK, testAccountID.String(),
		},
	}
	executeTests(t, tests)
}

func TestGetAccount(t *testing.T) {
	tests := []testCase{
		{
			"error due to auth", http.MethodGet, "/v1/accounts/invalid-uuid", false, nil,
			nil, nil,
			http.StatusUnauthorized, "Unauthorized",
		},
		{
			"error due to bad account id", http.MethodGet, "/v1/accounts/invalid-uuid", true, nil,
			nil, nil,
			http.StatusBadRequest, "invalid UUID",
		},
		{
			"error getting account in database", http.MethodGet, "/v1/accounts/" + testAccountID.String(), true,
			nil,
			nil,
			func(mock pgxmock.PgxPoolIface) {
				mock.ExpectExec("SELECT *").WithArgs(testAccountID).WillReturnError(pgx.ErrTxClosed)
			},
			http.StatusInternalServerError, "tx is closed",
		},
		{
			"success getting account", http.MethodGet, "/v1/accounts/" + testAccountID.String(), true,
			nil,
			nil,
			func(mock pgxmock.PgxPoolIface) {
				mock.ExpectQuery("SELECT *").WithArgs(testAccountID).WillReturnRows(pgxmock.NewRows(accountsRowCols).AddRow(testAccountID.String(), testAccountName, &testOffBudget, testCategory,
					testAdapter, testAccountTime, testAccountTime, testAmount))
			},
			http.StatusOK, testAccountID.String(),
		},
	}
	executeTests(t, tests)
}

func TestGetAdapters(t *testing.T) {
	tests := []testCase{
		{
			"error due to auth", http.MethodGet, "/v1/adapters", false, nil,
			nil, nil,
			http.StatusUnauthorized, "Unauthorized",
		},
		{
			"success", http.MethodGet, "/v1/adapters", true, nil,
			nil,
			nil,
			http.StatusOK, "icici",
		},
	}
	executeTests(t, tests)
}
