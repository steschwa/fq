package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/steschwa/fst/firestore"
	"github.com/steschwa/fst/firestore/parser"
)

const (
	timeoutCreateFirestoreClient = 10
	timeoutQueryDocuments        = 30
)

var queryCommand = &cobra.Command{
	Use:   "query",
	Short: "Query Firestore",
	Run: func(*cobra.Command, []string) {
		config, err := initQueryConfig()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		client, err := firestore.NewClient(config.ProjectID)
		if err != nil {
			fmt.Printf("failed to create firestore client: %v\n", err)
			os.Exit(1)
		}
		defer client.Close()

		if firestore.IsCollectionPath(config.Path) {
			builder := firestore.NewQueryBuilder(client, config.Path)
			builder.SetWheres(config.Wheres)

			if config.Count {
				count, err := builder.GetCount()
				if err != nil {
					fmt.Printf("loading documents count: %v\n", err)
					os.Exit(1)
				}

				fmt.Print(count)
			} else {
				docs, err := builder.GetDocs()
				if err != nil {
					fmt.Printf("loading documents: %v\n", err)
					os.Exit(1)
				}

				j, err := json.Marshal(docs)
				if err != nil {
					fmt.Printf("marshalling documents to json: %v\n", err)
					os.Exit(1)
				}

				fmt.Print(string(j))
			}

		}
	},
}

var (
	count bool
	where []string
)

func init() {
	queryCommand.Flags().BoolVar(&count, "count", false, "count documents instead of returning json")
	queryCommand.Flags().StringArrayVarP(&where, "where", "w", nil, "documents filter in format {KEY} {OPERATOR} {VALUE}. can be used multiple times")
}

type (
	QueryConfig struct {
		ProjectID string
		Path      string
		Count     bool
		Wheres    []firestore.Where
	}
)

var (
	errEmptyProjectID = errors.New("empty project id")
	errEmptyPath      = errors.New("empty path")
)

func initQueryConfig() (QueryConfig, error) {
	if ProjectID == "" {
		return QueryConfig{}, errEmptyProjectID
	}

	if Path == "" {
		return QueryConfig{}, errEmptyPath
	}
	err := firestore.ValidatePath(Path)
	if err != nil {
		return QueryConfig{}, fmt.Errorf("invalid firestore path")
	}

	wheres := make([]firestore.Where, len(where))
	for i, wRaw := range where {
		w, err := parser.Parse(wRaw)
		if err != nil {
			return QueryConfig{}, fmt.Errorf("failed to parse firestore where: %s", err.Error())
		}

		wheres[i] = w
	}

	return QueryConfig{
		ProjectID: ProjectID,
		Path:      Path,
		Count:     count,
		Wheres:    wheres,
	}, nil
}

func (c QueryConfig) DebugPrint() {
	fmt.Printf("ProjectID: %s\n", c.ProjectID)
	fmt.Printf("Path: %s\n", c.Path)
	fmt.Printf("Count: %t\n", c.Count)
	for i, w := range c.Wheres {
		fmt.Printf("Where (%d): %s\n", i+1, w.String())
	}
}
