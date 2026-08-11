.PHONY: help tidy fmt vet lint build run test cover clean

BINARY := bin/qio-server
MAIN   := ./cmd/server

help: ## 显示可用命令
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-10s\033[0m %s\n", $$1, $$2}'

tidy: ## 整理依赖
	go mod tidy

fmt: ## 格式化代码
	gofmt -l -w .

vet: ## 静态检查
	go vet ./...

lint: ## golangci-lint 检查（需先安装 golangci-lint）
	golangci-lint run ./...

build: ## 编译二进制
	go build -o $(BINARY) $(MAIN)

run: ## 本地启动
	go run $(MAIN)

test: ## 运行测试
	go test ./...

cover: ## 运行测试并输出覆盖率
	go test -coverprofile=coverage.out ./...
	go tool cover -func=coverage.out

clean: ## 清理产物
	rm -rf bin coverage.out
