.PHONY: help test lint fmt vet build clean examples

help: ## 显示帮助信息
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-15s\033[0m %s\n", $$1, $$2}'

test: ## 运行测试
	go test -v -race -cover ./...

test-coverage: ## 生成测试覆盖率报告
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

lint: ## 运行 golangci-lint
	@which golangci-lint > /dev/null || (echo "请先安装 golangci-lint: https://golangci-lint.run/usage/install/" && exit 1)
	golangci-lint run ./...

fmt: ## 格式化代码
	go fmt ./...
	gofmt -s -w .

vet: ## 运行 go vet
	go vet ./...

build: ## 构建示例程序
	go build -o bin/basic ./examples/basic
	go build -o bin/advanced ./examples/advanced

clean: ## 清理构建产物
	rm -rf bin/
	rm -f coverage.out coverage.html

tidy: ## 整理依赖
	go mod tidy

check: fmt vet lint test ## 运行所有检查

.DEFAULT_GOAL := help
