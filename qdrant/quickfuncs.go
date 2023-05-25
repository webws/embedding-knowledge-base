package qdrant

// 包含 一些快速函数
var FastQdrantClient *QdrantClient

func init() {
	FastQdrantClient = NewQdrantClient()
}

type Collection string

func (c Collection) Create(size uint64) error {
	return FastQdrantClient.CreateCollection(string(c), size)
}

func (c Collection) Delete() error {
	return FastQdrantClient.DeleteCollection(string(c))
}
