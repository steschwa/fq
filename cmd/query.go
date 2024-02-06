package cmd

import (
	"context"
	"errors"
	"fmt"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/steschwa/fq/firebase"
	"github.com/steschwa/fq/utils"
	"github.com/urfave/cli/v2"
)

type (
	queryCmdParams struct {
		projectID string
		path      string
		where     []*firebase.FirestoreWhere
		limit     uint
		orderBy   string
		orderDir  firestore.Direction
		timeout   uint
	}
)

const (
	firebaseProjectFlagName = "project"
	firestorePathFlagName   = "path"

	defaultQueryCmdTimeout uint = 30
)

var (
	firebaseProjectIDFlag = &cli.StringFlag{
		Name:     firebaseProjectFlagName,
		Aliases:  []string{"p"},
		Usage:    "gcloud project id",
		EnvVars:  []string{"GCLOUD_PROJECT", "GCLOUD_PROJECT_ID"},
		Required: true,
	}

	firestorePathFlag = &cli.StringFlag{
		Name:  firestorePathFlagName,
		Usage: "`path` to firestore collection or document separated with dashes (/)",
		Action: func(ctx *cli.Context, s string) error {
			pathType, err := firebase.GetFirestorePathType(s)
			if err != nil {
				return err
			}

			if pathType == firebase.FirestorePathTypeCollection {
				return nil
			}
			if pathType == firebase.FirestorePathTypeDocument {
				return nil
			}
			return errors.New("could not validate path. please provide a) a collection path (containing an uneven amount of parts separated by /) or b) document path (containing an even amount of parts separated by /)")
		},
	}

	firestoreWhereFlag = &cli.StringSliceFlag{
		Name:    "where",
		Aliases: []string{"w"},
		Usage:   "documents `filter`. must be in format '{property-path} {operator} {value}'. can be used multiple times",
		Action: func(ctx *cli.Context, s []string) error {
			for _, where := range s {
				err := firebase.ValidateFirestoreWhere(where)
				if err != nil {
					return fmt.Errorf("invalid where clause (%s)", where)
				}
			}

			return nil
		},
	}

	firestoreLimitFlag = &cli.UintFlag{
		Name:        "limit",
		DefaultText: "no limit",
	}

	QueryCmd = &cli.Command{
		Name:  "query",
		Usage: "load data from firestore",
		Flags: []cli.Flag{
			firebaseProjectIDFlag,
			firestorePathFlag,
			firestoreWhereFlag,
			firestoreLimitFlag,
			&cli.StringFlag{
				Name:        "orderby",
				DefaultText: "no ordering",
			},
			&cli.BoolFlag{
				Name:        "desc",
				Usage:       "order in descending order",
				DefaultText: "false - ascending",
			},
			&cli.UintFlag{
				Name:        "timeout",
				Usage:       "timeout in `seconds`",
				DefaultText: fmt.Sprintf("%d s", defaultQueryCmdTimeout),
			},
		},
		Action: func(c *cli.Context) error {
			params := createQueryCmdParams(c)

			client, err := firebase.InitFirestoreClient(params.projectID)
			if err != nil {
				return err
			}
			defer client.Close()

			ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(params.timeout))
			defer cancel()

			var serializable utils.JSONSerializable

			if firebase.IsFirestoreCollectionPath(params.path) {
				qb := firebase.NewQueryCollectionBuilder(client)
				serializable, err = qb.Collection(params.path).
					WithWheres(params.where).
					WithLimit(params.limit).
					WithOrderBy(params.orderBy, params.orderDir).
					GetAll(ctx)
			} else if firebase.IsFirestoreDocumentPath(params.path) {
				qb := firebase.NewQueryDocumentBuilder(client)
				serializable, err = qb.Document(params.path).Get(ctx)
			}

			if err != nil {
				return err
			}

			j, err := serializable.ToJSON()
			if err != nil {
				return err
			}

			fmt.Println(j)
			return nil
		},
	}
)

func createQueryCmdParams(c *cli.Context) queryCmdParams {
	projectID := c.String(firebaseProjectFlagName)
	path := c.String(firestorePathFlagName)

	wheresRaw := c.StringSlice("where")
	wheres := make([]*firebase.FirestoreWhere, len(wheresRaw))
	for i, where := range wheresRaw {
		if w, err := firebase.ParseFirestoreWhere(where); err == nil {
			wheres[i] = w
		}
	}

	limit := c.Uint("limit")
	order := c.String("orderby")

	var orderDir firestore.Direction
	if desc := c.Bool("desc"); desc {
		orderDir = firestore.Desc
	} else {
		orderDir = firestore.Asc
	}

	timeout := c.Uint("timeout")
	if timeout == 0 {
		timeout = defaultQueryCmdTimeout
	}

	return queryCmdParams{
		projectID: projectID,
		path:      path,
		where:     wheres,
		limit:     limit,
		orderBy:   order,
		orderDir:  orderDir,
		timeout:   timeout,
	}
}
