/*
 Navicat Premium Data Transfer

 Source Server         : kunpeng数据库
 Source Server Type    : MySQL
 Source Server Version : 80042 (8.0.42)
 Source Host           : 127.0.0.1:3306
 Source Schema         : kunpeng

 Target Server Type    : MySQL
 Target Server Version : 80042 (8.0.42)
 File Encoding         : 65001

 Date: 07/07/2025 09:45:20
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for admin_configs
-- ----------------------------
DROP TABLE IF EXISTS `admin_configs`;
CREATE TABLE `admin_configs` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `group_name` varchar(50) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '配置分组',
  `config_key` varchar(100) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '配置键',
  `config_value` text COLLATE utf8mb4_unicode_ci COMMENT '配置值',
  `config_type` varchar(20) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT 'string' COMMENT '配置类型: string, int, bool, json',
  `description` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '配置描述',
  `is_system` tinyint(1) NOT NULL DEFAULT '0' COMMENT '是否系统配置: 0-否, 1-是',
  `sort_order` int NOT NULL DEFAULT '0' COMMENT '排序',
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_group_key` (`group_name`,`config_key`),
  KEY `idx_group_name` (`group_name`),
  KEY `idx_group_system` (`group_name`,`is_system`),
  KEY `idx_config_type` (`config_type`)
) ENGINE=InnoDB AUTO_INCREMENT=9 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='系统配置表';

-- ----------------------------
-- Records of admin_configs
-- ----------------------------
BEGIN;
INSERT INTO `admin_configs` (`id`, `group_name`, `config_key`, `config_value`, `config_type`, `description`, `is_system`, `sort_order`, `created_at`, `updated_at`) VALUES (1, 'system', 'site_name', '后台管理系统', 'string', '网站名称', 1, 0, '2025-07-06 05:49:16', '2025-07-06 05:49:16');
INSERT INTO `admin_configs` (`id`, `group_name`, `config_key`, `config_value`, `config_type`, `description`, `is_system`, `sort_order`, `created_at`, `updated_at`) VALUES (2, 'system', 'site_logo', '/static/images/logo.png', 'string', '网站Logo', 1, 0, '2025-07-06 05:49:16', '2025-07-06 05:49:16');
INSERT INTO `admin_configs` (`id`, `group_name`, `config_key`, `config_value`, `config_type`, `description`, `is_system`, `sort_order`, `created_at`, `updated_at`) VALUES (3, 'system', 'login_captcha', 'true', 'bool', '是否开启登录验证码', 1, 0, '2025-07-06 05:49:16', '2025-07-06 05:49:16');
INSERT INTO `admin_configs` (`id`, `group_name`, `config_key`, `config_value`, `config_type`, `description`, `is_system`, `sort_order`, `created_at`, `updated_at`) VALUES (4, 'system', 'password_min_length', '6', 'int', '密码最小长度', 1, 0, '2025-07-06 05:49:16', '2025-07-06 05:49:16');
INSERT INTO `admin_configs` (`id`, `group_name`, `config_key`, `config_value`, `config_type`, `description`, `is_system`, `sort_order`, `created_at`, `updated_at`) VALUES (5, 'system', 'session_timeout', '7200', 'int', '会话超时时间(秒)', 1, 0, '2025-07-06 05:49:16', '2025-07-06 05:49:16');
INSERT INTO `admin_configs` (`id`, `group_name`, `config_key`, `config_value`, `config_type`, `description`, `is_system`, `sort_order`, `created_at`, `updated_at`) VALUES (6, 'security', 'max_login_attempts', '5', 'int', '最大登录尝试次数', 1, 0, '2025-07-06 05:49:16', '2025-07-06 05:49:16');
INSERT INTO `admin_configs` (`id`, `group_name`, `config_key`, `config_value`, `config_type`, `description`, `is_system`, `sort_order`, `created_at`, `updated_at`) VALUES (7, 'security', 'lockout_duration', '1800', 'int', '账户锁定时间(秒)', 1, 0, '2025-07-06 05:49:16', '2025-07-06 05:49:16');
INSERT INTO `admin_configs` (`id`, `group_name`, `config_key`, `config_value`, `config_type`, `description`, `is_system`, `sort_order`, `created_at`, `updated_at`) VALUES (8, 'security', 'password_expire_days', '90', 'int', '密码过期天数', 1, 0, '2025-07-06 05:49:16', '2025-07-06 05:49:16');
COMMIT;

-- ----------------------------
-- Table structure for admin_departments
-- ----------------------------
DROP TABLE IF EXISTS `admin_departments`;
CREATE TABLE `admin_departments` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `parent_id` bigint unsigned NOT NULL DEFAULT '0' COMMENT '父级部门ID',
  `name` varchar(100) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '部门名称',
  `code` varchar(50) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '部门编码',
  `level` int NOT NULL DEFAULT '1' COMMENT '层级',
  `path` varchar(500) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '层级路径(如: 1,2,3)',
  `manager_id` bigint unsigned DEFAULT NULL COMMENT '部门负责人ID',
  `phone` varchar(20) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '部门电话',
  `email` varchar(100) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '部门邮箱',
  `address` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '部门地址',
  `sort_order` int NOT NULL DEFAULT '0' COMMENT '排序',
  `status` tinyint(1) NOT NULL DEFAULT '1' COMMENT '状态: 0-禁用, 1-启用',
  `description` text COLLATE utf8mb4_unicode_ci COMMENT '部门描述',
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `deleted_at` timestamp NULL DEFAULT NULL COMMENT '删除时间(软删除)',
  PRIMARY KEY (`id`),
  KEY `idx_parent_id` (`parent_id`),
  KEY `idx_manager_id` (`manager_id`),
  KEY `idx_status` (`status`),
  KEY `idx_deleted_at` (`deleted_at`),
  KEY `idx_parent_level` (`parent_id`,`level`),
  KEY `idx_code_status` (`code`,`status`)
) ENGINE=InnoDB AUTO_INCREMENT=5 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='部门表';

-- ----------------------------
-- Records of admin_departments
-- ----------------------------
BEGIN;
INSERT INTO `admin_departments` (`id`, `parent_id`, `name`, `code`, `level`, `path`, `manager_id`, `phone`, `email`, `address`, `sort_order`, `status`, `description`, `created_at`, `updated_at`, `deleted_at`) VALUES (1, 0, '总公司', 'ROOT', 1, '1', NULL, NULL, NULL, NULL, 1, 1, '公司总部', '2025-07-06 05:49:16', '2025-07-06 05:49:16', NULL);
INSERT INTO `admin_departments` (`id`, `parent_id`, `name`, `code`, `level`, `path`, `manager_id`, `phone`, `email`, `address`, `sort_order`, `status`, `description`, `created_at`, `updated_at`, `deleted_at`) VALUES (2, 1, '技术部', 'TECH', 2, '1,2', NULL, NULL, NULL, NULL, 1, 1, '技术研发部门', '2025-07-06 05:49:16', '2025-07-06 05:49:16', NULL);
INSERT INTO `admin_departments` (`id`, `parent_id`, `name`, `code`, `level`, `path`, `manager_id`, `phone`, `email`, `address`, `sort_order`, `status`, `description`, `created_at`, `updated_at`, `deleted_at`) VALUES (3, 1, '运营部', 'OPERATION', 2, '1,3', NULL, NULL, NULL, NULL, 2, 1, '运营管理部门', '2025-07-06 05:49:16', '2025-07-06 05:49:16', NULL);
INSERT INTO `admin_departments` (`id`, `parent_id`, `name`, `code`, `level`, `path`, `manager_id`, `phone`, `email`, `address`, `sort_order`, `status`, `description`, `created_at`, `updated_at`, `deleted_at`) VALUES (4, 1, '财务部', 'FINANCE', 2, '1,4', NULL, NULL, NULL, NULL, 3, 1, '财务管理部门', '2025-07-06 05:49:16', '2025-07-06 05:49:16', NULL);
COMMIT;

-- ----------------------------
-- Table structure for admin_login_logs
-- ----------------------------
DROP TABLE IF EXISTS `admin_login_logs`;
CREATE TABLE `admin_login_logs` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `user_id` bigint unsigned DEFAULT NULL COMMENT '用户ID',
  `username` varchar(50) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '用户名',
  `login_type` tinyint(1) NOT NULL DEFAULT '1' COMMENT '登录类型: 1-用户名, 2-手机号, 3-邮箱',
  `login_method` varchar(20) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT 'password' COMMENT '登录方式: password-密码, sms-短信, qrcode-二维码',
  `ip_address` varchar(45) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT 'IP地址',
  `user_agent` text COLLATE utf8mb4_unicode_ci COMMENT '用户代理',
  `device_type` varchar(20) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '设备类型: web, mobile, tablet',
  `browser` varchar(50) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '浏览器',
  `os` varchar(50) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '操作系统',
  `location` varchar(100) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '登录地点',
  `status` tinyint(1) NOT NULL DEFAULT '1' COMMENT '登录状态: 0-失败, 1-成功',
  `failure_reason` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '失败原因',
  `login_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '登录时间',
  `logout_time` timestamp NULL DEFAULT NULL COMMENT '退出时间',
  PRIMARY KEY (`id`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_username` (`username`),
  KEY `idx_ip_address` (`ip_address`),
  KEY `idx_status` (`status`),
  KEY `idx_login_time` (`login_time`),
  KEY `idx_login_time_status` (`login_time`,`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='管理员登录日志表';

-- ----------------------------
-- Records of admin_login_logs
-- ----------------------------
BEGIN;
COMMIT;

-- ----------------------------
-- Table structure for admin_operation_logs
-- ----------------------------
DROP TABLE IF EXISTS `admin_operation_logs`;
CREATE TABLE `admin_operation_logs` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `user_id` bigint unsigned DEFAULT NULL COMMENT '操作用户ID',
  `username` varchar(50) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '操作用户名',
  `module` varchar(50) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '操作模块',
  `action` varchar(50) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '操作动作',
  `description` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '操作描述',
  `method` varchar(10) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT 'HTTP方法',
  `url` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '请求URL',
  `params` text COLLATE utf8mb4_unicode_ci COMMENT '请求参数',
  `result` text COLLATE utf8mb4_unicode_ci COMMENT '操作结果',
  `ip_address` varchar(45) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT 'IP地址',
  `user_agent` text COLLATE utf8mb4_unicode_ci COMMENT '用户代理',
  `status` tinyint(1) NOT NULL DEFAULT '1' COMMENT '操作状态: 0-失败, 1-成功',
  `error_message` text COLLATE utf8mb4_unicode_ci COMMENT '错误信息',
  `execution_time` int DEFAULT NULL COMMENT '执行时间(毫秒)',
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '操作时间',
  PRIMARY KEY (`id`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_username` (`username`),
  KEY `idx_module` (`module`),
  KEY `idx_action` (`action`),
  KEY `idx_ip_address` (`ip_address`),
  KEY `idx_status` (`status`),
  KEY `idx_created_at` (`created_at`),
  KEY `idx_created_status` (`created_at`,`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='管理员操作日志表';

-- ----------------------------
-- Records of admin_operation_logs
-- ----------------------------
BEGIN;
COMMIT;

-- ----------------------------
-- Table structure for admin_permissions
-- ----------------------------
DROP TABLE IF EXISTS `admin_permissions`;
CREATE TABLE `admin_permissions` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `parent_id` bigint unsigned NOT NULL DEFAULT '0' COMMENT '父级权限ID',
  `name` varchar(100) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '权限名称',
  `code` varchar(100) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '权限编码',
  `type` tinyint(1) NOT NULL DEFAULT '1' COMMENT '权限类型: 1-菜单, 2-按钮, 3-接口',
  `path` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '路由路径',
  `component` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '组件路径',
  `icon` varchar(100) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '图标',
  `method` varchar(10) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT 'HTTP方法(GET,POST,PUT,DELETE等)',
  `url` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT 'API接口地址',
  `level` int NOT NULL DEFAULT '1' COMMENT '层级',
  `sort_order` int NOT NULL DEFAULT '0' COMMENT '排序',
  `status` tinyint(1) NOT NULL DEFAULT '1' COMMENT '状态: 0-禁用, 1-启用',
  `is_hidden` tinyint(1) NOT NULL DEFAULT '0' COMMENT '是否隐藏: 0-否, 1-是',
  `description` text COLLATE utf8mb4_unicode_ci COMMENT '权限描述',
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `deleted_at` timestamp NULL DEFAULT NULL COMMENT '删除时间(软删除)',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_code` (`code`),
  KEY `idx_parent_id` (`parent_id`),
  KEY `idx_type` (`type`),
  KEY `idx_status` (`status`),
  KEY `idx_deleted_at` (`deleted_at`),
  KEY `idx_parent_type_status` (`parent_id`,`type`,`status`)
) ENGINE=InnoDB AUTO_INCREMENT=17 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='管理员权限表';

-- ----------------------------
-- Records of admin_permissions
-- ----------------------------
BEGIN;
INSERT INTO `admin_permissions` (`id`, `parent_id`, `name`, `code`, `type`, `path`, `component`, `icon`, `method`, `url`, `level`, `sort_order`, `status`, `is_hidden`, `description`, `created_at`, `updated_at`, `deleted_at`) VALUES (1, 0, '系统管理', 'system', 1, '/system', NULL, 'system', NULL, NULL, 1, 1, 1, 0, NULL, '2025-07-06 05:49:16', '2025-07-06 05:49:16', NULL);
INSERT INTO `admin_permissions` (`id`, `parent_id`, `name`, `code`, `type`, `path`, `component`, `icon`, `method`, `url`, `level`, `sort_order`, `status`, `is_hidden`, `description`, `created_at`, `updated_at`, `deleted_at`) VALUES (2, 1, '用户管理', 'system:user', 1, '/system/user', 'system/user/index', 'user', NULL, NULL, 2, 1, 1, 0, NULL, '2025-07-06 05:49:16', '2025-07-06 05:49:16', NULL);
INSERT INTO `admin_permissions` (`id`, `parent_id`, `name`, `code`, `type`, `path`, `component`, `icon`, `method`, `url`, `level`, `sort_order`, `status`, `is_hidden`, `description`, `created_at`, `updated_at`, `deleted_at`) VALUES (3, 2, '用户查询', 'system:user:query', 2, NULL, NULL, NULL, NULL, NULL, 3, 1, 1, 0, NULL, '2025-07-06 05:49:16', '2025-07-06 05:49:16', NULL);
INSERT INTO `admin_permissions` (`id`, `parent_id`, `name`, `code`, `type`, `path`, `component`, `icon`, `method`, `url`, `level`, `sort_order`, `status`, `is_hidden`, `description`, `created_at`, `updated_at`, `deleted_at`) VALUES (4, 2, '用户新增', 'system:user:add', 2, NULL, NULL, NULL, NULL, NULL, 3, 2, 1, 0, NULL, '2025-07-06 05:49:16', '2025-07-06 05:49:16', NULL);
INSERT INTO `admin_permissions` (`id`, `parent_id`, `name`, `code`, `type`, `path`, `component`, `icon`, `method`, `url`, `level`, `sort_order`, `status`, `is_hidden`, `description`, `created_at`, `updated_at`, `deleted_at`) VALUES (5, 2, '用户修改', 'system:user:edit', 2, NULL, NULL, NULL, NULL, NULL, 3, 3, 1, 0, NULL, '2025-07-06 05:49:16', '2025-07-06 05:49:16', NULL);
INSERT INTO `admin_permissions` (`id`, `parent_id`, `name`, `code`, `type`, `path`, `component`, `icon`, `method`, `url`, `level`, `sort_order`, `status`, `is_hidden`, `description`, `created_at`, `updated_at`, `deleted_at`) VALUES (6, 2, '用户删除', 'system:user:delete', 2, NULL, NULL, NULL, NULL, NULL, 3, 4, 1, 0, NULL, '2025-07-06 05:49:16', '2025-07-06 05:49:16', NULL);
INSERT INTO `admin_permissions` (`id`, `parent_id`, `name`, `code`, `type`, `path`, `component`, `icon`, `method`, `url`, `level`, `sort_order`, `status`, `is_hidden`, `description`, `created_at`, `updated_at`, `deleted_at`) VALUES (7, 1, '角色管理', 'system:role', 1, '/system/role', 'system/role/index', 'role', NULL, NULL, 2, 2, 1, 0, NULL, '2025-07-06 05:49:16', '2025-07-06 05:49:16', NULL);
INSERT INTO `admin_permissions` (`id`, `parent_id`, `name`, `code`, `type`, `path`, `component`, `icon`, `method`, `url`, `level`, `sort_order`, `status`, `is_hidden`, `description`, `created_at`, `updated_at`, `deleted_at`) VALUES (8, 7, '角色查询', 'system:role:query', 2, NULL, NULL, NULL, NULL, NULL, 3, 1, 1, 0, NULL, '2025-07-06 05:49:16', '2025-07-06 05:49:16', NULL);
INSERT INTO `admin_permissions` (`id`, `parent_id`, `name`, `code`, `type`, `path`, `component`, `icon`, `method`, `url`, `level`, `sort_order`, `status`, `is_hidden`, `description`, `created_at`, `updated_at`, `deleted_at`) VALUES (9, 7, '角色新增', 'system:role:add', 2, NULL, NULL, NULL, NULL, NULL, 3, 2, 1, 0, NULL, '2025-07-06 05:49:16', '2025-07-06 05:49:16', NULL);
INSERT INTO `admin_permissions` (`id`, `parent_id`, `name`, `code`, `type`, `path`, `component`, `icon`, `method`, `url`, `level`, `sort_order`, `status`, `is_hidden`, `description`, `created_at`, `updated_at`, `deleted_at`) VALUES (10, 7, '角色修改', 'system:role:edit', 2, NULL, NULL, NULL, NULL, NULL, 3, 3, 1, 0, NULL, '2025-07-06 05:49:16', '2025-07-06 05:49:16', NULL);
INSERT INTO `admin_permissions` (`id`, `parent_id`, `name`, `code`, `type`, `path`, `component`, `icon`, `method`, `url`, `level`, `sort_order`, `status`, `is_hidden`, `description`, `created_at`, `updated_at`, `deleted_at`) VALUES (11, 7, '角色删除', 'system:role:delete', 2, NULL, NULL, NULL, NULL, NULL, 3, 4, 1, 0, NULL, '2025-07-06 05:49:16', '2025-07-06 05:49:16', NULL);
INSERT INTO `admin_permissions` (`id`, `parent_id`, `name`, `code`, `type`, `path`, `component`, `icon`, `method`, `url`, `level`, `sort_order`, `status`, `is_hidden`, `description`, `created_at`, `updated_at`, `deleted_at`) VALUES (12, 1, '权限管理', 'system:permission', 1, '/system/permission', 'system/permission/index', 'permission', NULL, NULL, 2, 3, 1, 0, NULL, '2025-07-06 05:49:16', '2025-07-06 05:49:16', NULL);
INSERT INTO `admin_permissions` (`id`, `parent_id`, `name`, `code`, `type`, `path`, `component`, `icon`, `method`, `url`, `level`, `sort_order`, `status`, `is_hidden`, `description`, `created_at`, `updated_at`, `deleted_at`) VALUES (13, 1, '部门管理', 'system:dept', 1, '/system/dept', 'system/dept/index', 'dept', NULL, NULL, 2, 4, 1, 0, NULL, '2025-07-06 05:49:16', '2025-07-06 05:49:16', NULL);
INSERT INTO `admin_permissions` (`id`, `parent_id`, `name`, `code`, `type`, `path`, `component`, `icon`, `method`, `url`, `level`, `sort_order`, `status`, `is_hidden`, `description`, `created_at`, `updated_at`, `deleted_at`) VALUES (14, 0, '日志管理', 'log', 1, '/log', NULL, 'log', NULL, NULL, 1, 2, 1, 0, NULL, '2025-07-06 05:49:16', '2025-07-06 05:49:16', NULL);
INSERT INTO `admin_permissions` (`id`, `parent_id`, `name`, `code`, `type`, `path`, `component`, `icon`, `method`, `url`, `level`, `sort_order`, `status`, `is_hidden`, `description`, `created_at`, `updated_at`, `deleted_at`) VALUES (15, 14, '登录日志', 'log:login', 1, '/log/login', 'log/login/index', 'login-log', NULL, NULL, 2, 1, 1, 0, NULL, '2025-07-06 05:49:16', '2025-07-06 05:49:16', NULL);
INSERT INTO `admin_permissions` (`id`, `parent_id`, `name`, `code`, `type`, `path`, `component`, `icon`, `method`, `url`, `level`, `sort_order`, `status`, `is_hidden`, `description`, `created_at`, `updated_at`, `deleted_at`) VALUES (16, 14, '操作日志', 'log:operation', 1, '/log/operation', 'log/operation/index', 'operation-log', NULL, NULL, 2, 2, 1, 0, NULL, '2025-07-06 05:49:16', '2025-07-06 05:49:16', NULL);
COMMIT;

-- ----------------------------
-- Table structure for admin_role_permissions
-- ----------------------------
DROP TABLE IF EXISTS `admin_role_permissions`;
CREATE TABLE `admin_role_permissions` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `role_id` bigint unsigned NOT NULL COMMENT '角色ID',
  `permission_id` bigint unsigned NOT NULL COMMENT '权限ID',
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_role_permission` (`role_id`,`permission_id`),
  KEY `idx_role_id` (`role_id`),
  KEY `idx_permission_id` (`permission_id`),
  KEY `idx_role_permission` (`role_id`,`permission_id`),
  CONSTRAINT `fk_role_permissions_permission_id` FOREIGN KEY (`permission_id`) REFERENCES `admin_permissions` (`id`) ON DELETE CASCADE,
  CONSTRAINT `fk_role_permissions_role_id` FOREIGN KEY (`role_id`) REFERENCES `admin_roles` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB AUTO_INCREMENT=57 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='角色权限关联表';

-- ----------------------------
-- Records of admin_role_permissions
-- ----------------------------
BEGIN;
INSERT INTO `admin_role_permissions` (`id`, `role_id`, `permission_id`, `created_at`) VALUES (1, 1, 1, '2025-07-06 05:49:16');
INSERT INTO `admin_role_permissions` (`id`, `role_id`, `permission_id`, `created_at`) VALUES (2, 1, 2, '2025-07-06 05:49:16');
INSERT INTO `admin_role_permissions` (`id`, `role_id`, `permission_id`, `created_at`) VALUES (3, 1, 7, '2025-07-06 05:49:16');
INSERT INTO `admin_role_permissions` (`id`, `role_id`, `permission_id`, `created_at`) VALUES (4, 1, 12, '2025-07-06 05:49:16');
INSERT INTO `admin_role_permissions` (`id`, `role_id`, `permission_id`, `created_at`) VALUES (5, 1, 13, '2025-07-06 05:49:16');
INSERT INTO `admin_role_permissions` (`id`, `role_id`, `permission_id`, `created_at`) VALUES (6, 1, 14, '2025-07-06 05:49:16');
INSERT INTO `admin_role_permissions` (`id`, `role_id`, `permission_id`, `created_at`) VALUES (7, 1, 15, '2025-07-06 05:49:16');
INSERT INTO `admin_role_permissions` (`id`, `role_id`, `permission_id`, `created_at`) VALUES (8, 1, 16, '2025-07-06 05:49:16');
INSERT INTO `admin_role_permissions` (`id`, `role_id`, `permission_id`, `created_at`) VALUES (9, 1, 3, '2025-07-06 05:49:16');
INSERT INTO `admin_role_permissions` (`id`, `role_id`, `permission_id`, `created_at`) VALUES (10, 1, 4, '2025-07-06 05:49:16');
INSERT INTO `admin_role_permissions` (`id`, `role_id`, `permission_id`, `created_at`) VALUES (11, 1, 5, '2025-07-06 05:49:16');
INSERT INTO `admin_role_permissions` (`id`, `role_id`, `permission_id`, `created_at`) VALUES (12, 1, 6, '2025-07-06 05:49:16');
INSERT INTO `admin_role_permissions` (`id`, `role_id`, `permission_id`, `created_at`) VALUES (13, 1, 8, '2025-07-06 05:49:16');
INSERT INTO `admin_role_permissions` (`id`, `role_id`, `permission_id`, `created_at`) VALUES (14, 1, 9, '2025-07-06 05:49:16');
INSERT INTO `admin_role_permissions` (`id`, `role_id`, `permission_id`, `created_at`) VALUES (15, 1, 10, '2025-07-06 05:49:16');
INSERT INTO `admin_role_permissions` (`id`, `role_id`, `permission_id`, `created_at`) VALUES (16, 1, 11, '2025-07-06 05:49:16');
INSERT INTO `admin_role_permissions` (`id`, `role_id`, `permission_id`, `created_at`) VALUES (32, 2, 1, '2025-07-06 05:49:16');
INSERT INTO `admin_role_permissions` (`id`, `role_id`, `permission_id`, `created_at`) VALUES (33, 2, 2, '2025-07-06 05:49:16');
INSERT INTO `admin_role_permissions` (`id`, `role_id`, `permission_id`, `created_at`) VALUES (34, 2, 3, '2025-07-06 05:49:16');
INSERT INTO `admin_role_permissions` (`id`, `role_id`, `permission_id`, `created_at`) VALUES (35, 2, 4, '2025-07-06 05:49:16');
INSERT INTO `admin_role_permissions` (`id`, `role_id`, `permission_id`, `created_at`) VALUES (36, 2, 5, '2025-07-06 05:49:16');
INSERT INTO `admin_role_permissions` (`id`, `role_id`, `permission_id`, `created_at`) VALUES (37, 2, 6, '2025-07-06 05:49:16');
INSERT INTO `admin_role_permissions` (`id`, `role_id`, `permission_id`, `created_at`) VALUES (38, 2, 7, '2025-07-06 05:49:16');
INSERT INTO `admin_role_permissions` (`id`, `role_id`, `permission_id`, `created_at`) VALUES (39, 2, 8, '2025-07-06 05:49:16');
INSERT INTO `admin_role_permissions` (`id`, `role_id`, `permission_id`, `created_at`) VALUES (40, 2, 9, '2025-07-06 05:49:16');
INSERT INTO `admin_role_permissions` (`id`, `role_id`, `permission_id`, `created_at`) VALUES (41, 2, 10, '2025-07-06 05:49:16');
INSERT INTO `admin_role_permissions` (`id`, `role_id`, `permission_id`, `created_at`) VALUES (42, 2, 11, '2025-07-06 05:49:16');
INSERT INTO `admin_role_permissions` (`id`, `role_id`, `permission_id`, `created_at`) VALUES (43, 2, 12, '2025-07-06 05:49:16');
INSERT INTO `admin_role_permissions` (`id`, `role_id`, `permission_id`, `created_at`) VALUES (44, 2, 13, '2025-07-06 05:49:16');
INSERT INTO `admin_role_permissions` (`id`, `role_id`, `permission_id`, `created_at`) VALUES (45, 2, 14, '2025-07-06 05:49:16');
INSERT INTO `admin_role_permissions` (`id`, `role_id`, `permission_id`, `created_at`) VALUES (46, 2, 15, '2025-07-06 05:49:16');
INSERT INTO `admin_role_permissions` (`id`, `role_id`, `permission_id`, `created_at`) VALUES (47, 2, 16, '2025-07-06 05:49:16');
INSERT INTO `admin_role_permissions` (`id`, `role_id`, `permission_id`, `created_at`) VALUES (48, 3, 2, '2025-07-06 05:49:16');
INSERT INTO `admin_role_permissions` (`id`, `role_id`, `permission_id`, `created_at`) VALUES (49, 3, 3, '2025-07-06 05:49:16');
INSERT INTO `admin_role_permissions` (`id`, `role_id`, `permission_id`, `created_at`) VALUES (50, 3, 5, '2025-07-06 05:49:16');
INSERT INTO `admin_role_permissions` (`id`, `role_id`, `permission_id`, `created_at`) VALUES (51, 3, 8, '2025-07-06 05:49:16');
INSERT INTO `admin_role_permissions` (`id`, `role_id`, `permission_id`, `created_at`) VALUES (52, 3, 13, '2025-07-06 05:49:16');
INSERT INTO `admin_role_permissions` (`id`, `role_id`, `permission_id`, `created_at`) VALUES (53, 4, 3, '2025-07-06 05:49:16');
INSERT INTO `admin_role_permissions` (`id`, `role_id`, `permission_id`, `created_at`) VALUES (54, 4, 8, '2025-07-06 05:49:16');
INSERT INTO `admin_role_permissions` (`id`, `role_id`, `permission_id`, `created_at`) VALUES (55, 4, 13, '2025-07-06 05:49:16');
INSERT INTO `admin_role_permissions` (`id`, `role_id`, `permission_id`, `created_at`) VALUES (56, 5, 3, '2025-07-06 05:49:16');
COMMIT;

-- ----------------------------
-- Table structure for admin_roles
-- ----------------------------
DROP TABLE IF EXISTS `admin_roles`;
CREATE TABLE `admin_roles` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `name` varchar(50) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '角色名称',
  `code` varchar(50) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '角色编码',
  `description` text COLLATE utf8mb4_unicode_ci COMMENT '角色描述',
  `level` int NOT NULL DEFAULT '1' COMMENT '角色级别(数字越小权限越高)',
  `status` tinyint(1) NOT NULL DEFAULT '1' COMMENT '状态: 0-禁用, 1-启用',
  `is_system` tinyint(1) NOT NULL DEFAULT '0' COMMENT '是否系统角色: 0-否, 1-是',
  `sort_order` int NOT NULL DEFAULT '0' COMMENT '排序',
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `deleted_at` timestamp NULL DEFAULT NULL COMMENT '删除时间(软删除)',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_code` (`code`),
  KEY `idx_status` (`status`),
  KEY `idx_level` (`level`),
  KEY `idx_deleted_at` (`deleted_at`),
  KEY `idx_code_status` (`code`,`status`),
  KEY `idx_level_status` (`level`,`status`)
) ENGINE=InnoDB AUTO_INCREMENT=6 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='管理员角色表';

-- ----------------------------
-- Records of admin_roles
-- ----------------------------
BEGIN;
INSERT INTO `admin_roles` (`id`, `name`, `code`, `description`, `level`, `status`, `is_system`, `sort_order`, `created_at`, `updated_at`, `deleted_at`) VALUES (1, '超级管理员', 'super_admin', '系统超级管理员，拥有所有权限', 1, 1, 1, 1, '2025-07-06 05:49:16', '2025-07-06 05:49:16', NULL);
INSERT INTO `admin_roles` (`id`, `name`, `code`, `description`, `level`, `status`, `is_system`, `sort_order`, `created_at`, `updated_at`, `deleted_at`) VALUES (2, '系统管理员', 'admin', '系统管理员，拥有大部分权限', 2, 1, 1, 2, '2025-07-06 05:49:16', '2025-07-06 05:49:16', NULL);
INSERT INTO `admin_roles` (`id`, `name`, `code`, `description`, `level`, `status`, `is_system`, `sort_order`, `created_at`, `updated_at`, `deleted_at`) VALUES (3, '部门管理员', 'dept_admin', '部门管理员，管理本部门事务', 3, 1, 0, 3, '2025-07-06 05:49:16', '2025-07-06 05:49:16', NULL);
INSERT INTO `admin_roles` (`id`, `name`, `code`, `description`, `level`, `status`, `is_system`, `sort_order`, `created_at`, `updated_at`, `deleted_at`) VALUES (4, '普通用户', 'user', '普通用户，基础权限', 4, 1, 0, 4, '2025-07-06 05:49:16', '2025-07-06 05:49:16', NULL);
INSERT INTO `admin_roles` (`id`, `name`, `code`, `description`, `level`, `status`, `is_system`, `sort_order`, `created_at`, `updated_at`, `deleted_at`) VALUES (5, '访客', 'guest', '访客用户，只有基本查看权限', 5, 1, 0, 5, '2025-07-06 05:49:16', '2025-07-06 05:49:16', NULL);
COMMIT;

-- ----------------------------
-- Table structure for admin_user_roles
-- ----------------------------
DROP TABLE IF EXISTS `admin_user_roles`;
CREATE TABLE `admin_user_roles` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `user_id` bigint unsigned NOT NULL COMMENT '用户ID',
  `role_id` bigint unsigned NOT NULL COMMENT '角色ID',
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_user_role` (`user_id`,`role_id`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_role_id` (`role_id`),
  KEY `idx_user_role` (`user_id`,`role_id`),
  CONSTRAINT `fk_user_roles_role_id` FOREIGN KEY (`role_id`) REFERENCES `admin_roles` (`id`) ON DELETE CASCADE,
  CONSTRAINT `fk_user_roles_user_id` FOREIGN KEY (`user_id`) REFERENCES `admin_users` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB AUTO_INCREMENT=6 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户角色关联表';

-- ----------------------------
-- Records of admin_user_roles
-- ----------------------------
BEGIN;
INSERT INTO `admin_user_roles` (`id`, `user_id`, `role_id`, `created_at`) VALUES (1, 1, 1, '2025-07-06 05:49:16');
INSERT INTO `admin_user_roles` (`id`, `user_id`, `role_id`, `created_at`) VALUES (2, 2, 2, '2025-07-06 05:49:16');
INSERT INTO `admin_user_roles` (`id`, `user_id`, `role_id`, `created_at`) VALUES (3, 3, 3, '2025-07-06 05:49:16');
INSERT INTO `admin_user_roles` (`id`, `user_id`, `role_id`, `created_at`) VALUES (4, 4, 4, '2025-07-06 05:49:16');
INSERT INTO `admin_user_roles` (`id`, `user_id`, `role_id`, `created_at`) VALUES (5, 5, 5, '2025-07-06 05:49:16');
COMMIT;

-- ----------------------------
-- Table structure for admin_users
-- ----------------------------
DROP TABLE IF EXISTS `admin_users`;
CREATE TABLE `admin_users` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `username` varchar(50) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '用户名/账号',
  `phone` varchar(20) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '手机号',
  `email` varchar(100) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '邮箱',
  `password` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '密码(bcrypt加密)',
  `real_name` varchar(50) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '真实姓名',
  `nickname` varchar(50) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '昵称',
  `avatar` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '头像URL',
  `gender` tinyint(1) DEFAULT '0' COMMENT '性别: 0-未知, 1-男, 2-女',
  `birthday` date DEFAULT NULL COMMENT '生日',
  `department_id` bigint unsigned DEFAULT NULL COMMENT '部门ID',
  `position` varchar(100) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '职位',
  `status` tinyint(1) NOT NULL DEFAULT '1' COMMENT '状态: 0-禁用, 1-启用',
  `is_super_admin` tinyint(1) NOT NULL DEFAULT '0' COMMENT '是否超级管理员: 0-否, 1-是',
  `last_login_time` timestamp NULL DEFAULT NULL COMMENT '最后登录时间',
  `last_login_ip` varchar(45) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '最后登录IP',
  `login_count` int NOT NULL DEFAULT '0' COMMENT '登录次数',
  `password_changed_at` timestamp NULL DEFAULT NULL COMMENT '密码修改时间',
  `remark` text COLLATE utf8mb4_unicode_ci COMMENT '备注',
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
  KEY `idx_deleted_at` (`deleted_at`),
  KEY `idx_status_deleted` (`status`,`deleted_at`),
  KEY `idx_department_status` (`department_id`,`status`)
) ENGINE=InnoDB AUTO_INCREMENT=6 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='管理员用户表';

-- ----------------------------
-- Records of admin_users
-- ----------------------------
BEGIN;
INSERT INTO `admin_users` (`id`, `username`, `phone`, `email`, `password`, `real_name`, `nickname`, `avatar`, `gender`, `birthday`, `department_id`, `position`, `status`, `is_super_admin`, `last_login_time`, `last_login_ip`, `login_count`, `password_changed_at`, `remark`, `created_at`, `updated_at`, `deleted_at`) VALUES (1, 'admin', '13800138000', 'admin@example.com', '$2a$10$N.zmdr9k7uOCQb376NoUnuTJ8iKyqODYOUwjUxgvlMfCdWiw3Ca4a', '系统管理员', '超级管理员', NULL, 0, NULL, 1, 'CTO', 1, 1, NULL, NULL, 0, NULL, NULL, '2025-07-06 05:49:16', '2025-07-06 05:49:16', NULL);
INSERT INTO `admin_users` (`id`, `username`, `phone`, `email`, `password`, `real_name`, `nickname`, `avatar`, `gender`, `birthday`, `department_id`, `position`, `status`, `is_super_admin`, `last_login_time`, `last_login_ip`, `login_count`, `password_changed_at`, `remark`, `created_at`, `updated_at`, `deleted_at`) VALUES (2, 'user1', '13800138001', 'user1@example.com', '$2a$10$N.zmdr9k7uOCQb376NoUnuTJ8iKyqODYOUwjUxgvlMfCdWiw3Ca4a', '管理员1', '系统管理员', NULL, 0, NULL, 2, '技术经理', 1, 0, NULL, NULL, 0, NULL, NULL, '2025-07-06 05:49:16', '2025-07-06 05:49:16', NULL);
INSERT INTO `admin_users` (`id`, `username`, `phone`, `email`, `password`, `real_name`, `nickname`, `avatar`, `gender`, `birthday`, `department_id`, `position`, `status`, `is_super_admin`, `last_login_time`, `last_login_ip`, `login_count`, `password_changed_at`, `remark`, `created_at`, `updated_at`, `deleted_at`) VALUES (3, 'user2', '13800138002', 'user2@example.com', '$2a$10$N.zmdr9k7uOCQb376NoUnuTJ8iKyqODYOUwjUxgvlMfCdWiw3Ca4a', '部门经理1', '部门管理员', NULL, 0, NULL, 2, '部门经理', 1, 0, NULL, NULL, 0, NULL, NULL, '2025-07-06 05:49:16', '2025-07-06 05:49:16', NULL);
INSERT INTO `admin_users` (`id`, `username`, `phone`, `email`, `password`, `real_name`, `nickname`, `avatar`, `gender`, `birthday`, `department_id`, `position`, `status`, `is_super_admin`, `last_login_time`, `last_login_ip`, `login_count`, `password_changed_at`, `remark`, `created_at`, `updated_at`, `deleted_at`) VALUES (4, 'user3', '13800138003', 'user3@example.com', '$2a$10$N.zmdr9k7uOCQb376NoUnuTJ8iKyqODYOUwjUxgvlMfCdWiw3Ca4a', '普通员工1', '普通用户', NULL, 0, NULL, 3, '开发工程师', 1, 0, NULL, NULL, 0, NULL, NULL, '2025-07-06 05:49:16', '2025-07-06 05:49:16', NULL);
INSERT INTO `admin_users` (`id`, `username`, `phone`, `email`, `password`, `real_name`, `nickname`, `avatar`, `gender`, `birthday`, `department_id`, `position`, `status`, `is_super_admin`, `last_login_time`, `last_login_ip`, `login_count`, `password_changed_at`, `remark`, `created_at`, `updated_at`, `deleted_at`) VALUES (5, 'guest1', '13800138004', 'guest1@example.com', '$2a$10$N.zmdr9k7uOCQb376NoUnuTJ8iKyqODYOUwjUxgvlMfCdWiw3Ca4a', '访客1', '访客', NULL, 0, NULL, 4, '实习生', 1, 0, NULL, NULL, 0, NULL, NULL, '2025-07-06 05:49:16', '2025-07-06 05:49:16', NULL);
COMMIT;

-- ----------------------------
-- Table structure for casbin_rule
-- ----------------------------
DROP TABLE IF EXISTS `casbin_rule`;
CREATE TABLE `casbin_rule` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `ptype` varchar(100) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `v0` varchar(100) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `v1` varchar(100) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `v2` varchar(100) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `v3` varchar(100) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `v4` varchar(100) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `v5` varchar(100) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_casbin_rule` (`ptype`,`v0`,`v1`,`v2`,`v3`),
  UNIQUE KEY `idx_casbin_rule` (`ptype`,`v0`,`v1`,`v2`,`v3`,`v4`,`v5`),
  KEY `idx_ptype` (`ptype`),
  KEY `idx_v0` (`v0`),
  KEY `idx_v1` (`v1`),
  KEY `idx_v2` (`v2`)
) ENGINE=InnoDB AUTO_INCREMENT=70 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='Casbin权限策略表';

-- ----------------------------
-- Records of casbin_rule
-- ----------------------------
BEGIN;
INSERT INTO `casbin_rule` (`id`, `ptype`, `v0`, `v1`, `v2`, `v3`, `v4`, `v5`, `created_at`, `updated_at`) VALUES (1, 'p', 'super_admin', '/api/admin/*', '*', NULL, NULL, NULL, '2025-07-06 05:58:28', '2025-07-06 05:58:28');
INSERT INTO `casbin_rule` (`id`, `ptype`, `v0`, `v1`, `v2`, `v3`, `v4`, `v5`, `created_at`, `updated_at`) VALUES (2, 'p', 'super_admin', '/api/users/*', '*', NULL, NULL, NULL, '2025-07-06 05:58:28', '2025-07-06 05:58:28');
INSERT INTO `casbin_rule` (`id`, `ptype`, `v0`, `v1`, `v2`, `v3`, `v4`, `v5`, `created_at`, `updated_at`) VALUES (3, 'p', 'super_admin', '/api/roles/*', '*', NULL, NULL, NULL, '2025-07-06 05:58:28', '2025-07-06 05:58:28');
INSERT INTO `casbin_rule` (`id`, `ptype`, `v0`, `v1`, `v2`, `v3`, `v4`, `v5`, `created_at`, `updated_at`) VALUES (4, 'p', 'super_admin', '/api/permissions/*', '*', NULL, NULL, NULL, '2025-07-06 05:58:28', '2025-07-06 05:58:28');
INSERT INTO `casbin_rule` (`id`, `ptype`, `v0`, `v1`, `v2`, `v3`, `v4`, `v5`, `created_at`, `updated_at`) VALUES (5, 'p', 'super_admin', '/api/departments/*', '*', NULL, NULL, NULL, '2025-07-06 05:58:28', '2025-07-06 05:58:28');
INSERT INTO `casbin_rule` (`id`, `ptype`, `v0`, `v1`, `v2`, `v3`, `v4`, `v5`, `created_at`, `updated_at`) VALUES (6, 'p', 'super_admin', '/api/configs/*', '*', NULL, NULL, NULL, '2025-07-06 05:58:28', '2025-07-06 05:58:28');
INSERT INTO `casbin_rule` (`id`, `ptype`, `v0`, `v1`, `v2`, `v3`, `v4`, `v5`, `created_at`, `updated_at`) VALUES (7, 'p', 'super_admin', '/api/logs/*', '*', NULL, NULL, NULL, '2025-07-06 05:58:28', '2025-07-06 05:58:28');
INSERT INTO `casbin_rule` (`id`, `ptype`, `v0`, `v1`, `v2`, `v3`, `v4`, `v5`, `created_at`, `updated_at`) VALUES (8, 'p', 'admin', '/api/users', 'GET', NULL, NULL, NULL, '2025-07-06 05:58:28', '2025-07-06 05:58:28');
INSERT INTO `casbin_rule` (`id`, `ptype`, `v0`, `v1`, `v2`, `v3`, `v4`, `v5`, `created_at`, `updated_at`) VALUES (9, 'p', 'admin', '/api/users', 'POST', NULL, NULL, NULL, '2025-07-06 05:58:28', '2025-07-06 05:58:28');
INSERT INTO `casbin_rule` (`id`, `ptype`, `v0`, `v1`, `v2`, `v3`, `v4`, `v5`, `created_at`, `updated_at`) VALUES (10, 'p', 'admin', '/api/users/*', 'GET', NULL, NULL, NULL, '2025-07-06 05:58:28', '2025-07-06 05:58:28');
INSERT INTO `casbin_rule` (`id`, `ptype`, `v0`, `v1`, `v2`, `v3`, `v4`, `v5`, `created_at`, `updated_at`) VALUES (11, 'p', 'admin', '/api/users/*', 'PUT', NULL, NULL, NULL, '2025-07-06 05:58:28', '2025-07-06 05:58:28');
INSERT INTO `casbin_rule` (`id`, `ptype`, `v0`, `v1`, `v2`, `v3`, `v4`, `v5`, `created_at`, `updated_at`) VALUES (12, 'p', 'admin', '/api/users/*', 'DELETE', NULL, NULL, NULL, '2025-07-06 05:58:28', '2025-07-06 05:58:28');
INSERT INTO `casbin_rule` (`id`, `ptype`, `v0`, `v1`, `v2`, `v3`, `v4`, `v5`, `created_at`, `updated_at`) VALUES (13, 'p', 'admin', '/api/roles', 'GET', NULL, NULL, NULL, '2025-07-06 05:58:28', '2025-07-06 05:58:28');
INSERT INTO `casbin_rule` (`id`, `ptype`, `v0`, `v1`, `v2`, `v3`, `v4`, `v5`, `created_at`, `updated_at`) VALUES (14, 'p', 'admin', '/api/roles/*', 'GET', NULL, NULL, NULL, '2025-07-06 05:58:28', '2025-07-06 05:58:28');
INSERT INTO `casbin_rule` (`id`, `ptype`, `v0`, `v1`, `v2`, `v3`, `v4`, `v5`, `created_at`, `updated_at`) VALUES (15, 'p', 'admin', '/api/permissions', 'GET', NULL, NULL, NULL, '2025-07-06 05:58:28', '2025-07-06 05:58:28');
INSERT INTO `casbin_rule` (`id`, `ptype`, `v0`, `v1`, `v2`, `v3`, `v4`, `v5`, `created_at`, `updated_at`) VALUES (16, 'p', 'admin', '/api/permissions/*', 'GET', NULL, NULL, NULL, '2025-07-06 05:58:28', '2025-07-06 05:58:28');
INSERT INTO `casbin_rule` (`id`, `ptype`, `v0`, `v1`, `v2`, `v3`, `v4`, `v5`, `created_at`, `updated_at`) VALUES (17, 'p', 'admin', '/api/departments', 'GET', NULL, NULL, NULL, '2025-07-06 05:58:28', '2025-07-06 05:58:28');
INSERT INTO `casbin_rule` (`id`, `ptype`, `v0`, `v1`, `v2`, `v3`, `v4`, `v5`, `created_at`, `updated_at`) VALUES (18, 'p', 'admin', '/api/departments/*', 'GET', NULL, NULL, NULL, '2025-07-06 05:58:28', '2025-07-06 05:58:28');
INSERT INTO `casbin_rule` (`id`, `ptype`, `v0`, `v1`, `v2`, `v3`, `v4`, `v5`, `created_at`, `updated_at`) VALUES (19, 'p', 'admin', '/api/logs', 'GET', NULL, NULL, NULL, '2025-07-06 05:58:28', '2025-07-06 05:58:28');
INSERT INTO `casbin_rule` (`id`, `ptype`, `v0`, `v1`, `v2`, `v3`, `v4`, `v5`, `created_at`, `updated_at`) VALUES (20, 'p', 'admin', '/api/logs/*', 'GET', NULL, NULL, NULL, '2025-07-06 05:58:28', '2025-07-06 05:58:28');
INSERT INTO `casbin_rule` (`id`, `ptype`, `v0`, `v1`, `v2`, `v3`, `v4`, `v5`, `created_at`, `updated_at`) VALUES (21, 'p', 'dept_admin', '/api/users', 'GET', NULL, NULL, NULL, '2025-07-06 05:58:28', '2025-07-06 05:58:28');
INSERT INTO `casbin_rule` (`id`, `ptype`, `v0`, `v1`, `v2`, `v3`, `v4`, `v5`, `created_at`, `updated_at`) VALUES (22, 'p', 'dept_admin', '/api/users/dept/*', 'GET', NULL, NULL, NULL, '2025-07-06 05:58:28', '2025-07-06 05:58:28');
INSERT INTO `casbin_rule` (`id`, `ptype`, `v0`, `v1`, `v2`, `v3`, `v4`, `v5`, `created_at`, `updated_at`) VALUES (23, 'p', 'dept_admin', '/api/users/dept/*', 'PUT', NULL, NULL, NULL, '2025-07-06 05:58:28', '2025-07-06 05:58:28');
INSERT INTO `casbin_rule` (`id`, `ptype`, `v0`, `v1`, `v2`, `v3`, `v4`, `v5`, `created_at`, `updated_at`) VALUES (24, 'p', 'dept_admin', '/api/departments', 'GET', NULL, NULL, NULL, '2025-07-06 05:58:28', '2025-07-06 05:58:28');
INSERT INTO `casbin_rule` (`id`, `ptype`, `v0`, `v1`, `v2`, `v3`, `v4`, `v5`, `created_at`, `updated_at`) VALUES (25, 'p', 'dept_admin', '/api/departments/*', 'GET', NULL, NULL, NULL, '2025-07-06 05:58:28', '2025-07-06 05:58:28');
INSERT INTO `casbin_rule` (`id`, `ptype`, `v0`, `v1`, `v2`, `v3`, `v4`, `v5`, `created_at`, `updated_at`) VALUES (26, 'p', 'dept_admin', '/api/roles', 'GET', NULL, NULL, NULL, '2025-07-06 05:58:28', '2025-07-06 05:58:28');
INSERT INTO `casbin_rule` (`id`, `ptype`, `v0`, `v1`, `v2`, `v3`, `v4`, `v5`, `created_at`, `updated_at`) VALUES (27, 'p', 'dept_admin', '/api/permissions', 'GET', NULL, NULL, NULL, '2025-07-06 05:58:28', '2025-07-06 05:58:28');
INSERT INTO `casbin_rule` (`id`, `ptype`, `v0`, `v1`, `v2`, `v3`, `v4`, `v5`, `created_at`, `updated_at`) VALUES (28, 'p', 'user', '/api/users/profile', 'GET', NULL, NULL, NULL, '2025-07-06 05:58:28', '2025-07-06 05:58:28');
INSERT INTO `casbin_rule` (`id`, `ptype`, `v0`, `v1`, `v2`, `v3`, `v4`, `v5`, `created_at`, `updated_at`) VALUES (29, 'p', 'user', '/api/users/profile', 'PUT', NULL, NULL, NULL, '2025-07-06 05:58:28', '2025-07-06 05:58:28');
INSERT INTO `casbin_rule` (`id`, `ptype`, `v0`, `v1`, `v2`, `v3`, `v4`, `v5`, `created_at`, `updated_at`) VALUES (30, 'p', 'user', '/api/departments', 'GET', NULL, NULL, NULL, '2025-07-06 05:58:28', '2025-07-06 05:58:28');
INSERT INTO `casbin_rule` (`id`, `ptype`, `v0`, `v1`, `v2`, `v3`, `v4`, `v5`, `created_at`, `updated_at`) VALUES (31, 'p', 'user', '/api/departments/*', 'GET', NULL, NULL, NULL, '2025-07-06 05:58:28', '2025-07-06 05:58:28');
INSERT INTO `casbin_rule` (`id`, `ptype`, `v0`, `v1`, `v2`, `v3`, `v4`, `v5`, `created_at`, `updated_at`) VALUES (32, 'p', 'guest', '/api/users/profile', 'GET', NULL, NULL, NULL, '2025-07-06 05:58:28', '2025-07-06 05:58:28');
INSERT INTO `casbin_rule` (`id`, `ptype`, `v0`, `v1`, `v2`, `v3`, `v4`, `v5`, `created_at`, `updated_at`) VALUES (33, 'p', 'super_admin', '/api/admin/*', '*', NULL, NULL, NULL, '2025-07-06 05:58:37', '2025-07-06 05:58:37');
INSERT INTO `casbin_rule` (`id`, `ptype`, `v0`, `v1`, `v2`, `v3`, `v4`, `v5`, `created_at`, `updated_at`) VALUES (34, 'p', 'super_admin', '/api/users/*', '*', NULL, NULL, NULL, '2025-07-06 05:58:37', '2025-07-06 05:58:37');
INSERT INTO `casbin_rule` (`id`, `ptype`, `v0`, `v1`, `v2`, `v3`, `v4`, `v5`, `created_at`, `updated_at`) VALUES (35, 'p', 'super_admin', '/api/roles/*', '*', NULL, NULL, NULL, '2025-07-06 05:58:37', '2025-07-06 05:58:37');
INSERT INTO `casbin_rule` (`id`, `ptype`, `v0`, `v1`, `v2`, `v3`, `v4`, `v5`, `created_at`, `updated_at`) VALUES (36, 'p', 'super_admin', '/api/permissions/*', '*', NULL, NULL, NULL, '2025-07-06 05:58:37', '2025-07-06 05:58:37');
INSERT INTO `casbin_rule` (`id`, `ptype`, `v0`, `v1`, `v2`, `v3`, `v4`, `v5`, `created_at`, `updated_at`) VALUES (37, 'p', 'super_admin', '/api/departments/*', '*', NULL, NULL, NULL, '2025-07-06 05:58:37', '2025-07-06 05:58:37');
INSERT INTO `casbin_rule` (`id`, `ptype`, `v0`, `v1`, `v2`, `v3`, `v4`, `v5`, `created_at`, `updated_at`) VALUES (38, 'p', 'super_admin', '/api/configs/*', '*', NULL, NULL, NULL, '2025-07-06 05:58:37', '2025-07-06 05:58:37');
INSERT INTO `casbin_rule` (`id`, `ptype`, `v0`, `v1`, `v2`, `v3`, `v4`, `v5`, `created_at`, `updated_at`) VALUES (39, 'p', 'super_admin', '/api/logs/*', '*', NULL, NULL, NULL, '2025-07-06 05:58:37', '2025-07-06 05:58:37');
INSERT INTO `casbin_rule` (`id`, `ptype`, `v0`, `v1`, `v2`, `v3`, `v4`, `v5`, `created_at`, `updated_at`) VALUES (40, 'p', 'admin', '/api/users', 'GET', NULL, NULL, NULL, '2025-07-06 05:58:37', '2025-07-06 05:58:37');
INSERT INTO `casbin_rule` (`id`, `ptype`, `v0`, `v1`, `v2`, `v3`, `v4`, `v5`, `created_at`, `updated_at`) VALUES (41, 'p', 'admin', '/api/users', 'POST', NULL, NULL, NULL, '2025-07-06 05:58:37', '2025-07-06 05:58:37');
INSERT INTO `casbin_rule` (`id`, `ptype`, `v0`, `v1`, `v2`, `v3`, `v4`, `v5`, `created_at`, `updated_at`) VALUES (42, 'p', 'admin', '/api/users/*', 'GET', NULL, NULL, NULL, '2025-07-06 05:58:37', '2025-07-06 05:58:37');
INSERT INTO `casbin_rule` (`id`, `ptype`, `v0`, `v1`, `v2`, `v3`, `v4`, `v5`, `created_at`, `updated_at`) VALUES (43, 'p', 'admin', '/api/users/*', 'PUT', NULL, NULL, NULL, '2025-07-06 05:58:37', '2025-07-06 05:58:37');
INSERT INTO `casbin_rule` (`id`, `ptype`, `v0`, `v1`, `v2`, `v3`, `v4`, `v5`, `created_at`, `updated_at`) VALUES (44, 'p', 'admin', '/api/users/*', 'DELETE', NULL, NULL, NULL, '2025-07-06 05:58:37', '2025-07-06 05:58:37');
INSERT INTO `casbin_rule` (`id`, `ptype`, `v0`, `v1`, `v2`, `v3`, `v4`, `v5`, `created_at`, `updated_at`) VALUES (45, 'p', 'admin', '/api/roles', 'GET', NULL, NULL, NULL, '2025-07-06 05:58:37', '2025-07-06 05:58:37');
INSERT INTO `casbin_rule` (`id`, `ptype`, `v0`, `v1`, `v2`, `v3`, `v4`, `v5`, `created_at`, `updated_at`) VALUES (46, 'p', 'admin', '/api/roles/*', 'GET', NULL, NULL, NULL, '2025-07-06 05:58:37', '2025-07-06 05:58:37');
INSERT INTO `casbin_rule` (`id`, `ptype`, `v0`, `v1`, `v2`, `v3`, `v4`, `v5`, `created_at`, `updated_at`) VALUES (47, 'p', 'admin', '/api/permissions', 'GET', NULL, NULL, NULL, '2025-07-06 05:58:37', '2025-07-06 05:58:37');
INSERT INTO `casbin_rule` (`id`, `ptype`, `v0`, `v1`, `v2`, `v3`, `v4`, `v5`, `created_at`, `updated_at`) VALUES (48, 'p', 'admin', '/api/permissions/*', 'GET', NULL, NULL, NULL, '2025-07-06 05:58:37', '2025-07-06 05:58:37');
INSERT INTO `casbin_rule` (`id`, `ptype`, `v0`, `v1`, `v2`, `v3`, `v4`, `v5`, `created_at`, `updated_at`) VALUES (49, 'p', 'admin', '/api/departments', 'GET', NULL, NULL, NULL, '2025-07-06 05:58:37', '2025-07-06 05:58:37');
INSERT INTO `casbin_rule` (`id`, `ptype`, `v0`, `v1`, `v2`, `v3`, `v4`, `v5`, `created_at`, `updated_at`) VALUES (50, 'p', 'admin', '/api/departments/*', 'GET', NULL, NULL, NULL, '2025-07-06 05:58:37', '2025-07-06 05:58:37');
INSERT INTO `casbin_rule` (`id`, `ptype`, `v0`, `v1`, `v2`, `v3`, `v4`, `v5`, `created_at`, `updated_at`) VALUES (51, 'p', 'admin', '/api/logs', 'GET', NULL, NULL, NULL, '2025-07-06 05:58:37', '2025-07-06 05:58:37');
INSERT INTO `casbin_rule` (`id`, `ptype`, `v0`, `v1`, `v2`, `v3`, `v4`, `v5`, `created_at`, `updated_at`) VALUES (52, 'p', 'admin', '/api/logs/*', 'GET', NULL, NULL, NULL, '2025-07-06 05:58:37', '2025-07-06 05:58:37');
INSERT INTO `casbin_rule` (`id`, `ptype`, `v0`, `v1`, `v2`, `v3`, `v4`, `v5`, `created_at`, `updated_at`) VALUES (53, 'p', 'dept_admin', '/api/users', 'GET', NULL, NULL, NULL, '2025-07-06 05:58:37', '2025-07-06 05:58:37');
INSERT INTO `casbin_rule` (`id`, `ptype`, `v0`, `v1`, `v2`, `v3`, `v4`, `v5`, `created_at`, `updated_at`) VALUES (54, 'p', 'dept_admin', '/api/users/dept/*', 'GET', NULL, NULL, NULL, '2025-07-06 05:58:37', '2025-07-06 05:58:37');
INSERT INTO `casbin_rule` (`id`, `ptype`, `v0`, `v1`, `v2`, `v3`, `v4`, `v5`, `created_at`, `updated_at`) VALUES (55, 'p', 'dept_admin', '/api/users/dept/*', 'PUT', NULL, NULL, NULL, '2025-07-06 05:58:37', '2025-07-06 05:58:37');
INSERT INTO `casbin_rule` (`id`, `ptype`, `v0`, `v1`, `v2`, `v3`, `v4`, `v5`, `created_at`, `updated_at`) VALUES (56, 'p', 'dept_admin', '/api/departments', 'GET', NULL, NULL, NULL, '2025-07-06 05:58:37', '2025-07-06 05:58:37');
INSERT INTO `casbin_rule` (`id`, `ptype`, `v0`, `v1`, `v2`, `v3`, `v4`, `v5`, `created_at`, `updated_at`) VALUES (57, 'p', 'dept_admin', '/api/departments/*', 'GET', NULL, NULL, NULL, '2025-07-06 05:58:37', '2025-07-06 05:58:37');
INSERT INTO `casbin_rule` (`id`, `ptype`, `v0`, `v1`, `v2`, `v3`, `v4`, `v5`, `created_at`, `updated_at`) VALUES (58, 'p', 'dept_admin', '/api/roles', 'GET', NULL, NULL, NULL, '2025-07-06 05:58:37', '2025-07-06 05:58:37');
INSERT INTO `casbin_rule` (`id`, `ptype`, `v0`, `v1`, `v2`, `v3`, `v4`, `v5`, `created_at`, `updated_at`) VALUES (59, 'p', 'dept_admin', '/api/permissions', 'GET', NULL, NULL, NULL, '2025-07-06 05:58:37', '2025-07-06 05:58:37');
INSERT INTO `casbin_rule` (`id`, `ptype`, `v0`, `v1`, `v2`, `v3`, `v4`, `v5`, `created_at`, `updated_at`) VALUES (60, 'p', 'user', '/api/users/profile', 'GET', NULL, NULL, NULL, '2025-07-06 05:58:37', '2025-07-06 05:58:37');
INSERT INTO `casbin_rule` (`id`, `ptype`, `v0`, `v1`, `v2`, `v3`, `v4`, `v5`, `created_at`, `updated_at`) VALUES (61, 'p', 'user', '/api/users/profile', 'PUT', NULL, NULL, NULL, '2025-07-06 05:58:37', '2025-07-06 05:58:37');
INSERT INTO `casbin_rule` (`id`, `ptype`, `v0`, `v1`, `v2`, `v3`, `v4`, `v5`, `created_at`, `updated_at`) VALUES (62, 'p', 'user', '/api/departments', 'GET', NULL, NULL, NULL, '2025-07-06 05:58:37', '2025-07-06 05:58:37');
INSERT INTO `casbin_rule` (`id`, `ptype`, `v0`, `v1`, `v2`, `v3`, `v4`, `v5`, `created_at`, `updated_at`) VALUES (63, 'p', 'user', '/api/departments/*', 'GET', NULL, NULL, NULL, '2025-07-06 05:58:37', '2025-07-06 05:58:37');
INSERT INTO `casbin_rule` (`id`, `ptype`, `v0`, `v1`, `v2`, `v3`, `v4`, `v5`, `created_at`, `updated_at`) VALUES (64, 'p', 'guest', '/api/users/profile', 'GET', NULL, NULL, NULL, '2025-07-06 05:58:37', '2025-07-06 05:58:37');
INSERT INTO `casbin_rule` (`id`, `ptype`, `v0`, `v1`, `v2`, `v3`, `v4`, `v5`, `created_at`, `updated_at`) VALUES (65, 'g', 'admin', 'super_admin', '', NULL, NULL, NULL, '2025-07-06 05:59:41', '2025-07-06 05:59:41');
INSERT INTO `casbin_rule` (`id`, `ptype`, `v0`, `v1`, `v2`, `v3`, `v4`, `v5`, `created_at`, `updated_at`) VALUES (66, 'g', 'user1', 'admin', '', NULL, NULL, NULL, '2025-07-06 05:59:41', '2025-07-06 05:59:41');
INSERT INTO `casbin_rule` (`id`, `ptype`, `v0`, `v1`, `v2`, `v3`, `v4`, `v5`, `created_at`, `updated_at`) VALUES (67, 'g', 'user2', 'dept_admin', '', NULL, NULL, NULL, '2025-07-06 05:59:41', '2025-07-06 05:59:41');
INSERT INTO `casbin_rule` (`id`, `ptype`, `v0`, `v1`, `v2`, `v3`, `v4`, `v5`, `created_at`, `updated_at`) VALUES (68, 'g', 'user3', 'user', '', NULL, NULL, NULL, '2025-07-06 05:59:41', '2025-07-06 05:59:41');
INSERT INTO `casbin_rule` (`id`, `ptype`, `v0`, `v1`, `v2`, `v3`, `v4`, `v5`, `created_at`, `updated_at`) VALUES (69, 'g', 'guest1', 'guest', '', NULL, NULL, NULL, '2025-07-06 05:59:41', '2025-07-06 05:59:41');
COMMIT;

SET FOREIGN_KEY_CHECKS = 1;
