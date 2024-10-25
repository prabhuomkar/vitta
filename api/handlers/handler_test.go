package handlers

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"vitta/adapters"
	"vitta/config"

	"github.com/pashagolub/pgxmock/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type testCase struct {
	name                 string
	method               string
	path                 string
	auth                 bool
	body                 io.Reader
	headers              http.Header
	mockDBFunc           func(pgxmock.PgxPoolIface)
	expectedCode         int
	expectedBodyContains string
}

func executeTests(t *testing.T, tests []testCase) {
	t.Helper()
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockDB, err := pgxmock.NewPool()
			require.NoError(t, err)
			defer mockDB.Close()

			if tc.mockDBFunc != nil {
				tc.mockDBFunc(mockDB)
			}

			req, err := http.NewRequestWithContext(context.TODO(), tc.method, tc.path, tc.body)
			if err != nil {
				require.NoError(t, err)
			}

			for key, value := range tc.headers {
				req.Header.Add(key, value[0])
			}

			if tc.auth {
				req.SetBasicAuth("vitta", "vittaT3st!")
			}

			h := New(&config.Config{AdminUsername: "vitta", AdminPassword: "vittaT3st!",
				UploadMemoryLimit: 1048576}, mockDB, map[string]adapters.Config{"icici-CC": {
				DateName:        "Transaction Date",
				DateFormats:     []string{"02/01/2006"},
				Remarks:         "Details",
				Credit:          "Amount (INR)",
				Debit:           "Amount (INR)",
				TransactionDiff: []string{"cr.", "dr."},
			}})

			res := httptest.NewRecorder()

			handler := h
			handler.ServeHTTP(res, req)

			assert.Equal(t, tc.expectedCode, res.Code)
			assert.Contains(t, res.Body.String(), tc.expectedBodyContains)
		})
	}
}
