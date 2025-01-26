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
	Use:   "fq",
	Short: "CLI tool to interact with Firestore",
	RunE: func(*cobra.Command, []string) error {
		if PrintVersion {
			fmt.Printf("version: %v\n", Version)
			fmt.Printf("revision: %v", Revision)
			return nil
		}

		return errors.New("please specify a subcommand to run")
	},
}

var (
	ProjectID    string
	Path         string
	PrintVersion bool
)

var (
	errEmptyProjectID = errors.New("empty project id")
)

func init() {
	rootCmd.AddCommand(queryCommand)
	rootCmd.AddCommand(setCommand)
	rootCmd.AddCommand(deleteCommand)

	rootCmd.PersistentFlags().BoolVarP(&PrintVersion, "version", "v", false, "print the version")

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
