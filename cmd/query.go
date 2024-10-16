package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/steschwa/fst/firestore"
	"github.com/steschwa/fst/firestore/parser"
)

var queryCommand = &cobra.Command{
	Use:   "query",
	Short: "query firestore",
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
			builder.SetWheres(config.Wheres).
				SetOrderBy(config.OrderBy, firestore.GetFirestoreDirection(config.OrderDescending)).
				SetLimit(config.Limit)

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

		} else if firestore.IsDocumentPath(config.Path) {
			loader := firestore.NewDocLoader(client, config.Path)

			doc, err := loader.GetDoc()
			if err != nil {
				fmt.Printf("loading document: %v\n", err)
				os.Exit(1)
			}

			j, err := json.Marshal(doc)
			if err != nil {
				fmt.Printf("marshalling document to json: %v\n", err)
				os.Exit(1)
			}

			fmt.Print(string(j))
		}
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
