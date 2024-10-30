package cmd

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/carapace-sh/carapace"
	"github.com/spf13/cobra"
	"github.com/steschwa/fq/firestore"
	"github.com/steschwa/fq/firestore/parser"
)

var queryCommand = &cobra.Command{
	Use:   "query",
	Short: "query firestore",
	RunE: func(*cobra.Command, []string) error {
		config, err := initQueryConfig()
		if err != nil {
			return err
		}

		client, err := firestore.NewClient(config.ProjectID)
		if err != nil {
			return fmt.Errorf("creating firestore client: %v", err)
		}
		defer client.Close()

		if firestore.IsCollectionPath(config.Path) {
			queryClient := firestore.NewQueryClient(client, config.Path)
			queryClient.SetWheres(config.Wheres).
				SetOrderBy(config.OrderBy, firestore.GetFirestoreDirection(config.OrderDescending)).
				SetLimit(config.Limit)

			if config.Count {
				count, err := queryClient.GetCount()
				if err != nil {
					return fmt.Errorf("loading documents count: %v", err)
				}

				fmt.Print(count)
			} else {
				docs, err := queryClient.GetDocs()
				if err != nil {
					return fmt.Errorf("loading documents: %v", err)
				}

				if docs == nil {
					docs = []any{}
				}

				j, err := json.Marshal(docs)
				if err != nil {
					return fmt.Errorf("marshalling documents to json: %v", err)
				}

				fmt.Print(string(j))
			}

		} else if firestore.IsDocumentPath(config.Path) {
			docClient := firestore.NewDocClient(client, config.Path)

			doc, err := docClient.GetDoc()
			if errors.Is(err, firestore.ErrDocumentNotFound) {
				fmt.Print("null")
				return nil
			}
			if err != nil {
				return fmt.Errorf("loading document: %v", err)
			}

			j, err := json.Marshal(doc)
			if err != nil {
				return fmt.Errorf("marshalling document to json: %v", err)
			}

			fmt.Print(string(j))
		}

		return nil
	},
}

var (
	count   bool
	where   []string
	orderBy string
	desc    bool
	limit   int
)

func init() {
	queryCommand.Flags().BoolVar(&count, "count", false, "count documents instead of returning json")
	queryCommand.Flags().StringArrayVarP(&where, "where", "w", nil, "documents filter in format {KEY} {OPERATOR} {VALUE}. can be used multiple times")
	queryCommand.Flags().StringVar(&orderBy, "order-by", "", "set column to order by")
	queryCommand.Flags().BoolVar(&desc, "desc", false, "order documents in descending order (only used if --order-by is set)")
	queryCommand.Flags().IntVar(&limit, "limit", -1, "limit number of returned documents")

	addProjectFlag(queryCommand)
	addPathFlag(queryCommand)

	c := carapace.Gen(queryCommand)
	c.Standalone()
}

type QueryConfig struct {
	ProjectID       string
	Path            string
	Count           bool
	Wheres          []firestore.Where
	OrderBy         string
	OrderDescending bool
	Limit           int
}

func initQueryConfig() (config QueryConfig, err error) {
	if ProjectID == "" {
		return config, errEmptyProjectID
	}
	config.ProjectID = ProjectID

	err = firestore.ValidatePath(Path)
	if err != nil {
		return config, fmt.Errorf("invalid firestore path")
	}
	config.Path = Path

	config.Wheres = make([]firestore.Where, len(where))
	for i, wRaw := range where {
		w, err := parser.Parse(wRaw)
		if err != nil {
			return config, fmt.Errorf("failed to parse firestore where: %s", err.Error())
		}

		config.Wheres[i] = w
	}

	config.Count = count
	config.OrderBy = orderBy
	config.OrderDescending = desc
	config.Limit = limit

	return config, nil
}

func (c QueryConfig) DebugPrint() {
	fmt.Printf("ProjectID: %s\n", c.ProjectID)
	fmt.Printf("Path: %s\n", c.Path)
	fmt.Printf("Count: %t\n", c.Count)
	for i, w := range c.Wheres {
		fmt.Printf("Where (%d): %s\n", i+1, w.String())
	}
	fmt.Printf("Order-By: %s\n", c.OrderBy)
	fmt.Printf("Order Descending: %t\n", c.OrderDescending)
	fmt.Printf("Limit: %d\n", c.Limit)
}
