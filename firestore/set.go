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

func (c SetClient) SetMany(data JSONArray, options SetOptions) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*timeoutRunQuery)
	defer cancel()

	var setOptions []firestore.SetOption
	if !options.ReplaceDocument {
		setOptions = append(setOptions, firestore.MergeAll)
	}

	collection := c.client.Collection(c.path)
	writer := c.client.BulkWriter(ctx)
	queued := 0

	for _, obj := range data.Values {
		var doc *firestore.DocumentRef

		if id, found := obj["id"]; found {
			idStr, ok := id.(string)
			if !ok {
				return fmt.Errorf("id must be of type string. got: %T", id)
			}

			doc = collection.Doc(idStr)
		}

		if doc == nil {
			doc = collection.NewDoc()
		}

		_, err := writer.Set(doc, obj, setOptions...)
		if errors.Is(err, context.Canceled) {
			return fmt.Errorf("timed-out after %d seconds", timeoutRunQuery)
		}
		if err != nil {
			return err
		}

		queued++
		if queued > 10 {
			writer.Flush()
			queued = 0
		}
	}

	writer.End()

	return nil
}

func (c SetClient) Set(data JSONObject, options SetOptions) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*timeoutRunQuery)
	defer cancel()

	var setOptions []firestore.SetOption
	if !options.ReplaceDocument {
		setOptions = append(setOptions, firestore.MergeAll)
	}

	_, err := c.client.Doc(c.path).Set(ctx, data.Value, setOptions...)
	if errors.Is(err, context.Canceled) {
		return fmt.Errorf("timed-out after %d seconds", timeoutRunQuery)
	}
	if err != nil {
		return err
	}

	return nil
}
