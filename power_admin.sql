/*
 Navicat Premium Dump SQL

 Source Server         : localhost
 Source Server Type    : MySQL
 Source Server Version : 80012 (8.0.12)
 Source Host           : localhost:3306
 Source Schema         : power_admin

 Target Server Type    : MySQL
 Target Server Version : 80012 (8.0.12)
 File Encoding         : 65001

 Date: 01/12/2025 00:17:57
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for apis
-- ----------------------------
DROP TABLE IF EXISTS `apis`;
CREATE TABLE `apis`  (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT 'API ID',
  `api_name` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT 'API名称',
  `api_path` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT 'API路径',
  `api_method` varchar(10) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT 'HTTP方法 GET, POST, PUT, DELETE等',
  `group` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT 'API分组',
  `description` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT 'API描述',
  `permission_id` bigint(20) NULL DEFAULT NULL COMMENT '关联权限ID',
  `status` bigint(20) NULL DEFAULT 1 COMMENT '状态 1:启用 0:禁用',
  `remark` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '备注',
  `created_at` datetime(3) NULL DEFAULT NULL COMMENT '创建时间',
  `updated_at` datetime(3) NULL DEFAULT NULL COMMENT '更新时间',
  `deleted_at` datetime(3) NULL DEFAULT NULL COMMENT '删除时间',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `uk_path_method`(`api_path` ASC, `api_method` ASC) USING BTREE,
  INDEX `idx_api_method`(`api_method` ASC) USING BTREE,
  INDEX `idx_created_at`(`created_at` ASC) USING BTREE,
  INDEX `fk_apis_permission`(`permission_id` ASC) USING BTREE,
  CONSTRAINT `apis_ibfk_1` FOREIGN KEY (`permission_id`) REFERENCES `permissions` (`id`) ON DELETE SET NULL ON UPDATE RESTRICT,
  CONSTRAINT `fk_apis_permission` FOREIGN KEY (`permission_id`) REFERENCES `permissions` (`id`) ON DELETE RESTRICT ON UPDATE RESTRICT
) ENGINE = InnoDB AUTO_INCREMENT = 39 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = 'API?????' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of apis
-- ----------------------------
INSERT INTO `apis` VALUES (1, '用户登录', '/api/admin/auth/login', 'POST', '用户组', '用户登录接口', NULL, 1, NULL, '2025-11-29 21:25:50.000', '2025-11-29 22:30:29.149', NULL);
INSERT INTO `apis` VALUES (2, '用户登出', '/api/admin/auth/logout', 'POST', '用户组', '用户登出接口', NULL, 1, NULL, '2025-11-29 21:25:50.000', '2025-11-29 21:25:50.000', NULL);
INSERT INTO `apis` VALUES (3, '获取用户信息', '/api/admin/auth/info', 'GET', '用户组', '获取当前登录用户信息', NULL, 1, NULL, '2025-11-29 21:25:50.000', '2025-11-29 21:25:50.000', NULL);
INSERT INTO `apis` VALUES (4, '获取用户列表', '/api/admin/system/users', 'GET', '用户组', '分页获取用户列表', NULL, 1, NULL, '2025-11-29 21:25:50.000', '2025-11-29 21:25:50.000', NULL);
INSERT INTO `apis` VALUES (5, '创建用户', '/api/admin/system/users', 'POST', '用户组', '创建新用户', NULL, 1, NULL, '2025-11-29 21:25:50.000', '2025-11-29 21:25:50.000', NULL);
INSERT INTO `apis` VALUES (6, '编辑用户', '/api/admin/system/users', 'PUT', '用户组', '编辑用户信息', NULL, 1, NULL, '2025-11-29 21:25:50.000', '2025-11-29 21:25:50.000', NULL);
INSERT INTO `apis` VALUES (7, '删除用户', '/api/admin/system/users', 'DELETE', '用户组', '删除用户', NULL, 1, NULL, '2025-11-29 21:25:50.000', '2025-11-29 21:25:50.000', NULL);
INSERT INTO `apis` VALUES (8, '获取用户详情', '/api/admin/system/users/:id', 'GET', '用户组', '根据ID获取用户详情', NULL, 1, NULL, '2025-11-29 21:25:50.000', '2025-11-29 21:25:50.000', NULL);
INSERT INTO `apis` VALUES (9, '为用户分配角色', '/api/admin/system/users/:id/roles', 'POST', '角色组', '为用户分配角色', NULL, 1, NULL, '2025-11-29 21:25:50.000', '2025-11-29 21:25:50.000', NULL);
INSERT INTO `apis` VALUES (10, '获取用户角色', '/api/admin/system/users/:id/roles', 'GET', '角色组', '获取用户已分配的角色', NULL, 1, NULL, '2025-11-29 21:25:50.000', '2025-11-29 21:25:50.000', NULL);
INSERT INTO `apis` VALUES (11, '获取角色列表', '/api/admin/system/roles', 'GET', '角色组', '分页获取角色列表', NULL, 1, NULL, '2025-11-29 21:25:50.000', '2025-11-29 21:25:50.000', NULL);
INSERT INTO `apis` VALUES (12, '创建角色', '/api/admin/system/roles', 'POST', '角色组', '创建新角色', NULL, 1, NULL, '2025-11-29 21:25:50.000', '2025-11-29 21:25:50.000', NULL);
INSERT INTO `apis` VALUES (13, '编辑角色', '/api/admin/system/roles', 'PUT', '角色组', '编辑角色信息', NULL, 1, NULL, '2025-11-29 21:25:50.000', '2025-11-29 21:25:50.000', NULL);
INSERT INTO `apis` VALUES (14, '删除角色', '/api/admin/system/roles', 'DELETE', '角色组', '删除角色', NULL, 1, NULL, '2025-11-29 21:25:50.000', '2025-11-29 21:25:50.000', NULL);
INSERT INTO `apis` VALUES (15, '获取角色详情', '/api/admin/system/roles/:id', 'GET', '角色组', '根据ID获取角色详情', NULL, 1, NULL, '2025-11-29 21:25:50.000', '2025-11-29 21:25:50.000', NULL);
INSERT INTO `apis` VALUES (16, '为角色分配权限', '/api/admin/system/roles/:id/permissions', 'POST', '角色组', '为角色分配权限', NULL, 1, NULL, '2025-11-29 21:25:50.000', '2025-11-29 21:25:50.000', NULL);
INSERT INTO `apis` VALUES (17, '获取角色权限', '/api/admin/system/roles/:id/permissions', 'GET', '角色组', '获取角色已分配的权限', NULL, 1, NULL, '2025-11-29 21:25:50.000', '2025-11-29 21:25:50.000', NULL);
INSERT INTO `apis` VALUES (18, '获取菜单列表', '/api/admin/system/menus', 'GET', '菜单组', '分页获取菜单列表', NULL, 1, NULL, '2025-11-29 21:25:50.000', '2025-11-29 21:25:50.000', NULL);
INSERT INTO `apis` VALUES (19, '创建菜单', '/api/admin/system/menus', 'POST', '菜单组', '创建新菜单', NULL, 1, NULL, '2025-11-29 21:25:50.000', '2025-11-29 21:25:50.000', NULL);
INSERT INTO `apis` VALUES (20, '编辑菜单', '/api/admin/system/menus', 'PUT', '菜单组', '编辑菜单信息', NULL, 1, NULL, '2025-11-29 21:25:50.000', '2025-11-29 21:25:50.000', NULL);
INSERT INTO `apis` VALUES (21, '删除菜单', '/api/admin/system/menus', 'DELETE', '菜单组', '删除菜单', NULL, 1, NULL, '2025-11-29 21:25:50.000', '2025-11-29 21:25:50.000', NULL);
INSERT INTO `apis` VALUES (22, '获取菜单详情', '/api/admin/system/menus/:id', 'GET', '菜单组', '根据ID获取菜单详情', NULL, 1, NULL, '2025-11-29 21:25:50.000', '2025-11-29 21:25:50.000', NULL);
INSERT INTO `apis` VALUES (23, '获取权限列表', '/api/admin/system/permissions', 'GET', '权限组', '分页获取权限列表', NULL, 1, NULL, '2025-11-29 21:25:50.000', '2025-11-29 21:25:50.000', NULL);
INSERT INTO `apis` VALUES (24, '创建权限', '/api/admin/system/permissions', 'POST', '权限组', '创建新权限', NULL, 1, NULL, '2025-11-29 21:25:50.000', '2025-11-29 21:25:50.000', NULL);
INSERT INTO `apis` VALUES (25, '编辑权限', '/api/admin/system/permissions', 'PUT', '权限组', '编辑权限信息', NULL, 1, NULL, '2025-11-29 21:25:50.000', '2025-11-29 21:25:50.000', NULL);
INSERT INTO `apis` VALUES (26, '删除权限', '/api/admin/system/permissions', 'DELETE', '权限组', '删除权限', NULL, 1, NULL, '2025-11-29 21:25:50.000', '2025-11-29 21:25:50.000', NULL);
INSERT INTO `apis` VALUES (27, '获取权限详情', '/api/admin/system/permissions/:id', 'GET', '权限组', '根据ID获取权限详情', NULL, 1, NULL, '2025-11-29 21:25:50.000', '2025-11-29 21:25:50.000', NULL);
INSERT INTO `apis` VALUES (28, '获取字典列表', '/api/admin/content/dicts', 'GET', '字典组', '分页获取字典列表', NULL, 1, NULL, '2025-11-29 21:25:50.000', '2025-11-29 21:25:50.000', NULL);
INSERT INTO `apis` VALUES (29, '创建字典', '/api/admin/content/dicts', 'POST', '字典组', '创建新字典项', NULL, 1, NULL, '2025-11-29 21:25:50.000', '2025-11-29 21:25:50.000', NULL);
INSERT INTO `apis` VALUES (30, '编辑字典', '/api/admin/content/dicts', 'PUT', '字典组', '编辑字典项', NULL, 1, NULL, '2025-11-29 21:25:50.000', '2025-11-29 21:25:50.000', NULL);
INSERT INTO `apis` VALUES (31, '删除字典', '/api/admin/content/dicts', 'DELETE', '字典组', '删除字典项', NULL, 1, NULL, '2025-11-29 21:25:50.000', '2025-11-29 21:25:50.000', NULL);
INSERT INTO `apis` VALUES (32, '获取字典详情', '/api/admin/content/dicts/:id', 'GET', '字典组', '根据ID获取字典项详情', NULL, 1, NULL, '2025-11-29 21:25:50.000', '2025-11-29 21:25:50.000', NULL);
INSERT INTO `apis` VALUES (33, '获取API列表', '/api/admin/system/apis', 'GET', 'Api组', '分页获取API列表', NULL, 1, NULL, '2025-11-29 21:25:50.000', '2025-11-29 21:25:50.000', NULL);
INSERT INTO `apis` VALUES (34, '创建API', '/api/admin/system/apis', 'POST', 'Api组', '创建新API接口记录', NULL, 1, NULL, '2025-11-29 21:25:50.000', '2025-11-29 21:25:50.000', NULL);
INSERT INTO `apis` VALUES (35, '编辑API', '/api/admin/system/apis', 'PUT', 'Api组', '编辑API接口信息', NULL, 1, NULL, '2025-11-29 21:25:50.000', '2025-11-29 21:25:50.000', NULL);
INSERT INTO `apis` VALUES (36, '删除API', '/api/admin/system/apis', 'DELETE', 'Api组', '删除API接口记录', NULL, 1, NULL, '2025-11-29 21:25:50.000', '2025-11-29 21:25:50.000', NULL);
INSERT INTO `apis` VALUES (37, '获取API详情', '/api/admin/system/apis/:id', 'GET', 'Api组', '根据ID获取API接口详情', NULL, 1, NULL, '2025-11-29 21:25:50.000', '2025-11-29 21:25:50.000', NULL);
INSERT INTO `apis` VALUES (38, 'test', '', 'GET', 'Api组', '', NULL, 1, '', '2025-11-29 21:37:39.000', '2025-11-29 21:38:24.000', NULL);

-- ----------------------------
-- Table structure for apps
-- ----------------------------
DROP TABLE IF EXISTS `apps`;
CREATE TABLE `apps`  (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '应用ID',
  `app_name` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '应用名称',
  `app_key` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '应用标识',
  `version` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '应用版本',
  `author` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '作者',
  `description` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL COMMENT '应用描述',
  `icon` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '应用图标',
  `download_url` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '下载地址',
  `demo_url` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '演示地址',
  `category` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '应用分类',
  `tags` json NULL COMMENT '应用标签',
  `rating` double NULL DEFAULT NULL COMMENT '应用评分',
  `downloads` bigint(20) NULL DEFAULT NULL COMMENT '下载次数',
  `status` bigint(20) NULL DEFAULT NULL COMMENT '应用状态',
  `published` bigint(20) NULL DEFAULT NULL COMMENT '发布状态',
  `created_at` datetime(3) NULL DEFAULT NULL COMMENT '创建时间',
  `updated_at` datetime(3) NULL DEFAULT NULL COMMENT '更新时间',
  `deleted_at` datetime(3) NULL DEFAULT NULL COMMENT '删除时间',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `idx_apps_app_key`(`app_key`) USING BTREE
) ENGINE = MyISAM AUTO_INCREMENT = 1 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of apps
-- ----------------------------

-- ----------------------------
-- Table structure for casbin_rule
-- ----------------------------
DROP TABLE IF EXISTS `casbin_rule`;
CREATE TABLE `casbin_rule`  (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `ptype` varchar(10) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `v0` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `v1` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `v2` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `v3` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `v4` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `v5` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `idx_ptype`(`ptype`) USING BTREE,
  INDEX `idx_v0`(`v0`) USING BTREE
) ENGINE = MyISAM AUTO_INCREMENT = 1 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of casbin_rule
-- ----------------------------

-- ----------------------------
-- Table structure for dictionaries
-- ----------------------------
DROP TABLE IF EXISTS `dictionaries`;
CREATE TABLE `dictionaries`  (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '?ֵ?ID',
  `dict_type` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '字典类型',
  `dict_label` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '字典标签',
  `dict_value` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '字典值',
  `description` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '字典描述',
  `sort` bigint(20) NULL DEFAULT 0 COMMENT '排序号',
  `status` bigint(20) NULL DEFAULT 1 COMMENT '状态 1:启用 0:禁用',
  `remark` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '备注',
  `created_at` datetime(3) NULL DEFAULT NULL COMMENT '创建时间',
  `updated_at` datetime(3) NULL DEFAULT NULL COMMENT '更新时间',
  `deleted_at` datetime(3) NULL DEFAULT NULL COMMENT '删除时间',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `uk_type_value`(`dict_type` ASC, `dict_value` ASC) USING BTREE,
  INDEX `idx_dict_type`(`dict_type` ASC) USING BTREE,
  INDEX `idx_sort`(`sort` ASC) USING BTREE,
  INDEX `idx_dictionaries_dict_type`(`dict_type` ASC) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 11 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '?ֵ??' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of dictionaries
-- ----------------------------
INSERT INTO `dictionaries` VALUES (1, 'gender', '男', '1', NULL, 1, 1, NULL, '2025-11-29 10:31:32.000', '2025-11-29 10:31:32.000', NULL);
INSERT INTO `dictionaries` VALUES (2, 'gender', '女', '2', NULL, 2, 1, NULL, '2025-11-29 10:31:32.000', '2025-11-29 10:31:32.000', NULL);
INSERT INTO `dictionaries` VALUES (3, 'gender', '未知', '0', NULL, 3, 1, NULL, '2025-11-29 10:31:32.000', '2025-11-29 10:31:32.000', NULL);
INSERT INTO `dictionaries` VALUES (4, 'status', '启用', '1', NULL, 1, 1, NULL, '2025-11-29 10:31:32.000', '2025-11-29 10:31:32.000', NULL);
INSERT INTO `dictionaries` VALUES (5, 'status', '禁用', '0', NULL, 2, 1, NULL, '2025-11-29 10:31:32.000', '2025-11-29 10:31:32.000', NULL);
INSERT INTO `dictionaries` VALUES (6, 'menu_type', '菜单', '1', NULL, 1, 1, NULL, '2025-11-29 10:31:32.000', '2025-11-29 10:31:32.000', NULL);
INSERT INTO `dictionaries` VALUES (7, 'menu_type', '按钮', '2', NULL, 2, 1, NULL, '2025-11-29 10:31:32.000', '2025-11-29 10:31:32.000', NULL);
INSERT INTO `dictionaries` VALUES (8, 'user_status', '正常', '1', NULL, 1, 1, NULL, '2025-11-29 10:31:32.000', '2025-11-29 10:31:32.000', NULL);
INSERT INTO `dictionaries` VALUES (9, 'user_status', '禁用', '0', NULL, 2, 1, NULL, '2025-11-29 10:31:32.000', '2025-11-29 10:31:32.000', NULL);
INSERT INTO `dictionaries` VALUES (10, 'user_status', '锁定', '2', NULL, 3, 1, NULL, '2025-11-29 10:31:32.000', '2025-11-29 10:31:32.000', NULL);

-- ----------------------------
-- Table structure for logs
-- ----------------------------
DROP TABLE IF EXISTS `logs`;
CREATE TABLE `logs`  (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '??־ID',
  `user_id` bigint(20) NULL DEFAULT NULL COMMENT '用户ID',
  `username` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '用户名',
  `operation` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '操作名称',
  `method` varchar(10) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '请求方法',
  `path` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '请求路径',
  `ip` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '请求IP',
  `status` bigint(20) NULL DEFAULT NULL COMMENT '响应状态码',
  `error_msg` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL COMMENT '错误信息',
  `request_body` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL COMMENT '请求体',
  `response_body` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL COMMENT '响应体',
  `duration` bigint(20) NULL DEFAULT NULL COMMENT '耗时(ms)',
  `created_at` datetime(3) NULL DEFAULT NULL COMMENT '创建时间',
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `idx_user_id`(`user_id` ASC) USING BTREE,
  INDEX `idx_method_path`(`method` ASC, `path` ASC) USING BTREE,
  INDEX `idx_created_at`(`created_at` ASC) USING BTREE,
  INDEX `idx_status`(`status` ASC) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 1 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = 'ϵͳ??־?' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of logs
-- ----------------------------

-- ----------------------------
-- Table structure for menus
-- ----------------------------
DROP TABLE IF EXISTS `menus`;
CREATE TABLE `menus`  (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '菜单ID',
  `parent_id` bigint(20) NULL DEFAULT 0 COMMENT '父菜单ID',
  `menu_name` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '菜单名称',
  `menu_path` varchar(200) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '菜单路径',
  `component` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '组件路径',
  `icon` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '菜单图标',
  `sort` bigint(20) NULL DEFAULT 0 COMMENT '排序号',
  `status` bigint(20) NULL DEFAULT 1 COMMENT '状态 1:显示 0:隐藏',
  `menu_type` bigint(20) NULL DEFAULT 1 COMMENT '菜单类型 1:菜单 2:按钮',
  `permission_id` bigint(20) NULL DEFAULT NULL COMMENT '关联权限ID',
  `remark` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '备注',
  `created_at` datetime(3) NULL DEFAULT NULL COMMENT '创建时间',
  `updated_at` datetime(3) NULL DEFAULT NULL COMMENT '更新时间',
  `deleted_at` datetime(3) NULL DEFAULT NULL COMMENT '删除时间',
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `idx_parent_id`(`parent_id` ASC) USING BTREE,
  INDEX `idx_sort`(`sort` ASC) USING BTREE,
  INDEX `idx_created_at`(`created_at` ASC) USING BTREE,
  INDEX `fk_menus_permission`(`permission_id` ASC) USING BTREE,
  CONSTRAINT `fk_menus_permission` FOREIGN KEY (`permission_id`) REFERENCES `permissions` (`id`) ON DELETE RESTRICT ON UPDATE RESTRICT,
  CONSTRAINT `menus_ibfk_1` FOREIGN KEY (`permission_id`) REFERENCES `permissions` (`id`) ON DELETE SET NULL ON UPDATE RESTRICT
) ENGINE = InnoDB AUTO_INCREMENT = 20 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '?˵??' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of menus
-- ----------------------------
INSERT INTO `menus` VALUES (1, 0, '系统管理', '/system', '', 'setting', 10, 1, 1, NULL, NULL, '2025-11-29 10:31:27.000', '2025-11-29 10:31:27.000', NULL);
INSERT INTO `menus` VALUES (2, 1, '用户管理', '/system/users', 'system/user/UserList', 'user', 1, 1, 1, NULL, NULL, '2025-11-29 10:31:27.000', '2025-11-29 10:31:27.000', NULL);
INSERT INTO `menus` VALUES (3, 1, '角色管理', '/system/roles', 'system/role/RoleList', 'admin', 2, 1, 1, NULL, NULL, '2025-11-29 10:31:27.000', '2025-11-29 10:31:27.000', NULL);
INSERT INTO `menus` VALUES (4, 1, '菜单管理', '/system/menus', 'system/menu/MenuList', 'menu', 3, 1, 1, NULL, NULL, '2025-11-29 10:31:27.000', '2025-11-29 10:31:27.000', NULL);
INSERT INTO `menus` VALUES (5, 1, '权限管理', '/system/permissions', 'system/permission/PermissionList', 'lock', 4, 1, 1, NULL, NULL, '2025-11-29 10:31:27.000', '2025-11-29 10:31:27.000', NULL);
INSERT INTO `menus` VALUES (6, 1, 'API管理', '/system/apis', 'system/api/ApiList', 'link', 5, 1, 1, NULL, NULL, '2025-11-29 10:31:27.000', '2025-11-29 10:31:27.000', NULL);
INSERT INTO `menus` VALUES (7, 0, '内容管理', '/content', '', 'document', 20, 1, 1, NULL, NULL, '2025-11-29 10:31:27.000', '2025-11-29 10:31:27.000', NULL);
INSERT INTO `menus` VALUES (8, 7, '字典管理', '/content/dicts', 'content/dict/DictList', 'list', 1, 1, 1, NULL, NULL, '2025-11-29 10:31:27.000', '2025-11-29 10:31:27.000', NULL);
INSERT INTO `menus` VALUES (9, 0, '应用中心', '/market', '', 'shopping', 30, 1, 1, NULL, NULL, '2025-11-29 10:31:27.000', '2025-11-29 10:31:27.000', NULL);
INSERT INTO `menus` VALUES (10, 9, '应用市场', '/market/apps', 'market/AppMarket', 'shop', 1, 1, 1, NULL, NULL, '2025-11-29 10:31:27.000', '2025-11-29 10:31:27.000', NULL);
INSERT INTO `menus` VALUES (11, 0, '系统设置', '/system-config', '', 'setting', 40, 1, 1, NULL, NULL, '2025-11-29 10:31:27.000', '2025-11-29 10:31:27.000', NULL);
INSERT INTO `menus` VALUES (12, 11, '日志管理', '/logs', 'logs/LogList', 'monitor', 1, 1, 1, NULL, NULL, '2025-11-29 10:31:27.000', '2025-11-29 10:31:27.000', NULL);
INSERT INTO `menus` VALUES (13, 1, 'test', 'test', '', 'setting', 0, 1, 1, NULL, '', '2025-11-29 15:14:19.000', '2025-11-29 15:14:19.000', NULL);

-- ----------------------------
-- Table structure for permissions
-- ----------------------------
DROP TABLE IF EXISTS `permissions`;
CREATE TABLE `permissions`  (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT 'Ȩ??ID',
  `name` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '权限名称',
  `description` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '权限描述',
  `resource` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '资源 如 users:view',
  `action` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '操作 如 view, create, update, delete',
  `status` bigint(20) NULL DEFAULT 1 COMMENT '状态 1:启用 0:禁用',
  `created_at` datetime(3) NULL DEFAULT NULL COMMENT '创建时间',
  `updated_at` datetime(3) NULL DEFAULT NULL COMMENT '更新时间',
  `deleted_at` datetime(3) NULL DEFAULT NULL COMMENT '删除时间',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `idx_permissions_name`(`name` ASC) USING BTREE,
  INDEX `idx_resource_action`(`resource` ASC, `action` ASC) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 1 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = 'Ȩ?ޱ' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of permissions
-- ----------------------------

-- ----------------------------
-- Table structure for plugins
-- ----------------------------
DROP TABLE IF EXISTS `plugins`;
CREATE TABLE `plugins`  (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '???',
  `plugin_name` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '插件名称',
  `plugin_key` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '插件标识',
  `description` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '插件描述',
  `version` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '插件版本',
  `author` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '作者',
  `status` bigint(20) NULL DEFAULT 1 COMMENT '状态 1:启用 0:禁用',
  `config` json NULL COMMENT '插件配置JSON',
  `created_at` datetime(3) NULL DEFAULT NULL COMMENT '创建时间',
  `updated_at` datetime(3) NULL DEFAULT NULL COMMENT '更新时间',
  `deleted_at` datetime(3) NULL DEFAULT NULL COMMENT '删除时间',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `idx_plugins_plugin_name`(`plugin_name` ASC) USING BTREE,
  UNIQUE INDEX `idx_plugins_plugin_key`(`plugin_key` ASC) USING BTREE,
  INDEX `idx_plugin_key`(`plugin_key` ASC) USING BTREE,
  INDEX `idx_status`(`status` ASC) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 1 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '???' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of plugins
-- ----------------------------

-- ----------------------------
-- Table structure for reviews
-- ----------------------------
DROP TABLE IF EXISTS `reviews`;
CREATE TABLE `reviews`  (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '评价ID',
  `app_id` bigint(20) NULL DEFAULT NULL COMMENT '应用ID',
  `user_id` bigint(20) NULL DEFAULT NULL COMMENT '用户ID',
  `rating` bigint(20) NULL DEFAULT NULL COMMENT '评分',
  `comment` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL COMMENT '评价内容',
  `created_at` datetime(3) NULL DEFAULT NULL COMMENT '创建时间',
  `updated_at` datetime(3) NULL DEFAULT NULL COMMENT '更新时间',
  `deleted_at` datetime(3) NULL DEFAULT NULL COMMENT '删除时间',
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `idx_reviews_app_id`(`app_id`) USING BTREE
) ENGINE = MyISAM AUTO_INCREMENT = 1 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of reviews
-- ----------------------------

-- ----------------------------
-- Table structure for role_menus
-- ----------------------------
DROP TABLE IF EXISTS `role_menus`;
CREATE TABLE `role_menus`  (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '???',
  `role_id` bigint(20) NULL DEFAULT NULL COMMENT '角色ID',
  `menu_id` bigint(20) NULL DEFAULT NULL COMMENT '菜单ID',
  `created_at` datetime(3) NULL DEFAULT NULL COMMENT '创建时间',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `uk_role_menu`(`role_id` ASC, `menu_id` ASC) USING BTREE,
  INDEX `idx_role_id`(`role_id` ASC) USING BTREE,
  INDEX `idx_menu_id`(`menu_id` ASC) USING BTREE,
  INDEX `idx_role_menus_role_id`(`role_id` ASC) USING BTREE,
  INDEX `idx_role_menus_menu_id`(`menu_id` ASC) USING BTREE,
  CONSTRAINT `fk_role_menus_menu` FOREIGN KEY (`menu_id`) REFERENCES `menus` (`id`) ON DELETE RESTRICT ON UPDATE RESTRICT,
  CONSTRAINT `fk_role_menus_role` FOREIGN KEY (`role_id`) REFERENCES `roles` (`id`) ON DELETE RESTRICT ON UPDATE RESTRICT,
  CONSTRAINT `role_menus_ibfk_1` FOREIGN KEY (`role_id`) REFERENCES `roles` (`id`) ON DELETE CASCADE ON UPDATE RESTRICT,
  CONSTRAINT `role_menus_ibfk_2` FOREIGN KEY (`menu_id`) REFERENCES `menus` (`id`) ON DELETE CASCADE ON UPDATE RESTRICT
) ENGINE = InnoDB AUTO_INCREMENT = 66 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '??ɫ-?˵??????' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of role_menus
-- ----------------------------
INSERT INTO `role_menus` VALUES (50, 1, 13, '2025-11-30 11:19:25.256');
INSERT INTO `role_menus` VALUES (51, 1, 2, '2025-11-30 11:19:25.271');
INSERT INTO `role_menus` VALUES (52, 1, 3, '2025-11-30 11:19:25.284');
INSERT INTO `role_menus` VALUES (53, 1, 4, '2025-11-30 11:19:25.300');
INSERT INTO `role_menus` VALUES (54, 1, 5, '2025-11-30 11:19:25.316');
INSERT INTO `role_menus` VALUES (55, 1, 6, '2025-11-30 11:19:25.331');
INSERT INTO `role_menus` VALUES (56, 1, 1, '2025-11-30 11:19:25.347');
INSERT INTO `role_menus` VALUES (57, 1, 7, '2025-11-30 11:19:25.362');
INSERT INTO `role_menus` VALUES (58, 1, 8, '2025-11-30 11:19:25.378');
INSERT INTO `role_menus` VALUES (59, 2, 1, '2025-11-30 11:39:44.085');
INSERT INTO `role_menus` VALUES (60, 2, 13, '2025-11-30 11:39:44.143');
INSERT INTO `role_menus` VALUES (61, 2, 2, '2025-11-30 11:39:44.156');
INSERT INTO `role_menus` VALUES (62, 2, 3, '2025-11-30 11:39:44.173');
INSERT INTO `role_menus` VALUES (63, 2, 4, '2025-11-30 11:39:44.188');
INSERT INTO `role_menus` VALUES (64, 2, 5, '2025-11-30 11:39:44.203');
INSERT INTO `role_menus` VALUES (65, 2, 6, '2025-11-30 11:39:44.219');

-- ----------------------------
-- Table structure for role_permissions
-- ----------------------------
DROP TABLE IF EXISTS `role_permissions`;
CREATE TABLE `role_permissions`  (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '???',
  `role_id` bigint(20) NOT NULL COMMENT '??ɫID',
  `permission_id` bigint(20) NOT NULL COMMENT 'Ȩ??ID',
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '????ʱ?',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `uk_role_permission`(`role_id` ASC, `permission_id` ASC) USING BTREE,
  INDEX `idx_role_id`(`role_id` ASC) USING BTREE,
  INDEX `idx_permission_id`(`permission_id` ASC) USING BTREE,
  CONSTRAINT `fk_role_permissions_permission` FOREIGN KEY (`permission_id`) REFERENCES `permissions` (`id`) ON DELETE RESTRICT ON UPDATE RESTRICT,
  CONSTRAINT `fk_role_permissions_role` FOREIGN KEY (`role_id`) REFERENCES `roles` (`id`) ON DELETE RESTRICT ON UPDATE RESTRICT,
  CONSTRAINT `role_permissions_ibfk_1` FOREIGN KEY (`role_id`) REFERENCES `roles` (`id`) ON DELETE CASCADE ON UPDATE RESTRICT,
  CONSTRAINT `role_permissions_ibfk_2` FOREIGN KEY (`permission_id`) REFERENCES `permissions` (`id`) ON DELETE CASCADE ON UPDATE RESTRICT
) ENGINE = InnoDB AUTO_INCREMENT = 1 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '??ɫ-Ȩ?޹????' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of role_permissions
-- ----------------------------

-- ----------------------------
-- Table structure for roles
-- ----------------------------
DROP TABLE IF EXISTS `roles`;
CREATE TABLE `roles`  (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '??ɫID',
  `name` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '角色名称',
  `description` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '角色描述',
  `status` bigint(20) NULL DEFAULT 1 COMMENT '状态 1:启用 0:禁用',
  `remark` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '备注',
  `created_at` datetime(3) NULL DEFAULT NULL COMMENT '创建时间',
  `updated_at` datetime(3) NULL DEFAULT NULL COMMENT '更新时间',
  `deleted_at` datetime(3) NULL DEFAULT NULL COMMENT '删除时间',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `idx_roles_name`(`name` ASC) USING BTREE,
  INDEX `idx_created_at`(`created_at` ASC) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 5 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '??ɫ?' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of roles
-- ----------------------------
INSERT INTO `roles` VALUES (1, 'admin', '管理员', 1, NULL, '2025-11-29 10:31:18.000', '2025-11-29 10:31:18.000', NULL);
INSERT INTO `roles` VALUES (2, 'editor', '编辑', 1, NULL, '2025-11-29 10:31:18.000', '2025-11-29 10:31:18.000', NULL);
INSERT INTO `roles` VALUES (4, 'user', '普通用户', 1, '', '2025-11-29 21:59:07.597', '2025-11-29 21:59:07.597', NULL);

-- ----------------------------
-- Table structure for user_roles
-- ----------------------------
DROP TABLE IF EXISTS `user_roles`;
CREATE TABLE `user_roles`  (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '???',
  `user_id` bigint(20) NOT NULL COMMENT '?û?ID',
  `role_id` bigint(20) NOT NULL COMMENT '??ɫID',
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '????ʱ?',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `uk_user_role`(`user_id` ASC, `role_id` ASC) USING BTREE,
  INDEX `idx_user_id`(`user_id` ASC) USING BTREE,
  INDEX `idx_role_id`(`role_id` ASC) USING BTREE,
  CONSTRAINT `fk_user_roles_role` FOREIGN KEY (`role_id`) REFERENCES `roles` (`id`) ON DELETE RESTRICT ON UPDATE RESTRICT,
  CONSTRAINT `fk_user_roles_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE RESTRICT ON UPDATE RESTRICT,
  CONSTRAINT `user_roles_ibfk_1` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE ON UPDATE RESTRICT,
  CONSTRAINT `user_roles_ibfk_2` FOREIGN KEY (`role_id`) REFERENCES `roles` (`id`) ON DELETE CASCADE ON UPDATE RESTRICT
) ENGINE = InnoDB AUTO_INCREMENT = 6 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '?û?-??ɫ?????' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of user_roles
-- ----------------------------
INSERT INTO `user_roles` VALUES (1, 1, 1, '2025-11-29 10:31:21');
INSERT INTO `user_roles` VALUES (2, 2, 2, '2025-11-29 10:31:21');

-- ----------------------------
-- Table structure for users
-- ----------------------------
DROP TABLE IF EXISTS `users`;
CREATE TABLE `users`  (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '?û?ID',
  `username` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '用户名',
  `phone` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '手机号',
  `email` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '邮箱',
  `password` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '密码（加密存储）',
  `nickname` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '昵称',
  `avatar` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '头像URL',
  `gender` bigint(20) NULL DEFAULT NULL COMMENT '性别 1:男 2:女 0:未知',
  `status` bigint(20) NULL DEFAULT 1 COMMENT '状态 1:激活 0:禁用',
  `last_login_at` datetime(3) NULL DEFAULT NULL COMMENT '最后登录时间',
  `created_at` datetime(3) NULL DEFAULT NULL COMMENT '创建时间',
  `updated_at` datetime(3) NULL DEFAULT NULL COMMENT '更新时间',
  `deleted_at` datetime(3) NULL DEFAULT NULL COMMENT '删除时间',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `idx_users_username`(`username` ASC) USING BTREE,
  UNIQUE INDEX `idx_users_phone`(`phone` ASC) USING BTREE,
  UNIQUE INDEX `idx_users_email`(`email` ASC) USING BTREE,
  INDEX `idx_phone`(`phone` ASC) USING BTREE,
  INDEX `idx_email`(`email` ASC) USING BTREE,
  INDEX `idx_created_at`(`created_at` ASC) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 4 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '?û??' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of users
-- ----------------------------
INSERT INTO `users` VALUES (1, 'admin', '13800000000', 'admin@example.com', '$2a$10$oXf.6jZOW.gq/hm5hBR/OOk7uAmou3z6L89GkfriPRQGL5IIciEP6', '管理员', NULL, NULL, 1, NULL, '2025-11-29 10:31:14.000', '2025-11-30 11:17:11.402', NULL);
INSERT INTO `users` VALUES (2, 'editor', '13800000001', 'editor@example.com', '$2a$10$oXf.6jZOW.gq/hm5hBR/OOk7uAmou3z6L89GkfriPRQGL5IIciEP6', '编辑', NULL, NULL, 1, NULL, '2025-11-29 10:31:14.000', '2025-11-29 21:59:07.579', NULL);
INSERT INTO `users` VALUES (3, 'user', '13800000002', 'user@example.com', '$2a$10$oXf.6jZOW.gq/hm5hBR/OOk7uAmou3z6L89GkfriPRQGL5IIciEP6', '普通用户', NULL, NULL, 1, NULL, '2025-11-29 10:31:14.000', '2025-12-01 00:16:28.458', NULL);

SET FOREIGN_KEY_CHECKS = 1;
