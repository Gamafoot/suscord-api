package middleware

import (
	"net/http"
	"suscord/internal/domain/entity"
	"testing"

	domainErrors "suscord/internal/domain/errors"

	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestRequiredAuth_Success(t *testing.T) {
	mwCtx := newMiddlewareTest(t)

	mwCtx.c.Request().AddCookie(&http.Cookie{
		Name:  "session",
		Value: "uuid-123",
	})

	mwCtx.session.EXPECT().
		GetByUUID(mock.Anything, "uuid-123").
		Return(&entity.Session{UserID: 42}, nil)

	nextCalled := false

	handler := mwCtx.mw.RequiredAuth()(func(c echo.Context) error {
		nextCalled = true
		result := c.Get("user_id").(uint)
		assert.Equal(t, uint(42), result)
		return c.NoContent(http.StatusOK)
	})

	err := handler(mwCtx.c)
	require.NoError(t, err)

	assert.True(t, nextCalled)
	assert.Equal(t, http.StatusOK, mwCtx.c.Response().Status)
}

func TestRequiredAuth_NoCookie(t *testing.T) {
	mwCtx := newMiddlewareTest(t)

	handler := mwCtx.mw.RequiredAuth()(func(c echo.Context) error {
		return nil
	})

	err := handler(mwCtx.c)
	require.NoError(t, err)

	assert.Equal(t, http.StatusFound, mwCtx.c.Response().Status)
	assert.Equal(t, "/auth", mwCtx.c.Response().Header().Get("Location"))
}

func TestRequiredAuth_ErrRecordNotFound(t *testing.T) {
	mwCtx := newMiddlewareTest(t)

	mwCtx.c.Request().AddCookie(&http.Cookie{
		Name:  "session",
		Value: "uuid-123",
	})

	mwCtx.session.EXPECT().
		GetByUUID(mock.Anything, "uuid-123").
		Return(nil, domainErrors.ErrRecordNotFound)

	handler := mwCtx.mw.RequiredAuth()(func(c echo.Context) error {
		return nil
	})

	err := handler(mwCtx.c)
	require.NoError(t, err)

	assert.Equal(t, http.StatusFound, mwCtx.c.Response().Status)
	assert.Equal(t, "/auth", mwCtx.c.Response().Header().Get("Location"))
}

func TestRequiredAuth_UnknownError(t *testing.T) {
	mwCtx := newMiddlewareTest(t)

	mwCtx.c.Request().AddCookie(&http.Cookie{
		Name:  "session",
		Value: "uuid-123",
	})

	dbErr := errors.New("db closed")
	mwCtx.session.EXPECT().
		GetByUUID(mock.Anything, "uuid-123").
		Return(nil, dbErr)

	nextCalled := false
	handler := mwCtx.mw.RequiredAuth()(func(c echo.Context) error {
		nextCalled = true
		return nil
	})

	err := handler(mwCtx.c)
	require.Error(t, err)
	assert.False(t, nextCalled)
	assert.EqualError(t, err, dbErr.Error())
}
