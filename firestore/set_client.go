package firestore

import (
	"context"
	"errors"
	"fmt"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/steschwa/fq/utils"
)

type (
	SetClient struct {
		client *firestore.Client
		path   string
	}

	SetOptions struct {
		ReplaceDocument bool
		ShowProgress    bool
		Delay           int
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

	if len(data.Values) == 0 {
		fmt.Println("empty input data")
		return nil
	}

	var setOptions []firestore.SetOption
	if !options.ReplaceDocument {
		setOptions = append(setOptions, firestore.MergeAll)
	}

	collection := c.client.Collection(c.path)
	writer := c.client.BulkWriter(ctx)
	defer writer.End()

	for i, obj := range data.Values {
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
			return fmt.Errorf("setting documents timed out")
		}
		if err != nil {
			return err
		}

		if options.ShowProgress {
			utils.ClearLine()
			fmt.Printf("%d/%d", i+1, len(data.Values))
		}

		if options.Delay > 0 {
			time.Sleep(time.Millisecond * time.Duration(options.Delay))
		}
	}

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
		return fmt.Errorf("setting document timed out")
	}
	if err != nil {
		return err
	}

	if options.ShowProgress {
		fmt.Printf("1/1")
	}

	return nil
}
