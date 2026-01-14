package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"

	database "suscord/internal/domain/storage/database/mocks"
	storage "suscord/internal/domain/storage/mocks"
)

type middlewareTest struct {
	mw      *Middleware
	storage *storage.MockStorage
	db      *database.MockStorage
	session *database.MockSessionStorage
	c       echo.Context
}

func newMiddlewareTest(t *testing.T) *middlewareTest {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	storage := storage.NewMockStorage(t)
	db := database.NewMockStorage(t)
	session := database.NewMockSessionStorage(t)

	storage.EXPECT().
		Database().
		Return(db).
		Maybe()

	db.EXPECT().
		Session().
		Return(session).
		Maybe()

	mw := &Middleware{
		storage: storage,
	}

	return &middlewareTest{
		mw:      mw,
		storage: storage,
		db:      db,
		session: session,
		c:       c,
	}
}
