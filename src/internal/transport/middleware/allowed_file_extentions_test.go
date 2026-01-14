package middleware

import (
	"bytes"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"suscord/internal/config"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newMultipartContext(t *testing.T, method, fieldName, filename string, content []byte) echo.Context {
	t.Helper()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	fileWriter, err := writer.CreateFormFile(fieldName, filename)
	require.NoError(t, err)

	_, err = fileWriter.Write(content)
	require.NoError(t, err)

	require.NoError(t, writer.Close())

	e := echo.New()
	req := httptest.NewRequest(method, "/", body)
	req.Header.Set(echo.HeaderContentType, writer.FormDataContentType())
	rec := httptest.NewRecorder()

	return e.NewContext(req, rec)
}

func TestAllowedFileExtentions_OK(t *testing.T) {
	mwCtx := newMiddlewareTest(t)
	mwCtx.mw.config = &config.Config{}
	mwCtx.mw.config.Media.AllowedMedia = []string{".png"}

	c := newMultipartContext(t, http.MethodPost, "file", "test.png", []byte("fake bytes"))

	nextCalled := false

	handler := mwCtx.mw.AllowedFileExtentions()(func(c echo.Context) error {
		nextCalled = true
		return c.NoContent(http.StatusOK)
	})

	err := handler(c)
	require.NoError(t, err)
	require.True(t, nextCalled)

	assert.Equal(t, http.StatusOK, c.Response().Status)
}

func TestAllowedFileExtentions_NotAllowFileExtiontion(t *testing.T) {
	mwCtx := newMiddlewareTest(t)
	mwCtx.mw.config = &config.Config{}
	mwCtx.mw.config.Media.AllowedMedia = []string{".png"}

	c := newMultipartContext(t, http.MethodPost, "file", "test.exe", []byte("fake bytes"))

	handler := mwCtx.mw.AllowedFileExtentions()(func(c echo.Context) error {
		return c.NoContent(http.StatusOK)
	})

	err := handler(c)
	require.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, c.Response().Status)
}

func TestAllowedFileExtentions_NoFile(t *testing.T) {
	mwCtx := newMiddlewareTest(t)
	mwCtx.mw.config = &config.Config{}
	mwCtx.mw.config.Media.AllowedMedia = []string{".png"}

	c := newMultipartContext(t, http.MethodPost, "", "", nil)

	handler := mwCtx.mw.AllowedFileExtentions()(func(c echo.Context) error {
		return c.NoContent(http.StatusOK)
	})

	err := handler(c)
	require.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, c.Response().Status)
}
