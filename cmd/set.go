package cmd

import (
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

		fmt.Println(string(config.Data))
	},
}

var (
	dataPath string
)

func init() {
	setCommand.Flags().StringVar(&dataPath, "data", "--", "input data. can be -- to read from stdin")
}

type SetConfig struct {
	ProjectID string
	Path      string
	Data      []byte
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

	// TODO: investige if maye reading json object-by-object is possible
	config.Data, err = io.ReadAll(r)
	if err != nil {
		return config, fmt.Errorf("failed to read from %s", dataPath)
	}

	return config, nil
}
