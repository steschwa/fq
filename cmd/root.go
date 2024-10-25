package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/carapace-sh/carapace"
	"github.com/spf13/cobra"
	"github.com/steschwa/fq/completion"
)

var rootCmd = &cobra.Command{
	Use:   "fq",
	Short: "CLI tool to interact with Firestore",
	Run: func(*cobra.Command, []string) {
		fmt.Println("please specify a subcommand to run")
		os.Exit(1)
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

	rootCmd.PersistentFlags().StringVar(&ProjectID, "project", "", "firebase project id")
	rootCmd.MarkPersistentFlagRequired("project")

	rootCmd.PersistentFlags().StringVar(&Path, "path", "", "collection or document path")
	rootCmd.MarkPersistentFlagRequired("path")

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
