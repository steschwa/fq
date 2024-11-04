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
	DeleteClient struct {
		client *firestore.Client
		path   string
		wheres []Where
	}

	DeleteOptions struct {
		ShowProgress bool
		Delay        int
	}
)

func NewDeleteClient(client *firestore.Client, path string) *DeleteClient {
	return &DeleteClient{
		client: client,
		path:   path,
	}
}

func (c *DeleteClient) SetWheres(wheres []Where) {
	c.wheres = wheres
}

func (c DeleteClient) Exec(options DeleteOptions) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*timeoutRunQuery)
	defer cancel()

	if IsCollectionPath(c.path) {
		return c.deleteMany(ctx, options)
	} else if IsDocumentPath(c.path) {
		return c.deleteOne(ctx, options)
	}

	return nil
}

func (c DeleteClient) deleteMany(ctx context.Context, options DeleteOptions) error {
	q := c.client.Collection(c.path).Query
	for _, where := range c.wheres {
		q = q.Where(string(where.Key), where.Operator.String(), where.Value.Value())
	}

	iter := q.Documents(ctx)
	snapshots, err := iter.GetAll()
	if errors.Is(err, context.Canceled) {
		return fmt.Errorf("loading document timed out")
	}
	if err != nil {
		return fmt.Errorf("loading document: %v", err)
	}

	if len(snapshots) == 0 {
		fmt.Println("no documents to delete")
		return nil
	}

	writer := c.client.BulkWriter(ctx)
	defer writer.End()

	for i, snapshot := range snapshots {
		if !snapshot.Exists() {
			continue
		}

		_, err = writer.Delete(snapshot.Ref)
		if err != nil {
			return fmt.Errorf("deleting document ref: %v", err)
		}

		if options.ShowProgress {
			utils.ClearLine()
			fmt.Printf("%d/%d", i+1, len(snapshots))
		}

		if options.Delay > 0 {
			time.Sleep(time.Millisecond * time.Duration(options.Delay))
		}
	}

	return nil
}

func (c DeleteClient) deleteOne(ctx context.Context, options DeleteOptions) error {
	_, err := c.client.Doc(c.path).Delete(ctx)
	if errors.Is(err, context.Canceled) {
		return fmt.Errorf("deleting document timed out")
	}
	if err != nil {
		return fmt.Errorf("deleting doc: %v", err)
	}

	if options.ShowProgress {
		fmt.Printf("1/1")
	}

	return nil
}
