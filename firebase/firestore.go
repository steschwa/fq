package firebase

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
	"github.com/steschwa/fq/utils"
)

func InitFirestoreClient(projectID string) (*firestore.Client, error) {
	setupFirebase(projectID)

	c, err := firestore.NewClient(context.Background(), projectID)
	if err != nil {
		log.Println(err)
		return nil, errors.New("failed to create firestore client")
	}

	return c, nil
}

type (
	FirestoreWhere struct {
		Path     string               `parser:"@(Ident (Dot Ident)*)"`
		Operator string               `parser:"@Operator"`
		Value    *FirestoreWhereValue `parser:"@@"`
	}

	FirestoreWhereValue struct {
		String *string                `parser:"@String"`
		Number *float64               `parser:"| @Number"`
		True   bool                   `parser:"| @'true'"`
		False  bool                   `parser:"| @'false'"`
		Null   bool                   `parser:"| @'null'"`
		List   []*FirestoreWhereValue `parser:"| '[' ( @@ ( ',' @@ )* )? ']'"`
	}
)

var (
	firestoreWhereLexer = lexer.MustSimple([]lexer.SimpleRule{
		{Name: "Operator", Pattern: `\s([!=]=|[<>]=?|in|not-in|array-contains-any|array-contains)\s`},
		{Name: "Ident", Pattern: `[a-zA-Z][a-zA-Z0-9_]*`},
		{Name: "Whitespace", Pattern: `\s+`},
		{Name: "Bracket", Pattern: `[\[\]]`},
		{Name: "Comma", Pattern: `,`},
		{Name: "Dot", Pattern: `\.`},
		{Name: "Number", Pattern: `[-+]?\d+(?:\.\d+)?`},
		{Name: "String", Pattern: `"[^"]*"`},
	})

	firestoreWhereParser = participle.MustBuild[FirestoreWhere](
		participle.Lexer(firestoreWhereLexer),
		participle.Elide("Whitespace"),
		participle.Unquote("String"),
	)
)

func (w *FirestoreWhere) GetOperator() string {
	return strings.Trim(w.Operator, " ")
}

func (w *FirestoreWhere) GetValue() any {
	if w.Value == nil {
		return nil
	}

	return w.Value.getValue()
}

func (v *FirestoreWhereValue) getValue() any {
	if v.Null {
		return nil
	} else if v.False {
		return false
	} else if v.True {
		return true
	} else if v.Number != nil {
		return *v.Number
	} else if v.String != nil {
		return *v.String
	} else if v.List != nil {
		l := make([]any, len(v.List))
		for i, listValue := range v.List {
			l[i] = listValue.getValue()
		}
		return l
	}

	return nil
}

func (w *FirestoreWhere) Debug() {
	fmt.Println(fmt.Sprintf("Path: %s, Op: %s, Value: %v (%T)", w.Path, w.GetOperator(), w.GetValue(), w.GetValue()))
}

func ParseFirestoreWhere(where string) (*FirestoreWhere, error) {
	return firestoreWhereParser.ParseString("", where)
}

func ValidateFirestoreCollectionPath(path string) error {
	if path == "" {
		return errors.New("path is empty")
	}

	parts := strings.Split(path, "/")
	if len(parts)%2 == 0 {
		return errors.New("collection paths must contain an uneven amount of parts")
	}

	return nil
}

type (
	FirestoreQueryBuilder struct {
		c *firestore.Client
		q firestore.Query
	}

	FirestoreDocs []map[string]any
)

func NewQueryBuilder(client *firestore.Client) *FirestoreQueryBuilder {
	return &FirestoreQueryBuilder{
		c: client,
	}
}

func (qb *FirestoreQueryBuilder) Collection(path string) *FirestoreQueryBuilder {
	qb.q = qb.c.Collection(path).Query
	return qb
}

func (qb *FirestoreQueryBuilder) WithWheres(wheres []*FirestoreWhere) *FirestoreQueryBuilder {
	for _, where := range wheres {
		qb.q = qb.q.Where(where.Path, where.GetOperator(), where.GetValue())
	}
	return qb
}

func (qb *FirestoreQueryBuilder) WithLimit(limit int) *FirestoreQueryBuilder {
	if limit <= 0 {
		return qb
	}

	qb.q = qb.q.Limit(limit)
	return qb
}

func (qb *FirestoreQueryBuilder) WithOrderBy(orderBy string, dir firestore.Direction) *FirestoreQueryBuilder {
	if orderBy == "" {
		return qb
	}

	qb.q = qb.q.OrderBy(orderBy, dir)
	return qb
}

func (qb *FirestoreQueryBuilder) Execute() (FirestoreDocs, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	docs := qb.q.Documents(ctx)
	snapshots, err := docs.GetAll()
	if err != nil {
		log.Println(err)
		return nil, errors.New("failed to retrieve documents")
	}

	data := make(FirestoreDocs, len(snapshots))
	for i, snapshot := range snapshots {
		if !snapshot.Exists() {
			continue
		}
		data[i] = snapshot.Data()
	}

	return data, nil
}

func (d FirestoreDocs) ToJSON() (string, error) {
	jsonData, err := json.Marshal(d)
	if err != nil {
		log.Println(err)
		return "", errors.New("failed to serialize firestore docs to json")
	}

	if utils.IsInteractiveTTY() {
		jsonString, err := utils.PrettifyJSON(jsonData)
		if err != nil {
			jsonString = string(jsonData)
		}

		return jsonString, nil
	}

	return string(jsonData), nil
}
