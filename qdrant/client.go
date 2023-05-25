package qdrant

import (
	"context"
	"log"

	pb "github.com/qdrant/go-client/qdrant"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const QdrantAddr = "127.0.0.1:6334"

type QdrantClient struct {
	grpcConn *grpc.ClientConn
}

func (qc *QdrantClient) Close() {
	qc.grpcConn.Close()
}

func (qc *QdrantClient) Collection() pb.CollectionsClient {
	return pb.NewCollectionsClient(qc.grpcConn)
}

func NewQdrantClient() *QdrantClient {
	conn, err := grpc.Dial(QdrantAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	return &QdrantClient{grpcConn: conn}
}

func toPayload(payload map[string]string) map[string]*pb.Value {
	ret := make(map[string]*pb.Value)
	for k, v := range payload {
		ret[k] = &pb.Value{Kind: &pb.Value_StringValue{StringValue: v}}
	}
	return ret
}

func (qc *QdrantClient) DeleteCollection(name string) error {
	cc := pb.NewCollectionsClient(qc.grpcConn)
	_, err := cc.Delete(context.TODO(), &pb.DeleteCollection{
		CollectionName: name,
	})
	return err
}

func (qc *QdrantClient) CreateCollection(name string, size uint64) error {
	cc := pb.NewCollectionsClient(qc.grpcConn)

	req := &pb.CreateCollection{
		CollectionName: name,
		VectorsConfig: &pb.VectorsConfig{
			Config: &pb.VectorsConfig_Params{
				Params: &pb.VectorParams{
					Size:     size,
					Distance: pb.Distance_Cosine, // 余弦相似性
				},
			},
		},
	}
	_, err := cc.Create(context.Background(), req)
	if err != nil {
		panic(err)
	}
	return nil
}

// 批量创建Point
func (qc *QdrantClient) CreatePoints(collection string, points []*pb.PointStruct) error {
	pc := pb.NewPointsClient(qc.grpcConn)

	wait := true
	pointsReq := pb.UpsertPoints{
		CollectionName: collection,
		Points:         points,
		Wait:           &wait,
	}

	_, err := pc.Upsert(context.TODO(), &pointsReq)
	if err != nil {
		return err
	}
	return nil
}

// 创建Point的函数。如果不懂， 请1对1 提问。 就是凑数据而已。没啥技术含量

func (qc *QdrantClient) CreatePoint(uuid string, collection string, vector []float32, payload map[string]string) error {
	point := &pb.PointStruct{}
	point.Id = &pb.PointId{
		PointIdOptions: &pb.PointId_Uuid{
			Uuid: uuid,
		},
	}
	point.Vectors = &pb.Vectors{
		VectorsOptions: &pb.Vectors_Vector{
			Vector: &pb.Vector{
				Data: vector,
			},
		},
	}
	point.Payload = toPayload(payload)

	pc := pb.NewPointsClient(qc.grpcConn)

	wait := true
	points := pb.UpsertPoints{
		CollectionName: collection,
		Points:         []*pb.PointStruct{point},
		Wait:           &wait,
	}

	_, err := pc.Upsert(context.TODO(), &points)
	if err != nil {
		return err
	}
	return nil
}

func (qc *QdrantClient) Search(collection string, vector []float32) ([]*pb.ScoredPoint, error) {
	sc := pb.NewPointsClient(qc.grpcConn)
	rsp, err := sc.Search(context.Background(), &pb.SearchPoints{
		CollectionName: collection,
		Vector:         vector,
		Limit:          3, // 只取 3条
		WithPayload: &pb.WithPayloadSelector{
			SelectorOptions: &pb.WithPayloadSelector_Include{
				Include: &pb.PayloadIncludeSelector{
					Fields: []string{"title", "question", "answers"}, // 暴露 title 和 question 字段
				},
			},
		},
	})
	if err != nil {
		return nil, err
	}
	return rsp.Result, nil
}
