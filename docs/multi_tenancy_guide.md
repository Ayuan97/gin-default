## 多租户文档（Multi-Tenancy Guide）

### 目标与模式

- **目标**: 数据与权限按租户隔离，平台具备跨租户运维能力
- **模式**: 单库单表 + `tenant_id` 字段的逻辑隔离（模型与 Scope 支持）

### 关键组件

- **模型**: `internal/models/tenant.go`（租户定义），业务模型需包含 `TenantID`
- **Scope**: `internal/models/scope.go`（GORM Scope，按租户过滤）
- **中间件**: `internal/middleware/tenant/resolve.go`（解析并注入租户上下文）
- **权限**: `internal/middleware/admin/auth.go`（RBAC，控制跨租户能力）

### 租户解析（建议策略）

优先级（可在 `resolve.go` 中实现/拓展）：

1. 请求头 `X-Tenant-ID`
2. 子域名/路径（如 `xxx.example.com` → tenant `xxx`）
3. JWT Claims 中的 `tenant_id`

解析成功后放入 Gin Context，例如：`c.Set("tenant_id", id)`。

### GORM 租户隔离（Scope）

- 查询时统一应用 Scope：

```go
db.Scopes(models.WithTenant(c)).Find(&list)
```

- 创建/更新需写入 `TenantID`：

```go
entity.TenantID = ctxTenantID
db.Create(&entity)
```

- 事务内同样需要携带 Scope：

```go
err := db.Transaction(func(tx *gorm.DB) error {
    return tx.Scopes(models.WithTenant(c)).Save(&entity).Error
})
```

### 缓存与日志的租户维度

- **Redis Key**: 必须带租户，如 `justus:{tenant}:{module}:{biz}:{id}`
- **日志字段**: 必须包含 `tenant_id`
- **指标/审计**: 同样记录 `tenant_id`

### RBAC 与租户边界

- **默认**: 角色/权限在租户内生效
- **平台/超级管理员**: 具备跨租户能力（需专门权限）
- **跨租户操作**: 必须显式指定目标租户（建议头：`X-Cross-Tenant-ID`）并校验权限

### 平台级与租户级资源

- **平台级（无 `tenant_id`）**: 系统配置、全局字典、任务定义等
- **租户级（有 `tenant_id`）**: 用户、角色、菜单、业务数据等
- **约定**: 除平台级表外，业务表须包含 `tenant_id` 并纳入 Scope

### 多租户与路由/中间件

- API 与 Admin 均应在路由链路前段完成租户解析（在 JWT/RBAC 之前）
- 重要操作（如跨租户导入/迁移）需记录审计日志

### 后台任务与多租户

- **租户级任务**: 枚举租户逐个执行并注入上下文
- **平台级任务**: 禁止使用租户 Scope

### 常见问题

- 忘记应用 Scope：导致越权数据读写（务必在 Repository 统一注入）
- Goroutine 丢失上下文：启动新协程时显式传递 `tenant_id`
- 缓存 Key 未带租户：跨租户数据污染（按规范修正）
- 索引缺失：`tenant_id` + 业务维度建立联合索引
- 数据迁移：新增租户时初始化必要租户级数据（角色、菜单、默认配置）

### 测试与调试

- 本地调试：请求头设置 `X-Tenant-ID: <tenant>`；API 可加 `skip-signature: true`
- JWT 调试：解码查看 `tenant_id` 是否正确携带
- Redis：按前缀 `justus:` + 租户定位缓存
- 日志：过滤字段 `tenant_id` 定位问题

### 落地清单（新增业务检查）

- 模型含 `TenantID` 且建索引
- Repository 统一 `WithTenant` Scope
- 控制器/服务写入 `TenantID`
- 缓存 Key 按规范带租户
- 审计日志记录 `tenant_id`
- 跨租户能力有明确权限点与审计
