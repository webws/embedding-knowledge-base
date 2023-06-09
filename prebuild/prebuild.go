package main

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"log"

	myai "github.com/webws/embedding-knowledge-base/ai"
	"github.com/webws/embedding-knowledge-base/qdrant"

	pb "github.com/qdrant/go-client/qdrant"
	"github.com/sashabaranov/go-openai"
)

// 模拟数据集 question:answer
var questions = []string{
	"什么是Kubernetes中的Deployment？",
	"Kubernetes中的Service有什么作用？",
	"如何扩展Kubernetes中的Pod副本数量？",
	"什么是Kubernetes中的命名空间（Namespace）？",
	"Kubernetes中的ConfigMap和Secret有什么区别？",
	"如何在Kubernetes中进行水平扩展（Horizontal Pod Autoscaling）？",
	"Kubernetes中的Ingress是什么？",
	"如何在Kubernetes中进行滚动更新（Rolling Update）？",
	"Kubernetes中的PersistentVolume和PersistentVolumeClaim有什么关系？",
	"什么是Kubernetes中的DaemonSet？",
	"什么是Kubernetes中的Nginx Ingress？",
	"Kubernetes（K8s）中的网关是什么？",
}

var answers = []string{
	"Deployment是Kubernetes中用于管理应用程序副本的资源对象。它提供了副本的声明性定义，可以实现应用程序的部署、扩展和更新。",
	"Service用于定义一组Pod的访问方式和网络策略。它为Pod提供了一个稳定的网络地址，并可以实现负载均衡、服务发现和内部通信。",
	"可以通过更新Deployment的副本数量来扩展Kubernetes中的Pod。可以通过kubectl命令或修改Deployment的YAML文件来指定所需的副本数量。",
	"命名空间是Kubernetes中用于隔离和组织资源的一种机制。它可以将不同的资源划分到不同的命名空间中，实现资源的逻辑隔离和管理。",
	"ConfigMap用于存储应用程序的配置数据，而Secret用于存储敏感的密钥和凭证信息。它们的区别在于Secret的数据会被加密存储，并且可以安全地用于敏感信息的传递。",
	"可以使用Horizontal Pod Autoscaler（HPA）来实现在Kubernetes中的水平扩展。HPA可以根据Pod的CPU利用率或自定义指标自动调整副本数量。",
	"Ingress是Kubernetes中用于暴露HTTP和HTTPS服务的一种资源对象。它可以实现负载均衡、路由和TLS终止等功能。",
	"可以通过更新Deployment的版本或修改Deployment的YAML文件来实现在Kubernetes中的滚动更新。滚动更新可以确保应用程序在更新过程中保持可用性。",
	"PersistentVolume（PV）和PersistentVolumeClaim（PVC）是Kubernetes中用于持久化存储的两个关键概念。",
	"DaemonSet是Kubernetes中一种特殊的控制器，用于在集群中的节点上运行一个Pod副本。它确保每个节点上都有一个副本在运行，用于执行特定的任务或服务",
	"Kubernetes中的Nginx Ingress是一种用于暴露HTTP和HTTPS服务的Ingress控制器。它基于Nginx反向代理实现负载均衡、路由和TLS终止等功能。Nginx Ingress可以根据请求的域名、路径和其他规则将流量转发到相应的后端服务。通过使用Nginx Ingress，可以实现灵活的流量管理和应用程序的访问控制。",
	"在Kubernetes中，网关通常指的是Ingress（入口）资源对象。Ingress是一种Kubernetes API对象，用于配置和管理集群中的HTTP和HTTPS流量入口。它充当了从集群外部访问集群内部服务的入口点",
}

func main() {
	// 第一步：自己创建 一个collection:  kubernetes
	var err error
	err = qdrant.Collection("kubernetes").Create(1536)
	if err != nil {
		log.Fatalln("创建collection出错:", err.Error())
	}

	points := []*pb.PointStruct{}
	// 批量 进行BuildQdrantPoint
	for index, question := range questions {
		if index < 9 {
			continue
		}
		p, err := BuildQdrantPoint(question, answers[index])
		if err != nil {
			log.Fatalln("创建point出错:", err.Error())
		}
		fmt.Println(p.Id)
		points = append(points, p)

	}
	err = qdrant.FastQdrantClient.CreatePoints("kubernetes", points)
	if err != nil {
		log.Fatalln("批量创建point出错:", err.Error())
	}
}

func BuildQdrantPoint(
	question string,
	answers string,
) (*pb.PointStruct, error) {
	point := &pb.PointStruct{}
	// 纳秒
	// uuid := fmt.Sprintf("%s-%d", question, time.Now().UnixNano())
	uuid := question
	point.Id = &pb.PointId{
		PointIdOptions: &pb.PointId_Uuid{
			Uuid: md5str(uuid), // 保证 id的唯一性
		},
	}
	c := myai.NewOpenAiClient()
	req := openai.EmbeddingRequest{
		Input: []string{question},
		Model: openai.AdaEmbeddingV2,
	}
	rsp, err := c.CreateEmbeddings(context.TODO(), req)
	if err != nil {
		return nil, err
	}

	// vector
	point.Vectors = &pb.Vectors{
		VectorsOptions: &pb.Vectors_Vector{
			Vector: &pb.Vector{
				Data: rsp.Data[0].Embedding,
			},
		},
	}
	// payload
	ret := make(map[string]*pb.Value)
	// ret["title"] = &pb.Value{Kind: &pb.Value_StringValue{StringValue: title}}
	ret["question"] = &pb.Value{Kind: &pb.Value_StringValue{StringValue: question}}
	ret["answers"] = &pb.Value{Kind: &pb.Value_StringValue{StringValue: answers}}
	point.Payload = ret
	return point, nil
}

func md5str(s string) string {
	h := md5.New()
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
}
