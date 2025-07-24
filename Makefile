# Makefile for Justus Go API Project

# 变量定义
APP_NAME=justus-go
VERSION=$(shell git describe --tags --always)
BUILD_TIME=$(shell date +%Y-%m-%d\ %H:%M:%S)
GO_VERSION=$(shell go version | awk '{print $$3}')

# 编译标志
LDFLAGS=-ldflags "-X main.Version=${VERSION} -X main.BuildTime=${BUILD_TIME} -X main.GoVersion=${GO_VERSION}"

# 默认目标
.PHONY: help
help: ## 显示帮助信息
	@echo "Available commands:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-15s\033[0m %s\n", $$1, $$2}'

# 开发相关命令
.PHONY: dev
dev: ## 启动开发服务器（带热重载）
	@echo "Starting development server..."
	air -c .air.toml

.PHONY: run
run: ## 运行应用
	@echo "Running application..."
	go run cmd/justus-go.go

.PHONY: build
build: ## 编译应用
	@echo "Building application..."
	go build ${LDFLAGS} -o bin/${APP_NAME} cmd/justus-go.go

.PHONY: test
test: ## 运行测试
	@echo "Running tests..."
	go test -v ./...

.PHONY: benchmark
benchmark: ## 运行性能测试
	@echo "Running benchmarks..."
	go test -bench=. -benchmem ./...

# 代码质量检查
.PHONY: lint
lint: ## 运行代码检查
	@echo "Running linter..."
	golangci-lint run

.PHONY: fmt
fmt: ## 格式化代码
	@echo "Formatting code..."
	go fmt ./...
	goimports -w .

.PHONY: vet
vet: ## 运行 go vet
	@echo "Running go vet..."
	go vet ./...

# 依赖管理
.PHONY: deps
deps: ## 下载依赖
	@echo "Downloading dependencies..."
	go mod download

.PHONY: tidy
tidy: ## 整理依赖
	@echo "Tidying dependencies..."
	go mod tidy

.PHONY: vendor
vendor: ## 创建 vendor 目录
	@echo "Creating vendor directory..."
	go mod vendor

# 数据库相关命令
.PHONY: db-init
db-init: ## 使用SQL文件初始化数据库
	@echo "Initializing database with SQL..."
	@echo "注意：需要手动输入MySQL密码"
	mysql -u root -p < scripts/database-init.sql



# 清理命令
.PHONY: clean
clean: ## 清理构建文件
	@echo "Cleaning up..."
	rm -rf bin/
	rm -rf dist/
	rm -rf vendor/
	go clean

.PHONY: clean-logs
clean-logs: ## 清理日志文件
	@echo "Cleaning log files..."
	rm -rf storage/logs/*
	rm -rf runtime/logs/*

# 项目初始化命令
.PHONY: init
init: ## 初始化项目（首次运行）
	@echo "Initializing project..."
	go mod tidy
	mkdir -p bin/ storage/logs uploads runtime/logs
	@echo "Project initialized successfully!"

.PHONY: template
template: ## 将项目转换为新模板（指定 PROJECT_NAME=your-project）
	@echo "Converting to template..."
	@if [ -z "$(PROJECT_NAME)" ]; then \
		echo "Error: Please specify PROJECT_NAME"; \
		echo "Usage: make template PROJECT_NAME=your-project-name"; \
		exit 1; \
	fi
	chmod +x scripts/init_project.sh
	./scripts/init_project.sh $(PROJECT_NAME)

# 生产部署相关
.PHONY: build-prod
build-prod: ## 生产环境构建
	@echo "Building for production..."
	CGO_ENABLED=0 GOOS=linux go build ${LDFLAGS} -o bin/${APP_NAME} cmd/justus-go.go

.PHONY: deploy
deploy: build-prod ## 部署到生产环境
	@echo "Deploying to production..."
	# 在这里添加您的部署脚本

# 开发工具安装
.PHONY: install-tools
install-tools: ## 安装开发工具
	@echo "Installing development tools..."
	go install github.com/cosmtrek/air@latest
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install golang.org/x/tools/cmd/goimports@latest

# 版本管理
.PHONY: version
version: ## 显示版本信息
	@echo "App Name: ${APP_NAME}"
	@echo "Version: ${VERSION}"
	@echo "Build Time: ${BUILD_TIME}"
	@echo "Go Version: ${GO_VERSION}"

.DEFAULT_GOAL := help 