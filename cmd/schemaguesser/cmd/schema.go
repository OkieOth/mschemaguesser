package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var schemaCmd = &cobra.Command{
	Use:   "schema",
	Short: "functions around the schemas",
	Long:  "With this command you can create schemas out of mongodb collection",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("To learn about the possible options type:\n`schemaguesser schema --help`")
	},
}
