package cmd

import (
	"github.com/spf13/cobra"
	"github.com/webws/go-moda/logger"
)

var ackCmd = &cobra.Command{
	Use:   "ask",
	Short: "ask the knowledge base example: kbai ask --msg 'First, the chicken or the egg'",
	Long:  "ask the knowledge base example: kbai ask --msg 'First, the chicken or the egg'",
	RunE: func(cmd *cobra.Command, args []string) error {
		logger.Infow("ackCmd")
		return nil
	},
}
