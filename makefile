# 构建qdrant 映射端口 6333 6334
build:
	docker pull qdrant/qdrant && \
	docker run -p 6333:6333 -p 6334:6334 qdrant/qdrant
