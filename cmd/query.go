package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/steschwa/fst/firestore"
	"github.com/steschwa/fst/firestore/parser"
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

		fmt.Printf("ProjectID: %s\n", config.ProjectID)
		fmt.Printf("Path: %s\n", config.Path)
		fmt.Printf("Count: %t\n", config.Count)
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
