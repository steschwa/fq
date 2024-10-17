package firestore

import (
	"context"
	"errors"
	"fmt"
	"time"

	"cloud.google.com/go/firestore"
)

type (
	SetClient struct {
		client *firestore.Client
		path   string
	}

	SetOptions struct {
		ReplaceDocument bool
	}
)

func NewSetClient(client *firestore.Client, path string) *SetClient {
	return &SetClient{
		client: client,
		path:   path,
	}
}

func (c SetClient) Set(data any, options SetOptions) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*timeoutRunQuery)
	defer cancel()

	var setOptions []firestore.SetOption
	if !options.ReplaceDocument {
		setOptions = append(setOptions, firestore.MergeAll)
	}

	_, err := c.client.Doc(c.path).Set(ctx, data, setOptions...)
	if errors.Is(err, context.Canceled) {
		return fmt.Errorf("timed-out after %d seconds", timeoutRunQuery)
	}
	if err != nil {
		return err
	}

	return nil
}
