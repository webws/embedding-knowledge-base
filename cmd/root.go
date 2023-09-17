package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

const (
	DefaultQdrantAddr = "127.0.0.1:6334"
	DefaultSocksProxy = "socks5://127.0.0.1:1080"
)

var (
	qdrantFlag   = "qdrant"
	proxyFlag    = "proxy"
	dataFileFlag = "data-file" // cmd import flag
	msgFlag      = "msg"       // cmd ask flag
	rootCmd      = &cobra.Command{
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
)

func Execute() {
	fmt.Println(os.Args[1:])

	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(importCmd)
	rootCmd.AddCommand(ackCmd)

	// local flag
	importCmd.Flags().String(dataFileFlag, "", "import data-file")
	importCmd.MarkFlagRequired(dataFileFlag)
	ackCmd.Flags().String(msgFlag, "", "example: kbai ask --msg 'First, the chicken or the egg'")
	ackCmd.MarkFlagRequired(msgFlag)
	// PersistentFlags
	rootCmd.PersistentFlags().String(qdrantFlag, DefaultQdrantAddr, fmt.Sprintf("qdrant address default: %s", DefaultQdrantAddr))
	rootCmd.PersistentFlags().String(proxyFlag, DefaultSocksProxy, fmt.Sprintf("http client proxy default:%s ", DefaultSocksProxy))
}
