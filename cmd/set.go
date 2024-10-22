package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"
	"github.com/steschwa/fst/firestore"
	"github.com/steschwa/fst/utils"
)

var setCommand = &cobra.Command{
	Use:   "set",
	Short: "insert / update firestore documents",
	Run: func(*cobra.Command, []string) {
		config, err := initSetConfig()
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

		setClient := firestore.NewSetClient(client, config.Path)
		options := firestore.SetOptions{
			ReplaceDocument: config.ReplaceDoc,
		}

		if firestore.IsCollectionPath(config.Path) {
			err := setClient.SetMany(config.CollectionData, options)
			if err != nil {
				fmt.Printf("failed to set documents: %v\n", err)
				os.Exit(1)
			}
		} else if firestore.IsDocumentPath(config.Path) {
			err := setClient.Set(config.DocumentData, options)
			if err != nil {
				fmt.Printf("failed to set document: %v\n", err)
				os.Exit(1)
			}
		}
	},
}

var (
	dataPath   string
	replaceDoc bool
)

func init() {
	setCommand.Flags().StringVar(&dataPath, "data", "--", "input data json file. can be -- to read from stdin")
	setCommand.Flags().BoolVar(&replaceDoc, "replace", false, "replace documents instead of merging")
}

type SetConfig struct {
	ProjectID      string
	Path           string
	ReplaceDoc     bool
	DocumentData   firestore.JSONObject
	CollectionData firestore.JSONArray
}

func initSetConfig() (config SetConfig, err error) {
	if ProjectID == "" {
		return config, errEmptyProjectID
	}
	config.ProjectID = ProjectID

	err = firestore.ValidatePath(Path)
	if err != nil {
		return config, fmt.Errorf("invalid firestore path")
	}
	config.Path = Path
	config.ReplaceDoc = replaceDoc

	var (
		r            io.Reader
		dataPathName string
	)
	if dataPath == "" || dataPath == "--" {
		if utils.IsStdinEmpty() {
			return config, fmt.Errorf("no data from stdin")
		}

		dataPathName = "stdin"
		r = os.Stdin
	} else {
		dataPathName = dataPath

		r, err = os.Open(dataPath)
		if errors.Is(err, os.ErrNotExist) {
			return config, fmt.Errorf("file %s does not exist", dataPath)
		}
		if err != nil {
			return config, fmt.Errorf("file %s can't be opened for reading", dataPath)
		}
	}

	if firestore.IsDocumentPath(config.Path) {
		if err := json.NewDecoder(r).Decode(&config.DocumentData); err != nil {
			return config, fmt.Errorf("failed to decode json from %s: %v", dataPathName, err)
		}
	} else if firestore.IsCollectionPath(config.Path) {
		if err := json.NewDecoder(r).Decode(&config.CollectionData); err != nil {
			return config, fmt.Errorf("failed to decode json from %s: %v", dataPathName, err)
		}
	}

	return config, nil
}
