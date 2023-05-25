#  使用golang 基于 OpenAI Embedding + qdrant 实现k8s本地知识库
## 流程
![](http://qiniu.taoluyuan.com/2023/blog20230526000739.png?imageMogr2/auto-orient/interlace/1/blur/1x0/quality/70%7Cwatermark/2/text/YmxvZy50YW9sdXl1YW4uY29t/font/5a6L5L2T/fontsize/500/fill/I0E4QTBBMA==/dissolve/100/gravity/NorthWest/dx/10/dy/10)
1. 将数据集 通过 openai embedding 得到向量+组装payload,存入 qdrant
2. 用户进行问题搜索,通过 openai embedding 得到向量,从 qdrant 中搜索相似度大于0.8的数据
3. 从 qdrant 中取出第一条数据
4. 将问题标题,问题描述,问题回答,组装成promot向gpt进行提问,得到回答


## 向量数据库
qdrant 是一个开源的向量搜索引擎,支持多种向量距离计算方式
官方文档:https://qdrant.tech/documentation/quick_start/
本节 介绍 qdrant 都是基于官方文档的例子,如已熟悉可以直接阅读下一节 [数据集导入qdrant]
### 安装 qdrant
docker 安装
```
docker pull qdrant/qdrant && \
docker run -p 6333:6333 -p 6334:6334 qdrant/qdrant
```
### collection 说明
collection 是 qdrant 中的一个概念,类似于 mysql 中的 database,用于区分不同的数据集合 
官方文档:https://qdrant.tech/documentation/collections/#collections 
collection 下面是 collection 字段说明,以创建 collection 为例 
```bash
PUT /collections/{collection_name}
{
    "name": "example_collection",
    "vectors": {
      "size": 300,
      "distance": "Cosine"
    }
}
```
name: collection 名称 
vectors: 向量的配置  
size: 向量的维度 
distance: 向量的距离计算方式,Cosine(余弦距离), Euclidean(欧式距离),Dot product(点积)  
如果需要将 openai embedding 后 存入 qdrant，需要将 size 设置为 1536[openai embedding](https://openai.com/blog/new-and-improved-embedding-model) 
 
### 插入数据
这个是官网 http add point 的例子,可以看到 payload 是可以存储任意的 json 数据,这个数据可以用于后续的过滤
```bash
curl -L -X PUT 'http://localhost:6333/collections/test_collection/points?wait=true' \
    -H 'Content-Type: application/json' \
    --data-raw '{
        "points": [
          {"id": 1, "vector": [0.05, 0.61, 0.76, 0.74], "payload": {"city": "Berlin" }},
          {"id": 2, "vector": [0.19, 0.81, 0.75, 0.11], "payload": {"city": ["Berlin", "London"] }},
          {"id": 3, "vector": [0.36, 0.55, 0.47, 0.94], "payload": {"city": ["Berlin", "Moscow"] }},
          {"id": 4, "vector": [0.18, 0.01, 0.85, 0.80], "payload": {"city": ["London", "Moscow"] }},
          {"id": 5, "vector": [0.24, 0.18, 0.22, 0.44], "payload": {"count": [0] }},
          {"id": 6, "vector": [0.35, 0.08, 0.11, 0.44]}
        ]
    }'
```
* id:唯一
* vector:向量,可在HuggingFace 找相应的模型训练,获取,也可以 openai embedding 得到
* payload:任意的自定义 json 数据
### 搜索数据
这是 qdrant 官方搜索数据的例子,可以看到返回的数据中包含了 payload 中的数据
``` bash
curl -L -X POST 'http://localhost:6333/collections/test_collection/points/search' \
    -H 'Content-Type: application/json' \
    --data-raw '{
        "vector": [0.2,0.1,0.9,0.7],
        "limit": 3
    }'
```
vector:向量,通过 openai embedding 得到
limit:返回的数据条数
## 数据集导入k8s知识数据库
```golang
// 模拟数据集 question:answer
var questions = []string{
	"什么是Kubernetes中的Deployment？",
	"Kubernetes中的Service有什么作用？",
}

var answers = []string{
	"Deployment是Kubernetes中用于管理应用程序副本的资源对象。它提供了副本的声明性定义，可以实现应用程序的部署、扩展和更新。",
	"Service用于定义一组Pod的访问方式和网络策略。它为Pod提供了一个稳定的网络地址，并可以实现负载均衡、服务发现和内部通信。",
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
```
* 上面代码 通过模拟数据集,将数据集导入到 k8s 知识数据库中,具体的实现可以参考 prebuild/prebuild.go 的代码
* BuildQdrantPoint 函数是将问题和答案转换成 qdrant 的 point 
* 其中 vector 是通过 openai embedding 得到的,这里使用的是 [openai embedding](https://openai.com/blog/new-and-improved-embedding-model) 

## 搜索数据 
### 代码实现
```golang
import (
	"fmt"

	myai "embedding-knowledge-base/ai"
	"embedding-knowledge-base/qdrant"
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
```
* 上面代码是通过 prompt 搜索qdrant 知识库,如果相似度小于 0.8,有可能是用户乱提问,或问知识库无关的问题,直接返回
* 如果相似度大于 0.8,则取第一条数据,将问题标题,问题描述,问题回答,组装成promot向gpt进行提问,得到回答
* 具体的实现可以参考 main.go 的代码
### 示例
1. 问无关的问题,比如:苹果不削皮能吃吗 
![](http://qiniu.taoluyuan.com/2023/blog20230526002303.png?imageMogr2/auto-orient/interlace/1/blur/1x0/quality/70%7Cwatermark/2/text/YmxvZy50YW9sdXl1YW4uY29t/font/5a6L5L2T/fontsize/500/fill/I0E4QTBBMA==/dissolve/100/gravity/NorthWest/dx/10/dy/10)
可以看到 相似度太低,提示违规问题或者超纲问题
2. 问k8s 本地知识库的问题,比如:什么是Kubernetes中的Deployment？
   ![](http://qiniu.taoluyuan.com/2023/blog20230526002640.png?imageMogr2/auto-orient/interlace/1/blur/1x0/quality/70%7Cwatermark/2/text/YmxvZy50YW9sdXl1YW4uY29t/font/5a6L5L2T/fontsize/500/fill/I0E4QTBBMA==/dissolve/100/gravity/NorthWest/dx/10/dy/1
3. 问k8s本地知识库的问题,但问题单独向chatgpt提问,并不能得到k8s相关问题.体现qdrant 本地知识库 辅助的重要性,比如问"网关是什么"
  ![](http://qiniu.taoluyuan.com/2023/blog20230526003436.png?imageMogr2/auto-orient/interlace/1/blur/1x0/quality/70%7Cwatermark/2/text/YmxvZy50YW9sdXl1YW4uY29t/font/5a6L5L2T/fontsize/500/fill/I0E4QTBBMA==/dissolve/100/gravity/NorthWest/dx/10/dy/10) 
 可以看到,红线部分,直接问chatgpt,得到的答案可能跟k8s无关,结合k8s本地知识库,可以让回答偏向 数据集设定的主题,比如k8s
## 示例源码地址及使用
源码地址:
进入根目录,将目录 ai/common.go 的 以下 const改成自己的
```golang
    SocksProxy = "socks5://127.0.0.1:1080"
	AIKey      = "your api key"
```
### docker 安装 qdrant
```shell
make install-qdrant
```
### 数据集导入qdrant
* 导入 adrant,我这边就是模拟 了十几条k8s相关的问题,在 prebuild/prebuild.go 
* 更多的数据集,需要自己用脚本抓取,然后导入qdrant
```shell
make import-qdrant
```
### 搜索
```shell
make search
```