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
	Version   = "0.0.1"
	CommitSHA = "dev"
)

var rootCmd = &cobra.Command{
	Use:   "fq",
	Short: "CLI tool to interact with Firestore",
	RunE: func(*cobra.Command, []string) error {
		if PrintVersion {
			fmt.Printf("version: %v\n", Version)
			fmt.Printf("commit: %v", CommitSHA)
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

	rootCmd.PersistentFlags().StringVar(&ProjectID, "project", "", "firebase project id")
	rootCmd.MarkPersistentFlagRequired("project")

	rootCmd.PersistentFlags().StringVar(&Path, "path", "", "collection or document path")
	rootCmd.MarkPersistentFlagRequired("path")

	rootCmd.PersistentFlags().BoolVarP(&PrintVersion, "version", "v", false, "print the version")

	c := carapace.Gen(rootCmd)
	c.Standalone()
	c.FlagCompletion(carapace.ActionMap{
		"project": completion.ActionGCloudProjects(),
	})
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
