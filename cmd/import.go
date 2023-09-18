package cmd

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"
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
		logger.Infow("import cmd")
		qdrantAddr, err := cmd.Flags().GetString(qdrantFlag)
		if err != nil {
			return err
		}
		qdrantClient := qdrant.NewQdrantClient(qdrantAddr)
		defer qdrantClient.Close()

		if err = qdrantClient.CreateCollection("kubernetes", 1536); err != nil {
			return err
		}
		// TODO config to QA
		var qas []*QA
		points := []*pb.PointStruct{}
		for _, qa := range qas {
			// TODO openai CreateEmbeddings
			var embedding []float32
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
