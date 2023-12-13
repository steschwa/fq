package cmd

import (
	"errors"
	"fmt"

	"cloud.google.com/go/firestore"
	"github.com/steschwa/fq/firebase"
	"github.com/urfave/cli/v2"
)

type (
	queryCmdParams struct {
		projectID string
		path      string
		where     []*firebase.FirestoreWhere
		limit     int
		orderBy   string
		orderDir  firestore.Direction
	}
)

var (
	QueryCmd = &cli.Command{
		Name:  "query",
		Usage: "load data from firestore",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "project",
				Aliases:  []string{"p"},
				Usage:    "gcloud project id",
				EnvVars:  []string{"GCLOUD_PROJECT", "GCLOUD_PROJECT_ID"},
				Required: true,
			},
			&cli.StringSliceFlag{
				Name:    "where",
				Aliases: []string{"w"},
				Usage:   "documents filter. must be in format '{property-path} {operator} {value}'. can be used multiple times",
			},
			&cli.IntFlag{
				Name:        "limit",
				DefaultText: "no limit",
			},
			&cli.StringFlag{
				Name:        "orderby",
				DefaultText: "no ordering",
			},
			&cli.BoolFlag{
				Name:        "desc",
				Usage:       "order in descending order",
				DefaultText: "false - ascending",
			},
		},
		Action: func(c *cli.Context) error {
			params, err := createQueryCmdParams(c)
			if err != nil {
				return err
			}

			client, err := firebase.InitFirestoreClient(params.projectID)
			if err != nil {
				return err
			}
			defer client.Close()

			qb := firebase.NewQueryBuilder(client)
			docs, err := qb.Collection(params.path).
				WithWheres(params.where).
				WithLimit(params.limit).
				WithOrderBy(params.orderBy, params.orderDir).
				Execute()
			if err != nil {
				return err
			}

			j, err := docs.ToJSON()
			if err != nil {
				return err
			}

			fmt.Println(j)
			return nil
		},
	}
)

func createQueryCmdParams(c *cli.Context) (queryCmdParams, error) {
	path := c.Args().Get(0)
	if err := firebase.ValidateFirestoreCollectionPath(path); err != nil {
		return queryCmdParams{}, err
	}

	projectID := c.String("project")
	if projectID == "" {
		return queryCmdParams{}, errors.New("project is required")
	}

	wheresRaw := c.StringSlice("where")
	wheres := make([]*firebase.FirestoreWhere, len(wheresRaw))
	for i, where := range wheresRaw {
		w, err := firebase.ParseFirestoreWhere(where)
		if err != nil {
			return queryCmdParams{}, fmt.Errorf("invalid where clause %s", where)
		}

		wheres[i] = w
	}

	limit := c.Int("limit")
	order := c.String("orderby")

	var orderDir firestore.Direction
	if desc := c.Bool("desc"); desc {
		orderDir = firestore.Desc
	} else {
		orderDir = firestore.Asc
	}

	return queryCmdParams{
		projectID: projectID,
		path:      path,
		where:     wheres,
		limit:     limit,
		orderBy:   order,
		orderDir:  orderDir,
	}, nil
}
