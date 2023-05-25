package main

import (
	"fmt"

	myai "github.com/webws/embedding-knowledge-base/ai"
	"github.com/webws/embedding-knowledge-base/qdrant"
)

func main() {
	prompt := "什么是Kubernetes中的DaemonSet？"
	// prompt := "苹果不削皮能吃吗"
	p_vec, err := myai.SimpleGetVec(prompt)
	if err != nil {
		panic(err)
	}
	points, err := qdrant.FastQdrantClient.Search("kubernetes", p_vec)
	if err != nil {
		panic(err)
	}

	fmt.Printf("用户的问题是:%s\n", prompt)
	if points[0].Score < 0.8 {
		fmt.Println("违规问题或者超纲问题")
		return
	}
	answer := points[0].Payload["answers"].GetStringValue()
	fmt.Printf("知识库答案是:%s\n", answer)
	tmpl := "question: %s\n" + "reference answer: %s\n"
	finalPrompt := fmt.Sprintf(tmpl, prompt, points[0].Payload["question"].GetStringValue(), answer)
	fmt.Println("------------------------")
	fmt.Printf("结合知识库参考答案:chatgpt的回答是:%s\n", myai.K8sChat(finalPrompt))
	// 不结合知识库参考答案
	fmt.Printf("不依赖本地知识库, chatgpt的直接回答是:%s\n", myai.K8sChat(prompt))
}
