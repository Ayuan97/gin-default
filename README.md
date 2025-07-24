# Justus Go API 项目模版

这是一个基于 Gin 框架的 Go Web API 项目模版，包含了 **Admin 管理端** 和 **API 前端** 两个独立模块，以及常用的中间件和基础功能。

## 功能特性

- ✅ **双模块架构** - Admin 管理端 + API 前端，权限分离
- ✅ **Gin Web 框架** - 高性能 HTTP Web 框架
- ✅ **JWT 认证** - 基于 JWT 的用户认证系统
- ✅ **角色权限管理** - 完整的 RBAC 权限控制系统
- ✅ **数据库支持** - GORM ORM + MySQL 数据库
- ✅ **Redis 缓存** - Redis 缓存支持
- ✅ **请求签名验证** - API 请求签名中间件
- ✅ **日志系统** - 结构化日志记录
- ✅ **错误恢复** - 优雅的错误处理和恢复
- ✅ **请求日志** - 详细的请求响应日志
- ✅ **CORS 支持** - 跨域资源共享配置
- ✅ **国际化** - 多语言支持
- ✅ **配置管理** - 基于环境的配置文件
- ✅ **定时任务** - Cron 定时任务支持

## 模块架构

### API 模块 (`/api/v1/*`)

- **面向对象**: 前端用户、移动端应用
- **功能**: 用户注册登录、个人信息管理、业务功能接口
- **权限**: 基于 JWT 的用户认证，普通用户权限
- **特点**: 包含签名验证，适合对外开放的接口

### Admin 模块 (`/admin/v1/*`)

- **面向对象**: 管理员、运营人员
- **功能**: 用户管理、角色权限管理、系统监控、数据统计
- **权限**: 管理员权限验证，支持细粒度权限控制
- **特点**: 无签名验证要求，内部管理使用

## 项目结构

```
justus-go/
├── cmd/                    # 程序入口
│   └── justus-go.go       # 主程序文件
├── conf/                   # 配置文件
│   └── app.dev.yaml       # 开发环境配置
├── internal/               # 内部包
│   ├── dao/               # 数据访问层
│   ├── global/            # 全局变量
│   ├── middleware/        # 中间件
│   │   ├── admin/         # 管理员权限中间件
│   │   │   └── auth.go    # 权限验证中间件
│   │   ├── api_require/   # API请求中间件
│   │   ├── bodyLog/       # 请求日志中间件
│   │   ├── cors/          # 跨域中间件
│   │   ├── jwt/           # JWT认证中间件
│   │   ├── logging/       # 日志中间件
│   │   ├── recovers/      # 错误恢复中间件
│   │   └── sign/          # 签名验证中间件
│   ├── models/            # 数据模型
│   │   ├── models.go      # 基础模型
│   │   ├── user.go        # 用户模型
│   │   └── role.go        # 角色权限模型
│   └── service/           # 业务逻辑层
│       └── user.go        # 用户服务
├── pkg/                    # 工具包
│   ├── app/               # 应用工具
│   ├── e/                 # 错误码定义
│   │   ├── code.go        # 错误码常量
│   │   └── msg.go         # 错误消息
│   ├── gredis/            # Redis工具
│   ├── logger/            # 日志工具
│   └── util/              # 通用工具
└── routers/               # 路由定义
    ├── router.go          # 主路由配置
    ├── api/               # API模块路由 (前端用户)
    │   ├── test_controller.go    # 测试接口
    │   └── user_controller.go    # 用户接口
    └── admin/             # Admin模块路由 (管理员)
        ├── system_controller.go  # 系统管理接口
        ├── user_controller.go    # 用户管理接口
        └── role_controller.go    # 角色管理接口
```

### 目录说明

#### `internal/middleware/`

- **admin/**: 管理员专用中间件
  - `auth.go`: 管理员权限验证和特定权限检查
- **其他中间件**: 通用中间件，可被多个模块使用

#### `routers/`

- **api/**: API 模块控制器 (面向前端用户)
  - 命名规范: `*_controller.go`
  - 包含用户认证、个人信息管理等接口
- **admin/**: Admin 模块控制器 (面向管理员)
  - 命名规范: `*_controller.go`
  - 包含用户管理、系统监控、角色权限管理等接口

#### `internal/models/`

- 数据模型定义，支持完整的 RBAC 权限系统
- 包含用户、角色、权限的关联关系

#### `pkg/e/`

- 统一的错误码管理系统
- 分类清晰的错误码和消息定义

## 快速开始

### 环境要求

- Go 1.24.3+
- MySQL 5.7+
- Redis 6.0+

### 安装步骤

1. **克隆项目**

   ```bash
   git clone <repository-url>
   cd justus-go
   ```

2. **安装依赖**

   ```bash
   go mod tidy
   ```

3. **配置环境**

   ```bash
   # 复制环境变量配置文件
   cp .env.example .env

   # 编辑环境变量文件（可选，有默认值）
   vim .env

   # 或者直接编辑配置文件
   vim conf/app.dev.yaml
   ```

4. **启动数据库和 Redis**

   ```bash
   # 启动MySQL和Redis服务
   # 创建数据库：justus

   # 初始化数据库（推荐，自动使用配置文件连接）
   make db-init

   # 或者使用SQL文件初始化（需要输入密码）
   make db-init-sql
   ```

5. **运行项目**

   ```bash
   # 开发模式（热重载）
   make dev

   # 或者直接运行
   make run
   ```

服务器将在 `http://localhost:8787` 启动

## API 接口文档

### API 模块接口 (`/api/v1/*`)

```bash
# 测试接口
curl -X POST http://localhost:8787/api/v1/test \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <your-jwt-token>" \
  -H "skip-signature: true"

# 获取个人信息
curl -X GET http://localhost:8787/api/v1/profile \
  -H "Authorization: Bearer <your-jwt-token>" \
  -H "skip-signature: true"

# 更新个人信息
curl -X PUT http://localhost:8787/api/v1/profile \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <your-jwt-token>" \
  -H "skip-signature: true" \
  -d '{"first_name": "John", "last_name": "Doe", "phone": "1234567890"}'

# 获取用户列表
curl -X GET http://localhost:8787/api/v1/users \
  -H "Authorization: Bearer <your-jwt-token>" \
  -H "skip-signature: true"
```

### Admin 模块接口 (`/admin/v1/*`)

```bash
# 管理员测试接口
curl -X POST http://localhost:8787/admin/v1/test \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <admin-jwt-token>"

# 获取用户列表（管理员）
curl -X GET "http://localhost:8787/admin/v1/users?page=1&limit=20&keyword=john&status=1&role=admin" \
  -H "Authorization: Bearer <admin-jwt-token>"

# 获取单个用户详情
curl -X GET http://localhost:8787/admin/v1/users/1 \
  -H "Authorization: Bearer <admin-jwt-token>"

# 更新用户状态
curl -X PUT http://localhost:8787/admin/v1/users/1/status \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <admin-jwt-token>" \
  -d '{"status": 0}'

# 获取系统信息
curl -X GET http://localhost:8787/admin/v1/system/info \
  -H "Authorization: Bearer <admin-jwt-token>"

# 获取系统统计
curl -X GET http://localhost:8787/admin/v1/system/stats \
  -H "Authorization: Bearer <admin-jwt-token>"

# 获取角色列表
curl -X GET http://localhost:8787/admin/v1/roles \
  -H "Authorization: Bearer <admin-jwt-token>"
```

## 权限系统

### 角色定义

- **admin**: 系统管理员，拥有所有权限
- **moderator**: 内容管理员，拥有用户管理和内容管理权限
- **user**: 普通用户，只能管理自己的资料

### 权限定义

权限采用 `resource.action` 格式：

- `user.read`: 查看用户信息
- `user.write`: 修改用户信息
- `user.delete`: 删除用户
- `system.read`: 查看系统信息
- `system.write`: 修改系统配置
- `role.read`: 查看角色信息
- `role.write`: 管理角色权限

### 权限验证

```go
// 检查管理员权限
adminGroup.Use(admin.AdminAuth())

// 检查特定权限
adminGroup.GET("/sensitive", admin.RequirePermission("system.write"), handler)
```

## 配置说明

### 主要配置文件

- `conf/app.dev.yaml` - 开发环境配置
- `.env` - 环境变量配置

### 关键配置项

**配置文件方式** (`conf/app.dev.yaml`):

```yaml
app:
  PageSize: 20
  JwtSecret: "your-jwt-secret"

server:
  RunMode: debug
  HttpPort: 8787

database:
  Host: 127.0.0.1:3306
  Name: justus
  User: root
  Password: root

redis:
  Host: 127.0.0.1:6379
  Prefix: "justus:"
```

**环境变量方式** (`.env`):

```bash
# 应用配置
JWT_SECRET=your-super-secret-key
APP_PORT=8787

# 数据库配置
DB_HOST=127.0.0.1:3306
DB_USER=root
DB_PASSWORD=root
DB_NAME=justus

# Redis配置
REDIS_HOST=127.0.0.1:6379
REDIS_PREFIX=justus:
```

## 作为模版使用

要将此项目用作新项目的模版：

1. **使用模版化脚本** (推荐)

   ```bash
   chmod +x scripts/init_project.sh
   ./scripts/init_project.sh your-new-project
   ```

2. **手动替换**
   - 修改 `go.mod` 中的模块名称
   - 替换所有代码中的 `justus` 导入路径
   - 更新配置文件中的项目相关配置

## 开发指南

### 添加新 API 接口

1. **API 模块**:

   - 在 `routers/api/` 中添加处理函数
   - 在 `routers/router.go` 的 `apiGroup` 中注册路由
   - 遵循 RESTful 设计原则

2. **Admin 模块**:
   - 在 `routers/admin/` 中添加处理函数
   - 在 `routers/router.go` 的 `adminGroup` 中注册路由
   - 考虑权限控制和数据安全

### 添加新权限

1. 在 `models/role.go` 中定义新权限
2. 更新角色权限关联
3. 在需要的接口中使用 `RequirePermission` 中间件

### 添加新角色

1. 在数据库中创建新角色记录
2. 分配相应的权限
3. 更新用户角色关联

## 技术栈

- **Web 框架**: Gin
- **ORM**: GORM
- **数据库**: MySQL
- **缓存**: Redis
- **认证**: JWT
- **权限**: RBAC (Role-Based Access Control)
- **日志**: Logrus
- **配置**: Viper
- **定时任务**: Cron

## 安全注意事项

- 确保 JWT 密钥在生产环境中使用强密码
- 数据库连接信息请勿提交到版本控制系统
- Redis 前缀建议使用项目名，避免键冲突
- 生产环境建议使用环境变量管理敏感配置
- Admin 接口应该部署在内网或使用 VPN 访问
- 定期审查用户权限，及时清理无效账户
- 对管理员操作进行日志记录和审计

## 错误码说明

项目使用统一的错误码系统：

- `2xxxx`: 认证相关错误
- `3xxxx`: 文件上传相关错误
- `4xxxx`: 用户相关错误
- `41xxx`: 权限和角色相关错误
- `42xxx`: 管理员相关错误
- `5xxxx`: 数据库相关错误
- `6xxxx`: 缓存相关错误
- `7xxxx`: 文件相关错误
- `8xxxx`: 网络相关错误
- `9xxxx`: 业务逻辑和系统相关错误

## License

MIT License
