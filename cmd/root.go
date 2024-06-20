package cmd

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/carapace-sh/carapace"
	"github.com/spf13/cobra"
	"github.com/steschwa/fq/cmp/gcloud/projects"
	"github.com/steschwa/fq/firebase"
	"github.com/steschwa/fq/utils"
)

var (
	rootCmd = &cobra.Command{
		Use:   "fq",
		Short: "firestore query tool",
		Args:  cobra.NoArgs,
		Run:   run,
	}

	projectID         string
	path              string
	where             []string
	limit             uint
	orderBy           string
	orderByDescending bool
	timeout           uint
	count             bool
)

const (
	defaultTimeout uint = 30
)

func init() {
	carapace.Gen(rootCmd).Standalone()

	rootCmd.Flags().StringVar(&projectID, "project", "", "gcloud project id")
	rootCmd.Flags().StringVar(&path, "path", "", "path to firestore collection or document separated with dashes (/)")
	rootCmd.Flags().StringArrayVar(&where, "where", nil, "format: '{property-path} {operator} {value}'")
	rootCmd.Flags().UintVar(&limit, "limit", 0, "maximum documents to return")
	rootCmd.Flags().StringVar(&orderBy, "orderby", "", "order-by column")
	rootCmd.Flags().BoolVar(&orderByDescending, "desc", false, "reverse sort direction")
	rootCmd.Flags().UintVar(&timeout, "timeout", defaultTimeout, "timeout in seconds")
	rootCmd.Flags().BoolVar(&count, "count", false, "return count instead of documents")

	rootCmd.MarkFlagRequired("project")
	rootCmd.MarkFlagRequired("path")

	carapace.Gen(rootCmd).FlagCompletion(carapace.ActionMap{
		"project": projects.ActionProjects(),
	})
}

func validateFlags() error {
	if err := validatePath(); err != nil {
		return err
	}

	return validateWhere()
}

func run(*cobra.Command, []string) {
	if err := validateFlags(); err != nil {
		fmt.Printf("invalid flags: %v\n", err)
		os.Exit(1)
	}

	client, err := firebase.InitFirestoreClient(projectID)
	if err != nil {
		fmt.Printf("failed to create firestore client: %v\n", err)
		os.Exit(1)
	}

	defer client.Close()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(timeout))
	defer cancel()

	var serializable utils.JSONSerializable

	if firebase.IsFirestoreCollectionPath(path) {
		queryBuilder := firebase.NewQueryCollectionBuilder(client)

		var orderDirection firestore.Direction
		if orderByDescending {
			orderDirection = firestore.Desc
		} else {
			orderDirection = firestore.Asc
		}

		firestoreWheres := make([]*firebase.FirestoreWhere, len(where))
		for i, where := range where {
			if w, err := firebase.ParseFirestoreWhere(where); err == nil {
				firestoreWheres[i] = w
			}
		}

		queryBuilder.Collection(path).
			WithWheres(firestoreWheres).
			WithLimit(limit).
			WithOrderBy(orderBy, orderDirection)

		if count {
			serializable, err = queryBuilder.Count(ctx)
		} else {
			serializable, err = queryBuilder.GetAll(ctx)
		}

	} else if firebase.IsFirestoreDocumentPath(path) {
		queryBuilder := firebase.NewQueryDocumentBuilder(client)
		serializable, err = queryBuilder.Document(path).Get(ctx)
	}

	if err != nil {
		fmt.Printf("firestore query failed: %v\n", err)
		os.Exit(1)
	}

	j, err := serializable.ToJSON()
	if err != nil {
		fmt.Printf("failed to serialize data: %v\n", err)
		os.Exit(1)
	}

	fmt.Print(j)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func validatePath() error {
	pathType, err := firebase.GetFirestorePathType(path)
	if err != nil {
		return err
	}

	if pathType == firebase.FirestorePathTypeCollection {
		return nil
	}
	if pathType == firebase.FirestorePathTypeDocument {
		return nil
	}

	return errors.New("invalid path. please provide\na) a collection path (containing an uneven amount of parts separated by /)\nor b) document path (containing an even amount of parts separated by /)")
}

func validateWhere() error {
	for _, where := range where {
		err := firebase.ValidateFirestoreWhere(where)
		if err != nil {
			return fmt.Errorf("invalid where clause: %s", where)
		}
	}

	return nil
}
