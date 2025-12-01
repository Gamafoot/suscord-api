package utils

import (
	"path/filepath"
	"strings"
	"suscord/internal/config"
)

func FileExtensionValidate(filename string) bool {
	cfg := config.GetConfig()

	ok := false
	ext := strings.ToLower(filepath.Ext(filename))

	for _, allowExt := range cfg.Media.AllowedExtentions {
		if ext == allowExt {
			ok = true
		}
	}

	return ok
}
