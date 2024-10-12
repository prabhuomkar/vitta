package handlers

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
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
	body                 string
	mockDBFunc           func(pgxmock.PgxPoolIface)
	expectedCode         int
	expectedBodyContains string
}

func (tc *testCase) Run(t *testing.T) {
	t.Helper()
	t.Run(tc.name, func(t *testing.T) {
		mockDB, err := pgxmock.NewPool()
		require.NoError(t, err)
		defer mockDB.Close()

		if tc.mockDBFunc != nil {
			tc.mockDBFunc(mockDB)
		}

		req, err := http.NewRequestWithContext(context.TODO(), tc.method, tc.path, strings.NewReader(tc.body))
		if err != nil {
			require.NoError(t, err)
		}

		if tc.auth {
			req.SetBasicAuth("vitta", "vittaT3st!")
		}

		h := New(&config.Config{AdminUsername: "vitta", AdminPassword: "vittaT3st!"}, mockDB)

		res := httptest.NewRecorder()

		handler := h
		handler.ServeHTTP(res, req)

		assert.Equal(t, tc.expectedCode, res.Code)
		assert.Contains(t, res.Body.String(), tc.expectedBodyContains)
	})
}
