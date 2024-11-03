package handlers

import (
	"errors"
	"net/http"
	"strings"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/pashagolub/pgxmock/v4"
)

var (
	testPayeeName = "Swiggy"
	payeeRowCols  = []string{"id", "name", "created_at", "updated_at"}
)

func TestCreatePayee(t *testing.T) {
	tests := []testCase{
		{
			"error due to auth", http.MethodPost, "/v1/payees", false, strings.NewReader("invalid-body"),
			nil, nil,
			http.StatusUnauthorized, "Unauthorized",
		},
		{
			"error due to bad request", http.MethodPost, "/v1/payees", true, strings.NewReader("invalid-body"),
			nil, nil,
			http.StatusBadRequest, "invalid character",
		},
		{
			"error inserting payee to database", http.MethodPost, "/v1/payees", true,
			strings.NewReader(`{"name":"` + testPayeeName + `"}`),
			nil,
			func(mock pgxmock.PgxPoolIface) {
				mock.ExpectExec("INSERT INTO payees").WithArgs(pgxmock.AnyArg(), testPayeeName,
					pgxmock.AnyArg(), pgxmock.AnyArg()).WillReturnError(pgx.ErrTxClosed)
			},
			http.StatusInternalServerError, "tx is closed",
		},
		{
			"success creating payee", http.MethodPost, "/v1/payees", true,
			strings.NewReader(`{"name":"` + testPayeeName + `"}`),
			nil,
			func(mock pgxmock.PgxPoolIface) {
				mock.ExpectExec("INSERT INTO payees").WithArgs(pgxmock.AnyArg(), testPayeeName,
					pgxmock.AnyArg(), pgxmock.AnyArg()).WillReturnResult(pgxmock.NewResult("INSERT", 1))
			},
			http.StatusCreated, testPayeeName,
		},
	}
	executeTests(t, tests)
}

func TestUpdatePayee(t *testing.T) {
	tests := []testCase{
		{
			"error due to auth", http.MethodPatch, "/v1/payees/invalid-uuid", false, strings.NewReader("invalid-body"),
			nil, nil,
			http.StatusUnauthorized, "Unauthorized",
		},
		{
			"error due to bad payee id", http.MethodPatch, "/v1/payees/invalid-uuid", true, strings.NewReader("invalid-body"),
			nil, nil,
			http.StatusBadRequest, "invalid UUID",
		},
		{
			"error due to bad request", http.MethodPatch, "/v1/payees/" + testPayeeID.String(), true, strings.NewReader("invalid-body"),
			nil, nil,
			http.StatusBadRequest, "invalid character",
		},
		{
			"error updating payee in database", http.MethodPatch, "/v1/payees/" + testPayeeID.String(), true,
			strings.NewReader(`{"name":"` + testPayeeName + `"}`),
			nil,
			func(mock pgxmock.PgxPoolIface) {
				mock.ExpectExec("UPDATE payees").WithArgs(
					testPayeeName, pgxmock.AnyArg(), testPayeeID,
				).WillReturnError(pgx.ErrTxClosed)
			},
			http.StatusInternalServerError, "tx is closed",
		},
		{
			"success updating payee", http.MethodPatch, "/v1/payees/" + testPayeeID.String(), true,
			strings.NewReader(`{"name":"` + testPayeeName + `"}`),
			nil,
			func(mock pgxmock.PgxPoolIface) {
				mock.ExpectExec("UPDATE payees").WithArgs(
					testPayeeName, pgxmock.AnyArg(), testPayeeID,
				).WillReturnResult(pgxmock.NewResult("UPDATE", 1))
			},
			http.StatusNoContent, "",
		},
	}
	executeTests(t, tests)
}

func TestDeletePayee(t *testing.T) {
	tests := []testCase{
		{
			"error due to auth", http.MethodDelete, "/v1/payees/invalid-uuid", false, nil,
			nil, nil,
			http.StatusUnauthorized, "Unauthorized",
		},
		{
			"error due to bad payee id", http.MethodDelete, "/v1/payees/invalid-uuid", true, nil,
			nil, nil,
			http.StatusBadRequest, "invalid UUID",
		},
		{
			"error deleting payee in database", http.MethodDelete, "/v1/payees/" + testPayeeID.String(), true,
			nil,
			nil,
			func(mock pgxmock.PgxPoolIface) {
				mock.ExpectExec("DELETE FROM payees").WithArgs(testPayeeID).WillReturnError(pgx.ErrTxClosed)
			},
			http.StatusInternalServerError, "tx is closed",
		},
		{
			"success deleting payee", http.MethodDelete, "/v1/payees/" + testPayeeID.String(), true,
			nil,
			nil,
			func(mock pgxmock.PgxPoolIface) {
				mock.ExpectExec("DELETE FROM payees").WithArgs(testPayeeID).WillReturnResult(pgxmock.NewResult("DELETE", 1))
			},
			http.StatusNoContent, "",
		},
	}
	executeTests(t, tests)
}

func TestGetPayees(t *testing.T) {
	tests := []testCase{
		{
			"error due to auth", http.MethodGet, "/v1/payees", false, nil,
			nil, nil,
			http.StatusUnauthorized, "Unauthorized",
		},
		{
			"error getting payees from db", http.MethodGet, "/v1/payees?q=some", true, nil,
			nil,
			func(mock pgxmock.PgxPoolIface) {
				mock.ExpectQuery("SELECT *").WithArgs("some").WillReturnError(pgx.ErrNoRows)
			},
			http.StatusInternalServerError, "no rows",
		},
		{
			"error scanning payees rows from db", http.MethodGet, "/v1/payees?q=some", true, nil,
			nil,
			func(mock pgxmock.PgxPoolIface) {
				mock.ExpectQuery("SELECT *").WithArgs("some").WillReturnRows(pgxmock.NewRows(payeeRowCols).
					AddRow("invalid", "ok", "bad-time", "bad-time"))
			},
			http.StatusInternalServerError, "Scanning value error",
		},
		{
			"error reading payees rows from db", http.MethodGet, "/v1/payees?q=some", true, nil,
			nil,
			func(mock pgxmock.PgxPoolIface) {
				mock.ExpectQuery("SELECT *").WithArgs("some").WillReturnRows(pgxmock.NewRows(payeeRowCols).
					RowError(0, errors.New("some error in db")))
			},
			http.StatusInternalServerError, "some error in db",
		},
		{
			"success", http.MethodGet, "/v1/payees?q=some", true, nil,
			nil,
			func(mock pgxmock.PgxPoolIface) {
				mock.ExpectQuery("SELECT *").WithArgs("some").WillReturnRows(pgxmock.NewRows(payeeRowCols).
					AddRow(testPayeeID.String(), testPayeeName, testAccountTime, testAccountTime))
			},
			http.StatusOK, testPayeeID.String(),
		},
	}
	executeTests(t, tests)
}
