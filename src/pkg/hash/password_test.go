package hash

import (
	"crypto/sha1"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHash(t *testing.T) {
	tests := []struct {
		name     string
		salt     string
		password string
	}{
		{
			name:     "normal",
			salt:     "salt",
			password: "password",
		},
		{
			name:     "empty values",
			salt:     "",
			password: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hasher := NewSHA1Hasher(tt.salt)
			expected := expectedResult(tt.salt, tt.password)
			result, err := hasher.Hash(tt.password)
			require.NoError(t, err)
			assert.Equal(t, expected, result)
		})
	}
}

func expectedResult(salt, password string) string {
	h := sha1.New()
	h.Write([]byte(password))
	return fmt.Sprintf("%x", h.Sum([]byte(salt)))
}
