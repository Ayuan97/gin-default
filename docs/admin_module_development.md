## 后台功能模块开发指南（Admin Module Development）

本文面向在 Admin 侧（`/admin/v1/*`）新增业务模块的标准流程，覆盖模型、仓储、服务、控制器、路由、RBAC、多租户、缓存与审计日志等关键点，确保快速且一致地落地新能力。

### 1. 适用范围与目标

- 面向管理员端的 REST API（前缀：`/admin/v1/`）
- 默认启用 JWT + RBAC 权限控制
- 统一响应、统一错误码、多租户隔离、可观测（日志/审计）

### 2. 新建模块 Checklist（一次性通看）

- 模型：在 `internal/models` 定义结构体与索引（含 `TenantID`）
- 迁移：确保新表结构纳入迁移流程
- 仓储：在 `internal/repository` 封装数据访问（统一租户 Scope）
- 服务：在 `internal/service` 实现业务逻辑、事务、缓存策略
- 控制器：在 `internal/controllers/admin` 编写入参校验、调用服务、统一响应
- 路由：在 `internal/routers` 注册 REST 路由并挂接中间件链
- 权限与菜单：定义权限码，关联菜单，角色赋权
- 缓存：键包含租户维度，读写与失效策略明确
- 审计：重要操作记录管理员操作日志

### 3. 模型（`internal/models`）

- 规范：除平台级表外必须包含 `TenantID`，并建立必要索引/唯一约束

```go
package models

import "gorm.io/gorm"

type Project struct {
    gorm.Model
    TenantID uint   `gorm:"index;not null"`
    Name     string `gorm:"type:varchar(128);index"`
    Code     string `gorm:"type:varchar(64);uniqueIndex:uk_tenant_code"`
    Status   int    `gorm:"default:0"`
}
```

### 4. 仓储（`internal/repository`）

- 职责：只处理数据访问；所有查询统一应用多租户 Scope；提供事务友好方法

```go
package repository

import (
    "github.com/gin-gonic/gin"
    "justus-go/internal/models"
    "gorm.io/gorm"
)

type ProjectRepository struct { DB *gorm.DB }

func (r *ProjectRepository) List(c *gin.Context, page, size int) (items []models.Project, total int64, err error) {
    db := r.DB.Scopes(models.WithTenant(c))
    if err = db.Model(&models.Project{}).Count(&total).Error; err != nil { return }
    err = db.Order("id desc").Offset((page-1)*size).Limit(size).Find(&items).Error
    return
}

func (r *ProjectRepository) Get(c *gin.Context, id uint) (*models.Project, error) {
    var m models.Project
    if err := r.DB.Scopes(models.WithTenant(c)).First(&m, id).Error; err != nil { return nil, err }
    return &m, nil
}
```

### 5. 服务（`internal/service`）

- 职责：业务编排、参数与状态校验、事务、缓存策略统一

```go
package service

import (
    "github.com/gin-gonic/gin"
    "justus-go/internal/models"
    "justus-go/internal/repository"
    "gorm.io/gorm"
)

type ProjectService struct {
    DB   *gorm.DB
    Repo *repository.ProjectRepository
}

func (s *ProjectService) Create(c *gin.Context, p *models.Project) error {
    p.TenantID = models.GetTenantIDFromContext(c)
    return s.DB.Transaction(func(tx *gorm.DB) error {
        return tx.Create(p).Error
    })
}
```

### 6. 控制器（`internal/controllers/admin`）

- 职责：入参绑定与校验、调用服务、使用 `pkg/app` 统一响应

```go
package admin

import (
    "github.com/gin-gonic/gin"
    "justus-go/internal/models"
    "justus-go/internal/service"
    "justus-go/pkg/app"
    "justus-go/pkg/e"
)

type ProjectController struct { Svc *service.ProjectService }

type createProjectReq struct {
    Name string `json:"name" binding:"required,min=2,max=128"`
    Code string `json:"code" binding:"required,alphanum,max=64"`
}

func (ctl *ProjectController) Create(c *gin.Context) {
    appG := app.Gin{C: c}
    var req createProjectReq
    if err := c.ShouldBindJSON(&req); err != nil { appG.InvalidParams(); return }
    if err := ctl.Svc.Create(c, &models.Project{Name: req.Name, Code: req.Code}); err != nil {
        appG.Error(e.ERROR)
        return
    }
    appG.Success(nil)
}
```

### 7. 路由注册（`internal/routers`）

- 中间件顺序建议：`tenant.Resolve()` → `jwt.Middleware()` → `admin.Auth()`

```go
// 伪代码：在 internal/routers/admin/*.go 中
admin := r.Group("/admin/v1")
admin.Use(tenant.Resolve(), jwt.Middleware(), adminmw.Auth())
{
    project := admin.Group("/projects")
    project.POST("", projectCtl.Create)
    project.GET("", projectCtl.List)
    project.GET(":id", projectCtl.Get)
    project.PUT(":id", projectCtl.Update)
    project.DELETE(":id", projectCtl.Delete)
}
```

### 8. RBAC 权限与菜单

- 权限码命名：`project:list|get|create|update|delete|export` 等
- 将路由与权限码在菜单管理中关联，角色赋权后生效
- 在 `admin.Auth()` 中检查权限（基于角色-权限映射）

### 9. 多租户接入

- 解析：`internal/middleware/tenant/resolve.go` 注入 `tenant_id` 到上下文
- 查询：统一使用 `models.WithTenant(c)` Scope
- 写入：显式设置 `TenantID`（来自上下文）
- 缓存：Key 必含租户，如 `justus:{tenant}:{module}:{biz}:{id}`

### 10. 缓存策略（按需）

- 查询缓存：Key 包含租户与分页/过滤要素，设置合理 TTL
- 写操作：删除或更新相关缓存 Key（建议在服务层集中封装）

### 11. 审计与日志

- 重要操作记录管理员操作日志（操作者、租户、资源、动作、对象 ID）
- 日志字段：`request_id, tenant_id, admin_id, module, action`
- 错误需带上下文与堆栈（由 Recover 捕获）

### 12. 错误码与统一响应

- 使用 `pkg/app`：`Success / Error / InvalidParams`
- 错误码集中在 `pkg/e` 维护，按模块细化

### 13. 验证、分页与返回结构

- 参数校验：`ShouldBindJSON` / `ShouldBindQuery` + `binding` 标签
- 分页参数：`page`、`page_size`；排序：`sort_by`、`order`
- 返回分页：`items + total`，由控制器统一封装

### 14. 最小落地清单（Checklist）

- 模型含 `TenantID`、索引完善；迁移已执行
- Repository 查询均使用 `WithTenant`；事务边界在 Service
- Controller 入参校验完备；统一响应
- 路由按 REST 注册并挂接 `tenant.resolve`、`jwt`、`admin.Auth()`
- 权限码定义并生效；菜单与权限绑定
- 缓存 Key 含租户维度；写操作正确失效缓存
- 审计日志记录关键操作

### 15. 调试与排障

- 令牌：使用管理员登录获取 `Authorization: Bearer <token>`
- 多租户：请求头 `X-Tenant-ID: <tenant>`（若使用该策略）
- 日志：查看 `runtime/logs/`，筛选 `tenant_id`、`request_id`
- 数据：Redis 以 `justus:` 前缀定位，数据库检查 `tenant_id`
