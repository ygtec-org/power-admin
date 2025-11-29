-- 创建数据库
CREATE DATABASE IF NOT EXISTS power_admin CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

USE power_admin;

-- 用户表
CREATE TABLE IF NOT EXISTS users (
    id BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '用户ID',
    username VARCHAR(50) UNIQUE NOT NULL COMMENT '用户名',
    phone VARCHAR(20) UNIQUE COMMENT '手机号',
    email VARCHAR(100) UNIQUE COMMENT '邮箱',
    password VARCHAR(255) NOT NULL COMMENT '密码（加密存储）',
    nickname VARCHAR(100) COMMENT '昵称',
    avatar VARCHAR(255) COMMENT '头像URL',
    gender INT COMMENT '性别 1:男 2:女 0:未知',
    status INT DEFAULT 1 COMMENT '状态 1:激活 0:禁用',
    last_login_at TIMESTAMP COMMENT '最后登录时间',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    deleted_at TIMESTAMP NULL COMMENT '删除时间',
    INDEX idx_phone (phone),
    INDEX idx_email (email),
    INDEX idx_created_at (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户表';

-- 角色表
CREATE TABLE IF NOT EXISTS roles (
    id BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '角色ID',
    name VARCHAR(50) UNIQUE NOT NULL COMMENT '角色名称',
    description VARCHAR(255) COMMENT '角色描述',
    status INT DEFAULT 1 COMMENT '状态 1:启用 0:禁用',
    remark VARCHAR(255) COMMENT '备注',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    deleted_at TIMESTAMP NULL COMMENT '删除时间',
    UNIQUE KEY uk_name (name),
    INDEX idx_created_at (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='角色表';

-- 用户-角色关联表
CREATE TABLE IF NOT EXISTS user_roles (
    id BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '主键ID',
    user_id BIGINT NOT NULL COMMENT '用户ID',
    role_id BIGINT NOT NULL COMMENT '角色ID',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    UNIQUE KEY uk_user_role (user_id, role_id),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (role_id) REFERENCES roles(id) ON DELETE CASCADE,
    INDEX idx_user_id (user_id),
    INDEX idx_role_id (role_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户-角色关联表';

-- 权限表
CREATE TABLE IF NOT EXISTS permissions (
    id BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '权限ID',
    name VARCHAR(100) UNIQUE NOT NULL COMMENT '权限名称',
    description VARCHAR(255) COMMENT '权限描述',
    resource VARCHAR(100) COMMENT '资源',
    action VARCHAR(50) COMMENT '操作',
    status INT DEFAULT 1 COMMENT '状态 1:启用 0:禁用',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    deleted_at TIMESTAMP NULL COMMENT '删除时间',
    INDEX idx_resource_action (resource, action)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='权限表';

-- 角色-权限关联表
CREATE TABLE IF NOT EXISTS role_permissions (
    id BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '主键ID',
    role_id BIGINT NOT NULL COMMENT '角色ID',
    permission_id BIGINT NOT NULL COMMENT '权限ID',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    UNIQUE KEY uk_role_permission (role_id, permission_id),
    FOREIGN KEY (role_id) REFERENCES roles(id) ON DELETE CASCADE,
    FOREIGN KEY (permission_id) REFERENCES permissions(id) ON DELETE CASCADE,
    INDEX idx_role_id (role_id),
    INDEX idx_permission_id (permission_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='角色-权限关联表';

-- 菜单表
CREATE TABLE IF NOT EXISTS menus (
    id BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '菜单ID',
    parent_id BIGINT DEFAULT 0 COMMENT '父菜单ID',
    menu_name VARCHAR(100) NOT NULL COMMENT '菜单名称',
    menu_path VARCHAR(200) COMMENT '菜单路径',
    component VARCHAR(255) COMMENT '组件路径',
    icon VARCHAR(100) COMMENT '菜单图标',
    sort INT DEFAULT 0 COMMENT '排序号',
    status INT DEFAULT 1 COMMENT '状态 1:显示 0:隐藏',
    menu_type INT DEFAULT 1 COMMENT '菜单类型 1:菜单 2:按钮',
    permission_id BIGINT COMMENT '关联权限ID',
    remark VARCHAR(255) COMMENT '备注',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    deleted_at TIMESTAMP NULL COMMENT '删除时间',
    INDEX idx_parent_id (parent_id),
    INDEX idx_sort (sort),
    INDEX idx_created_at (created_at),
    FOREIGN KEY (permission_id) REFERENCES permissions(id) ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='菜单表';

-- 角色-菜单关联表
CREATE TABLE IF NOT EXISTS role_menus (
    id BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '主键ID',
    role_id BIGINT NOT NULL COMMENT '角色ID',
    menu_id BIGINT NOT NULL COMMENT '菜单ID',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    UNIQUE KEY uk_role_menu (role_id, menu_id),
    FOREIGN KEY (role_id) REFERENCES roles(id) ON DELETE CASCADE,
    FOREIGN KEY (menu_id) REFERENCES menus(id) ON DELETE CASCADE,
    INDEX idx_role_id (role_id),
    INDEX idx_menu_id (menu_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='角色-菜单关联表';

-- 字典表
CREATE TABLE IF NOT EXISTS dictionaries (
    id BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '字典ID',
    dict_type VARCHAR(100) NOT NULL COMMENT '字典类型',
    dict_label VARCHAR(100) NOT NULL COMMENT '字典标签',
    dict_value VARCHAR(255) NOT NULL COMMENT '字典值',
    description VARCHAR(255) COMMENT '字典描述',
    sort INT DEFAULT 0 COMMENT '排序号',
    status INT DEFAULT 1 COMMENT '状态 1:启用 0:禁用',
    remark VARCHAR(255) COMMENT '备注',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    deleted_at TIMESTAMP NULL COMMENT '删除时间',
    UNIQUE KEY uk_type_value (dict_type, dict_value),
    INDEX idx_dict_type (dict_type),
    INDEX idx_sort (sort)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='字典表';

-- API管理表
CREATE TABLE IF NOT EXISTS apis (
    id BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT 'API ID',
    api_name VARCHAR(100) NOT NULL COMMENT 'API名称',
    api_path VARCHAR(255) NOT NULL COMMENT 'API路径',
    api_method VARCHAR(10) NOT NULL COMMENT 'HTTP方法',
    description VARCHAR(255) COMMENT 'API描述',
    permission_id BIGINT COMMENT '关联权限ID',
    status INT DEFAULT 1 COMMENT '状态 1:启用 0:禁用',
    remark VARCHAR(255) COMMENT '备注',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    deleted_at TIMESTAMP NULL COMMENT '删除时间',
    UNIQUE KEY uk_path_method (api_path, api_method),
    INDEX idx_api_method (api_method),
    INDEX idx_created_at (created_at),
    FOREIGN KEY (permission_id) REFERENCES permissions(id) ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='API管理表';

-- 插件表
CREATE TABLE IF NOT EXISTS plugins (
    id BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '插件ID',
    plugin_name VARCHAR(100) UNIQUE NOT NULL COMMENT '插件名称',
    plugin_key VARCHAR(100) UNIQUE NOT NULL COMMENT '插件标识',
    description VARCHAR(255) COMMENT '插件描述',
    version VARCHAR(50) COMMENT '插件版本',
    author VARCHAR(100) COMMENT '作者',
    status INT DEFAULT 1 COMMENT '状态 1:启用 0:禁用',
    config JSON COMMENT '插件配置JSON',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    deleted_at TIMESTAMP NULL COMMENT '删除时间',
    INDEX idx_plugin_key (plugin_key),
    INDEX idx_status (status)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='插件表';

-- 系统日志表
CREATE TABLE IF NOT EXISTS logs (
    id BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '日志ID',
    user_id BIGINT COMMENT '用户ID',
    username VARCHAR(50) COMMENT '用户名',
    operation VARCHAR(100) COMMENT '操作名称',
    method VARCHAR(10) COMMENT '请求方法',
    path VARCHAR(255) COMMENT '请求路径',
    ip VARCHAR(50) COMMENT '请求IP',
    status INT COMMENT '响应状态码',
    error_msg TEXT COMMENT '错误信息',
    request_body LONGTEXT COMMENT '请求体',
    response_body LONGTEXT COMMENT '响应体',
    duration BIGINT COMMENT '耗时(ms)',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    INDEX idx_user_id (user_id),
    INDEX idx_method_path (method, path),
    INDEX idx_created_at (created_at),
    INDEX idx_status (status)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='系统日志表';

-- Casbin RBAC 策略表
CREATE TABLE IF NOT EXISTS casbin_rule (
    id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY COMMENT 'ID',
    ptype VARCHAR(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL,
    v0 VARCHAR(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL,
    v1 VARCHAR(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL,
    v2 VARCHAR(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL,
    v3 VARCHAR(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL,
    v4 VARCHAR(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL,
    v5 VARCHAR(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL,
    UNIQUE KEY idx_casbin_rule (ptype, v0, v1, v2) USING BTREE,
    INDEX idx_v0 (v0),
    INDEX idx_v1 (v1)
) ENGINE=InnoDB CHARACTER SET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='Casbin策略规则表';

-- ========== 初始化数据 ==========

-- 插入用户
INSERT INTO users (username, password, nickname, email, phone, status) VALUES
('admin', '$2a$10$9L89bPEx.1S4DBsv0blEgu9rK3MmSWmqtd/LbOWBxPi3iuXG3UwxW', '管理员', 'admin@example.com', '13800000000', 1),
('editor', '$2a$10$9L89bPEx.1S4DBsv0blEgu9rK3MmSWmqtd/LbOWBxPi3iuXG3UwxW', '编辑', 'editor@example.com', '13800000001', 1),
('user', '$2a$10$9L89bPEx.1S4DBsv0blEgu9rK3MmSWmqtd/LbOWBxPi3iuXG3UwxW', '普通用户', 'user@example.com', '13800000002', 1);

-- 插入角色
INSERT INTO roles (name, description, status) VALUES
('admin', '管理员', 1),
('editor', '编辑', 1),
('user', '普通用户', 1);

-- 关联用户到角色
INSERT INTO user_roles (user_id, role_id) VALUES
(1, 1),
(2, 2),
(3, 3);

-- 插入权限
INSERT INTO permissions (name, description, resource, action, status) VALUES
('user_view', '查看用户', 'users', 'view', 1),
('user_create', '创建用户', 'users', 'create', 1),
('user_update', '编辑用户', 'users', 'update', 1),
('user_delete', '删除用户', 'users', 'delete', 1),
('role_view', '查看角色', 'roles', 'view', 1),
('role_create', '创建角色', 'roles', 'create', 1),
('role_update', '编辑角色', 'roles', 'update', 1),
('role_delete', '删除角色', 'roles', 'delete', 1),
('menu_view', '查看菜单', 'menus', 'view', 1),
('menu_create', '创建菜单', 'menus', 'create', 1),
('menu_update', '编辑菜单', 'menus', 'update', 1),
('menu_delete', '删除菜单', 'menus', 'delete', 1),
('permission_view', '查看权限', 'permissions', 'view', 1),
('permission_create', '创建权限', 'permissions', 'create', 1),
('permission_update', '编辑权限', 'permissions', 'update', 1),
('permission_delete', '删除权限', 'permissions', 'delete', 1),
('api_view', '查看API', 'apis', 'view', 1),
('api_create', '创建API', 'apis', 'create', 1),
('api_update', '编辑API', 'apis', 'update', 1),
('api_delete', '删除API', 'apis', 'delete', 1),
('dict_view', '查看字典', 'dicts', 'view', 1),
('dict_create', '创建字典', 'dicts', 'create', 1),
('dict_update', '编辑字典', 'dicts', 'update', 1),
('dict_delete', '删除字典', 'dicts', 'delete', 1),
('app_view', '查看应用', 'apps', 'view', 1),
('app_install', '安装应用', 'apps', 'install', 1);

-- 为管理员角色分配所有权限
INSERT INTO role_permissions (role_id, permission_id)
SELECT 1, id FROM permissions WHERE status = 1;

-- 为编辑角色分配内容管理权限
INSERT INTO role_permissions (role_id, permission_id)
SELECT 2, id FROM permissions WHERE resource IN ('menus', 'dicts') AND action IN ('view', 'create', 'update');

-- 为用户角色分配基本权限
INSERT INTO role_permissions (role_id, permission_id)
SELECT 3, id FROM permissions WHERE resource IN ('users') AND action = 'view';

-- 插入菜单数据（支持无限级层级）
INSERT INTO menus (parent_id, menu_name, menu_path, component, icon, sort, status, menu_type) VALUES
(0, '系统管理', '/system', '', 'setting', 10, 1, 1),
(1, '用户管理', '/system/users', 'system/user/UserList', 'user', 1, 1, 1),
(1, '角色管理', '/system/roles', 'system/role/RoleList', 'admin', 2, 1, 1),
(1, '菜单管理', '/system/menus', 'system/menu/MenuList', 'menu', 3, 1, 1),
(1, '权限管理', '/system/permissions', 'system/permission/PermissionList', 'lock', 4, 1, 1),
(1, 'API管理', '/system/apis', 'system/api/ApiList', 'link', 5, 1, 1),
(0, '内容管理', '/content', '', 'document', 20, 1, 1),
(7, '字典管理', '/content/dicts', 'content/dict/DictList', 'list', 1, 1, 1),
(0, '应用中心', '/market', '', 'shopping', 30, 1, 1),
(9, '应用市场', '/market/apps', 'market/AppMarket', 'shop', 1, 1, 1),
(0, '系统设置', '/system-config', '', 'setting', 40, 1, 1),
(11, '日志管理', '/logs', 'logs/LogList', 'monitor', 1, 1, 1);

-- 为管理员角色分配所有菜单
INSERT INTO role_menus (role_id, menu_id)
SELECT 1, id FROM menus WHERE status = 1;

-- 为编辑角色分配内容管理菜单
INSERT INTO role_menus (role_id, menu_id)
SELECT 2, id FROM menus WHERE parent_id = 7 AND status = 1;

-- 为用户角色分配部分菜单
INSERT INTO role_menus (role_id, menu_id)
SELECT 3, id FROM menus WHERE menu_name IN ('系统管理', '用户管理') AND status = 1;

-- 插入字典数据
INSERT INTO dictionaries (dict_type, dict_label, dict_value, sort, status) VALUES
('gender', '男', '1', 1, 1),
('gender', '女', '2', 2, 1),
('gender', '未知', '0', 3, 1),
('status', '启用', '1', 1, 1),
('status', '禁用', '0', 2, 1),
('menu_type', '菜单', '1', 1, 1),
('menu_type', '按钮', '2', 2, 1),
('user_status', '正常', '1', 1, 1),
('user_status', '禁用', '0', 2, 1),
('user_status', '锁定', '2', 3, 1);

-- 插入Casbin权限规则（简化版，只保留主要API端点）
INSERT INTO casbin_rule (ptype, v0, v1, v2) VALUES
('p', '1', '/api/v1/admin/login', 'POST'),
('p', '1', '/api/v1/admin/register', 'POST'),
('p', '1', '/api/v1/user/list', 'GET'),
('p', '1', '/api/v1/user/list', 'POST'),
('p', '1', '/api/v1/user/:id', 'GET'),
('p', '1', '/api/v1/user/create', 'POST'),
('p', '1', '/api/v1/user/:id', 'PUT'),
('p', '1', '/api/v1/user/:id', 'DELETE'),
('p', '1', '/api/v1/role/list', 'GET'),
('p', '1', '/api/v1/role/list', 'POST'),
('p', '1', '/api/v1/role/:id', 'GET'),
('p', '1', '/api/v1/role/create', 'POST'),
('p', '1', '/api/v1/role/:id', 'PUT'),
('p', '1', '/api/v1/role/:id', 'DELETE'),
('p', '1', '/api/v1/menu/list', 'GET'),
('p', '1', '/api/v1/menu/list', 'POST'),
('p', '1', '/api/v1/menu/tree', 'GET'),
('p', '1', '/api/v1/menu/:id', 'GET'),
('p', '1', '/api/v1/menu/create', 'POST'),
('p', '1', '/api/v1/menu/:id', 'PUT'),
('p', '1', '/api/v1/menu/:id', 'DELETE'),
('p', '1', '/api/v1/permission/list', 'GET'),
('p', '1', '/api/v1/permission/create', 'POST'),
('p', '1', '/api/v1/permission/:id', 'PUT'),
('p', '1', '/api/v1/permission/:id', 'DELETE'),
('p', '1', '/api/v1/dict/list', 'GET'),
('p', '1', '/api/v1/dict/create', 'POST'),
('p', '1', '/api/v1/dict/:id', 'PUT'),
('p', '1', '/api/v1/dict/:id', 'DELETE'),
('p', '2', '/api/v1/admin/login', 'POST'),
('p', '2', '/api/v1/menu/list', 'GET'),
('p', '2', '/api/v1/menu/tree', 'GET'),
('p', '2', '/api/v1/dict/list', 'GET'),
('p', '2', '/api/v1/dict/create', 'POST'),
('p', '2', '/api/v1/dict/:id', 'PUT'),
('p', '3', '/api/v1/admin/login', 'POST'),
('p', '3', '/api/v1/menu/list', 'GET'),
('p', '3', '/api/v1/menu/tree', 'GET');
