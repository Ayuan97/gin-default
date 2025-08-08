## 开发文档（Dev Guide）

### 概览

- **Web 框架**: Gin
- **数据库**: MySQL + GORM（表前缀 `ay_`）
- **缓存**: Redis（键前缀 `justus:`）
- **认证**: JWT
- **权限**: RBAC
- **日志**: Logrus（文件输出在 `runtime/logs/`）
- **国际化**: 50+ 语言（`pkg/langefile`）
- **端口**: 8787（见 `conf/app.*.yaml` 或环境变量）

### 目录结构（核心）

- **入口**: `cmd/justus-go.go`
- **配置**: `conf/app.dev.yaml` / `conf/app.production.yaml`
- **控制器**: `internal/controllers/{api,admin,common}`
- **中间件**: `internal/middleware/**`
- **模型**: `internal/models/**`
- **仓储**: `internal/repository/**`
- **服务**: `internal/service/**`
- **路由**: `internal/routers/**`
- **定时任务**: `internal/cron/cron.go`
- **基础设施**: `internal/infrastructure/**`
- **工具包**: `pkg/**`

### 启动与环境

- 依赖: Go 1.21+、MySQL、Redis
- 初始化数据库
  - 本地初始化: `make db-init`
  - 仅迁移: `make migrate`
  - 仅种子: `make seed`
- 运行
  - 开发热更: `make dev`
  - 直接运行: `make run`
- 健康检查: `GET http://localhost:8787/health`
- 配置优先级: 环境变量 > `.env` > `conf/app.*.yaml`
  - 常用环境变量: `JWT_SECRET, APP_PORT, DB_HOST, DB_USER, DB_PASSWORD, REDIS_HOST`

### 中间件顺序（建议）

Logger → Recover → BodyLog → Common → Sign → JWT → Auth

- **API 签名**: 开发时可用请求头 `skip-signature: true` 跳过
- **管理端权限**: 使用 `internal/middleware/admin/auth.go`

### 路由规范

- **API 模块**: `/api/v1/*`（签名验证）
- **Admin 模块**: `/admin/v1/*`（JWT + RBAC）
- **路由注册位置**: `internal/routers/*`（按模块拆分）
- **入口**: 在 `cmd/justus-go.go` 中挂载路由树

### 统一响应

控制器中统一使用 `pkg/app`：

```go
appG := app.Gin{C: c}
appG.Success(data)
appG.Error(e.ERROR_CODE)
appG.InvalidParams()
```

响应示例：

```json
{ "code": 200, "msg": "ok", "data": {} }
```

### 错误码

- **200**: 成功
- **400**: 参数错误
- **401**: 签名错误
- **2xxxx**: 认证相关错误
- **4xxxx**: 用户相关错误
- **5xxxx**: 数据库相关错误

### 数据库与迁移

- 表前缀: `ay_`
- 模型放置: `internal/models/*`
- 建议: 设置 `CreatedAt/UpdatedAt/DeletedAt`；必要索引与唯一约束
- 租户字段与范围查询参见多租户文档

### Redis 规范

- 键前缀: `justus:`
- 建议按模块与租户分层: `justus:{tenant}:{module}:{biz}:{id}`

### 日志

- 位置: `runtime/logs/`
- 组件: `internal/infrastructure/logger.go`、`pkg/logger/*`、`pkg/logging/*`
- 建议注入字段: `request_id, tenant_id, user_id/admin_id, path, method`

### 国际化

- 语言文件: `pkg/langefile/*.toml`
- 建议：集中封装获取文案的入口，控制器/服务层按需使用

### 定时任务

- 定义入口: `internal/cron/cron.go`
- 建议：任务中显式设置租户上下文或执行为“平台级任务”

### 文件上传与工具

- 上传: `pkg/upload/image.go`
- 工具: `pkg/util/{jwt,md5,pagination,password,util}.go`
- 二维码: `pkg/qrcode/qrcode.go`
- FFmpeg: `pkg/ffmpeg/*`（可选）

### 新增接口流程（范式）

1. 模型：`internal/models/{model}.go`
2. 仓储：`internal/repository/{model}_repository.go`
3. 服务：`internal/service/{model}_service_impl.go`
4. 控制器：`internal/controllers/{api|admin}/{xxx}_controller.go`
5. 路由注册：`internal/routers/{api|admin}/*.go`
6. 中间件：确保签名/JWT/RBAC 按需启用
7. 缓存/日志/国际化：按规范接入
8. 错误码与响应：统一出口 `app.Gin`

### 调试技巧

- API 开发: 头部 `skip-signature: true`
- JWT: 解码 token 检查 claims（含 `tenant_id`）
- Redis: 使用 `justus:` 前缀快速定位
- 健康检查: `curl http://localhost:8787/health`
