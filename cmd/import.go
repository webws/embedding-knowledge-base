package cmd

import (
	"github.com/spf13/cobra"
	"github.com/webws/go-moda/logger"
)

var importCmd = &cobra.Command{
	Use:   "import",
	Short: "import data to vector database",
	Long:  "import data to vector database",

	RunE: func(cmd *cobra.Command, args []string) error {
		// TODO:flg to conf
		// dataFile, err := cmd.Flags().GetString(dataFileFlag)
		// if err != nil {
		// 	logger.Errorw("importCmd GetString fail", "err", err)
		// 	return err
		// }
		logger.Infow("import cmd")
		return nil
	},
}
