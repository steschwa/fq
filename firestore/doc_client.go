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

type DocClient struct {
	doc *firestore.DocumentRef
}

func NewDocClient(client *firestore.Client, path string) *DocClient {
	return &DocClient{
		doc: client.Doc(path),
	}
}

func (l DocClient) GetDoc() (*FirestoreDoc, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*timeoutRunQuery)
	defer cancel()

	snapshot, err := l.doc.Get(ctx)
	if errors.Is(err, context.Canceled) {
		return nil, fmt.Errorf("getting document timed out")
	}
	if status.Code(err) == codes.NotFound {
		return nil, ErrDocumentNotFound
	}
	if err != nil {
		return nil, err
	}

	return NewFirestoreDoc(snapshot.Data()), nil
}
