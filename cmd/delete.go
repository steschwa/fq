package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/carapace-sh/carapace"
	"github.com/spf13/cobra"
	"github.com/steschwa/fq/firestore"
	"github.com/steschwa/fq/firestore/parser"
)

var deleteCommand = &cobra.Command{
	Use:   "delete",
	Short: "delete firestore documents",
	RunE: func(*cobra.Command, []string) error {
		config, err := initDeleteConfig()
		if errors.Is(err, errNonEmulatorProjectID) {
			fmt.Println("only emulator projects are supported (projects starting with demo-*).")
			fmt.Println("see https://firebase.google.com/docs/emulator-suite/connect_firestore#choose_a_firebase_project")
			os.Exit(1)
			return nil
		}
		if err != nil {
			return err
		}

		client, err := firestore.NewClient(config.ProjectID)
		if err != nil {
			return fmt.Errorf("failed to create firestore client: %v", err)
		}
		defer client.Close()

		deleteClient := firestore.NewDeleteClient(client, config.Path)
		deleteClient.SetWheres(config.Wheres)
		err = deleteClient.Exec(firestore.DeleteOptions{
			ShowProgress: config.ShowProgress,
			Delay:        config.Delay,
		})
		if err != nil {
			return fmt.Errorf("deleting documents: %v", err)
		}

		return nil
	},
}

var (
	deleteWhere        []string
	deleteShowProgress bool
	deleteDelay        int
)

func init() {
	deleteCommand.Flags().StringArrayVarP(&deleteWhere, "where", "w", nil, "documents filter in format {KEY} {OPERATOR} {VALUE}. can be used multiple times")
	deleteCommand.Flags().BoolVar(&deleteShowProgress, "progress", false, "show the progress")
	deleteCommand.Flags().IntVar(&deleteDelay, "delay", 0, "delay between operations in milliseconds")

	addProjectFlag(deleteCommand)
	addPathFlag(deleteCommand)

	c := carapace.Gen(deleteCommand)
	c.Standalone()
}

type DeleteConfig struct {
	ProjectID    string
	Path         string
	Wheres       []firestore.Where
	ShowProgress bool
	Delay        int
}

func initDeleteConfig() (config DeleteConfig, err error) {
	if ProjectID == "" {
		return config, errEmptyProjectID
	}
	if !firestore.IsEmulatorProject(ProjectID) {
		return config, errNonEmulatorProjectID
	}
	config.ProjectID = ProjectID

	err = firestore.ValidatePath(Path)
	if err != nil {
		return config, fmt.Errorf("invalid firestore path")
	}
	config.Path = Path

	config.Wheres = make([]firestore.Where, len(deleteWhere))
	for i, wRaw := range deleteWhere {
		w, err := parser.Parse(wRaw)
		if err != nil {
			return config, fmt.Errorf("failed to parse firestore where: %s", err.Error())
		}

		config.Wheres[i] = w
	}

	config.ShowProgress = deleteShowProgress

	if setDelay < 0 {
		return config, errNegativeDelay
	}
	config.Delay = setDelay

	return config, nil
}
