package cmd

import (
	"github.com/spf13/cobra"
	"github.com/webws/embedding-knowledge-base/options"
	"github.com/webws/go-moda/logger"
)

var flagmsg = "msg" // cmd ask flag  question

func NewSearchCmd(configFlags options.ConfigFlags) *cobra.Command {
	var msg string
	searchCmd := &cobra.Command{
		Use:   "search",
		Short: "ask the knowledge base example: kbai ask --msg 'First, the chicken or the egg'",
		Long:  "ask the knowledge base example: kbai ask --msg 'First, the chicken or the egg'",
		RunE: func(cmd *cobra.Command, args []string) error {
			logger.Infow("searchCmd")
			return nil
		},
	}
	searchCmd.Flags().StringVar(&msg, flagmsg, "", "example: kbai search --msg 'First, the chicken or the egg'")
	searchCmd.MarkFlagRequired(flagmsg)

	return searchCmd
}
