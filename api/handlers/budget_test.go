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
	for _, tc := range tests {
		tc.Run(t)
	}
}

func TestUpdateGroup(t *testing.T) {
	tests := []testCase{
		{
			"error due to auth", http.MethodPatch, "/v1/groups/invalid-uuid", false, strings.NewReader("invalid-body"),
			nil, nil,
			http.StatusUnauthorized, "Unauthorized",
		},
		{
			"error due to bad request", http.MethodPatch, "/v1/groups/invalid-uuid", true, strings.NewReader("invalid-body"),
			nil, nil,
			http.StatusBadRequest, "invalid UUID",
		},
		// TODO(omkar): Add more unit test cases
	}
	for _, tc := range tests {
		tc.Run(t)
	}
}

func TestDeleteGroup(t *testing.T) {
	tests := []testCase{
		{
			"error due to auth", http.MethodDelete, "/v1/groups/invalid-uuid", false, strings.NewReader("invalid-body"),
			nil, nil,
			http.StatusUnauthorized, "Unauthorized",
		},
		{
			"error due to bad request", http.MethodDelete, "/v1/groups/invalid-uuid", true, strings.NewReader("invalid-body"),
			nil, nil,
			http.StatusBadRequest, "invalid UUID",
		},
		// TODO(omkar): Add more unit test cases
	}
	for _, tc := range tests {
		tc.Run(t)
	}
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
	for _, tc := range tests {
		tc.Run(t)
	}
}

func TestUpdateCategory(t *testing.T) {
	tests := []testCase{
		{
			"error due to auth", http.MethodPatch, "/v1/categories/invalid-uuid", false, strings.NewReader("invalid-body"),
			nil, nil,
			http.StatusUnauthorized, "Unauthorized",
		},
		{
			"error due to bad request", http.MethodPatch, "/v1/categories/invalid-uuid", true, strings.NewReader("invalid-body"),
			nil, nil,
			http.StatusBadRequest, "invalid UUID",
		},
		// TODO(omkar): Add more unit test cases
	}
	for _, tc := range tests {
		tc.Run(t)
	}
}

func TestDeleteCategory(t *testing.T) {
	tests := []testCase{
		{
			"error due to auth", http.MethodDelete, "/v1/categories/invalid-uuid", false, strings.NewReader("invalid-body"),
			nil, nil,
			http.StatusUnauthorized, "Unauthorized",
		},
		{
			"error due to bad request", http.MethodDelete, "/v1/categories/invalid-uuid", true, strings.NewReader("invalid-body"),
			nil, nil,
			http.StatusBadRequest, "invalid UUID",
		},
		// TODO(omkar): Add more unit test cases
	}
	for _, tc := range tests {
		tc.Run(t)
	}
}
