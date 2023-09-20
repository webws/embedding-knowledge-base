package cmd

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"io"
	"os"

	"github.com/spf13/cobra"
	"github.com/webws/embedding-knowledge-base/ai"
	"github.com/webws/embedding-knowledge-base/options"
	"github.com/webws/embedding-knowledge-base/qdrant"
	"github.com/webws/go-moda/logger"

	pb "github.com/qdrant/go-client/qdrant"
)

var FlagdataFile = "dataFile" // cmd import flag
type QA struct {
	Questions string `json:"questions" bson:"questions"`
	Answers   string `json:"answers" bson:"answers"`
}

func NewImportCmd(configFlags options.ConfigFlags) *cobra.Command {
	var dataFile string
	importCmd := &cobra.Command{
		Use:   "import",
		Short: "import data to vector database",
		Long:  "import data to vector database",
		RunE: func(cmd *cobra.Command, args []string) error {
			qdrantClient := qdrant.NewQdrantClient(configFlags.Qdrant, configFlags.Collection, configFlags.VectorSize)
			defer qdrantClient.Close()
			aiClient, err := ai.NewAiClient(configFlags.Proxy, configFlags.ApiKey)
			if err != nil {
				return err
			}
			if err = qdrantClient.CreateCollection(configFlags.Collection, configFlags.VectorSize); err != nil {
				return err
			}
			qas, err := convertToQAs(dataFile)
			if err != nil {
				return err
			}
			points := []*pb.PointStruct{}
			logger.Infow("import", "data", qas)
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
	importCmd.Flags().StringVarP(&dataFile, FlagdataFile, "p", "", "import dataFile")
	importCmd.MarkFlagRequired(FlagdataFile)

	return importCmd
}

func buildPoint(question string, answers string, embedding []float32) *pb.PointStruct {
	point := &pb.PointStruct{}
	// point id
	// uuid := fmt.Sprintf("%s%d", md5str(question), time.Now().UnixNano())
	point.Id = &pb.PointId{
		PointIdOptions: &pb.PointId_Uuid{
			Uuid: md5str(question),
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

func md5str(s string) string {
	h := md5.New()
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
}
