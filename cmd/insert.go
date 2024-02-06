package cmd

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/steschwa/fq/firebase"
	"github.com/urfave/cli/v2"
)

type (
	insertCmdParams struct {
		projectID  string
		collection string
		data       []firebase.FirestoreInsertData
		timeout    uint
	}
)

const (
	defaultInsertCmdTimeout uint = 30

	firestoreCollectionFlagName = "collection"
)

var (
	firestoreCollectionPathFlag = &cli.StringFlag{
		Name:    firestoreCollectionFlagName,
		Aliases: []string{"c"},
		Usage:   "`path` to firestore collection separated with dashes (/)",
		Action: func(ctx *cli.Context, s string) error {
			pathType, err := firebase.GetFirestorePathType(s)
			if err != nil {
				return err
			}

			if pathType == firebase.FirestorePathTypeCollection {
				return nil
			}
			return errors.New("only collection paths (containing an uneven amount of parts separated by /) are allowed")
		},
	}

	InsertCmd = &cli.Command{
		Name:  "insert",
		Usage: "insert data into firestore",
		Description: `this command expects json formatted data.
it should consist of an top-level array that contains objects.
if an item contains an property 'id' this is used as document path.
if there is no 'id' property, then an auto generated id is used`,
		Flags: []cli.Flag{
			firebaseProjectIDFlag,
			firestoreCollectionPathFlag,
			&cli.StringFlag{
				Name:    "file",
				Aliases: []string{"f"},
				Usage:   "read insert data from this file (leave empty to read data from stdin)",
				Action: func(ctx *cli.Context, s string) error {
					if _, err := os.Stat(s); errors.Is(err, os.ErrNotExist) {
						return fmt.Errorf("file %s does not exist", s)
					} else if err != nil {
						return err
					}
					return nil
				},
			},
			&cli.UintFlag{
				Name:        "timeout",
				Usage:       "timeout in `seconds`",
				DefaultText: fmt.Sprintf("%d s", defaultInsertCmdTimeout),
			},
		},
		Action: func(c *cli.Context) error {
			params, err := createInsertCmdParams(c)
			if err != nil {
				return err
			}

			slog.Info(fmt.Sprintf("using project %s", params.projectID))
			slog.Info(fmt.Sprintf("inserting into %s", params.collection))

			client, err := firebase.InitFirestoreClient(params.projectID)
			if err != nil {
				return err
			}
			defer client.Close()

			if firebase.IsRunningInEmulator() {
				slog.Info("disabling function triggers ...")
				err = firebase.EmulatorDisableFunctionTriggers()
				if err != nil {
					slog.Error(err.Error())
					slog.Error("disabling function triggers ... failed")
				} else {
					slog.Info("disabling function triggers ... done")
				}

				defer func() {
					slog.Info("enabling function triggers ...")
					err := firebase.EmulatorEnableFunctionTriggers()
					if err != nil {
						slog.Error(err.Error())
						slog.Error("enabling function triggers ... failed")
					} else {
						slog.Info("enabling function triggers ... done")
					}
				}()
			}

			ib := firebase.NewInsertBuilder(client).
				Collection(params.collection)

			slog.Info("inserting data ...")
			ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(params.timeout))
			defer cancel()
			errs := ib.InsertMany(ctx, params.data)
			errs.Log()
			slog.Info("inserting data ... done")

			return nil
		},
	}
)

func createInsertCmdParams(c *cli.Context) (insertCmdParams, error) {
	collection := c.String(firestoreCollectionFlagName)
	projectID := c.String(firebaseProjectFlagName)

	filePath := c.String("file")
	var insertData []firebase.FirestoreInsertData
	var insertDataErr error
	if filePath == "" {
		slog.Info("reading data from stdin")
		insertData, insertDataErr = loadInsertData(os.Stdin)
	} else {
		slog.Info(fmt.Sprintf("reading data from %s", filePath))
		file, err := os.Open(filePath)
		if err != nil {
			slog.Error(err.Error())
			return insertCmdParams{}, fmt.Errorf("failed to open %s for reading insert data", filePath)
		}
		defer file.Close()
		insertData, insertDataErr = loadInsertData(file)
	}

	if insertDataErr != nil {
		slog.Error(insertDataErr.Error())
		return insertCmdParams{}, errors.New("failed to load insert data")
	}

	timeout := c.Uint("timeout")
	if timeout == 0 {
		timeout = defaultInsertCmdTimeout
	}

	slog.Info(fmt.Sprintf("parsed %d items", len(insertData)))

	return insertCmdParams{
		projectID:  projectID,
		collection: collection,
		data:       insertData,
		timeout:    timeout,
	}, nil
}

func loadInsertData(file *os.File) ([]firebase.FirestoreInsertData, error) {
	fi, err := file.Stat()
	if err != nil {
		slog.Error(err.Error())
		return nil, errors.New("failed to get file stat")
	}

	if fi.Size() <= 0 {
		return nil, errors.New("no data")
	}

	var data []firebase.FirestoreInsertData

	if err := json.NewDecoder(file).Decode(&data); err != nil {
		return nil, err
	}
	return data, nil
}
