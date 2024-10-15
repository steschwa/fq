package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "fst",
	Short: "CLI tool to interact with Firestore",
	Run:   func(*cobra.Command, []string) {},
}

var (
	ProjectID string
	Path      string
)

func init() {
	rootCmd.AddCommand(queryCommand)

	rootCmd.PersistentFlags().StringVar(&ProjectID, "project", "", "firebase project id")
	rootCmd.MarkPersistentFlagRequired("project")

	rootCmd.PersistentFlags().StringVar(&Path, "path", "", "collection or document path")
	rootCmd.MarkPersistentFlagRequired("path")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
