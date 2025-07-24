# Justus Go API 项目模版

基于 Gin 框架的企业级 Go Web API 项目模版，采用双模块架构设计。

## 核心特性

- 🏗️ **双模块架构** - API 前端 + Admin 管理端，权限分离
- 🔐 **完整认证** - JWT 认证 + RBAC 权限控制 + 请求签名验证
- 💾 **数据支持** - MySQL + GORM ORM + Redis 缓存
- 🌍 **国际化** - 50+语言支持
- 📝 **日志系统** - 结构化日志记录和请求追踪
- ⚙️ **开发友好** - 热重载、统一错误码、配置管理

## 模块架构

| 模块           | 路径          | 面向对象         | 特点                    |
| -------------- | ------------- | ---------------- | ----------------------- |
| **API 模块**   | `/api/v1/*`   | 前端用户、移动端 | JWT 认证 + 签名验证     |
| **Admin 模块** | `/admin/v1/*` | 管理员、运营人员 | 管理员权限 + 细粒度控制 |

## 项目结构

```
justus-go/
├── cmd/                    # 程序入口
├── conf/                   # 配置文件
├── internal/               # 内部代码
│   ├── middleware/        # 中间件（认证、日志、签名等）
│   ├── models/            # 数据模型（用户、角色、权限）
│   └── service/           # 业务逻辑
├── pkg/                    # 工具包
└── routers/               # 路由控制器
    ├── api/               # API模块（前端用户）
    └── admin/             # Admin模块（管理员）
```

## 快速开始

### 环境要求

- Go 1.24.3+
- MySQL 5.7+
- Redis 6.0+

### 启动步骤

```bash
# 1. 安装依赖
go mod tidy

# 2. 配置环境（可选，有默认配置）
cp .env.example .env

# 3. 初始化数据库
make db-init

# 4. 启动开发服务（热重载）
make dev

# 或直接运行
make run
```

服务器启动在 `http://localhost:8787`

### 常用命令

- `make dev` - 热重载开发
- `make run` - 直接运行
- `make db-init` - 智能数据库初始化
- `make help` - 查看所有命令

## 配置说明

**默认配置** (`conf/app.dev.yaml`):

```yaml
app:
  PageSize: 20
  JwtSecret: "your-jwt-secret"
server:
  HttpPort: 8787
database:
  Host: 127.0.0.1:3306
  Name: justus
redis:
  Host: 127.0.0.1:6379
  Prefix: "justus:"
```

**环境变量覆盖** (`.env`):

```bash
JWT_SECRET=your-secret-key
APP_PORT=8787
DB_HOST=127.0.0.1:3306
DB_PASSWORD=your-password
```

## API 使用

### API 模块示例

```bash
# 健康检查
curl http://localhost:8787/health

# 获取个人信息（需要JWT token）
curl -H "Authorization: Bearer <token>" \
     -H "skip-signature: true" \
     http://localhost:8787/api/v1/profile
```

### Admin 模块示例

```bash
# 获取用户列表（需要管理员权限）
curl -H "Authorization: Bearer <admin-token>" \
     "http://localhost:8787/admin/v1/users?page=1&limit=20"

# 获取系统信息
curl -H "Authorization: Bearer <admin-token>" \
     http://localhost:8787/admin/v1/system/info
```

## 权限系统

### 内置角色

- **admin**: 系统管理员（所有权限）
- **moderator**: 内容管理员（用户和内容管理）
- **user**: 普通用户（个人资料管理）

### 权限格式

采用 `resource.action` 格式：

- `user.read/write/delete` - 用户管理权限
- `system.read/write` - 系统管理权限
- `role.read/write` - 角色管理权限

## 开发指南

### 添加 API 接口

1. 在 `routers/api/` 创建控制器
2. 在 `internal/service/` 添加业务逻辑
3. 在 `routers/router.go` 注册路由
4. 测试时使用 `skip-signature: true` 跳过签名

### 添加 Admin 接口

1. 在 `routers/admin/` 创建控制器
2. 使用 `admin.Auth()` 验证权限
3. 添加特定权限检查（可选）

### 统一响应格式

```go
appG := app.Gin{C: c}
appG.Success(data)           // 成功
appG.Error(e.ERROR_CODE)     // 错误
appG.InvalidParams()         // 参数错误
```

## 技术栈

- **Web 框架**: Gin
- **数据库**: MySQL + GORM ORM
- **缓存**: Redis
- **认证**: JWT + 请求签名
- **权限**: RBAC
- **日志**: Logrus
- **国际化**: 50+语言支持

## 作为模版使用

```bash
# 使用模版化脚本
chmod +x scripts/init_project.sh
./scripts/init_project.sh your-new-project
```

## License

MIT License
