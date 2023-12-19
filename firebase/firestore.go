package firebase

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strings"

	"cloud.google.com/go/firestore"
	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
	"github.com/steschwa/fq/utils"
)

func InitFirestoreClient(projectID string) (*firestore.Client, error) {
	setupFirebase(projectID)

	c, err := firestore.NewClient(context.Background(), projectID)
	if err != nil {
		slog.Error(err.Error())
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

func ValidateFirestoreWhere(where string) error {
	_, err := firestoreWhereParser.ParseString("", where)
	return err
}

func getFirestorePathParts(path string) []string {
	if path == "" {
		return []string{}
	}

	return strings.Split(path, "/")
}

func ValidateFirestoreCollectionPath(path string) error {
	parts := getFirestorePathParts(path)
	if len(parts) == 0 {
		return errors.New("path is empty")
	}
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

	FirestoreInsertBuilder struct {
		c   *firestore.Client
		ref *firestore.CollectionRef
	}

	FirestoreDoc struct {
		Path string
		Data map[string]any
	}
	FirestoreDocs []FirestoreDoc

	FirestoreInsertData     map[string]any
	FirestoreInsertErrorMap map[string]int
)

func NewQueryBuilder(client *firestore.Client) *FirestoreQueryBuilder {
	return &FirestoreQueryBuilder{
		c: client,
	}
}

func (qb *FirestoreQueryBuilder) Collection(collectionPath string) *FirestoreQueryBuilder {
	qb.q = qb.c.Collection(collectionPath).Query
	return qb
}

func (qb *FirestoreQueryBuilder) WithWheres(wheres []*FirestoreWhere) *FirestoreQueryBuilder {
	for _, where := range wheres {
		qb.q = qb.q.Where(where.Path, where.GetOperator(), where.GetValue())
	}
	return qb
}

func (qb *FirestoreQueryBuilder) WithLimit(limit uint) *FirestoreQueryBuilder {
	if limit == 0 {
		return qb
	}

	qb.q = qb.q.Limit(int(limit))
	return qb
}

func (qb *FirestoreQueryBuilder) WithOrderBy(orderBy string, dir firestore.Direction) *FirestoreQueryBuilder {
	if orderBy == "" {
		return qb
	}

	qb.q = qb.q.OrderBy(orderBy, dir)
	return qb
}

func (qb *FirestoreQueryBuilder) GetAll(ctx context.Context) (FirestoreDocs, error) {
	docs := qb.q.Documents(ctx)
	snapshots, err := docs.GetAll()
	if err != nil {
		slog.Error(err.Error())
		return nil, errors.New("failed to retrieve documents")
	}

	data := make(FirestoreDocs, len(snapshots))
	for i, snapshot := range snapshots {
		if !snapshot.Exists() {
			continue
		}
		data[i] = FirestoreDoc{
			Path: snapshot.Ref.Path,
			Data: snapshot.Data(),
		}
	}

	return data, nil
}

func NewInsertBuilder(client *firestore.Client) *FirestoreInsertBuilder {
	return &FirestoreInsertBuilder{
		c: client,
	}
}

func (ib *FirestoreInsertBuilder) Collection(collectionPath string) *FirestoreInsertBuilder {
	ib.ref = ib.c.Collection(collectionPath)
	return ib
}

func (ib *FirestoreInsertBuilder) InsertMany(ctx context.Context, data []FirestoreInsertData) FirestoreInsertErrorMap {
	bw := ib.c.BulkWriter(ctx)

	jobsCh := make(chan error)
	jobsCount := 0

	for _, item := range data {
		doc := item.GetDoc(ib.ref)
		job, err := bw.Set(doc, item)
		if err != nil {
			continue
		}

		jobsCount++
		go func(job *firestore.BulkWriterJob) {
			_, err := job.Results()
			jobsCh <- err
		}(job)
	}

	bw.Flush()

	errs := map[string]int{}

	finished := 0
	for {
		if finished == jobsCount {
			break
		}

		select {
		case err := <-jobsCh:
			finished++
			if err != nil {
				errs[err.Error()]++
			}
		}
	}

	return errs
}

func (d FirestoreDocs) GetData() []map[string]any {
	data := make([]map[string]any, len(d))
	for i, doc := range d {
		data[i] = doc.Data
	}

	return data
}

func (d FirestoreDocs) ToJSON() (string, error) {
	return utils.ToJSON(d.GetData())
}

func (d FirestoreInsertData) GetDoc(collectionRef *firestore.CollectionRef) *firestore.DocumentRef {
	id, ok := d["id"]
	if !ok {
		return collectionRef.NewDoc()
	}

	switch id := id.(type) {
	case string:
		return collectionRef.Doc(id)
	default:
		return collectionRef.NewDoc()
	}
}

func (errs FirestoreInsertErrorMap) Log() {
	for err, count := range errs {
		slog.Error(err, "count", fmt.Sprint(count))
	}
}
