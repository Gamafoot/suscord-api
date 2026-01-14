package utils

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	domainErrors "suscord/internal/domain/errors"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func Test_GetUIntFromParam(t *testing.T) {
	tests := []struct {
		name       string
		paramValue string
		want       uint
		wantErr    error
	}{
		{
			name:       "ok",
			paramValue: "1",
			want:       1,
			wantErr:    nil,
		},
		{
			name:       "less 0",
			paramValue: "-1",
			want:       0,
			wantErr:    domainErrors.ErrIsNotPositiveDigit,
		},
		{
			name:       "not digit",
			paramValue: "dhaida",
			want:       0,
			wantErr:    domainErrors.ErrIsNotDigit,
		},
		{
			name:       "empty",
			paramValue: "",
			want:       0,
			wantErr:    domainErrors.ErrIsNotDigit,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := echo.New()
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			c.SetParamNames("id")
			c.SetParamValues(tt.paramValue)

			result, err := GetUIntFromParam(c, "id")

			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, tt.want, result)
		})
	}
}

func Test_GetIntFromQuery(t *testing.T) {
	tests := []struct {
		name       string
		queryValue string
		result     int
		err        error
	}{
		{
			name:       "ok",
			queryValue: "1",
			result:     1,
			err:        nil,
		},
		{
			name:       "less 0",
			queryValue: "-1",
			result:     0,
			err:        domainErrors.ErrIsNotPositiveDigit,
		},
		{
			name:       "not digit",
			queryValue: "dhaida",
			result:     0,
			err:        domainErrors.ErrIsNotDigit,
		},
		{
			name:       "empty",
			queryValue: "",
			result:     0,
			err:        domainErrors.ErrIsNotDigit,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			values := url.Values{}
			values.Set("id", tt.queryValue)
			uri := values.Encode()

			req := httptest.NewRequest(http.MethodGet, "/?"+uri, nil)
			rec := httptest.NewRecorder()

			e := echo.New()
			c := e.NewContext(req, rec)

			result, err := GetIntFromQuery(c, "id")

			if tt.err != nil {
				assert.ErrorIs(t, err, tt.err, fmt.Sprintf("Err: %v; Want: %v; Result: %d", err, tt.err, result))
			} else {
				assert.NoError(t, err, fmt.Sprintf("Err: %v; Want: %v; Result: %d", err, tt.err, result))
			}

			assert.Equal(t, tt.result, result)
		})
	}
}
