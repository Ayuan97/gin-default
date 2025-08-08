## 项目规范文档（Project Conventions）

### 命名与代码风格

- **命名**: 语义化、使用完整单词；函数为动词/动宾短语
- **控制流**: 早返回；优先处理边界与错误；避免深层嵌套
- **错误处理**: 包装错误并加上下文；对外统一错误码
- **注释**: 解释“为什么”；避免无意义注释；不留 TODO（直接实现或建任务）
- **格式**: `gofmt`；避免无关重排
- **导出类型**: 为结构体/接口编写注释与字段说明

### 目录与分层职责

- **Controller**: 参数校验、调用服务、统一响应
- **Service**: 业务编排、事务、缓存策略
- **Repository**: 数据访问（GORM 查询与事务）
- **Model**: 结构定义与 GORM 标签、索引
- **Middleware**: 鉴权、审计、租户解析、日志等横切关注点

### Git 与提交

- **分支命名**: `feature/*`、`fix/*`、`chore/*`、`refactor/*`、`perf/*`
- **提交信息（Conventional Commits）**: `feat:`、`fix:`、`docs:`、`refactor:`、`perf:`、`chore:`
- **PR 要求**: 描述问题、方案、影响面、回滚策略；小步提交

### API 规范（REST）

- **资源**: 使用复数，如 `/users`、`/roles`
- **方法语义**: GET 查询、POST 创建、PUT 全量、PATCH 部分、DELETE 删除
- **分页参数**: `page`、`page_size`；排序: `sort_by`、`order`
- **认证**:
  - API：签名（开发可用 `skip-signature: true`）+ 可选 JWT
  - Admin：`Authorization: Bearer <token>`（JWT + RBAC）
- **幂等**: 需要时对 POST 提供 `Idempotency-Key`
- **速率限制**: 放在 Common 中间件之后（如接入）

### 统一响应与错误码

- **统一返回**: `pkg/app/response.go`
- **错误码定义**: `pkg/e/{code.go,msg.go}`
- **禁止**: 控制器直接将底层错误返回给前端

### 安全与权限

- **Admin 路由**: 必须接入 `admin.Auth()` 中间件
- **审计**: 重要操作记录管理员操作日志
- **签名**: API 模块默认开启签名（开发可跳过）
- **JWT 密钥**: 仅从配置/环境变量读取，禁止硬编码
- **敏感信息**: 密码与密钥不以明文写入日志

### 配置与环境

- **密钥**: 禁止提交真实密钥；使用 `.env`/环境变量
- **环境分离**: `conf/app.dev.yaml` 与 `conf/app.production.yaml`
- **本地开发**: `make dev`；数据库初始化 `make db-init`

### 数据库规范

- **表前缀**: `ay_`；除平台级表外，业务表均含 `tenant_id`
- **索引**: 建立必要索引与唯一约束；避免跨租户唯一冲突
- **事务**: 以 Service 为边界；Repository 提供可组合方法
- **删除**: 优先软删除；必要时外键约束
- **迁移/种子**: 使用 `make` 系列命令

### Redis 规范

- **Key 模式**: `justus:{tenant}:{module}:{biz}:{id}`
- **TTL**: 依据业务设置；避免无 TTL 的无限增长集合
- **值格式**: 推荐 JSON；必要时压缩；避免超大对象
- **清理策略**: 模块级 Key 前缀统一管理

### 国际化规范

- **文案位置**: `pkg/langefile/*.toml`
- **Key 命名**: `domain.action.result`（如 `user.login.failed`）
- **回落**: 优先回落英文或简体中文

### 日志规范

- **级别**: Debug/Info/Warn/Error/Fatal
- **结构化字段**: `request_id, tenant_id, user_id/admin_id, module, action`
- **错误日志**: Recover 捕获堆栈并输出上下文

### CI/CD 与依赖

- **依赖管理**: `go mod tidy`；避免未使用依赖
- **构建检查**: 构建前执行静态检查与基础单测（如接入）
