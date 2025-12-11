package urlpath

import (
	"path"
)

func GetMediaURL(mediaURL, filepath string) string {
	if len(filepath) == 0 {
		return ""
	}
	return path.Join(mediaURL, filepath)
}
