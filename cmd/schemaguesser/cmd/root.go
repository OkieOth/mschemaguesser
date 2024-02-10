package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "schemaguesser",
	Short: "Tool to generate JSON schema from mongodb content",
	Long: `Based on a given mongodb connection you can evaluate its general content
                and generate JSON schemas out of it.
                This sould support model driven development and documentation activities`,
	Run: func(cmd *cobra.Command, args []string) {
		// Do Stuff Here
	},
}

func init() {
	rootCmd.AddCommand(schemaCmd)
	rootCmd.AddCommand(listCmd)
}

// Execute executes the root command.
func Execute() error {
	return rootCmd.Execute()
}
