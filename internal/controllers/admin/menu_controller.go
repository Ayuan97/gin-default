package admin

import (
	"encoding/json"
	"net/http"
	"sort"
	"strconv"
	"time"

	"justus/internal/container"
	"justus/internal/models"
	"justus/pkg/app"
	"justus/pkg/e"
	"justus/pkg/rediskey"

	"github.com/gin-gonic/gin"
)

// MenuController 菜单控制器（租户感知）
type MenuController struct {
	logger container.Logger
	cache  container.Cache
}

func NewMenuController(logger container.Logger, cache container.Cache) *MenuController {
	return &MenuController{logger: logger, cache: cache}
}

// GetMyMenus 返回当前用户在当前租户下可见的菜单树
func (mc *MenuController) GetMyMenus(c *gin.Context) {
	appG := app.Gin{C: c}

	tenantIdVal, tenantExists := c.Get("tenantId")
	userIdVal, userExists := c.Get("userId")
	if !tenantExists || !userExists {
		appG.Unauthorized(e.ERROR_AUTH)
		return
	}
	tenantID, _ := tenantIdVal.(int)
	userID, _ := userIdVal.(int)

	// 先尝试读取用户维度菜单树缓存
	if mc.cache != nil {
		cacheKey := rediskey.TenantUserMenuTreeKey(uint(tenantID), uint(userID))
		if raw := mc.cache.Get(cacheKey); raw != "" {
			var cached struct {
				Menus []gin.H `json:"menus"`
			}
			if err := json.Unmarshal([]byte(raw), &cached); err == nil {
				appG.Success(gin.H{"menus": cached.Menus, "tenant_id": tenantID, "user_id": userID, "cached": true})
				return
			}
		}
	}

	// 读取租户白名单（如果没有配置，默认空=无菜单；可按需调整兼容策略）
	var whiteIDs []uint
	// 先读缓存
	if mc.cache != nil {
		cacheKey := rediskey.TenantMenuWhitelistKey(uint(tenantID))
		if raw := mc.cache.Get(cacheKey); raw != "" {
			_ = json.Unmarshal([]byte(raw), &whiteIDs)
		}
	}
	if len(whiteIDs) == 0 {
		ids, err := models.GetTenantPermissionIDs(uint(tenantID))
		if err != nil {
			mc.logger.Errorf("get tenant permission ids error: %v", err)
			c.Status(http.StatusInternalServerError)
			return
		}
		whiteIDs = ids
		if mc.cache != nil {
			b, _ := json.Marshal(whiteIDs)
			_ = mc.cache.Set(rediskey.TenantMenuWhitelistKey(uint(tenantID)), string(b), 10*time.Minute)
		}
	}
	whiteSet := map[uint]struct{}{}
	for _, id := range whiteIDs {
		whiteSet[id] = struct{}{}
	}

	// 用户在当前租户的权限ID集合
	userPermIDs, err := models.GetAdminUserPermissionIDsInTenant(userID, uint(tenantID))
	if err != nil {
		mc.logger.Errorf("get user permission ids error: %v", err)
		c.Status(http.StatusInternalServerError)
		return
	}

	// 求交集并筛选菜单项
	var menuPerms []models.Permission
	if len(userPermIDs) == 0 || len(whiteSet) == 0 {
		appG.Success(gin.H{"menus": []gin.H{}})
		return
	}
	// 查询菜单权限详情（不缓存，体量较小；已统一使用 ay_permissions）
	if err := models.GetDb().
		Table("ay_permissions").
		Where("id IN ?", userPermIDs).
		Where("id IN ?", whiteIDs).
		Where("is_menu = 1").
		Order("sort_order ASC, id ASC").
		Find(&menuPerms).Error; err != nil {
		mc.logger.Errorf("query menu perms error: %v", err)
		c.Status(http.StatusInternalServerError)
		return
	}

	// 构建树（可选缓存：按租户+用户缓存菜单树，后续如需可开启）
	nodeMap := map[uint]gin.H{}
	var roots []gin.H
	for _, p := range menuPerms {
		node := gin.H{
			"id":           p.ID,
			"name":         p.Name,
			"display_name": p.DisplayName,
			"route":        p.Route,
			"menu_icon":    p.MenuIcon,
			"parent_id":    p.ParentID,
			"sort_order":   p.SortOrder,
			"children":     []gin.H{},
		}
		nodeMap[p.ID] = node
	}
	for _, p := range menuPerms {
		node := nodeMap[p.ID]
		if p.ParentID == 0 {
			roots = append(roots, node)
		} else if parent, ok := nodeMap[p.ParentID]; ok {
			children := parent["children"].([]gin.H)
			parent["children"] = append(children, node)
		} else {
			roots = append(roots, node)
		}
	}
	// 根节点按排序
	sort.Slice(roots, func(i, j int) bool {
		return roots[i]["sort_order"].(int) < roots[j]["sort_order"].(int)
	})

	// 写入用户维度菜单树缓存（5分钟）
	if mc.cache != nil {
		cacheKey := rediskey.TenantUserMenuTreeKey(uint(tenantID), uint(userID))
		payload, _ := json.Marshal(gin.H{"menus": roots})
		_ = mc.cache.Set(cacheKey, string(payload), 5*time.Minute)
	}

	appG.Success(gin.H{"menus": roots, "tenant_id": tenantID, "user_id": userID})
}

// GetMyMenusVben 返回 Vben 期望结构的菜单（后端访问控制对接）
func (mc *MenuController) GetMyMenusVben(c *gin.Context) {
	appG := app.Gin{C: c}

	// 直接复用已有菜单树，然后映射字段
	// 为避免重复实现，调用当前方法的核心逻辑，或简单复制并改字段名称
	tenantVal, tenantExists := c.Get("tenantId")
	userVal, userExists := c.Get("userId")
	if !tenantExists || !userExists {
		appG.Unauthorized(e.ERROR_AUTH)
		return
	}
	tenantID := uint(tenantVal.(int))
	userID := userVal.(int)

	// 读取租户白名单
	var whiteIDs []uint
	if mc.cache != nil {
		if raw := mc.cache.Get(rediskey.TenantMenuWhitelistKey(tenantID)); raw != "" {
			_ = json.Unmarshal([]byte(raw), &whiteIDs)
		}
	}
	if len(whiteIDs) == 0 {
		ids, _ := models.GetTenantPermissionIDs(tenantID)
		whiteIDs = ids
	}
	// 用户在租户的权限ID
	userPermIDs, _ := models.GetAdminUserPermissionIDsInTenant(userID, tenantID)
	if len(userPermIDs) == 0 || len(whiteIDs) == 0 {
		appG.Success([]gin.H{})
		return
	}
	// 查询菜单权限
	var menuPerms []models.Permission
	if err := models.GetDb().
		Table("ay_permissions").
		Where("id IN ?", userPermIDs).
		Where("id IN ?", whiteIDs).
		Where("is_menu = 1").
		Order("sort_order ASC, id ASC").
		Find(&menuPerms).Error; err != nil {
		appG.Error(e.ERROR_DATABASE_QUERY)
		return
	}
	// 构建 Vben 结构
	nodeMap := map[uint]gin.H{}
	var roots []gin.H
	for _, p := range menuPerms {
		node := gin.H{
			"name":      p.Name,
			"path":      p.Route,
			"component": p.Component,
			"meta": gin.H{
				"title": p.DisplayName,
				"icon":  p.MenuIcon,
				"order": p.SortOrder,
			},
			"children":  []gin.H{},
			"parent_id": p.ParentID,
		}
		nodeMap[p.ID] = node
	}
	for _, p := range menuPerms {
		node := nodeMap[p.ID]
		if p.ParentID == 0 {
			roots = append(roots, node)
		} else if parent, ok := nodeMap[p.ParentID]; ok {
			parent["children"] = append(parent["children"].([]gin.H), node)
		} else {
			roots = append(roots, node)
		}
	}
	appG.Success(roots)
}

// GetTenantMenus 获取指定租户允许的菜单（仅超级管理员）
func (mc *MenuController) GetTenantMenus(c *gin.Context) {
	appG := app.Gin{C: c}
	if isSuper, _ := c.Get("isSuper"); isSuper != true {
		appG.Error(e.ERROR_PERMISSION_DENIED)
		return
	}
	idStr := c.Param("id")
	tid, _ := strconv.Atoi(idStr)
	ids, err := models.GetTenantPermissionIDs(uint(tid))
	if err != nil {
		appG.Error(e.ERROR_DATABASE_QUERY)
		return
	}
	appG.Success(gin.H{"permission_ids": ids})
}

// UpdateTenantMenus 更新指定租户允许的菜单（仅超级管理员）
func (mc *MenuController) UpdateTenantMenus(c *gin.Context) {
	appG := app.Gin{C: c}
	if isSuper, _ := c.Get("isSuper"); isSuper != true {
		appG.Error(e.ERROR_PERMISSION_DENIED)
		return
	}
	idStr := c.Param("id")
	tid, _ := strconv.Atoi(idStr)
	var req struct {
		PermissionIDs []uint `json:"permission_ids" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		appG.InvalidParams()
		return
	}
	// 先清空后插入（幂等）
	tx := models.GetDb().Begin()
	if err := tx.Where("tenant_id = ?", tid).Delete(&models.TenantPermission{}).Error; err != nil {
		tx.Rollback()
		appG.Error(e.ERROR_DATABASE_QUERY)
		return
	}
	for _, pid := range req.PermissionIDs {
		tp := models.TenantPermission{TenantID: uint(tid), PermissionID: pid}
		if err := tx.Create(&tp).Error; err != nil {
			tx.Rollback()
			appG.Error(e.ERROR_DATABASE_QUERY)
			return
		}
	}
	if err := tx.Commit().Error; err != nil {
		appG.Error(e.ERROR_DATABASE_QUERY)
		return
	}
	// 清理该租户相关菜单缓存
	if mc.cache != nil {
		_, _ = mc.cache.Del(rediskey.TenantMenuWhitelistKey(uint(tid)))
		// 用户菜单树可按需批量失效，这里不做扫描清理，交由上层约定或设置短过期时间
	}
	appG.Success(gin.H{"message": "更新成功"})
}
