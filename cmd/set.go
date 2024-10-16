package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"
	"github.com/steschwa/fst/firestore"
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

		if firestore.IsDocumentPath(config.Path) {
			setClient := firestore.NewSetClient(client, config.Path)
			err := setClient.Set(config.Data, firestore.SetOptions{
				ReplaceDoc: config.ReplaceDoc,
			})
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
	setCommand.Flags().StringVar(&dataPath, "data", "--", "input data. can be -- to read from stdin")
	setCommand.Flags().BoolVar(&replaceDoc, "replace", false, "replace documents instead of merging")
}

type SetConfig struct {
	ProjectID  string
	Path       string
	ReplaceDoc bool
	Data       map[string]any
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

	var r io.Reader
	if dataPath == "" || dataPath == "--" {
		if fi, err := os.Stdin.Stat(); err == nil {
			if fi.Size() <= 0 {
				return config, fmt.Errorf("no data from stdin")
			}
		}

		r = os.Stdin
	} else {
		r, err = os.Open(dataPath)
		if errors.Is(err, os.ErrNotExist) {
			return config, fmt.Errorf("file %s does not exist", dataPath)
		}
		if err != nil {
			return config, fmt.Errorf("file %s can't be opened for reading", dataPath)
		}
	}

	if err := json.NewDecoder(r).Decode(&config.Data); err != nil {
		return config, fmt.Errorf("failed to decode json from %s", dataPath)
	}

	return config, nil
}
