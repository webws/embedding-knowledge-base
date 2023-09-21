package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/webws/embedding-knowledge-base/ai"
	"github.com/webws/embedding-knowledge-base/options"
	"github.com/webws/embedding-knowledge-base/qdrant"
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
			// logger.Debugw("search", "config", configFlags)
			qdrantClient := qdrant.NewQdrantClient(configFlags.Qdrant, configFlags.Collection, configFlags.VectorSize)
			defer qdrantClient.Close()

			aiClient, err := ai.NewAiClient(configFlags.Proxy, configFlags.ApiKey)
			if err != nil {
				return err
			}
			vector, err := aiClient.SimpleGetVec(msg)
			if err != nil {
				return err
			}
			points, err := qdrantClient.Search(vector)
			if err != nil {
				logger.Errorw("qdrant search fail", "err", err)
				return err
			}
			if len(points) == 0 {
				fmt.Println("rearch term violation or exceeding category")
				return nil
				// return errors.New("rearch term violation or exceeding category")
			}
			// Score less than 0.8, rearch term violation or exceeding category
			if points[0].Score < 0.8 {
				fmt.Println("rearch term violation or exceeding category")
				return nil
				// return errors.New("rearch term violation or exceeding category")
			}

			answer := points[0].Payload["answers"].GetStringValue()
			fmt.Printf("The answer to the knowledge base:\n%s\n", answer)

			tmpl := "question:%s" + "reference answer: %s"
			finalPrompt := fmt.Sprintf(tmpl, points[0].Payload["question"].GetStringValue(), answer)

			chatgptAnswer, err := aiClient.Chat(finalPrompt)
			if err != nil {
				return err
			}
			fmt.Printf("Results of chatgpt answers  with reference answers:\n%s\n", chatgptAnswer)
			chatgptAnswer, err = aiClient.Chat(msg)
			if err != nil {
				return err
			}
			fmt.Printf("only chatgpt answers:\n%s\n", chatgptAnswer)

			return nil
		},
	}
	searchCmd.Flags().StringVar(&msg, flagmsg, "", "example: kbai search --msg 'First, the chicken or the egg'")
	searchCmd.MarkFlagRequired(flagmsg)
	return searchCmd
}
