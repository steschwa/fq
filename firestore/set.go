package firestore

import (
	"context"
	"errors"
	"fmt"
	"time"

	"cloud.google.com/go/firestore"
)

type SetClient struct {
	client *firestore.Client
	path   string
}

func NewSetClient(client *firestore.Client, path string) *SetClient {
	return &SetClient{
		client: client,
		path:   path,
	}
}

func (c SetClient) Set(data any) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*timeoutRunQuery)
	defer cancel()

	_, err := c.client.Doc(c.path).Set(ctx, data, firestore.MergeAll)
	if errors.Is(err, context.Canceled) {
		return fmt.Errorf("timed-out after %d seconds", timeoutRunQuery)
	}
	if err != nil {
		return err
	}

	return nil
}
