package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/spf13/cobra"
	"github.com/webws/embedding-knowledge-base/pkg/ai"
	"github.com/webws/embedding-knowledge-base/pkg/qdrant"
	"github.com/webws/go-moda/logger"

	pb "github.com/qdrant/go-client/qdrant"
)

type QA struct {
	Questions string `json:"questions" bson:"questions"`
	Answers   string `json:"answers" bson:"answers"`
}

var importCmd = &cobra.Command{
	Use:   "import",
	Short: "import data to vector database",
	Long:  "import data to vector database",
	RunE: func(cmd *cobra.Command, args []string) error {
		// TODO:flg to conf
		qdrantAddr, err := cmd.Flags().GetString(qdrantFlag)
		if err != nil {
			return err
		}
		qdrantClient := qdrant.NewQdrantClient(qdrantAddr)
		defer qdrantClient.Close()

		// TODO flg to conf apikey,proxy
		apiKey := os.Getenv(apiKeyFlag)
		if apiKey == "" {
			apiKey, err = cmd.Flags().GetString(apiKeyFlag)
			if err != nil {
				return err
			}
		}
		//
		proxy, err := cmd.Flags().GetString(proxyFlag)
		if err != nil {
			return err
		}
		aiClient, err := ai.NewAiClient(proxy, apiKey)
		if err != nil {
			return err
		}

		// TODO  collection,size use flg
		if err = qdrantClient.CreateCollection("kubernetes", 1536); err != nil {
			return err
		}

		// TODO flg json file to QA
		// TODO flg json string to QA
		dataFilePath, err := cmd.Flags().GetString(dataFileFlag)
		if err != nil {
			return err
		}
		qas, err := convertToQAs(dataFilePath)
		if err != nil {
			return err
		}
		points := []*pb.PointStruct{}
		qpsLenth := len(qas)
		for i, qa := range qas {
			embedding, err := aiClient.SimpleGetVec(qa.Questions)
			if err != nil {
				logger.Errorw("SimpleGetVec", "err", err, "question", qa.Questions, "index", i, "total", qpsLenth)
				return err
			}
			point := buildPoint(qa.Questions, qa.Answers, embedding)
			points = append(points, point)
		}
		return qdrantClient.CreatePoints("kubernetes", points)
	},
}

func buildPoint(question string, answers string, embedding []float32) *pb.PointStruct {
	point := &pb.PointStruct{}
	// point id
	uuid := fmt.Sprintf("%s-%d", question, time.Now().UnixNano())
	point.Id = &pb.PointId{
		PointIdOptions: &pb.PointId_Uuid{
			Uuid: uuid,
		},
	}

	// vector
	point.Vectors = &pb.Vectors{
		VectorsOptions: &pb.Vectors_Vector{
			Vector: &pb.Vector{
				Data: embedding,
			},
		},
	}

	// payload
	ret := make(map[string]*pb.Value)
	ret["question"] = &pb.Value{Kind: &pb.Value_StringValue{StringValue: question}}
	ret["answers"] = &pb.Value{Kind: &pb.Value_StringValue{StringValue: answers}}
	point.Payload = ret
	return point
}

func convertToQAs(dataFilePath string) ([]*QA, error) {
	f, err := os.Open(dataFilePath)
	if err != nil {
		return nil, err
	}
	// TODO :big file handle
	defer f.Close()
	content, err := io.ReadAll(f)
	if err != nil {
		return nil, err
	}
	var qas []*QA
	err = json.Unmarshal(content, &qas)
	if err != nil {
		return nil, err
	}
	return qas, nil
}
