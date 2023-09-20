package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/webws/embedding-knowledge-base/options"
)

func Execute() {
	rootCmd := &cobra.Command{
		Use:   "kbai",
		Short: "a local knowledge base, based on chatgpt and qdrant",
		Long:  "a local knowledge base, based on chatgpt and qdrant",
		// Long:  `Built a local smart search knowledge repository using Golang, OpenAI Embedding, and the qdrant vector database`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cobra.MinimumNArgs(1)(cmd, args); err != nil {
				return cmd.Help()
			}
			return nil
		},
	}
	configFlags := options.NewConfigFlags()
	configFlags.AddFlags(rootCmd.PersistentFlags())

	rootCmd.AddCommand(NewImportCmd(*configFlags))
	rootCmd.AddCommand(NewSearchCmd(*configFlags))
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
