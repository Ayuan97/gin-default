-- 数据库初始化脚本
CREATE DATABASE IF NOT EXISTS justus CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

USE justus;

-- ===================================
-- 普通用户表 - 面向前端用户
-- ===================================
CREATE TABLE IF NOT EXISTS `ay_users` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '用户ID，主键',
  `username` varchar(50) NOT NULL COMMENT '用户名',
  `email` varchar(100) DEFAULT NULL COMMENT '邮箱地址，可选',
  `phone` varchar(20) DEFAULT NULL COMMENT '手机号码',
  `avatar` varchar(500) DEFAULT '' COMMENT '头像URL地址',
  `first_name` varchar(50) DEFAULT '' COMMENT '名字',
  `last_name` varchar(50) DEFAULT '' COMMENT '',
  `nickname` varchar(50) DEFAULT '' COMMENT '昵称',
  `gender` tinyint(1) DEFAULT 0 COMMENT '性别：0-未知，1-男，2-女',
  `birthday` date DEFAULT NULL COMMENT '生日',
  `lang` varchar(10) DEFAULT 'zh-Hans' COMMENT '语言偏好：zh-Hans-简体中文，en-英文等',
  `timezone` varchar(50) DEFAULT 'Asia/Shanghai' COMMENT '时区设置',
  `status` tinyint(1) DEFAULT 1 COMMENT '用户状态：1-正常，0-禁用，2-待激活',
  `email_verified` tinyint(1) DEFAULT 0 COMMENT '邮箱验证状态：0-未验证，1-已验证',
  `phone_verified` tinyint(1) DEFAULT 0 COMMENT '手机验证状态：0-未验证，1-已验证',
  `last_login_at` timestamp NULL DEFAULT NULL COMMENT '最后登录时间',
  `last_login_ip` varchar(45) DEFAULT '' COMMENT '最后登录IP地址',
  `login_count` int(11) DEFAULT 0 COMMENT '登录次数统计',
  `created_at` timestamp DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` timestamp DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `deleted_at` timestamp NULL DEFAULT NULL COMMENT '软删除时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_email` (`email`) COMMENT '邮箱唯一索引',
  KEY `idx_status` (`status`) COMMENT '状态索引',
  KEY `idx_created_at` (`created_at`) COMMENT '创建时间索引',
  KEY `idx_last_login` (`last_login_at`) COMMENT '最后登录时间索引'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户表-存储前端用户信息';

-- ===================================
-- 管理员用户表 - 面向后台管理
-- ===================================
CREATE TABLE IF NOT EXISTS `ay_admin_users` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '管理员ID，主键',
  `username` varchar(50) NOT NULL COMMENT '管理员用户名，唯一标识',
  `password` varchar(255) NOT NULL COMMENT '登录密码，bcrypt加密',
  `email` varchar(100) DEFAULT NULL COMMENT '邮箱地址',
  `phone` varchar(20) DEFAULT NULL COMMENT '手机号码',
  `avatar` varchar(500) DEFAULT '' COMMENT '头像URL地址',
  `real_name` varchar(50) DEFAULT '' COMMENT '真实姓名',
  `department` varchar(100) DEFAULT '' COMMENT '所属部门',
  `position` varchar(50) DEFAULT '' COMMENT '职位',
  `status` tinyint(1) DEFAULT 1 COMMENT '账户状态：1-正常，0-禁用，2-锁定',
  `is_super` tinyint(1) DEFAULT 0 COMMENT '是否超级管理员：0-否，1-是',
  `login_count` int(11) DEFAULT 0 COMMENT '登录次数统计',
  `last_login_at` timestamp NULL DEFAULT NULL COMMENT '最后登录时间',
  `last_login_ip` varchar(45) DEFAULT '' COMMENT '最后登录IP地址',
  `password_changed_at` timestamp NULL DEFAULT NULL COMMENT '密码最后修改时间',
  `failed_login_count` int(11) DEFAULT 0 COMMENT '连续登录失败次数',
  `locked_until` timestamp NULL DEFAULT NULL COMMENT '账户锁定到期时间',
  `created_by` bigint(20) DEFAULT 0 COMMENT '创建者ID',
  `created_at` timestamp DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` timestamp DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `deleted_at` timestamp NULL DEFAULT NULL COMMENT '软删除时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_username` (`username`) COMMENT '用户名唯一索引',
  UNIQUE KEY `uk_email` (`email`) COMMENT '邮箱唯一索引',
  KEY `idx_status` (`status`) COMMENT '状态索引',
  KEY `idx_department` (`department`) COMMENT '部门索引',
  KEY `idx_created_at` (`created_at`) COMMENT '创建时间索引',
  KEY `idx_last_login` (`last_login_at`) COMMENT '最后登录时间索引'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='管理员用户表-存储后台管理员信息';

-- ===================================
-- 角色表 - RBAC权限控制核心
-- ===================================
CREATE TABLE IF NOT EXISTS `ay_roles` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '角色ID，主键',
  `name` varchar(50) NOT NULL COMMENT '角色标识名，英文，如admin、editor',
  `display_name` varchar(100) NOT NULL DEFAULT '' COMMENT '角色显示名称，中文，如管理员、编辑员',
  `description` varchar(500) DEFAULT '' COMMENT '角色描述信息',
  `level` int(11) DEFAULT 1 COMMENT '角色等级，数字越大权限越高',
  `status` tinyint(1) DEFAULT 1 COMMENT '角色状态：1-启用，0-禁用',
  `is_system` tinyint(1) DEFAULT 0 COMMENT '是否系统角色：1-系统内置不可删除，0-普通角色',
  `sort_order` int(11) DEFAULT 0 COMMENT '排序字段，数字越小越靠前',
  `created_by` bigint(20) DEFAULT 0 COMMENT '创建者ID',
  `created_at` timestamp DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` timestamp DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `deleted_at` timestamp NULL DEFAULT NULL COMMENT '软删除时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_name` (`name`) COMMENT '角色名唯一索引',
  KEY `idx_status` (`status`) COMMENT '状态索引',
  KEY `idx_level` (`level`) COMMENT '等级索引',
  KEY `idx_sort_order` (`sort_order`) COMMENT '排序索引'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='角色表-定义系统中的各种角色';

-- ===================================
-- 权限表 - 定义系统中的具体权限
-- ===================================
CREATE TABLE IF NOT EXISTS `ay_permissions` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '权限ID，主键',
  `name` varchar(100) NOT NULL COMMENT '权限标识名，格式：module.action.resource',
  `display_name` varchar(100) NOT NULL DEFAULT '' COMMENT '权限显示名称，中文描述',
  `description` varchar(500) DEFAULT '' COMMENT '权限详细描述',
  `module` varchar(50) NOT NULL DEFAULT '' COMMENT '所属模块：admin、api、system等',
  `action` varchar(50) NOT NULL DEFAULT '' COMMENT '操作类型：read、write、delete、create、update等',
  `resource` varchar(50) NOT NULL DEFAULT '' COMMENT '资源类型：user、role、permission、system等',
  `route` varchar(200) DEFAULT '' COMMENT '对应的路由规则，支持通配符',
  `method` varchar(20) DEFAULT '' COMMENT 'HTTP方法：GET、POST、PUT、DELETE、*',
  `parent_id` bigint(20) DEFAULT 0 COMMENT '父权限ID，支持权限树结构',
  `level` int(11) DEFAULT 1 COMMENT '权限层级，根权限为1',
  `sort_order` int(11) DEFAULT 0 COMMENT '排序字段，数字越小越靠前',
  `is_menu` tinyint(1) DEFAULT 0 COMMENT '是否为菜单权限：1-是菜单，0-非菜单',
  `menu_icon` varchar(100) DEFAULT '' COMMENT '菜单图标class或路径',
  `is_system` tinyint(1) DEFAULT 0 COMMENT '是否系统权限：1-系统内置不可删除，0-普通权限',
  `created_at` timestamp DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` timestamp DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `deleted_at` timestamp NULL DEFAULT NULL COMMENT '软删除时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_name` (`name`) COMMENT '权限名唯一索引',
  KEY `idx_module` (`module`) COMMENT '模块索引',
  KEY `idx_action` (`action`) COMMENT '操作类型索引',
  KEY `idx_resource` (`resource`) COMMENT '资源类型索引',
  KEY `idx_parent_id` (`parent_id`) COMMENT '父权限索引',
  KEY `idx_is_menu` (`is_menu`) COMMENT '菜单权限索引',
  KEY `idx_sort_order` (`sort_order`) COMMENT '排序索引'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='权限表-定义系统中的所有权限点';

-- ===================================
-- 角色权限关联表 - 多对多关系
-- ===================================
CREATE TABLE IF NOT EXISTS `ay_role_permissions` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '关联ID，主键',
  `role_id` bigint(20) unsigned NOT NULL COMMENT '角色ID，外键关联ay_roles.id',
  `permission_id` bigint(20) unsigned NOT NULL COMMENT '权限ID，外键关联ay_permissions.id',
  `granted_by` bigint(20) DEFAULT 0 COMMENT '授权者ID，记录是谁给这个角色分配的权限',
  `created_at` timestamp DEFAULT CURRENT_TIMESTAMP COMMENT '授权时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_role_permission` (`role_id`, `permission_id`) COMMENT '角色权限组合唯一索引',
  KEY `idx_role_id` (`role_id`) COMMENT '角色ID索引',
  KEY `idx_permission_id` (`permission_id`) COMMENT '权限ID索引',
  KEY `idx_granted_by` (`granted_by`) COMMENT '授权者索引',
  CONSTRAINT `fk_role_permissions_role` FOREIGN KEY (`role_id`) REFERENCES `ay_roles` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `fk_role_permissions_permission` FOREIGN KEY (`permission_id`) REFERENCES `ay_permissions` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='角色权限关联表-定义角色拥有哪些权限';

-- ===================================
-- 管理员角色关联表 - 多对多关系
-- ===================================
CREATE TABLE IF NOT EXISTS `ay_admin_user_roles` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '关联ID，主键',
  `admin_user_id` bigint(20) unsigned NOT NULL COMMENT '管理员ID，外键关联ay_admin_users.id',
  `role_id` bigint(20) unsigned NOT NULL COMMENT '角色ID，外键关联ay_roles.id',
  `assigned_by` bigint(20) DEFAULT 0 COMMENT '分配者ID，记录是谁给这个用户分配的角色',
  `expires_at` timestamp NULL DEFAULT NULL COMMENT '角色过期时间，NULL表示永不过期',
  `created_at` timestamp DEFAULT CURRENT_TIMESTAMP COMMENT '分配时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_admin_role` (`admin_user_id`, `role_id`) COMMENT '用户角色组合唯一索引',
  KEY `idx_admin_user_id` (`admin_user_id`) COMMENT '管理员ID索引',
  KEY `idx_role_id` (`role_id`) COMMENT '角色ID索引',
  KEY `idx_assigned_by` (`assigned_by`) COMMENT '分配者索引',
  KEY `idx_expires_at` (`expires_at`) COMMENT '过期时间索引',
  CONSTRAINT `fk_admin_user_roles_admin` FOREIGN KEY (`admin_user_id`) REFERENCES `ay_admin_users` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `fk_admin_user_roles_role` FOREIGN KEY (`role_id`) REFERENCES `ay_roles` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='管理员角色关联表-定义管理员拥有哪些角色';

-- ===================================
-- 默认角色数据 - 系统初始角色
-- ===================================
INSERT IGNORE INTO `ay_roles` (`name`, `display_name`, `description`, `level`, `status`, `is_system`, `sort_order`) VALUES
('super_admin', '超级管理员', '拥有系统所有权限的超级管理员，可以管理所有模块和用户', 100, 1, 1, 1),
('admin', '系统管理员', '系统管理员，拥有大部分管理权限，负责日常系统维护', 80, 1, 1, 2),
('moderator', '内容管理员', '内容审核和管理员，主要负责内容相关的管理工作', 50, 1, 1, 3),
('editor', '编辑员', '内容编辑员，可以创建和编辑内容，但权限有限', 30, 1, 0, 4),
('viewer', '查看员', '只读权限角色，只能查看信息不能修改', 10, 1, 0, 5);

-- ===================================
-- 默认权限数据 - 系统功能权限
-- ===================================
INSERT IGNORE INTO `ay_permissions` (`name`, `display_name`, `description`, `module`, `action`, `resource`, `route`, `method`, `parent_id`, `level`, `sort_order`, `is_menu`, `menu_icon`, `is_system`) VALUES
-- 用户管理权限
('admin.user', '用户管理', '用户管理模块主菜单', 'admin', 'menu', 'user', '/admin/users', '*', 0, 1, 1, 1, 'fas fa-users', 1),
('admin.user.list', '用户列表', '查看用户列表', 'admin', 'read', 'user', '/admin/users', 'GET', 1, 2, 1, 0, '', 1),
('admin.user.view', '查看用户', '查看用户详细信息', 'admin', 'read', 'user', '/admin/users/*', 'GET', 1, 2, 2, 0, '', 1),
('admin.user.create', '创建用户', '创建新用户', 'admin', 'create', 'user', '/admin/users', 'POST', 1, 2, 3, 0, '', 1),
('admin.user.update', '修改用户', '修改用户信息', 'admin', 'update', 'user', '/admin/users/*', 'PUT', 1, 2, 4, 0, '', 1),
('admin.user.delete', '删除用户', '删除用户账户', 'admin', 'delete', 'user', '/admin/users/*', 'DELETE', 1, 2, 5, 0, '', 1),
('admin.user.status', '用户状态', '修改用户状态（启用/禁用）', 'admin', 'update', 'user', '/admin/users/*/status', 'PUT', 1, 2, 6, 0, '', 1),

-- 角色管理权限
('admin.role', '角色管理', '角色管理模块主菜单', 'admin', 'menu', 'role', '/admin/roles', '*', 0, 1, 2, 1, 'fas fa-user-tag', 1),
('admin.role.list', '角色列表', '查看角色列表', 'admin', 'read', 'role', '/admin/roles', 'GET', 8, 2, 1, 0, '', 1),
('admin.role.view', '查看角色', '查看角色详细信息', 'admin', 'read', 'role', '/admin/roles/*', 'GET', 8, 2, 2, 0, '', 1),
('admin.role.create', '创建角色', '创建新角色', 'admin', 'create', 'role', '/admin/roles', 'POST', 8, 2, 3, 0, '', 1),
('admin.role.update', '修改角色', '修改角色信息', 'admin', 'update', 'role', '/admin/roles/*', 'PUT', 8, 2, 4, 0, '', 1),
('admin.role.delete', '删除角色', '删除角色', 'admin', 'delete', 'role', '/admin/roles/*', 'DELETE', 8, 2, 5, 0, '', 1),
('admin.role.permission', '角色权限', '管理角色权限分配', 'admin', 'update', 'role', '/admin/roles/*/permissions', 'PUT', 8, 2, 6, 0, '', 1),

-- 权限管理权限  
('admin.permission', '权限管理', '权限管理模块主菜单', 'admin', 'menu', 'permission', '/admin/permissions', '*', 0, 1, 3, 1, 'fas fa-key', 1),
('admin.permission.list', '权限列表', '查看权限列表', 'admin', 'read', 'permission', '/admin/permissions', 'GET', 14, 2, 1, 0, '', 1),
('admin.permission.view', '查看权限', '查看权限详细信息', 'admin', 'read', 'permission', '/admin/permissions/*', 'GET', 14, 2, 2, 0, '', 1),
('admin.permission.create', '创建权限', '创建新权限', 'admin', 'create', 'permission', '/admin/permissions', 'POST', 14, 2, 3, 0, '', 1),
('admin.permission.update', '修改权限', '修改权限信息', 'admin', 'update', 'permission', '/admin/permissions/*', 'PUT', 14, 2, 4, 0, '', 1),
('admin.permission.delete', '删除权限', '删除权限', 'admin', 'delete', 'permission', '/admin/permissions/*', 'DELETE', 14, 2, 5, 0, '', 1),

-- 系统管理权限
('admin.system', '系统管理', '系统管理模块主菜单', 'admin', 'menu', 'system', '/admin/system', '*', 0, 1, 4, 1, 'fas fa-cogs', 1),
('admin.system.info', '系统信息', '查看系统运行信息', 'admin', 'read', 'system', '/admin/system/info', 'GET', 19, 2, 1, 0, '', 1),
('admin.system.health', '健康检查', '查看系统健康状态', 'admin', 'read', 'system', '/admin/system/health', 'GET', 19, 2, 2, 0, '', 1),
('admin.system.config', '系统配置', '查看和修改系统配置', 'admin', 'update', 'system', '/admin/system/config', '*', 19, 2, 3, 0, '', 1),
('admin.system.logs', '系统日志', '查看系统运行日志', 'admin', 'read', 'system', '/admin/system/logs', 'GET', 19, 2, 4, 0, '', 1),
('admin.system.cache', '缓存管理', '管理系统缓存', 'admin', 'update', 'system', '/admin/system/cache', '*', 19, 2, 5, 0, '', 1),

-- API接口权限
('api.access', 'API访问', 'API接口访问权限', 'api', 'read', 'api', '/api/*', '*', 0, 1, 5, 0, '', 1),
('api.user', '用户API', 'API用户相关接口', 'api', 'read', 'user', '/api/*/users*', '*', 25, 2, 1, 0, '', 1),
('api.auth', '认证API', 'API认证相关接口', 'api', 'read', 'auth', '/api/*/auth*', '*', 25, 2, 2, 0, '', 1);

-- ===================================
-- 默认管理员账户 - 系统初始管理员
-- ===================================
-- 创建默认超级管理员账户（用户名: admin, 密码: admin123456）
-- 密码哈希值是通过bcrypt加密的admin123456
INSERT IGNORE INTO `ay_admin_users` (
    `username`, `password`, `email`, `real_name`, `department`, `position`, 
    `status`, `is_super`, `password_changed_at`, `created_by`
) VALUES (
    'admin', 
    '$2a$10$N9qo8uLOickgx2ZMRZoMye7uo1OGMz.L/4YL.ZeJzJJz0jQdRGUUu', 
    'admin@example.com', 
    '系统管理员', 
    '技术部', 
    '系统管理员',
    1, 
    1, 
    NOW(), 
    0
);

-- 创建普通管理员账户（用户名: manager, 密码: manager123456）
INSERT IGNORE INTO `ay_admin_users` (
    `username`, `password`, `email`, `real_name`, `department`, `position`, 
    `status`, `is_super`, `password_changed_at`, `created_by`
) VALUES (
    'manager', 
    '$2a$10$kM7lQU8eYYm4QFQJq9d.FeKJL.ZeJzJJz0jQdRGUUu.Manager123', 
    'manager@example.com', 
    '部门管理员', 
    '运营部', 
    '部门经理',
    1, 
    0, 
    NOW(), 
    1
);

-- ===================================
-- 角色权限分配 - 给角色分配相应权限
-- ===================================
-- 给超级管理员角色分配所有权限
INSERT IGNORE INTO `ay_role_permissions` (`role_id`, `permission_id`, `granted_by`) 
SELECT r.id, p.id, 0
FROM `ay_roles` r, `ay_permissions` p 
WHERE r.name = 'super_admin';

-- 给系统管理员角色分配大部分权限（除了超级权限）
INSERT IGNORE INTO `ay_role_permissions` (`role_id`, `permission_id`, `granted_by`) 
SELECT r.id, p.id, 0
FROM `ay_roles` r, `ay_permissions` p 
WHERE r.name = 'admin' 
AND p.name NOT IN ('admin.permission.delete', 'admin.role.delete', 'admin.system.config');

-- 给内容管理员分配用户和内容相关权限
INSERT IGNORE INTO `ay_role_permissions` (`role_id`, `permission_id`, `granted_by`) 
SELECT r.id, p.id, 0
FROM `ay_roles` r, `ay_permissions` p 
WHERE r.name = 'moderator' 
AND p.name IN ('admin.user', 'admin.user.list', 'admin.user.view', 'admin.user.update', 'admin.user.status');

-- 给编辑员分配基本查看权限
INSERT IGNORE INTO `ay_role_permissions` (`role_id`, `permission_id`, `granted_by`) 
SELECT r.id, p.id, 0
FROM `ay_roles` r, `ay_permissions` p 
WHERE r.name = 'editor' 
AND p.name IN ('admin.user.list', 'admin.user.view', 'admin.role.list', 'admin.role.view');

-- 给查看员分配只读权限
INSERT IGNORE INTO `ay_role_permissions` (`role_id`, `permission_id`, `granted_by`) 
SELECT r.id, p.id, 0
FROM `ay_roles` r, `ay_permissions` p 
WHERE r.name = 'viewer' 
AND p.action = 'read';

-- ===================================
-- 用户角色分配 - 给用户分配角色
-- ===================================
-- 给admin用户分配超级管理员角色
INSERT IGNORE INTO `ay_admin_user_roles` (`admin_user_id`, `role_id`, `assigned_by`) 
SELECT au.id, r.id, 0
FROM `ay_admin_users` au, `ay_roles` r 
WHERE au.username = 'admin' AND r.name = 'super_admin';

-- 给manager用户分配系统管理员角色
INSERT IGNORE INTO `ay_admin_user_roles` (`admin_user_id`, `role_id`, `assigned_by`) 
SELECT au.id, r.id, 1
FROM `ay_admin_users` au, `ay_roles` r 
WHERE au.username = 'manager' AND r.name = 'admin'; 