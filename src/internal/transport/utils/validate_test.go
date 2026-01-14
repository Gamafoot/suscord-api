package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_FilenameValidate(t *testing.T) {
	allowedExts := []string{
		"image",
		"video",
	}

	tests := []struct {
		name     string
		filename string
		want     bool
	}{
		{
			name:     "ok",
			filename: "test.png",
			want:     true,
		},
		{
			name:     "empty",
			filename: "",
			want:     false,
		},
		{
			name:     "filename without ext",
			filename: "test",
			want:     false,
		},
		{
			name:     "filename with not allowed ext",
			filename: "test.mp3",
			want:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FilenameValidate(tt.filename, allowedExts)
			assert.Equal(t, tt.want, result)
		})
	}
}
