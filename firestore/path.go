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
	parts := strings.Split(path, "/")
	return len(parts)%2 == 1
}

func IsDocumentPath(path string) bool {
	parts := strings.Split(path, "/")
	return len(parts)%2 == 0
}
