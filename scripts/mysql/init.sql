-- 鲲鹏后台管理系统数据库初始化脚本

-- 使用数据库
USE kunpeng;

-- 设置字符集
SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- 创建用户表
CREATE TABLE IF NOT EXISTS `kp_user` (
  `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '用户ID',
  `username` varchar(50) NOT NULL COMMENT '用户名',
  `password` varchar(100) NOT NULL COMMENT '密码',
  `nickname` varchar(50) DEFAULT NULL COMMENT '昵称',
  `real_name` varchar(50) DEFAULT NULL COMMENT '真实姓名',
  `avatar` varchar(255) DEFAULT NULL COMMENT '头像',
  `gender` tinyint(1) DEFAULT 0 COMMENT '性别(0:未知 1:男 2:女)',
  `email` varchar(100) DEFAULT NULL COMMENT '邮箱',
  `mobile` varchar(20) DEFAULT NULL COMMENT '手机号',
  `dept_id` bigint(20) UNSIGNED DEFAULT NULL COMMENT '部门ID',
  `post_id` bigint(20) UNSIGNED DEFAULT NULL COMMENT '岗位ID',
  `role_id` bigint(20) UNSIGNED DEFAULT NULL COMMENT '角色ID',
  `status` tinyint(1) DEFAULT 1 COMMENT '状态(0:禁用 1:启用)',
  `login_ip` varchar(50) DEFAULT NULL COMMENT '最后登录IP',
  `login_time` datetime DEFAULT NULL COMMENT '最后登录时间',
  `app_key` varchar(50) DEFAULT NULL COMMENT 'AppKey',
  `app_secret` varchar(100) DEFAULT NULL COMMENT 'AppSecret',
  `remark` varchar(255) DEFAULT NULL COMMENT '备注',
  `created_at` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `deleted_at` datetime DEFAULT NULL COMMENT '删除时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_username` (`username`),
  KEY `idx_dept_id` (`dept_id`),
  KEY `idx_post_id` (`post_id`),
  KEY `idx_role_id` (`role_id`),
  KEY `idx_status` (`status`),
  KEY `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户表';

-- 创建角色表
CREATE TABLE IF NOT EXISTS `kp_role` (
  `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '角色ID',
  `name` varchar(50) NOT NULL COMMENT '角色名称',
  `code` varchar(50) NOT NULL COMMENT '角色编码',
  `sort` int(11) DEFAULT 0 COMMENT '排序',
  `status` tinyint(1) DEFAULT 1 COMMENT '状态(0:禁用 1:启用)',
  `remark` varchar(255) DEFAULT NULL COMMENT '备注',
  `created_at` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `deleted_at` datetime DEFAULT NULL COMMENT '删除时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_code` (`code`),
  KEY `idx_status` (`status`),
  KEY `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='角色表';

-- 创建菜单表
CREATE TABLE IF NOT EXISTS `kp_menu` (
  `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '菜单ID',
  `parent_id` bigint(20) UNSIGNED DEFAULT 0 COMMENT '父菜单ID',
  `name` varchar(50) NOT NULL COMMENT '菜单名称',
  `type` tinyint(1) DEFAULT 0 COMMENT '类型(0:目录 1:菜单 2:按钮)',
  `path` varchar(100) DEFAULT NULL COMMENT '路由地址',
  `component` varchar(100) DEFAULT NULL COMMENT '组件路径',
  `permission` varchar(100) DEFAULT NULL COMMENT '权限标识',
  `icon` varchar(100) DEFAULT NULL COMMENT '图标',
  `sort` int(11) DEFAULT 0 COMMENT '排序',
  `visible` tinyint(1) DEFAULT 1 COMMENT '是否可见(0:隐藏 1:显示)',
  `status` tinyint(1) DEFAULT 1 COMMENT '状态(0:禁用 1:启用)',
  `is_cache` tinyint(1) DEFAULT 0 COMMENT '是否缓存(0:不缓存 1:缓存)',
  `is_frame` tinyint(1) DEFAULT 0 COMMENT '是否外链(0:否 1:是)',
  `remark` varchar(255) DEFAULT NULL COMMENT '备注',
  `created_at` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `deleted_at` datetime DEFAULT NULL COMMENT '删除时间',
  PRIMARY KEY (`id`),
  KEY `idx_parent_id` (`parent_id`),
  KEY `idx_status` (`status`),
  KEY `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='菜单表';

-- 创建角色菜单关联表
CREATE TABLE IF NOT EXISTS `kp_role_menu` (
  `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `role_id` bigint(20) UNSIGNED NOT NULL COMMENT '角色ID',
  `menu_id` bigint(20) UNSIGNED NOT NULL COMMENT '菜单ID',
  `created_at` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_role_menu` (`role_id`,`menu_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='角色菜单关联表';

-- 创建API表
CREATE TABLE IF NOT EXISTS `kp_api` (
  `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT 'API ID',
  `group` varchar(50) NOT NULL COMMENT 'API分组',
  `name` varchar(100) NOT NULL COMMENT 'API名称',
  `method` varchar(10) NOT NULL COMMENT '请求方法',
  `path` varchar(100) NOT NULL COMMENT '请求路径',
  `status` tinyint(1) DEFAULT 1 COMMENT '状态(0:禁用 1:启用)',
  `remark` varchar(255) DEFAULT NULL COMMENT '备注',
  `created_at` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `deleted_at` datetime DEFAULT NULL COMMENT '删除时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_method_path` (`method`,`path`),
  KEY `idx_group` (`group`),
  KEY `idx_status` (`status`),
  KEY `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='API表';

-- 创建角色API关联表
CREATE TABLE IF NOT EXISTS `kp_role_api` (
  `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `role_id` bigint(20) UNSIGNED NOT NULL COMMENT '角色ID',
  `api_id` bigint(20) UNSIGNED NOT NULL COMMENT 'API ID',
  `created_at` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_role_api` (`role_id`,`api_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='角色API关联表';

-- 创建部门表
CREATE TABLE IF NOT EXISTS `kp_dept` (
  `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '部门ID',
  `parent_id` bigint(20) UNSIGNED DEFAULT 0 COMMENT '父部门ID',
  `name` varchar(50) NOT NULL COMMENT '部门名称',
  `leader` varchar(50) DEFAULT NULL COMMENT '负责人',
  `phone` varchar(20) DEFAULT NULL COMMENT '联系电话',
  `email` varchar(100) DEFAULT NULL COMMENT '邮箱',
  `sort` int(11) DEFAULT 0 COMMENT '排序',
  `status` tinyint(1) DEFAULT 1 COMMENT '状态(0:禁用 1:启用)',
  `remark` varchar(255) DEFAULT NULL COMMENT '备注',
  `created_at` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `deleted_at` datetime DEFAULT NULL COMMENT '删除时间',
  PRIMARY KEY (`id`),
  KEY `idx_parent_id` (`parent_id`),
  KEY `idx_status` (`status`),
  KEY `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='部门表';

-- 创建岗位表
CREATE TABLE IF NOT EXISTS `kp_post` (
  `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '岗位ID',
  `name` varchar(50) NOT NULL COMMENT '岗位名称',
  `code` varchar(50) NOT NULL COMMENT '岗位编码',
  `sort` int(11) DEFAULT 0 COMMENT '排序',
  `status` tinyint(1) DEFAULT 1 COMMENT '状态(0:禁用 1:启用)',
  `remark` varchar(255) DEFAULT NULL COMMENT '备注',
  `created_at` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `deleted_at` datetime DEFAULT NULL COMMENT '删除时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_code` (`code`),
  KEY `idx_status` (`status`),
  KEY `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='岗位表';

-- 创建操作日志表
CREATE TABLE IF NOT EXISTS `kp_operation_log` (
  `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '日志ID',
  `user_id` bigint(20) UNSIGNED DEFAULT NULL COMMENT '用户ID',
  `username` varchar(50) DEFAULT NULL COMMENT '用户名',
  `module` varchar(50) DEFAULT NULL COMMENT '模块名称',
  `action` varchar(50) DEFAULT NULL COMMENT '操作类型',
  `method` varchar(10) DEFAULT NULL COMMENT '请求方法',
  `path` varchar(100) DEFAULT NULL COMMENT '请求路径',
  `ip` varchar(50) DEFAULT NULL COMMENT 'IP地址',
  `user_agent` varchar(255) DEFAULT NULL COMMENT '用户代理',
  `request` text DEFAULT NULL COMMENT '请求参数',
  `response` text DEFAULT NULL COMMENT '响应结果',
  `status` tinyint(1) DEFAULT 1 COMMENT '状态(0:失败 1:成功)',
  `error_message` text DEFAULT NULL COMMENT '错误信息',
  `duration` int(11) DEFAULT 0 COMMENT '执行时长(ms)',
  `created_at` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  PRIMARY KEY (`id`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_username` (`username`),
  KEY `idx_module` (`module`),
  KEY `idx_action` (`action`),
  KEY `idx_status` (`status`),
  KEY `idx_created_at` (`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='操作日志表';

-- 创建登录日志表
CREATE TABLE IF NOT EXISTS `kp_login_log` (
  `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '日志ID',
  `user_id` bigint(20) UNSIGNED DEFAULT NULL COMMENT '用户ID',
  `username` varchar(50) DEFAULT NULL COMMENT '用户名',
  `ip` varchar(50) DEFAULT NULL COMMENT 'IP地址',
  `location` varchar(100) DEFAULT NULL COMMENT '地理位置',
  `browser` varchar(50) DEFAULT NULL COMMENT '浏览器',
  `os` varchar(50) DEFAULT NULL COMMENT '操作系统',
  `device` varchar(50) DEFAULT NULL COMMENT '设备',
  `status` tinyint(1) DEFAULT 1 COMMENT '状态(0:失败 1:成功)',
  `message` varchar(255) DEFAULT NULL COMMENT '消息描述',
  `created_at` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  PRIMARY KEY (`id`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_username` (`username`),
  KEY `idx_ip` (`ip`),
  KEY `idx_status` (`status`),
  KEY `idx_created_at` (`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='登录日志表';

-- 创建登录尝试记录表
CREATE TABLE IF NOT EXISTS `kp_login_attempt` (
  `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '记录ID',
  `account` varchar(100) NOT NULL COMMENT '登录账号（用户名或手机号）',
  `ip` varchar(50) NOT NULL COMMENT '登录IP',
  `attempts` int(11) DEFAULT 1 COMMENT '失败次数',
  `last_try` datetime NOT NULL COMMENT '最后尝试时间',
  `blocked_at` datetime DEFAULT NULL COMMENT '拉黑时间',
  `created_at` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `deleted_at` datetime DEFAULT NULL COMMENT '删除时间',
  PRIMARY KEY (`id`),
  KEY `idx_account` (`account`),
  KEY `idx_ip` (`ip`),
  KEY `idx_account_ip` (`account`, `ip`),
  KEY `idx_blocked_at` (`blocked_at`),
  KEY `idx_created_at` (`created_at`),
  KEY `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='登录尝试记录表';

-- 创建系统配置表
CREATE TABLE IF NOT EXISTS `kp_config` (
  `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '配置ID',
  `name` varchar(50) NOT NULL COMMENT '配置名称',
  `key` varchar(50) NOT NULL COMMENT '配置键',
  `value` text NOT NULL COMMENT '配置值',
  `type` varchar(20) DEFAULT 'string' COMMENT '配置类型',
  `status` tinyint(1) DEFAULT 1 COMMENT '状态(0:禁用 1:启用)',
  `remark` varchar(255) DEFAULT NULL COMMENT '备注',
  `created_at` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `deleted_at` datetime DEFAULT NULL COMMENT '删除时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_key` (`key`),
  KEY `idx_status` (`status`),
  KEY `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='系统配置表';

-- 创建字典类型表
CREATE TABLE IF NOT EXISTS `kp_dict_type` (
  `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '字典类型ID',
  `name` varchar(50) NOT NULL COMMENT '字典名称',
  `type` varchar(50) NOT NULL COMMENT '字典类型',
  `status` tinyint(1) DEFAULT 1 COMMENT '状态(0:禁用 1:启用)',
  `remark` varchar(255) DEFAULT NULL COMMENT '备注',
  `created_at` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `deleted_at` datetime DEFAULT NULL COMMENT '删除时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_type` (`type`),
  KEY `idx_status` (`status`),
  KEY `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='字典类型表';

-- 创建字典数据表
CREATE TABLE IF NOT EXISTS `kp_dict_data` (
  `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '字典数据ID',
  `dict_type` varchar(50) NOT NULL COMMENT '字典类型',
  `label` varchar(50) NOT NULL COMMENT '字典标签',
  `value` varchar(50) NOT NULL COMMENT '字典值',
  `sort` int(11) DEFAULT 0 COMMENT '排序',
  `status` tinyint(1) DEFAULT 1 COMMENT '状态(0:禁用 1:启用)',
  `css_class` varchar(100) DEFAULT NULL COMMENT 'CSS类名',
  `list_class` varchar(100) DEFAULT NULL COMMENT '表格回显样式',
  `is_default` tinyint(1) DEFAULT 0 COMMENT '是否默认(0:否 1:是)',
  `remark` varchar(255) DEFAULT NULL COMMENT '备注',
  `created_at` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `deleted_at` datetime DEFAULT NULL COMMENT '删除时间',
  PRIMARY KEY (`id`),
  KEY `idx_dict_type` (`dict_type`),
  KEY `idx_status` (`status`),
  KEY `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='字典数据表';

-- 插入初始数据

-- 插入管理员用户
INSERT INTO `kp_user` (`id`, `username`, `password`, `nickname`, `real_name`, `avatar`, `gender`, `email`, `mobile`, `dept_id`, `post_id`, `role_id`, `status`, `app_key`, `app_secret`, `remark`) VALUES
(1, 'admin', '$2a$10$YEzOYVCz6jBhwgCHJEQXG.0/FROxhA/MxQYV0F1hUWvtQgV1CvZT.', '管理员', '系统管理员', NULL, 1, 'admin@example.com', '13800138000', 1, 1, 1, 1, 'admin', 'c5e330214fb33e2d485f207b33e4c92f', '系统管理员');

-- 插入角色
INSERT INTO `kp_role` (`id`, `name`, `code`, `sort`, `status`, `remark`) VALUES
(1, '超级管理员', 'admin', 1, 1, '超级管理员'),
(2, '普通用户', 'user', 2, 1, '普通用户');

-- 插入部门
INSERT INTO `kp_dept` (`id`, `parent_id`, `name`, `leader`, `phone`, `email`, `sort`, `status`, `remark`) VALUES
(1, 0, '总公司', '张三', '13800138001', 'zhangsan@example.com', 1, 1, '总公司'),
(2, 1, '研发部', '李四', '13800138002', 'lisi@example.com', 1, 1, '研发部'),
(3, 1, '市场部', '王五', '13800138003', 'wangwu@example.com', 2, 1, '市场部'),
(4, 1, '财务部', '赵六', '13800138004', 'zhaoliu@example.com', 3, 1, '财务部');

-- 插入岗位
INSERT INTO `kp_post` (`id`, `name`, `code`, `sort`, `status`, `remark`) VALUES
(1, '董事长', 'ceo', 1, 1, '董事长'),
(2, '技术总监', 'cto', 2, 1, '技术总监'),
(3, '项目经理', 'pm', 3, 1, '项目经理'),
(4, '高级工程师', 'se', 4, 1, '高级工程师'),
(5, '中级工程师', 'sde', 5, 1, '中级工程师'),
(6, '初级工程师', 'jde', 6, 1, '初级工程师');

-- 插入菜单
INSERT INTO `kp_menu` (`id`, `parent_id`, `name`, `type`, `path`, `component`, `permission`, `icon`, `sort`, `visible`, `status`, `is_cache`, `is_frame`) VALUES
(1, 0, '系统管理', 0, '/system', NULL, NULL, 'system', 1, 1, 1, 0, 0),
(2, 1, '用户管理', 1, 'user', 'system/user/index', 'system:user:list', 'user', 1, 1, 1, 0, 0),
(3, 1, '角色管理', 1, 'role', 'system/role/index', 'system:role:list', 'peoples', 2, 1, 1, 0, 0),
(4, 1, '菜单管理', 1, 'menu', 'system/menu/index', 'system:menu:list', 'tree-table', 3, 1, 1, 0, 0),
(5, 1, '部门管理', 1, 'dept', 'system/dept/index', 'system:dept:list', 'tree', 4, 1, 1, 0, 0),
(6, 1, '岗位管理', 1, 'post', 'system/post/index', 'system:post:list', 'post', 5, 1, 1, 0, 0),
(7, 1, '字典管理', 1, 'dict', 'system/dict/index', 'system:dict:list', 'dict', 6, 1, 1, 0, 0),
(8, 1, '参数设置', 1, 'config', 'system/config/index', 'system:config:list', 'edit', 7, 1, 1, 0, 0),
(9, 1, '日志管理', 0, 'log', NULL, NULL, 'log', 8, 1, 1, 0, 0),
(10, 9, '操作日志', 1, 'operlog', 'monitor/operlog/index', 'monitor:operlog:list', 'form', 1, 1, 1, 0, 0),
(11, 9, '登录日志', 1, 'logininfor', 'monitor/logininfor/index', 'monitor:logininfor:list', 'logininfor', 2, 1, 1, 0, 0);

-- 插入角色菜单关联
INSERT INTO `kp_role_menu` (`role_id`, `menu_id`) VALUES
(1, 1), (1, 2), (1, 3), (1, 4), (1, 5), (1, 6), (1, 7), (1, 8), (1, 9), (1, 10), (1, 11),
(2, 1), (2, 2), (2, 3), (2, 4), (2, 5), (2, 6), (2, 7);

-- 插入API
INSERT INTO `kp_api` (`id`, `group`, `name`, `method`, `path`, `status`, `remark`) VALUES
(1, '用户管理', '获取用户列表', 'GET', '/api/v1/users', 1, '获取用户列表'),
(2, '用户管理', '获取用户详情', 'GET', '/api/v1/users/:id', 1, '获取用户详情'),
(3, '用户管理', '创建用户', 'POST', '/api/v1/users', 1, '创建用户'),
(4, '用户管理', '更新用户', 'PUT', '/api/v1/users', 1, '更新用户'),
(5, '用户管理', '删除用户', 'DELETE', '/api/v1/users/:id', 1, '删除用户'),
(6, '角色管理', '获取角色列表', 'GET', '/api/v1/roles', 1, '获取角色列表'),
(7, '角色管理', '获取角色详情', 'GET', '/api/v1/roles/:id', 1, '获取角色详情'),
(8, '角色管理', '创建角色', 'POST', '/api/v1/roles', 1, '创建角色'),
(9, '角色管理', '更新角色', 'PUT', '/api/v1/roles', 1, '更新角色'),
(10, '角色管理', '删除角色', 'DELETE', '/api/v1/roles/:id', 1, '删除角色');

-- 插入角色API关联
INSERT INTO `kp_role_api` (`role_id`, `api_id`) VALUES
(1, 1), (1, 2), (1, 3), (1, 4), (1, 5), (1, 6), (1, 7), (1, 8), (1, 9), (1, 10),
(2, 1), (2, 2), (2, 6), (2, 7);

-- 插入字典类型
INSERT INTO `kp_dict_type` (`id`, `name`, `type`, `status`, `remark`) VALUES
(1, '用户性别', 'sys_user_gender', 1, '用户性别列表'),
(2, '菜单状态', 'sys_menu_status', 1, '菜单状态列表'),
(3, '系统开关', 'sys_switch', 1, '系统开关列表'),
(4, '任务状态', 'sys_job_status', 1, '任务状态列表'),
(5, '任务分组', 'sys_job_group', 1, '任务分组列表'),
(6, '系统是否', 'sys_yes_no', 1, '系统是否列表');

-- 插入字典数据
INSERT INTO `kp_dict_data` (`dict_type`, `label`, `value`, `sort`, `status`, `css_class`, `list_class`, `is_default`, `remark`) VALUES
('sys_user_gender', '男', '1', 1, 1, NULL, 'default', 1, '性别男'),
('sys_user_gender', '女', '2', 2, 1, NULL, 'default', 0, '性别女'),
('sys_user_gender', '未知', '0', 3, 1, NULL, 'default', 0, '性别未知'),
('sys_menu_status', '显示', '1', 1, 1, NULL, 'success', 1, '显示菜单'),
('sys_menu_status', '隐藏', '0', 2, 1, NULL, 'danger', 0, '隐藏菜单'),
('sys_switch', '开启', '1', 1, 1, NULL, 'success', 1, '开启状态'),
('sys_switch', '关闭', '0', 2, 1, NULL, 'danger', 0, '关闭状态'),
('sys_job_status', '正常', '1', 1, 1, NULL, 'success', 1, '正常状态'),
('sys_job_status', '暂停', '0', 2, 1, NULL, 'danger', 0, '暂停状态'),
('sys_job_group', '默认', 'DEFAULT', 1, 1, NULL, 'default', 1, '默认分组'),
('sys_job_group', '系统', 'SYSTEM', 2, 1, NULL, 'default', 0, '系统分组'),
('sys_yes_no', '是', '1', 1, 1, NULL, 'success', 1, '系统是'),
('sys_yes_no', '否', '0', 2, 1, NULL, 'danger', 0, '系统否');

-- 创建token黑名单表
CREATE TABLE IF NOT EXISTS `kp_token_blacklist` (
  `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `token` text NOT NULL COMMENT 'JWT token',
  `user_id` bigint(20) UNSIGNED NOT NULL COMMENT '用户ID',
  `username` varchar(50) NOT NULL COMMENT '用户名',
  `reason` varchar(255) DEFAULT NULL COMMENT '加入黑名单原因',
  `expires_at` datetime NOT NULL COMMENT 'token过期时间',
  `created_at` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `deleted_at` datetime DEFAULT NULL COMMENT '删除时间',
  PRIMARY KEY (`id`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_username` (`username`),
  KEY `idx_expires_at` (`expires_at`),
  KEY `idx_created_at` (`created_at`),
  KEY `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='token黑名单表';

-- 设置外键检查
SET FOREIGN_KEY_CHECKS = 1;