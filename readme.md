Translations: [English](readme.md) | [简体中文](README_zh.md)
## kbai
a local knowledge base, based on chatgpt and qdrant
### usage
local knowledge base based on chatgpt and qdrant, supporting data import and Q&A
```
❯ kbai -h
a local knowledge base, based on chatgpt and qdrant

usage:
  kbai [flags]
  kbai [command]

available commands:
  completion  generate the autocompletion script for the specified shell
  help        help about any command
  import      import data to vector database
  search      ask the knowledge base example: kbai ask --msg 'first, the chicken or the egg'

flags:
      --apikey string       openai apikey:default from env apikey
      --collection string   qdrant collection name default: kubernetes (default "kubernetes")
  -h, --help                help for kbai
      --proxy string        http client proxy default:socks5://127.0.0.1:1080  (default "socks5://127.0.0.1:1080")
      --qdrant string       qdrant address default: 127.0.0.1:6334 (default "127.0.0.1:6334")
      --vectorsize uint     qdrant vector size default: 1536 (default 1536)

use "kbai [command] --help" for more information about a command.
```
## install
go build install rename and move to the $PATH
```
sudo go build -o kbai  github.com/webws/embedding-knowledge-base && sudo  mv ./kbai /usr/local/bin
```
Or use golang to execute source code
```
git clone https://github.com/webws/embedding-knowledge-base.git && cd ./embedding-knowledge-base

```
or download binary file from release……(todo)
## use example
first, you must start the vector database qdrant
```
docker run --rm -p 6334:6334 qdrant/qdrant
```

set openai apikey
```
export apikey=xxx
```
### import
import the prepared JSON data into qdrant
```
kbai import --datafile ./example/data.json
```
example format of file data.json
```
[
    {
        "questions": "question",
        "answers": "answer"
    },
]
```
note:the imported sample data content is: k8s related knowledge, and it is in Chinese

### search

ask knowledge related to the knowledge base
```
kbai search --msg "What is a gateway"

The answer to the knowledge base:
在Kubernetes中，网关通常指的是Ingress（入 口）资源对象。Ingress是一种Kubernetes API对象，用于配置和管理集群中的HTTP和HTTPS流量入口。它充当了从集群外部访问集群内部服务的入口点

Results of chatgpt answers  with reference answers:
Ingress acts as a gateway in Kubernetes, allowing external traffic to access the internal services within the cluster. It provides a configuration layer for managing HTTP and HTTPS traffic routing rules, load balancing, and SSL termination. Ingress resources define the rules for how the traffic should be directed to the appropriate backend services based on the requested host or URL path. Ingress controllers, such as Nginx or Traefik, are responsible for implementing these rules and routing the traffic accordingly.

only chatgpt answers:
A gateway is a device or software that acts as an entry point or interface between two different networks, systems, or protocols. It serves as a connecting link that allows data to flow between different networks, translating between different protocols or formats if necessary. Gateways are commonly used in computer networks to connect local area networks (LANs) to wide area networks (WANs) or to bridge different types of networks. They can also be used in telecommunications to connect different types of networks, such as telephone networks and internet networks.
```
ps:
>It can be seen that directly asking chatgpt may result in an answer unrelated to k8s. Combining with the local knowledge base of k8s can bias the answer towards the theme set in the dataset


ask questions unrelated to the knowledge base,There won't be any results

```
kbai search --msg "Can apples be eaten without washing?"
rearch term violation or exceeding category
```
### reference project
* [https://github.com/spf13/cobra](https://github.com/spf13/cobra) 
* [https://github.com/kubernetes/kubernetes](https://github.com/kubernetes/kubernetes) 
* [https://github.com/gohugoio/hugo](https://github.com/gohugoio/hugo) 
* [https://github.com/qdrant/qdrant](https://github.com/qdrant/qdrant) 
