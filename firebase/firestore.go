package firebase

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strings"

	"cloud.google.com/go/firestore"
	"cloud.google.com/go/firestore/apiv1/firestorepb"
	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
	"github.com/steschwa/fq/utils"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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

type (
	FirestorePathType int
)

const (
	FirestorePathTypeCollection FirestorePathType = iota
	FirestorePathTypeDocument
)

func GetFirestorePathType(path string) (FirestorePathType, error) {
	if path == "" {
		return FirestorePathTypeCollection, errors.New("path can't be empty")
	}
	if strings.HasPrefix(path, "/") {
		return FirestorePathTypeCollection, errors.New("path can't start with a /")
	}
	if strings.HasSuffix(path, "/") {
		return FirestorePathTypeCollection, errors.New("path can't end with a /")
	}

	if IsFirestoreCollectionPath(path) {
		return FirestorePathTypeCollection, nil
	} else if IsFirestoreDocumentPath(path) {
		return FirestorePathTypeDocument, nil
	}

	return FirestorePathTypeCollection, errors.New("unknown firestore path type")
}

func IsFirestoreCollectionPath(path string) bool {
	parts := getFirestorePathParts(path)
	if len(parts) == 0 {
		return false
	}

	return len(parts)%2 == 1
}

func IsFirestoreDocumentPath(path string) bool {
	parts := getFirestorePathParts(path)
	if len(parts) == 0 {
		return false
	}

	return len(parts)%2 == 0
}

type (
	FirestoreQueryCollectionBuilder struct {
		c     *firestore.Client
		query firestore.Query
	}
	FirestoreQueryDocumentBuilder struct {
		c      *firestore.Client
		docRef *firestore.DocumentRef
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

	FirestoreCount int
)

func NewQueryCollectionBuilder(client *firestore.Client) *FirestoreQueryCollectionBuilder {
	return &FirestoreQueryCollectionBuilder{
		c: client,
	}
}

func (qb *FirestoreQueryCollectionBuilder) Collection(path string) *FirestoreQueryCollectionBuilder {
	qb.query = qb.c.Collection(path).Query
	return qb
}

func (qb *FirestoreQueryCollectionBuilder) WithWheres(wheres []*FirestoreWhere) *FirestoreQueryCollectionBuilder {
	for _, where := range wheres {
		qb.query = qb.query.Where(where.Path, where.GetOperator(), where.GetValue())
	}
	return qb
}

func (qb *FirestoreQueryCollectionBuilder) WithLimit(limit uint) *FirestoreQueryCollectionBuilder {
	if limit == 0 {
		return qb
	}

	qb.query = qb.query.Limit(int(limit))
	return qb
}

func (qb *FirestoreQueryCollectionBuilder) WithOrderBy(orderBy string, dir firestore.Direction) *FirestoreQueryCollectionBuilder {
	if orderBy == "" {
		return qb
	}

	qb.query = qb.query.OrderBy(orderBy, dir)
	return qb
}

func (qb *FirestoreQueryCollectionBuilder) Count(ctx context.Context) (FirestoreCount, error) {
	q := qb.query.NewAggregationQuery().WithCount("count")
	res, err := q.Get(ctx)
	if err != nil {
		return 0, fmt.Errorf("getting count query: %v", err)
	}

	count, ok := res["count"]
	if !ok {
		return 0, fmt.Errorf("extracting count from aggreagte result: %+v", res)
	}

	val, ok := count.(*firestorepb.Value)
	if !ok {
		return 0, fmt.Errorf("converting count to firestore value. received: %T", count)
	}

	return FirestoreCount(val.GetIntegerValue()), nil
}

func (qb *FirestoreQueryCollectionBuilder) GetAll(ctx context.Context) (FirestoreDocs, error) {
	docs := qb.query.Documents(ctx)
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

func NewQueryDocumentBuilder(client *firestore.Client) *FirestoreQueryDocumentBuilder {
	return &FirestoreQueryDocumentBuilder{
		c: client,
	}
}

func (b *FirestoreQueryDocumentBuilder) Document(path string) *FirestoreQueryDocumentBuilder {
	b.docRef = b.c.Doc(path)
	return b
}

func (b *FirestoreQueryDocumentBuilder) Get(ctx context.Context) (FirestoreDoc, error) {
	snapshot, err := b.docRef.Get(ctx)
	if status.Code(err) == codes.NotFound {
		return FirestoreDoc{}, fmt.Errorf("document %s does not exist", b.docRef.Path)
	} else if err != nil {
		slog.Error(err.Error())
		return FirestoreDoc{}, fmt.Errorf("failed to retrieve document %s", b.docRef.Path)
	}

	return FirestoreDoc{
		Path: snapshot.Ref.Path,
		Data: snapshot.Data(),
	}, nil
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

func (d FirestoreDoc) ToJSON() (string, error) {
	return utils.ToJSON(d.Data)
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

func (c FirestoreCount) ToJSON() (string, error) {
	return fmt.Sprintf("%d", c), nil
}
