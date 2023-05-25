#  chatgpt embedding 与 qdrant 向量数据库 实现k8s本地知识库
## 一.向量数据库
qdrant 是一个开源的向量搜索引擎,支持多种向量距离计算方式
官方文档:https://qdrant.tech/documentation/quick_start/
本节 介绍 qdrant 都是基于官方文档的例子,如熟悉可以直接阅读下一节 数据集导入k8s知识数据库
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
id:唯一
vector:向量,通过 openai embedding 得到
payload:任意的 json 数据
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
## 查询知识数据库
## 示例源码地址及使用
# 启动
1. 运行 qdrant:make build 启动qdrant
2. 往 qdrant 塞入数据集,包括 问题标题,问题描述,问题回答: prebuild/prebuild.go
3. 模拟根据prompt 搜索问题,mian.go 先通过 prompt 从向量数据库搜索数据,相似度要 〉0.8 ,防止问一些违规的问题.
4. 搜索到数据后,从向量数据库取第一条数据,将问题标题,问题描述,问题回答,组装成promot向gpt进行提问,得到回答