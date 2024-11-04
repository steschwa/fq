package firestore

import (
	"context"
	"errors"
	"fmt"
	"time"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
)

type DeleteClient struct {
	client *firestore.Client
	path   string
	wheres []Where
}

func NewDeleteClient(client *firestore.Client, path string) *DeleteClient {
	return &DeleteClient{
		client: client,
		path:   path,
	}
}

func (c *DeleteClient) SetWheres(wheres []Where) {
	c.wheres = wheres
}

func (c DeleteClient) Exec() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*timeoutRunQuery)
	defer cancel()

	if IsCollectionPath(c.path) {
		return c.deleteMany(ctx)
	} else if IsDocumentPath(c.path) {
		return c.deleteOne(ctx)
	}

	return nil
}

func (c DeleteClient) deleteOne(ctx context.Context) error {
	_, err := c.client.Doc(c.path).Delete(ctx)
	if errors.Is(err, context.Canceled) {
		return fmt.Errorf("deleting document timed-out")
	}
	if err != nil {
		return fmt.Errorf("deleting doc: %v", err)
	}

	return nil
}

func (c DeleteClient) deleteMany(ctx context.Context) error {
	q := c.client.Collection(c.path)
	for _, where := range c.wheres {
		q.Where(string(where.Key), where.Operator.String(), where.Value.Value())
	}

	iter := q.DocumentRefs(ctx)

	writer := c.client.BulkWriter(ctx)
	defer writer.End()

	for {
		ref, err := iter.Next()
		if errors.Is(err, iterator.Done) {
			break
		}
		if errors.Is(err, context.Canceled) {
			return fmt.Errorf("deleting documents timed-out")
		}
		if err != nil {
			return fmt.Errorf("reading next ref: %v", err)
		}

		_, err = writer.Delete(ref)
		if err != nil {
			return fmt.Errorf("deleting document ref: %v", err)
		}
	}

	return nil
}
