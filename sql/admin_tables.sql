-- =============================================
-- 后台管理系统数据表结构
-- 创建时间: 2024
-- 描述: 支持手机号和账号登录的后台管理员用户系统
-- =============================================

-- 设置字符集
SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- =============================================
-- 1. 管理员用户表
-- =============================================
DROP TABLE IF EXISTS `admin_users`;
CREATE TABLE `admin_users` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `username` varchar(50) NOT NULL COMMENT '用户名/账号',
  `phone` varchar(20) DEFAULT NULL COMMENT '手机号',
  `email` varchar(100) DEFAULT NULL COMMENT '邮箱',
  `password` varchar(255) NOT NULL COMMENT '密码(bcrypt加密)',
  `real_name` varchar(50) DEFAULT NULL COMMENT '真实姓名',
  `nickname` varchar(50) DEFAULT NULL COMMENT '昵称',
  `avatar` varchar(255) DEFAULT NULL COMMENT '头像URL',
  `gender` tinyint(1) DEFAULT 0 COMMENT '性别: 0-未知, 1-男, 2-女',
  `birthday` date DEFAULT NULL COMMENT '生日',
  `department_id` bigint(20) unsigned DEFAULT NULL COMMENT '部门ID',
  `position` varchar(100) DEFAULT NULL COMMENT '职位',
  `status` tinyint(1) NOT NULL DEFAULT 1 COMMENT '状态: 0-禁用, 1-启用',
  `is_super_admin` tinyint(1) NOT NULL DEFAULT 0 COMMENT '是否超级管理员: 0-否, 1-是',
  `last_login_time` timestamp NULL DEFAULT NULL COMMENT '最后登录时间',
  `last_login_ip` varchar(45) DEFAULT NULL COMMENT '最后登录IP',
  `login_count` int(11) NOT NULL DEFAULT 0 COMMENT '登录次数',
  `password_changed_at` timestamp NULL DEFAULT NULL COMMENT '密码修改时间',
  `remark` text COMMENT '备注',
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `deleted_at` timestamp NULL DEFAULT NULL COMMENT '删除时间(软删除)',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_username` (`username`),
  UNIQUE KEY `uk_phone` (`phone`),
  UNIQUE KEY `uk_email` (`email`),
  KEY `idx_department_id` (`department_id`),
  KEY `idx_status` (`status`),
  KEY `idx_created_at` (`created_at`),
  KEY `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='管理员用户表';

-- =============================================
-- 2. 角色表
-- =============================================
DROP TABLE IF EXISTS `admin_roles`;
CREATE TABLE `admin_roles` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `name` varchar(50) NOT NULL COMMENT '角色名称',
  `code` varchar(50) NOT NULL COMMENT '角色编码',
  `description` text COMMENT '角色描述',
  `level` int(11) NOT NULL DEFAULT 1 COMMENT '角色级别(数字越小权限越高)',
  `status` tinyint(1) NOT NULL DEFAULT 1 COMMENT '状态: 0-禁用, 1-启用',
  `is_system` tinyint(1) NOT NULL DEFAULT 0 COMMENT '是否系统角色: 0-否, 1-是',
  `sort_order` int(11) NOT NULL DEFAULT 0 COMMENT '排序',
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `deleted_at` timestamp NULL DEFAULT NULL COMMENT '删除时间(软删除)',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_code` (`code`),
  KEY `idx_status` (`status`),
  KEY `idx_level` (`level`),
  KEY `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='管理员角色表';

-- =============================================
-- 3. 权限表
-- =============================================
DROP TABLE IF EXISTS `admin_permissions`;
CREATE TABLE `admin_permissions` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `parent_id` bigint(20) unsigned NOT NULL DEFAULT 0 COMMENT '父级权限ID',
  `name` varchar(100) NOT NULL COMMENT '权限名称',
  `code` varchar(100) NOT NULL COMMENT '权限编码',
  `type` tinyint(1) NOT NULL DEFAULT 1 COMMENT '权限类型: 1-菜单, 2-按钮, 3-接口',
  `path` varchar(255) DEFAULT NULL COMMENT '路由路径',
  `component` varchar(255) DEFAULT NULL COMMENT '组件路径',
  `icon` varchar(100) DEFAULT NULL COMMENT '图标',
  `method` varchar(10) DEFAULT NULL COMMENT 'HTTP方法(GET,POST,PUT,DELETE等)',
  `url` varchar(255) DEFAULT NULL COMMENT 'API接口地址',
  `level` int(11) NOT NULL DEFAULT 1 COMMENT '层级',
  `sort_order` int(11) NOT NULL DEFAULT 0 COMMENT '排序',
  `status` tinyint(1) NOT NULL DEFAULT 1 COMMENT '状态: 0-禁用, 1-启用',
  `is_hidden` tinyint(1) NOT NULL DEFAULT 0 COMMENT '是否隐藏: 0-否, 1-是',
  `description` text COMMENT '权限描述',
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `deleted_at` timestamp NULL DEFAULT NULL COMMENT '删除时间(软删除)',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_code` (`code`),
  KEY `idx_parent_id` (`parent_id`),
  KEY `idx_type` (`type`),
  KEY `idx_status` (`status`),
  KEY `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='管理员权限表';

-- =============================================
-- 4. 用户角色关联表
-- =============================================
DROP TABLE IF EXISTS `admin_user_roles`;
CREATE TABLE `admin_user_roles` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `user_id` bigint(20) unsigned NOT NULL COMMENT '用户ID',
  `role_id` bigint(20) unsigned NOT NULL COMMENT '角色ID',
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_user_role` (`user_id`, `role_id`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_role_id` (`role_id`),
  CONSTRAINT `fk_user_roles_user_id` FOREIGN KEY (`user_id`) REFERENCES `admin_users` (`id`) ON DELETE CASCADE,
  CONSTRAINT `fk_user_roles_role_id` FOREIGN KEY (`role_id`) REFERENCES `admin_roles` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户角色关联表';

-- =============================================
-- 5. 角色权限关联表
-- =============================================
DROP TABLE IF EXISTS `admin_role_permissions`;
CREATE TABLE `admin_role_permissions` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `role_id` bigint(20) unsigned NOT NULL COMMENT '角色ID',
  `permission_id` bigint(20) unsigned NOT NULL COMMENT '权限ID',
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_role_permission` (`role_id`, `permission_id`),
  KEY `idx_role_id` (`role_id`),
  KEY `idx_permission_id` (`permission_id`),
  CONSTRAINT `fk_role_permissions_role_id` FOREIGN KEY (`role_id`) REFERENCES `admin_roles` (`id`) ON DELETE CASCADE,
  CONSTRAINT `fk_role_permissions_permission_id` FOREIGN KEY (`permission_id`) REFERENCES `admin_permissions` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='角色权限关联表';

-- =============================================
-- 6. 部门表
-- =============================================
DROP TABLE IF EXISTS `admin_departments`;
CREATE TABLE `admin_departments` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `parent_id` bigint(20) unsigned NOT NULL DEFAULT 0 COMMENT '父级部门ID',
  `name` varchar(100) NOT NULL COMMENT '部门名称',
  `code` varchar(50) DEFAULT NULL COMMENT '部门编码',
  `level` int(11) NOT NULL DEFAULT 1 COMMENT '层级',
  `path` varchar(500) DEFAULT NULL COMMENT '层级路径(如: 1,2,3)',
  `manager_id` bigint(20) unsigned DEFAULT NULL COMMENT '部门负责人ID',
  `phone` varchar(20) DEFAULT NULL COMMENT '部门电话',
  `email` varchar(100) DEFAULT NULL COMMENT '部门邮箱',
  `address` varchar(255) DEFAULT NULL COMMENT '部门地址',
  `sort_order` int(11) NOT NULL DEFAULT 0 COMMENT '排序',
  `status` tinyint(1) NOT NULL DEFAULT 1 COMMENT '状态: 0-禁用, 1-启用',
  `description` text COMMENT '部门描述',
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `deleted_at` timestamp NULL DEFAULT NULL COMMENT '删除时间(软删除)',
  PRIMARY KEY (`id`),
  KEY `idx_parent_id` (`parent_id`),
  KEY `idx_manager_id` (`manager_id`),
  KEY `idx_status` (`status`),
  KEY `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='部门表';

-- =============================================
-- 7. 登录日志表
-- =============================================
DROP TABLE IF EXISTS `admin_login_logs`;
CREATE TABLE `admin_login_logs` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `user_id` bigint(20) unsigned DEFAULT NULL COMMENT '用户ID',
  `username` varchar(50) DEFAULT NULL COMMENT '用户名',
  `login_type` tinyint(1) NOT NULL DEFAULT 1 COMMENT '登录类型: 1-用户名, 2-手机号, 3-邮箱',
  `login_method` varchar(20) NOT NULL DEFAULT 'password' COMMENT '登录方式: password-密码, sms-短信, qrcode-二维码',
  `ip_address` varchar(45) NOT NULL COMMENT 'IP地址',
  `user_agent` text COMMENT '用户代理',
  `device_type` varchar(20) DEFAULT NULL COMMENT '设备类型: web, mobile, tablet',
  `browser` varchar(50) DEFAULT NULL COMMENT '浏览器',
  `os` varchar(50) DEFAULT NULL COMMENT '操作系统',
  `location` varchar(100) DEFAULT NULL COMMENT '登录地点',
  `status` tinyint(1) NOT NULL DEFAULT 1 COMMENT '登录状态: 0-失败, 1-成功',
  `failure_reason` varchar(255) DEFAULT NULL COMMENT '失败原因',
  `login_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '登录时间',
  `logout_time` timestamp NULL DEFAULT NULL COMMENT '退出时间',
  PRIMARY KEY (`id`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_username` (`username`),
  KEY `idx_ip_address` (`ip_address`),
  KEY `idx_status` (`status`),
  KEY `idx_login_time` (`login_time`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='管理员登录日志表';

-- =============================================
-- 8. 操作日志表
-- =============================================
DROP TABLE IF EXISTS `admin_operation_logs`;
CREATE TABLE `admin_operation_logs` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `user_id` bigint(20) unsigned DEFAULT NULL COMMENT '操作用户ID',
  `username` varchar(50) DEFAULT NULL COMMENT '操作用户名',
  `module` varchar(50) NOT NULL COMMENT '操作模块',
  `action` varchar(50) NOT NULL COMMENT '操作动作',
  `description` varchar(255) DEFAULT NULL COMMENT '操作描述',
  `method` varchar(10) NOT NULL COMMENT 'HTTP方法',
  `url` varchar(255) NOT NULL COMMENT '请求URL',
  `params` text COMMENT '请求参数',
  `result` text COMMENT '操作结果',
  `ip_address` varchar(45) NOT NULL COMMENT 'IP地址',
  `user_agent` text COMMENT '用户代理',
  `status` tinyint(1) NOT NULL DEFAULT 1 COMMENT '操作状态: 0-失败, 1-成功',
  `error_message` text COMMENT '错误信息',
  `execution_time` int(11) DEFAULT NULL COMMENT '执行时间(毫秒)',
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '操作时间',
  PRIMARY KEY (`id`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_username` (`username`),
  KEY `idx_module` (`module`),
  KEY `idx_action` (`action`),
  KEY `idx_ip_address` (`ip_address`),
  KEY `idx_status` (`status`),
  KEY `idx_created_at` (`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='管理员操作日志表';

-- =============================================
-- 9. 系统配置表
-- =============================================
DROP TABLE IF EXISTS `admin_configs`;
CREATE TABLE `admin_configs` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `group_name` varchar(50) NOT NULL COMMENT '配置分组',
  `config_key` varchar(100) NOT NULL COMMENT '配置键',
  `config_value` text COMMENT '配置值',
  `config_type` varchar(20) NOT NULL DEFAULT 'string' COMMENT '配置类型: string, int, bool, json',
  `description` varchar(255) DEFAULT NULL COMMENT '配置描述',
  `is_system` tinyint(1) NOT NULL DEFAULT 0 COMMENT '是否系统配置: 0-否, 1-是',
  `sort_order` int(11) NOT NULL DEFAULT 0 COMMENT '排序',
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_group_key` (`group_name`, `config_key`),
  KEY `idx_group_name` (`group_name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='系统配置表';

-- =============================================
-- 初始化数据
-- =============================================

-- 插入默认部门
INSERT INTO `admin_departments` (`id`, `parent_id`, `name`, `code`, `level`, `path`, `sort_order`, `status`, `description`) VALUES
(1, 0, '总公司', 'ROOT', 1, '1', 1, 1, '公司总部'),
(2, 1, '技术部', 'TECH', 2, '1,2', 1, 1, '技术研发部门'),
(3, 1, '运营部', 'OPERATION', 2, '1,3', 2, 1, '运营管理部门'),
(4, 1, '财务部', 'FINANCE', 2, '1,4', 3, 1, '财务管理部门');

-- 插入默认角色
INSERT INTO `admin_roles` (`id`, `name`, `code`, `description`, `level`, `status`, `is_system`, `sort_order`) VALUES
(1, '超级管理员', 'super_admin', '系统超级管理员，拥有所有权限', 1, 1, 1, 1),
(2, '系统管理员', 'admin', '系统管理员，拥有大部分权限', 2, 1, 1, 2),
(3, '部门管理员', 'dept_admin', '部门管理员，管理本部门事务', 3, 1, 0, 3),
(4, '普通用户', 'user', '普通用户，基础权限', 4, 1, 0, 4),
(5, '访客', 'guest', '访客用户，只有基本查看权限', 5, 1, 0, 5);

-- 插入默认权限
INSERT INTO `admin_permissions` (`id`, `parent_id`, `name`, `code`, `type`, `path`, `component`, `icon`, `level`, `sort_order`, `status`) VALUES
(1, 0, '系统管理', 'system', 1, '/system', NULL, 'system', 1, 1, 1),
(2, 1, '用户管理', 'system:user', 1, '/system/user', 'system/user/index', 'user', 2, 1, 1),
(3, 2, '用户查询', 'system:user:query', 2, NULL, NULL, NULL, 3, 1, 1),
(4, 2, '用户新增', 'system:user:add', 2, NULL, NULL, NULL, 3, 2, 1),
(5, 2, '用户修改', 'system:user:edit', 2, NULL, NULL, NULL, 3, 3, 1),
(6, 2, '用户删除', 'system:user:delete', 2, NULL, NULL, NULL, 3, 4, 1),
(7, 1, '角色管理', 'system:role', 1, '/system/role', 'system/role/index', 'role', 2, 2, 1),
(8, 7, '角色查询', 'system:role:query', 2, NULL, NULL, NULL, 3, 1, 1),
(9, 7, '角色新增', 'system:role:add', 2, NULL, NULL, NULL, 3, 2, 1),
(10, 7, '角色修改', 'system:role:edit', 2, NULL, NULL, NULL, 3, 3, 1),
(11, 7, '角色删除', 'system:role:delete', 2, NULL, NULL, NULL, 3, 4, 1),
(12, 1, '权限管理', 'system:permission', 1, '/system/permission', 'system/permission/index', 'permission', 2, 3, 1),
(13, 1, '部门管理', 'system:dept', 1, '/system/dept', 'system/dept/index', 'dept', 2, 4, 1),
(14, 0, '日志管理', 'log', 1, '/log', NULL, 'log', 1, 2, 1),
(15, 14, '登录日志', 'log:login', 1, '/log/login', 'log/login/index', 'login-log', 2, 1, 1),
(16, 14, '操作日志', 'log:operation', 1, '/log/operation', 'log/operation/index', 'operation-log', 2, 2, 1);

-- 插入默认用户
INSERT INTO `admin_users` (`id`, `username`, `phone`, `email`, `password`, `real_name`, `nickname`, `department_id`, `position`, `status`, `is_super_admin`) VALUES
(1, 'admin', '13800138000', 'admin@example.com', '$2a$10$N.zmdr9k7uOCQb376NoUnuTJ8iKyqODYOUwjUxgvlMfCdWiw3Ca4a', '系统管理员', '超级管理员', 1, 'CTO', 1, 1),
(2, 'user1', '13800138001', 'user1@example.com', '$2a$10$N.zmdr9k7uOCQb376NoUnuTJ8iKyqODYOUwjUxgvlMfCdWiw3Ca4a', '管理员1', '系统管理员', 2, '技术经理', 1, 0),
(3, 'user2', '13800138002', 'user2@example.com', '$2a$10$N.zmdr9k7uOCQb376NoUnuTJ8iKyqODYOUwjUxgvlMfCdWiw3Ca4a', '部门经理1', '部门管理员', 2, '部门经理', 1, 0),
(4, 'user3', '13800138003', 'user3@example.com', '$2a$10$N.zmdr9k7uOCQb376NoUnuTJ8iKyqODYOUwjUxgvlMfCdWiw3Ca4a', '普通员工1', '普通用户', 3, '开发工程师', 1, 0),
(5, 'guest1', '13800138004', 'guest1@example.com', '$2a$10$N.zmdr9k7uOCQb376NoUnuTJ8iKyqODYOUwjUxgvlMfCdWiw3Ca4a', '访客1', '访客', 4, '实习生', 1, 0);
-- 密码: admin123

-- 分配用户角色
INSERT INTO `admin_user_roles` (`user_id`, `role_id`) VALUES 
(1, 1),  -- admin -> super_admin
(2, 2),  -- user1 -> admin
(3, 3),  -- user2 -> dept_admin
(4, 4),  -- user3 -> user
(5, 5);  -- guest1 -> guest

-- 分配角色权限
-- 超级管理员拥有所有权限
INSERT INTO `admin_role_permissions` (`role_id`, `permission_id`) 
SELECT 1, id FROM `admin_permissions`;

-- 系统管理员权限(除了超级管理员专有权限)
INSERT INTO `admin_role_permissions` (`role_id`, `permission_id`) VALUES
(2, 1), (2, 2), (2, 3), (2, 4), (2, 5), (2, 6),  -- 系统管理-用户管理
(2, 7), (2, 8), (2, 9), (2, 10), (2, 11),         -- 系统管理-角色管理
(2, 12), (2, 13),                                 -- 权限管理、部门管理
(2, 14), (2, 15), (2, 16);                       -- 日志管理

-- 部门管理员权限(本部门管理权限)
INSERT INTO `admin_role_permissions` (`role_id`, `permission_id`) VALUES
(3, 2), (3, 3), (3, 5),    -- 用户管理(查询、修改)
(3, 8), (3, 13);           -- 角色查询、部门管理

-- 普通用户权限(基础查看权限)
INSERT INTO `admin_role_permissions` (`role_id`, `permission_id`) VALUES
(4, 3), (4, 8), (4, 13);  -- 用户查询、角色查询、部门管理

-- 访客权限(最基本权限)
INSERT INTO `admin_role_permissions` (`role_id`, `permission_id`) VALUES
(5, 3);  -- 仅用户查询权限

-- 插入系统配置
INSERT INTO `admin_configs` (`group_name`, `config_key`, `config_value`, `config_type`, `description`, `is_system`) VALUES
('system', 'site_name', '后台管理系统', 'string', '网站名称', 1),
('system', 'site_logo', '/static/images/logo.png', 'string', '网站Logo', 1),
('system', 'login_captcha', 'true', 'bool', '是否开启登录验证码', 1),
('system', 'password_min_length', '6', 'int', '密码最小长度', 1),
('system', 'session_timeout', '7200', 'int', '会话超时时间(秒)', 1),
('security', 'max_login_attempts', '5', 'int', '最大登录尝试次数', 1),
('security', 'lockout_duration', '1800', 'int', '账户锁定时间(秒)', 1),
('security', 'password_expire_days', '90', 'int', '密码过期天数', 1);

-- =============================================
-- 10. Casbin权限策略表 (RBAC + RESTful)
-- =============================================
DROP TABLE IF EXISTS `casbin_rule`;
CREATE TABLE `casbin_rule` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `ptype` varchar(100) NOT NULL COMMENT '策略类型: p-权限策略, g-角色继承',
  `v0` varchar(100) NOT NULL COMMENT '主体(用户/角色)',
  `v1` varchar(100) NOT NULL COMMENT '对象(资源/权限)',
  `v2` varchar(100) NOT NULL COMMENT '动作(操作/HTTP方法)',
  `v3` varchar(100) DEFAULT NULL COMMENT '效果(allow/deny)',
  `v4` varchar(100) DEFAULT NULL COMMENT '扩展字段4',
  `v5` varchar(100) DEFAULT NULL COMMENT '扩展字段5',
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_casbin_rule` (`ptype`, `v0`, `v1`, `v2`, `v3`),
  KEY `idx_ptype` (`ptype`),
  KEY `idx_v0` (`v0`),
  KEY `idx_v1` (`v1`),
  KEY `idx_v2` (`v2`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='Casbin权限策略表';

-- =============================================
-- 11. Casbin权限策略初始化数据
-- =============================================

-- 角色权限策略 (p策略: 角色, 资源, 操作)
-- p策略需要4个字段: ptype, v0(角色), v1(资源), v2(操作)
INSERT INTO `casbin_rule` (`ptype`, `v0`, `v1`, `v2`) VALUES
-- 超级管理员权限 (所有资源的所有操作)
('p', 'super_admin', '/api/admin/*', '*'),
('p', 'super_admin', '/api/users/*', '*'),
('p', 'super_admin', '/api/roles/*', '*'),
('p', 'super_admin', '/api/permissions/*', '*'),
('p', 'super_admin', '/api/departments/*', '*'),
('p', 'super_admin', '/api/configs/*', '*'),
('p', 'super_admin', '/api/logs/*', '*'),

-- 系统管理员权限 (用户和角色管理)
('p', 'admin', '/api/users', 'GET'),
('p', 'admin', '/api/users', 'POST'),
('p', 'admin', '/api/users/*', 'GET'),
('p', 'admin', '/api/users/*', 'PUT'),
('p', 'admin', '/api/users/*', 'DELETE'),
('p', 'admin', '/api/roles', 'GET'),
('p', 'admin', '/api/roles/*', 'GET'),
('p', 'admin', '/api/permissions', 'GET'),
('p', 'admin', '/api/permissions/*', 'GET'),
('p', 'admin', '/api/departments', 'GET'),
('p', 'admin', '/api/departments/*', 'GET'),
('p', 'admin', '/api/logs', 'GET'),
('p', 'admin', '/api/logs/*', 'GET'),

-- 部门管理员权限 (本部门用户管理)
('p', 'dept_admin', '/api/users', 'GET'),
('p', 'dept_admin', '/api/users/dept/*', 'GET'),
('p', 'dept_admin', '/api/users/dept/*', 'PUT'),
('p', 'dept_admin', '/api/departments', 'GET'),
('p', 'dept_admin', '/api/departments/*', 'GET'),
('p', 'dept_admin', '/api/roles', 'GET'),
('p', 'dept_admin', '/api/permissions', 'GET'),

-- 普通用户权限 (个人信息和基础查看权限)
('p', 'user', '/api/users/profile', 'GET'),
('p', 'user', '/api/users/profile', 'PUT'),
('p', 'user', '/api/departments', 'GET'),
('p', 'user', '/api/departments/*', 'GET'),

-- 访客权限 (最基本的查看权限)
('p', 'guest', '/api/users/profile', 'GET');

-- 用户角色分配 (g策略: 用户, 角色)
-- g策略需要为v2字段提供空值，因为该字段定义为NOT NULL
INSERT INTO `casbin_rule` (`ptype`, `v0`, `v1`, `v2`) VALUES
('g', 'admin', 'super_admin', ''),
('g', 'user1', 'admin', ''),
('g', 'user2', 'dept_admin', ''),
('g', 'user3', 'user', ''),
('g', 'guest1', 'guest', '');

-- =============================================
-- 12. 索引优化和表结构说明
-- =============================================

-- Casbin表索引说明:
-- 1. uk_casbin_rule: 确保策略唯一性
-- 2. idx_ptype: 按策略类型查询优化
-- 3. idx_v0: 按主体(用户/角色)查询优化
-- 4. idx_v1: 按对象(资源)查询优化
-- 5. idx_v2: 按动作(操作)查询优化

-- RBAC + RESTful权限模型说明:
-- 1. p策略: 定义角色对资源的操作权限
--    格式: p, 角色, 资源路径, HTTP方法
-- 2. g策略: 定义用户与角色的继承关系
--    格式: g, 用户, 角色
-- 3. 支持通配符匹配: * 表示所有, /api/users/* 表示users下所有子路径
-- 4. HTTP方法: GET, POST, PUT, DELETE, PATCH等

-- 恢复外键检查
SET FOREIGN_KEY_CHECKS = 1;

-- =============================================
-- 创建索引优化查询性能
-- =============================================

-- 用户表复合索引
-- 注意: 如果索引已存在会报错，可以忽略或先检查索引是否存在
ALTER TABLE `admin_users` ADD INDEX `idx_status_deleted` (`status`, `deleted_at`);
ALTER TABLE `admin_users` ADD INDEX `idx_department_status` (`department_id`, `status`);

-- 权限表复合索引
ALTER TABLE `admin_permissions` ADD INDEX `idx_parent_type_status` (`parent_id`, `type`, `status`);

-- 日志表索引
ALTER TABLE `admin_login_logs` ADD INDEX `idx_login_time_status` (`login_time`, `status`);
ALTER TABLE `admin_operation_logs` ADD INDEX `idx_created_status` (`created_at`, `status`);

-- 角色表索引
ALTER TABLE `admin_roles` ADD INDEX `idx_code_status` (`code`, `status`);
ALTER TABLE `admin_roles` ADD INDEX `idx_level_status` (`level`, `status`);

-- 部门表索引
ALTER TABLE `admin_departments` ADD INDEX `idx_parent_level` (`parent_id`, `level`);
ALTER TABLE `admin_departments` ADD INDEX `idx_code_status` (`code`, `status`);

-- 关联表索引
ALTER TABLE `admin_user_roles` ADD INDEX `idx_user_role` (`user_id`, `role_id`);
ALTER TABLE `admin_role_permissions` ADD INDEX `idx_role_permission` (`role_id`, `permission_id`);

-- 配置表索引
ALTER TABLE `admin_configs` ADD INDEX `idx_group_system` (`group_name`, `is_system`);
ALTER TABLE `admin_configs` ADD INDEX `idx_config_type` (`config_type`);

-- =============================================
-- 说明文档
-- =============================================
/*
表结构说明:

1. admin_users: 管理员用户表
   - 支持用户名、手机号、邮箱登录
   - 包含用户基本信息、部门关联、状态管理
   - 支持软删除、登录统计

2. admin_roles: 角色表
   - 支持角色层级管理
   - 系统角色和自定义角色区分

3. admin_permissions: 权限表
   - 支持树形结构权限管理
   - 菜单、按钮、接口权限类型

4. admin_user_roles: 用户角色关联表
   - 多对多关系，一个用户可以有多个角色

5. admin_role_permissions: 角色权限关联表
   - 多对多关系，一个角色可以有多个权限

6. admin_departments: 部门表
   - 支持树形结构部门管理
   - 部门负责人关联

7. admin_login_logs: 登录日志表
   - 记录所有登录尝试
   - 支持多种登录方式统计

8. admin_operation_logs: 操作日志表
   - 记录所有管理操作
   - 支持操作审计

9. admin_configs: 系统配置表
   - 动态系统配置管理
   - 支持多种数据类型

安全特性:
- 密码bcrypt加密
- 软删除支持
- 登录日志记录
- 操作审计
- 权限细粒度控制
- 账户锁定机制

性能优化:
- 合理的索引设计
- 外键约束
- 字段长度优化
- 查询性能考虑
*/