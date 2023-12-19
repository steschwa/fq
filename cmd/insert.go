package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"os"

	"github.com/steschwa/fq/firebase"
	"github.com/urfave/cli/v2"
)

type (
	insertCmdParams struct {
		projectID  string
		collection string
		data       []firebase.FirestoreInsertData
	}
)

var (
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
			ib.InsertMany(params.data)
			slog.Info("inserting data ... done")

			return nil
		},
	}
)

func createInsertCmdParams(c *cli.Context) (insertCmdParams, error) {
	collection := c.String(firestoreCollectionName)
	projectID := c.String(firebaseProjectName)

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

	slog.Info(fmt.Sprintf("parsed %d items", len(insertData)))

	return insertCmdParams{
		projectID:  projectID,
		collection: collection,
		data:       insertData,
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
