package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/carapace-sh/carapace"
	"github.com/spf13/cobra"
	"github.com/steschwa/fq/completion"
)

var (
	Version  = "0.0.1"
	Revision = "dev"
)

var rootCmd = &cobra.Command{
	Use:     "fq",
	Short:   "CLI tool to interact with Firestore",
	Version: fmt.Sprintf("%s - revision %s", Version, Revision),
	RunE: func(*cobra.Command, []string) error {
		return errors.New("please specify a subcommand to run")
	},
}

var (
	ProjectID string
	Path      string
)

var (
	errEmptyProjectID = errors.New("empty project id")
)

func init() {
	rootCmd.AddCommand(queryCommand)
	rootCmd.AddCommand(setCommand)
	rootCmd.AddCommand(deleteCommand)

	carapace.Gen(rootCmd).Standalone()
}

func addProjectFlag(cmd *cobra.Command) {
	cmd.Flags().StringVar(&ProjectID, "project", "", "firebase project id")
	cmd.MarkFlagRequired("project")

	carapace.Gen(cmd).FlagCompletion(carapace.ActionMap{
		"project": completion.ActionGCloudProjects(),
	})
}

func addPathFlag(cmd *cobra.Command) {
	cmd.Flags().StringVar(&Path, "path", "", "collection or document path")
	cmd.MarkFlagRequired("path")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
