package firestore

import (
	"context"
	"errors"
	"fmt"
	"time"

	"cloud.google.com/go/firestore"
	"cloud.google.com/go/firestore/apiv1/firestorepb"
)

const (
	timeoutRunQuery = 30
)

type QueryClient struct {
	query firestore.Query
}

func NewQueryClient(client *firestore.Client, path string) *QueryClient {
	return &QueryClient{
		query: client.Collection(path).Query,
	}
}

func (b *QueryClient) SetWheres(wheres []Where) *QueryClient {
	for _, where := range wheres {
		b.applyWhere(where)
	}

	return b
}

func (b *QueryClient) SetOrderBy(orderBy string, dir firestore.Direction) *QueryClient {
	if orderBy == "" {
		return b
	}

	b.query = b.query.OrderBy(orderBy, dir)

	return b
}

func (b *QueryClient) SetLimit(limit int) *QueryClient {
	if limit <= 0 {
		return b
	}

	b.query = b.query.Limit(limit)

	return b
}

func (b *QueryClient) applyWhere(where Where) {
	b.query = b.query.Where(string(where.Key), where.Operator.String(), where.Value.Value())
}

func (b QueryClient) GetDocs() ([]any, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*timeoutRunQuery)
	defer cancel()

	iter := b.query.Documents(ctx)
	defer iter.Stop()

	docs, err := iter.GetAll()
	if errors.Is(err, context.Canceled) {
		return nil, fmt.Errorf("timed-out after %d seconds", timeoutRunQuery)
	}
	if err != nil {
		return nil, err
	}

	var out []any
	for _, doc := range docs {
		if doc == nil || !doc.Exists() {
			continue
		}

		out = append(out, doc.Data())
	}

	return out, nil
}

func (b QueryClient) GetCount() (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*timeoutRunQuery)
	defer cancel()

	aggr := b.query.NewAggregationQuery().WithCount("count")
	res, err := aggr.Get(ctx)
	if errors.Is(err, context.Canceled) {
		return 0, fmt.Errorf("timed-out after %d seconds", timeoutRunQuery)
	}
	if err != nil {
		return 0, err
	}

	value, ok := res["count"]
	if !ok {
		return 0, fmt.Errorf("missing 'count' key in response")
	}

	if v, ok := value.(*firestorepb.Value); ok {
		return int(v.GetIntegerValue()), nil
	}

	return 0, fmt.Errorf("converting to int. got %T", value)
}
