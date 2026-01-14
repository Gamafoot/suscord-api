package urlpath

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetMediaURL(t *testing.T) {
	t.Run("empty filepath", func(t *testing.T) {
		if got := GetMediaURL("/media/", ""); got != "" {
			assert.Empty(t, got)
		}
	})

	t.Run("joins paths", func(t *testing.T) {
		got := GetMediaURL("/media/", "a/b.png")
		want := "/media/a/b.png"
		if got != want {
			assert.Equal(t, want, got)
		}
	})
}
