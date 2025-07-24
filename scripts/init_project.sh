#!/bin/bash

# 项目模版化脚本
# 用法: ./scripts/init_project.sh <new-project-name>

set -e

if [ $# -eq 0 ]; then
    echo "错误: 请提供新项目名称"
    echo "用法: $0 <new-project-name>"
    echo "示例: $0 my-new-api"
    exit 1
fi

NEW_PROJECT_NAME=$1
OLD_PROJECT_NAME="justus"

echo "🚀 开始初始化新项目: $NEW_PROJECT_NAME"

# 检查项目名称是否有效
if [[ ! "$NEW_PROJECT_NAME" =~ ^[a-z][a-z0-9-]*$ ]]; then
    echo "❌ 错误: 项目名称必须以小写字母开头，只能包含小写字母、数字和连字符"
    exit 1
fi

# 备份原始文件
echo "📝 备份原始文件..."
cp go.mod go.mod.backup
cp conf/app.dev.yaml conf/app.dev.yaml.backup

# 1. 更新 go.mod 文件
echo "🔧 更新 go.mod 文件..."
sed -i "" "s/module $OLD_PROJECT_NAME/module $NEW_PROJECT_NAME/g" go.mod

# 2. 更新所有 Go 文件中的导入路径
echo "🔧 更新导入路径..."
find . -name "*.go" -type f -exec sed -i "" "s/$OLD_PROJECT_NAME\//$NEW_PROJECT_NAME\//g" {} \;

# 3. 更新配置文件中的项目相关配置
echo "🔧 更新配置文件..."
sed -i "" "s/justus:/$NEW_PROJECT_NAME:/g" conf/app.dev.yaml

# 4. 更新 README.md
echo "🔧 更新 README.md..."
sed -i "" "s/Justus Go/$NEW_PROJECT_NAME/g" README.md
sed -i "" "s/justus-go/$NEW_PROJECT_NAME/g" README.md

# 5. 清理依赖
echo "🧹 清理依赖..."
go mod tidy

# 6. 创建环境变量文件
echo "📄 创建环境变量文件..."
cat > .env.example << EOF
# 环境配置
APP_ENV=dev

# 数据库配置
DB_HOST=127.0.0.1:3306
DB_USER=root
DB_PASSWORD=root
DB_NAME=$NEW_PROJECT_NAME

# Redis配置
REDIS_HOST=127.0.0.1:6379
REDIS_PASSWORD=
REDIS_DB=1

# JWT配置
JWT_SECRET=M2lTP9APLdRY6TA5RC42rz5AkxsgoAZNdVx1bC2XQXlh2pxEJ5waB5EIKfp4CHfM

# 应用配置
APP_HOST=127.0.0.1
APP_PORT=8787
APP_DEBUG=true

# 文件上传配置
UPLOAD_PATH=uploads/
MAX_UPLOAD_SIZE=10485760

# 日志配置
LOG_LEVEL=info
LOG_PATH=storage/logs/
EOF

# 7. 创建 .env 文件
if [ ! -f .env ]; then
    echo "📄 创建 .env 文件..."
    cp .env.example .env
fi

# 8. 创建必要的目录
echo "📁 创建必要的目录..."
mkdir -p storage/logs
mkdir -p uploads
mkdir -p runtime

# 9. 编译测试
echo "🔍 编译测试..."
if go build -v cmd/${OLD_PROJECT_NAME}-go.go; then
    echo "✅ 编译成功!"
    rm -f ${OLD_PROJECT_NAME}-go  # 删除编译产物
else
    echo "❌ 编译失败，请检查错误信息"
    # 恢复备份文件
    mv go.mod.backup go.mod
    mv conf/app.dev.yaml.backup conf/app.dev.yaml
    exit 1
fi

# 清理备份文件
rm -f go.mod.backup conf/app.dev.yaml.backup

echo ""
echo "🎉 项目初始化完成!"
echo ""
echo "接下来的步骤:"
echo "1. 安装并启动 MySQL 和 Redis 服务"
echo "2. 创建数据库: CREATE DATABASE $NEW_PROJECT_NAME;"
echo "3. 编辑 .env 文件配置数据库和Redis连接信息（可选）"
echo "4. 执行数据库初始化: make db-init"
echo "5. 启动项目: make run 或 make dev"
echo "6. 访问健康检查: http://localhost:8787/health"
echo "7. 访问测试接口: http://localhost:8787/api/v1/test"
echo ""
echo "�� 更多信息请查看 README.md" 