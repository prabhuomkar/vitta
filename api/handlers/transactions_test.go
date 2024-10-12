package handlers

import (
	"net/http"
	"testing"

	uuid "github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/pashagolub/pgxmock/v4"
)

var (
	testTransactionName = "House Party Food"
	testPayeeName       = "Swiggy"
	testCategoryID      = uuid.MustParse("01927f3e-11b6-7a79-bbc7-affae59272ae")
	testPayeeID         = uuid.MustParse("01927f3e-5609-703b-b067-f9b9dd9d8ee2")
	testAccountID       = uuid.MustParse("01927f3e-6ecf-7091-987f-8aa23adcda09")
)

func TestCreateTransaction(t *testing.T) {
	tests := []testCase{
		{
			"error due to auth", http.MethodPost, "/v1/transactions", false, "invalid-body",
			nil,
			http.StatusUnauthorized, "Unauthorized",
		},
		{
			"error due to bad request", http.MethodPost, "/v1/transactions", true, "invalid-body",
			nil,
			http.StatusBadRequest, "invalid character",
		},
		{
			"error inserting transaction to database", http.MethodPost, "/v1/transactions", true,
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
			"success creating transaction", http.MethodPost, "/v1/transactions", true,
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
			"error due to auth", http.MethodPatch, "/v1/transactions/invalid-uuid", false, "invalid-body",
			nil,
			http.StatusUnauthorized, "Unauthorized",
		},
		{
			"error due to bad request", http.MethodPatch, "/v1/transactions/invalid-uuid", true, "invalid-body",
			nil,
			http.StatusBadRequest, "invalid UUID",
		},
		// TODO(omkar): Add more unit test cases
	}
	for _, tc := range tests {
		tc.Run(t)
	}
}

func TestDeleteTransaction(t *testing.T) {
	tests := []testCase{
		{
			"error due to auth", http.MethodDelete, "/v1/transactions/invalid-uuid", false, "invalid-body",
			nil,
			http.StatusUnauthorized, "Unauthorized",
		},
		{
			"error due to bad request", http.MethodDelete, "/v1/transactions/invalid-uuid", true, "invalid-body",
			nil,
			http.StatusBadRequest, "invalid UUID",
		},
		// TODO(omkar): Add more unit test cases
	}
	for _, tc := range tests {
		tc.Run(t)
	}
}

func TestGetTransactions(t *testing.T) {
	tests := []testCase{
		{
			"error due to auth", http.MethodGet, "/v1/transactions", false, "",
			nil,
			http.StatusUnauthorized, "Unauthorized",
		},
		// TODO(omkar): Add more unit test cases
	}
	for _, tc := range tests {
		tc.Run(t)
	}
}

func TestCreatePayee(t *testing.T) {
	tests := []testCase{
		{
			"error due to auth", http.MethodPost, "/v1/payees", false, "invalid-body",
			nil,
			http.StatusUnauthorized, "Unauthorized",
		},
		{
			"error due to bad request", http.MethodPost, "/v1/payees", true, "invalid-body",
			nil,
			http.StatusBadRequest, "invalid character",
		},
		{
			"error inserting payee to database", http.MethodPost, "/v1/payees", true,
			`{"name":"` + testPayeeName + `"}`,
			func(mock pgxmock.PgxPoolIface) {
				mock.ExpectExec("INSERT INTO payees").WithArgs(pgxmock.AnyArg(), testPayeeName,
					pgxmock.AnyArg(), pgxmock.AnyArg()).WillReturnError(pgx.ErrTxClosed)
			},
			http.StatusInternalServerError, "tx is closed",
		},
		{
			"success creating payee", http.MethodPost, "/v1/payees", true,
			`{"name":"` + testPayeeName + `"}`,
			func(mock pgxmock.PgxPoolIface) {
				mock.ExpectExec("INSERT INTO payees").WithArgs(pgxmock.AnyArg(), testPayeeName,
					pgxmock.AnyArg(), pgxmock.AnyArg()).WillReturnResult(pgxmock.NewResult("INSERT", 1))
			},
			http.StatusCreated, testPayeeName,
		},
	}
	for _, tc := range tests {
		tc.Run(t)
	}
}

func TestUpdatePayee(t *testing.T) {
	tests := []testCase{
		{
			"error due to auth", http.MethodPatch, "/v1/payees/invalid-uuid", false, "invalid-body",
			nil,
			http.StatusUnauthorized, "Unauthorized",
		},
		{
			"error due to bad request", http.MethodPatch, "/v1/payees/invalid-uuid", true, "invalid-body",
			nil,
			http.StatusBadRequest, "invalid UUID",
		},
		// TODO(omkar): Add more unit test cases
	}
	for _, tc := range tests {
		tc.Run(t)
	}
}

func TestDeletePayee(t *testing.T) {
	tests := []testCase{
		{
			"error due to auth", http.MethodDelete, "/v1/payees/invalid-uuid", false, "invalid-body",
			nil,
			http.StatusUnauthorized, "Unauthorized",
		},
		{
			"error due to bad request", http.MethodDelete, "/v1/payees/invalid-uuid", true, "invalid-body",
			nil,
			http.StatusBadRequest, "invalid UUID",
		},
		// TODO(omkar): Add more unit test cases
	}
	for _, tc := range tests {
		tc.Run(t)
	}
}

func TestGetPayees(t *testing.T) {
	tests := []testCase{
		{
			"error due to auth", http.MethodGet, "/v1/payees", false, "",
			nil,
			http.StatusUnauthorized, "Unauthorized",
		},
		// TODO(omkar): Add more unit test cases
	}
	for _, tc := range tests {
		tc.Run(t)
	}
}
