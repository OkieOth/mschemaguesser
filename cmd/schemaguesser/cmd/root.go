package cmd

import (
	"github.com/spf13/cobra"

	"okieoth/schemaguesser/internal/pkg/mongoHelper"
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
	rootCmd.AddCommand(versionCmd)

	rootCmd.PersistentFlags().StringVar(&mongoHelper.ConStr, "con_str", "mongodb://{MONGO_USER}:{MONGO_PASSWORD}@{MONGO_HOST}:{MONGO_PORT}/admin", "Connection string to mongodb")
}

// Execute executes the root command.
func Execute() error {
	return rootCmd.Execute()
}
