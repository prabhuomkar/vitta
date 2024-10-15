package handlers

import (
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
