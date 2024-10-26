package handlers

import (
	"errors"
	"net/http"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/pashagolub/pgxmock/v4"
)

var (
	testGroupName    = "Recurring Expenses"
	testCategoryName = "Swiggy Online"
	testGroupID      = uuid.MustParse("01927f36-44b8-7e62-a7a9-395eacab562b")
	testBudgetID     = uuid.MustParse("01928eeb-b65c-7840-beaa-a1adb8c75956")
	budgetRowCols    = []string{"id", "budgeted", "spent", "year", "month", "category_id",
		"category_name", "group_id", "group_name"}
)

func TestCreateGroup(t *testing.T) {
	tests := []testCase{
		{
			"error due to auth", http.MethodPost, "/v1/groups", false, strings.NewReader("invalid-body"),
			nil, nil,
			http.StatusUnauthorized, "Unauthorized",
		},
		{
			"error due to bad request", http.MethodPost, "/v1/groups", true, strings.NewReader("invalid-body"),
			nil, nil,
			http.StatusBadRequest, "invalid character",
		},
		{
			"error inserting group to database", http.MethodPost, "/v1/groups", true,
			strings.NewReader(`{"name":"` + testGroupName + `"}`),
			nil,
			func(mock pgxmock.PgxPoolIface) {
				mock.ExpectExec("INSERT INTO groups").WithArgs(pgxmock.AnyArg(), testGroupName,
					"", pgxmock.AnyArg(), pgxmock.AnyArg()).WillReturnError(pgx.ErrTxClosed)
			},
			http.StatusInternalServerError, "tx is closed",
		},
		{
			"success creating group", http.MethodPost, "/v1/groups", true,
			strings.NewReader(`{"name":"` + testGroupName + `"}`),
			nil,
			func(mock pgxmock.PgxPoolIface) {
				mock.ExpectExec("INSERT INTO groups").WithArgs(pgxmock.AnyArg(), testGroupName,
					"", pgxmock.AnyArg(), pgxmock.AnyArg()).WillReturnResult(pgxmock.NewResult("INSERT", 1))
			},
			http.StatusCreated, testGroupName,
		},
	}
	executeTests(t, tests)
}

func TestUpdateGroup(t *testing.T) {
	tests := []testCase{
		{
			"error due to auth", http.MethodPatch, "/v1/groups/invalid-uuid", false, strings.NewReader("invalid-body"),
			nil, nil,
			http.StatusUnauthorized, "Unauthorized",
		},
		{
			"error due to bad group id", http.MethodPatch, "/v1/groups/invalid-uuid", true, strings.NewReader("invalid-body"),
			nil, nil,
			http.StatusBadRequest, "invalid UUID",
		},
		{
			"error due to bad request", http.MethodPatch, "/v1/groups/" + testGroupID.String(), true, strings.NewReader("invalid-body"),
			nil, nil,
			http.StatusBadRequest, "invalid character",
		},
		{
			"error updating group in database", http.MethodPatch, "/v1/groups/" + testGroupID.String(), true,
			strings.NewReader(`{"name":"` + testGroupName + `","notes":"Some notes"}`),
			nil,
			func(mock pgxmock.PgxPoolIface) {
				mock.ExpectExec("UPDATE groups").WithArgs(
					testGroupName, "Some notes", pgxmock.AnyArg(), testGroupID,
				).WillReturnError(pgx.ErrTxClosed)
			},
			http.StatusInternalServerError, "tx is closed",
		},
		{
			"success updating group", http.MethodPatch, "/v1/groups/" + testGroupID.String(), true,
			strings.NewReader(`{"name":"` + testGroupName + `","notes":"Some notes"}`),
			nil,
			func(mock pgxmock.PgxPoolIface) {
				mock.ExpectExec("UPDATE groups").WithArgs(
					testGroupName, "Some notes", pgxmock.AnyArg(), testGroupID,
				).WillReturnResult(pgxmock.NewResult("UPDATE", 1))
			},
			http.StatusNoContent, "",
		},
	}
	executeTests(t, tests)
}

func TestDeleteGroup(t *testing.T) {
	tests := []testCase{
		{
			"error due to auth", http.MethodDelete, "/v1/groups/invalid-uuid", false, nil,
			nil, nil,
			http.StatusUnauthorized, "Unauthorized",
		},
		{
			"error due to bad group id", http.MethodDelete, "/v1/groups/invalid-uuid", true, nil,
			nil, nil,
			http.StatusBadRequest, "invalid UUID",
		},
		{
			"error deleting group in database", http.MethodDelete, "/v1/groups/" + testGroupID.String(), true,
			nil,
			nil,
			func(mock pgxmock.PgxPoolIface) {
				mock.ExpectExec("DELETE FROM groups").WithArgs(testGroupID).WillReturnError(pgx.ErrTxClosed)
			},
			http.StatusInternalServerError, "tx is closed",
		},
		{
			"success deleting group", http.MethodDelete, "/v1/groups/" + testGroupID.String(), true,
			nil,
			nil,
			func(mock pgxmock.PgxPoolIface) {
				mock.ExpectExec("DELETE FROM groups").WithArgs(testGroupID).WillReturnResult(pgxmock.NewResult("DELETE", 1))
			},
			http.StatusNoContent, "",
		},
	}
	executeTests(t, tests)
}

func TestCreateCategory(t *testing.T) {
	tests := []testCase{
		{
			"error due to auth", http.MethodPost, "/v1/categories", false, strings.NewReader("invalid-body"),
			nil, nil,
			http.StatusUnauthorized, "Unauthorized",
		},
		{
			"error due to bad request", http.MethodPost, "/v1/categories", true, strings.NewReader("invalid-body"),
			nil, nil,
			http.StatusBadRequest, "invalid character",
		},
		{
			"error inserting category to database", http.MethodPost, "/v1/categories", true,
			strings.NewReader(`{"name":"` + testCategoryName + `","groupId":"` + testGroupID.String() + `"}`),
			nil,
			func(mock pgxmock.PgxPoolIface) {
				mock.ExpectExec("INSERT INTO categories").WithArgs(pgxmock.AnyArg(),
					testGroupID, testCategoryName, "", pgxmock.AnyArg(),
					pgxmock.AnyArg()).WillReturnError(pgx.ErrTxClosed)
			},
			http.StatusInternalServerError, "tx is closed",
		},
		{
			"success creating category", http.MethodPost, "/v1/categories", true,
			strings.NewReader(`{"name":"` + testCategoryName + `","groupId":"` + testGroupID.String() + `"}`),
			nil,
			func(mock pgxmock.PgxPoolIface) {
				mock.ExpectExec("INSERT INTO categories").WithArgs(pgxmock.AnyArg(),
					testGroupID, testCategoryName, "", pgxmock.AnyArg(),
					pgxmock.AnyArg()).WillReturnResult(pgxmock.NewResult("INSERT", 1))
			},
			http.StatusCreated, testCategoryName,
		},
	}
	executeTests(t, tests)
}

func TestUpdateCategory(t *testing.T) {
	tests := []testCase{
		{
			"error due to auth", http.MethodPatch, "/v1/categories/invalid-uuid", false, strings.NewReader("invalid-body"),
			nil, nil,
			http.StatusUnauthorized, "Unauthorized",
		},
		{
			"error due to bad category id", http.MethodPatch, "/v1/categories/invalid-uuid", true, strings.NewReader("invalid-body"),
			nil, nil,
			http.StatusBadRequest, "invalid UUID",
		},
		{
			"error due to bad request", http.MethodPatch, "/v1/categories/" + testCategoryID.String(), true, strings.NewReader("invalid-body"),
			nil, nil,
			http.StatusBadRequest, "invalid character",
		},
		{
			"error updating category in database", http.MethodPatch, "/v1/categories/" + testCategoryID.String(), true,
			strings.NewReader(`{"name":"` + testCategoryName + `","notes":"Some notes","groupId":"` + testGroupID.String() + `"}`),
			nil,
			func(mock pgxmock.PgxPoolIface) {
				mock.ExpectExec("UPDATE categories").WithArgs(
					testCategoryName, "Some notes", testGroupID, pgxmock.AnyArg(), testCategoryID,
				).WillReturnError(pgx.ErrTxClosed)
			},
			http.StatusInternalServerError, "tx is closed",
		},
		{
			"success updating category", http.MethodPatch, "/v1/categories/" + testCategoryID.String(), true,
			strings.NewReader(`{"name":"` + testCategoryName + `","notes":"Some notes","groupId":"` + testGroupID.String() + `"}`),
			nil,
			func(mock pgxmock.PgxPoolIface) {
				mock.ExpectExec("UPDATE categories").WithArgs(
					testCategoryName, "Some notes", testGroupID, pgxmock.AnyArg(), testCategoryID,
				).WillReturnResult(pgxmock.NewResult("UPDATE", 1))
			},
			http.StatusNoContent, "",
		},
	}
	executeTests(t, tests)
}

func TestDeleteCategory(t *testing.T) {
	tests := []testCase{
		{
			"error due to auth", http.MethodDelete, "/v1/categories/invalid-uuid", false, nil,
			nil, nil,
			http.StatusUnauthorized, "Unauthorized",
		},
		{
			"error due to bad category id", http.MethodDelete, "/v1/categories/invalid-uuid", true, nil,
			nil, nil,
			http.StatusBadRequest, "invalid UUID",
		},
		{
			"error deleting category in database", http.MethodDelete, "/v1/categories/" + testCategoryID.String(), true,
			nil,
			nil,
			func(mock pgxmock.PgxPoolIface) {
				mock.ExpectExec("DELETE FROM categories").WithArgs(testCategoryID).WillReturnError(pgx.ErrTxClosed)
			},
			http.StatusInternalServerError, "tx is closed",
		},
		{
			"success deleting category", http.MethodDelete, "/v1/categories/" + testCategoryID.String(), true,
			nil,
			nil,
			func(mock pgxmock.PgxPoolIface) {
				mock.ExpectExec("DELETE FROM categories").WithArgs(testCategoryID).WillReturnResult(pgxmock.NewResult("DELETE", 1))
			},
			http.StatusNoContent, "",
		},
	}
	executeTests(t, tests)
}

func TestSetBudget(t *testing.T) {
	tests := []testCase{
		{
			"error due to auth", http.MethodPut, "/v1/budgets", false, strings.NewReader("invalid-body"),
			nil, nil,
			http.StatusUnauthorized, "Unauthorized",
		},
		{
			"error due to bad request", http.MethodPut, "/v1/budgets", true, strings.NewReader("invalid-body"),
			nil, nil,
			http.StatusBadRequest, "invalid character",
		},
		{
			"error setting budget in database", http.MethodPut, "/v1/budgets", true,
			strings.NewReader(`{"categoryId":"` + testCategoryID.String() + `","year":2024,"month":10,"budgeted":500.69}`),
			nil,
			func(mock pgxmock.PgxPoolIface) {
				mock.ExpectExec("INSERT INTO budgets").WithArgs(pgxmock.AnyArg(), testCategoryID,
					uint16(2024), uint8(10), 500.69, pgxmock.AnyArg(), pgxmock.AnyArg()).WillReturnError(pgx.ErrTxClosed)
			},
			http.StatusInternalServerError, "tx is closed",
		},
		{
			"success setting budget", http.MethodPut, "/v1/budgets", true,
			strings.NewReader(`{"categoryId":"` + testCategoryID.String() + `","year":2024,"month":10,"budgeted":500.69}`),
			nil,
			func(mock pgxmock.PgxPoolIface) {
				mock.ExpectExec("INSERT INTO budgets").WithArgs(pgxmock.AnyArg(), testCategoryID,
					uint16(2024), uint8(10), 500.69, pgxmock.AnyArg(), pgxmock.AnyArg()).WillReturnResult(pgxmock.NewResult("INSERT", 1))
			},
			http.StatusOK, "500.69",
		},
	}
	executeTests(t, tests)
}

func TestGetBudget(t *testing.T) {
	tests := []testCase{
		{
			"error due to auth", http.MethodGet, "/v1/budgets?year=2024&month=10", false, nil,
			nil, nil,
			http.StatusUnauthorized, "Unauthorized",
		},
		{
			"error due to bad year", http.MethodGet, "/v1/budgets?year=NA&month=10", true, nil,
			nil, nil,
			http.StatusBadRequest, "parsing",
		},
		{
			"error due to bad month", http.MethodGet, "/v1/budgets?year=2024&month=NA", true, nil,
			nil, nil,
			http.StatusBadRequest, "parsing",
		},
		{
			"error getting budgets from db", http.MethodGet, "/v1/budgets?year=2024&month=10", true, nil,
			nil,
			func(mock pgxmock.PgxPoolIface) {
				mock.ExpectQuery("SELECT *").WithArgs(2024, 10).WillReturnError(pgx.ErrNoRows)
			},
			http.StatusInternalServerError, "no rows",
		},
		{
			"error scanning budgets rows from db", http.MethodGet, "/v1/budgets?year=2024&month=10", true, nil,
			nil,
			func(mock pgxmock.PgxPoolIface) {
				mock.ExpectQuery("SELECT *").WithArgs(2024, 10).WillReturnRows(pgxmock.NewRows(budgetRowCols).
					AddRow("invalid", "budgeted",
						"spent", "year", "month", "category_id", "category_name", "group_id", "group_name"))
			},
			http.StatusInternalServerError, "Scanning value error",
		},
		{
			"error reading budgets rows from db", http.MethodGet, "/v1/budgets?year=2024&month=10", true, nil,
			nil,
			func(mock pgxmock.PgxPoolIface) {
				mock.ExpectQuery("SELECT *").WithArgs(2024, 10).WillReturnRows(pgxmock.NewRows(budgetRowCols).
					RowError(0, errors.New("some error in db")))
			},
			http.StatusInternalServerError, "some error in db",
		},
		{
			"success", http.MethodGet, "/v1/budgets?year=2024&month=10", true, nil,
			nil,
			func(mock pgxmock.PgxPoolIface) {
				mock.ExpectQuery("SELECT *").WithArgs(2024, 10).WillReturnRows(pgxmock.NewRows(budgetRowCols).
					AddRow(testBudgetID, 500.69, 4.20, uint16(2024), uint8(10), testCategoryID, testCategoryName,
						testGroupID, testGroupName))
			},
			http.StatusOK, testBudgetID.String(),
		},
	}
	executeTests(t, tests)
}
