package handlers

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/pashagolub/pgxmock/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	testPayeeName = "Swiggy"
	testRules     = Rules{Includes: []string{"abc"}, Excludes: []string{"xyz"}, StartsWith: []string{"abc"}, EndsWith: []string{"xyz"}}
	payeeRowCols  = []string{"id", "name", "rules", "auto_category_id", "created_at", "updated_at"}
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
			strings.NewReader(`{"name":"` + testPayeeName + `","rules":{"includes":["abc"],"excludes":["xyz"],` +
				`"startsWith":["abc"],"endsWith":["xyz"]},"autoCategoryId":"` + testCategoryID.String() + `"}`),
			nil,
			func(mock pgxmock.PgxPoolIface) {
				mock.ExpectExec("INSERT INTO payees").WithArgs(pgxmock.AnyArg(), testPayeeName,
					&testRules, &testCategoryID,
					pgxmock.AnyArg(), pgxmock.AnyArg()).WillReturnError(pgx.ErrTxClosed)
			},
			http.StatusInternalServerError, "tx is closed",
		},
		{
			"success creating payee", http.MethodPost, "/v1/payees", true,
			strings.NewReader(`{"name":"` + testPayeeName + `","rules":{"includes":["abc"],"excludes":["xyz"],` +
				`"startsWith":["abc"],"endsWith":["xyz"]},"autoCategoryId":"` + testCategoryID.String() + `"}`),
			nil,
			func(mock pgxmock.PgxPoolIface) {
				mock.ExpectExec("INSERT INTO payees").WithArgs(pgxmock.AnyArg(), testPayeeName,
					&testRules, &testCategoryID,
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
			strings.NewReader(`{"name":"` + testPayeeName + `","rules":{"includes":["abc"],"excludes":["xyz"],` +
				`"startsWith":["abc"],"endsWith":["xyz"]},"autoCategoryId":"` + testCategoryID.String() + `"}`),
			nil,
			func(mock pgxmock.PgxPoolIface) {
				mock.ExpectExec("UPDATE payees").WithArgs(
					testPayeeName, &Rules{Includes: []string{"abc"}, Excludes: []string{"xyz"}, StartsWith: []string{"abc"},
						EndsWith: []string{"xyz"}}, &testCategoryID, pgxmock.AnyArg(), testPayeeID,
				).WillReturnError(pgx.ErrTxClosed)
			},
			http.StatusInternalServerError, "tx is closed",
		},
		{
			"success updating payee", http.MethodPatch, "/v1/payees/" + testPayeeID.String(), true,
			strings.NewReader(`{"name":"` + testPayeeName + `","rules":{"includes":["abc"],"excludes":["xyz"],` +
				`"startsWith":["abc"],"endsWith":["xyz"]},"autoCategoryId":"` + testCategoryID.String() + `"}`),
			nil,
			func(mock pgxmock.PgxPoolIface) {
				mock.ExpectExec("UPDATE payees").WithArgs(
					testPayeeName, &Rules{Includes: []string{"abc"}, Excludes: []string{"xyz"}, StartsWith: []string{"abc"},
						EndsWith: []string{"xyz"}}, &testCategoryID, pgxmock.AnyArg(), testPayeeID,
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
			"error getting total payees from db", http.MethodGet, "/v1/payees?q=some", true, nil,
			nil,
			func(mock pgxmock.PgxPoolIface) {
				mock.ExpectQuery("SELECT *").WithArgs("some").WillReturnError(pgx.ErrNoRows)
			},
			http.StatusInternalServerError, "no rows",
		},
		{
			"error scanning total payees rows from db", http.MethodGet, "/v1/payees?q=some", true, nil,
			nil,
			func(mock pgxmock.PgxPoolIface) {
				mock.ExpectQuery("SELECT *").WithArgs("some").WillReturnRows(pgxmock.NewRows([]string{"total"}).AddRow(""))
			},
			http.StatusInternalServerError, "not supported",
		},
		{
			"error getting payees from db", http.MethodGet, "/v1/payees?q=some", true, nil,
			nil,
			func(mock pgxmock.PgxPoolIface) {
				mock.ExpectQuery("SELECT *").WithArgs("some").WillReturnRows(pgxmock.NewRows([]string{"total"}).AddRow(1))
				mock.ExpectQuery("SELECT *").WithArgs("some").WillReturnError(pgx.ErrNoRows)
			},
			http.StatusInternalServerError, "no rows",
		},
		{
			"error scanning payees rows from db", http.MethodGet, "/v1/payees?q=some", true, nil,
			nil,
			func(mock pgxmock.PgxPoolIface) {
				mock.ExpectQuery("SELECT *").WithArgs("some").WillReturnRows(pgxmock.NewRows([]string{"total"}).AddRow(1))
				mock.ExpectQuery("SELECT *").WithArgs("some").WillReturnRows(pgxmock.NewRows(payeeRowCols).
					AddRow("invalid", "ok", "ok", "invalid", "bad-time", "bad-time"))
			},
			http.StatusInternalServerError, "Scanning value error",
		},
		{
			"error reading payees rows from db", http.MethodGet, "/v1/payees?q=some", true, nil,
			nil,
			func(mock pgxmock.PgxPoolIface) {
				mock.ExpectQuery("SELECT *").WithArgs("some").WillReturnRows(pgxmock.NewRows([]string{"total"}).AddRow(1))
				mock.ExpectQuery("SELECT *").WithArgs("some").WillReturnRows(pgxmock.NewRows(payeeRowCols).
					RowError(0, errors.New("some error in db")))
			},
			http.StatusInternalServerError, "some error in db",
		},
		{
			"success", http.MethodGet, "/v1/payees?q=some", true, nil,
			nil,
			func(mock pgxmock.PgxPoolIface) {
				mock.ExpectQuery("SELECT *").WithArgs("some").WillReturnRows(pgxmock.NewRows([]string{"total"}).AddRow(1))
				mock.ExpectQuery("SELECT *").WithArgs("some").WillReturnRows(pgxmock.NewRows(payeeRowCols).
					AddRow(testPayeeID.String(), testPayeeName, &testRules, &testCategoryID, testAccountTime, testAccountTime))
			},
			http.StatusOK, testPayeeID.String(),
		},
	}
	executeTests(t, tests)
}

func TestAssignPayeeAndCategory(t *testing.T) {
	tests := []struct {
		name             string
		mockDBFunc       func(pgxmock.PgxPoolIface)
		errContains      string
		input            string
		expectedPayee    *uuid.UUID
		expectedCategory *uuid.UUID
	}{
		{
			"error getting payees",
			func(mock pgxmock.PgxPoolIface) {
				mock.ExpectQuery("SELECT *").WithArgs("some").WillReturnError(pgx.ErrNoRows)
			},
			"error getting payees", "", nil, nil,
		},
		{
			"error scanning payees row",
			func(mock pgxmock.PgxPoolIface) {
				mock.ExpectQuery("SELECT *").WithArgs("").WillReturnRows(pgxmock.NewRows(payeeRowCols).
					AddRow("invalid", "ok", "ok", "invalid", "bad-time", "bad-time"))
			},
			"error scanning payees row", "", nil, nil,
		},
		{
			"error reading payees rows",
			func(mock pgxmock.PgxPoolIface) {
				mock.ExpectQuery("SELECT *").WithArgs("").WillReturnRows(pgxmock.NewRows(payeeRowCols).
					RowError(0, errors.New("some error in db")))
			},
			"error reading payees rows", "", nil, nil,
		},
		{
			"match includes and excludes",
			func(mock pgxmock.PgxPoolIface) {
				mock.ExpectQuery("SELECT *").WithArgs("").WillReturnRows(pgxmock.NewRows(payeeRowCols).
					AddRow(testPayeeID.String(), testPayeeName, &Rules{Includes: []string{"ato"}, Excludes: []string{"swiggy"}}, &testCategoryID, testAccountTime, testAccountTime))
			},
			"", "ZOMATO", &testPayeeID, &testCategoryID,
		},
		{
			"match starts with",
			func(mock pgxmock.PgxPoolIface) {
				mock.ExpectQuery("SELECT *").WithArgs("").WillReturnRows(pgxmock.NewRows(payeeRowCols).
					AddRow(testPayeeID.String(), testPayeeName, &Rules{StartsWith: []string{"zoma"}}, &testCategoryID, testAccountTime, testAccountTime))
			},
			"", "ZOMATO", &testPayeeID, &testCategoryID,
		},
		{
			"match ends with",
			func(mock pgxmock.PgxPoolIface) {
				mock.ExpectQuery("SELECT *").WithArgs("").WillReturnRows(pgxmock.NewRows(payeeRowCols).
					AddRow(testPayeeID.String(), testPayeeName, &Rules{EndsWith: []string{"ato"}}, &testCategoryID, testAccountTime, testAccountTime))
			},
			"", "ZOMATO", &testPayeeID, &testCategoryID,
		},
		{
			"match with nothing and empty rules",
			func(mock pgxmock.PgxPoolIface) {
				mock.ExpectQuery("SELECT *").WithArgs("").WillReturnRows(pgxmock.NewRows(payeeRowCols).
					AddRow(testPayeeID.String(), testPayeeName, &Rules{}, &testCategoryID, testAccountTime, testAccountTime))
			},
			"", "ZOMATO", nil, nil,
		},
		{
			"match with nothing and single rule",
			func(mock pgxmock.PgxPoolIface) {
				mock.ExpectQuery("SELECT *").WithArgs("").WillReturnRows(pgxmock.NewRows(payeeRowCols).
					AddRow(testPayeeID.String(), testPayeeName, &Rules{Includes: []string{"wig"}, Excludes: []string{"zomato"}, StartsWith: []string{"swi"},
						EndsWith: []string{"gy"}}, &testCategoryID, testAccountTime, testAccountTime))
			},
			"", "ZOMATO", nil, nil,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockDB, err := pgxmock.NewPool()
			require.NoError(t, err)
			defer mockDB.Close()

			if tc.mockDBFunc != nil {
				tc.mockDBFunc(mockDB)
			}

			h := &Handler{cfg: nil, db: mockDB, adapters: nil}
			actFn, err := h.assignPayeeAndCategory(context.TODO())
			if len(tc.errContains) > 0 {
				assert.Error(t, err)
				assert.ErrorContains(t, err, tc.errContains)
			} else {
				assert.NoError(t, err)
				gotPayee, gotCategory := actFn(tc.input)
				assert.Equal(t, tc.expectedPayee, gotPayee)
				assert.Equal(t, tc.expectedCategory, gotCategory)
			}
		})
	}
}
