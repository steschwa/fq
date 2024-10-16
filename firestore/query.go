package firestore

import (
	"context"
	"errors"
	"fmt"
	"time"

	"cloud.google.com/go/firestore"
	"cloud.google.com/go/firestore/apiv1/firestorepb"
	"google.golang.org/api/iterator"
)

const (
	timeoutRunQuery = 30
)

type QueryBuilder struct {
	query firestore.Query
}

func NewQueryBuilder(client *firestore.Client, path string) *QueryBuilder {
	return &QueryBuilder{
		query: client.Collection(path).Query,
	}
}

func (b *QueryBuilder) SetWheres(wheres []Where) *QueryBuilder {
	for _, where := range wheres {
		b.applyWhere(where)
	}

	return b
}

func (b *QueryBuilder) SetOrderBy(orderBy string, dir firestore.Direction) *QueryBuilder {
	if orderBy == "" {
		return b
	}

	b.query = b.query.OrderBy(orderBy, dir)

	return b
}

func (b *QueryBuilder) SetLimit(limit int) *QueryBuilder {
	if limit <= 0 {
		return b
	}

	b.query = b.query.Limit(limit)

	return b
}

func (b *QueryBuilder) applyWhere(where Where) {
	b.query = b.query.Where(string(where.Key), where.Operator.String(), where.Value.Value())
}

func (b QueryBuilder) GetDocs() ([]any, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*timeoutRunQuery)
	defer cancel()

	iter := b.query.Documents(ctx)
	defer iter.Stop()

	var out []any
	for {
		doc, err := iter.Next()
		if errors.Is(err, iterator.Done) {
			break
		}
		if errors.Is(err, context.Canceled) {
			return nil, fmt.Errorf("timed-out after %d seconds", timeoutRunQuery)
		}
		if err != nil {
			return nil, err
		}

		out = append(out, doc.Data())
	}

	return out, nil
}

func (b QueryBuilder) GetCount() (int, error) {
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
