package firestore

import (
	"errors"
	"strings"
)

var (
	ErrEmptyPath = errors.New("empty path")
)

func ValidatePath(path string) error {
	if path == "" {
		return ErrEmptyPath
	}
	if IsCollectionPath(path) {
		return nil
	}
	if IsDocumentPath(path) {
		return nil
	}

	return nil
}

func IsCollectionPath(path string) bool {
	segments := strings.Split(path, "/")
	return len(segments)%2 == 1 && isValidPathSegments(segments)
}

func IsDocumentPath(path string) bool {
	segments := strings.Split(path, "/")
	return len(segments)%2 == 0 && isValidPathSegments(segments)
}

func isValidPathSegments(segments []string) bool {
	for _, segment := range segments {
		if segment == "" {
			return false
		}
	}

	return true
}
