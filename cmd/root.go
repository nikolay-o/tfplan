package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"tfplan/pkg/diff"
)

var rootCmd = &cobra.Command{
	Use:   "tfplan",
	Short: "Print diff for values.yaml",
	Long:  `Print diff for values.yaml`,
	Run: func(cmd *cobra.Command, args []string) {
		filePath, _ := cmd.Flags().GetString("file")
		fmt.Printf("DIFF: flag %v\n", filePath)
		diff.Diff(filePath)
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().StringP("file", "f", "", "Path to terraform plan file in JSON format")
}
