package handlers

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
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
	body                 io.Reader
	headers              http.Header
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

		h := New(&config.Config{AdminUsername: "vitta", AdminPassword: "vittaT3st!", UploadMemoryLimit: 1048576}, mockDB)

		res := httptest.NewRecorder()

		handler := h
		handler.ServeHTTP(res, req)

		assert.Equal(t, tc.expectedCode, res.Code)
		assert.Contains(t, res.Body.String(), tc.expectedBodyContains)
	})
}
