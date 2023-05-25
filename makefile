# 构建qdrant 映射端口 6333 6334
install-qdrant:
	docker pull qdrant/qdrant && \
	docker run -p 6333:6333 -p 6334:6334 qdrant/qdrant
# 将数据集导入qdrant
import-qdrant:
	go run ./pre-build/prebuild.go
# 搜索qdrant 的结果,获取结果中的问题和推荐答案,组装后对openai进行调用
search:
	go run ./main.go
	