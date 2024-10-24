package firestore

import (
	"context"
	"errors"
	"fmt"
	"time"

	"cloud.google.com/go/firestore"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	ErrDocumentNotFound = errors.New("document does not exist")
)

type DocLoader struct {
	doc *firestore.DocumentRef
}

func NewDocLoader(client *firestore.Client, path string) *DocLoader {
	return &DocLoader{
		doc: client.Doc(path),
	}
}

func (l DocLoader) GetDoc() (any, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*timeoutRunQuery)
	defer cancel()

	snapshot, err := l.doc.Get(ctx)
	if errors.Is(err, context.Canceled) {
		return nil, fmt.Errorf("timed-out after %d seconds", timeoutRunQuery)
	}
	if status.Code(err) == codes.NotFound {
		return nil, ErrDocumentNotFound
	}
	if err != nil {
		return nil, err
	}

	return snapshot.Data(), nil
}
