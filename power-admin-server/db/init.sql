-- 数据库初始化脚本
CREATE DATABASE IF NOT EXISTS power_admin CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
USE power_admin;

SET NAMES utf8mb4;
SET CHARACTER_SET_CLIENT = utf8mb4;
SET CHARACTER_SET_CONNECTION = utf8mb4;
SET CHARACTER_SET_RESULTS = utf8mb4;

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
    resource VARCHAR(100) COMMENT '资源 如 users:view',
    action VARCHAR(50) COMMENT '操作 如 view, create, update, delete',
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
    api_method VARCHAR(10) NOT NULL COMMENT 'HTTP方法 GET, POST, PUT, DELETE等',
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

-- 初始化数据：默认管理员账号
-- 用户名: admin 密码: 123456（已使用bcrypt加密）
INSERT INTO users (username, password, nickname, email, phone, status) VALUES
('admin', '$2a$10$9L89bPEx.1S4DBsv0blEgu9rK3MmSWmqtd/LbOWBxPi3iuXG3UwxW', '管理员', 'admin@example.com', '13800000000', 1),
('editor', '$2a$10$9L89bPEx.1S4DBsv0blEgu9rK3MmSWmqtd/LbOWBxPi3iuXG3UwxW', '编辑', 'editor@example.com', '13800000001', 1),
('user', '$2a$10$9L89bPEx.1S4DBsv0blEgu9rK3MmSWmqtd/LbOWBxPi3iuXG3UwxW', '普通用户', 'user@example.com', '13800000002', 1);

-- 初始化数据：角色
INSERT INTO roles (name, description, status) VALUES
('admin', '管理员', 1),
('editor', '编辑', 1),
('user', '普通用户', 1);

-- Casbin RBAC 策略表
CREATE TABLE IF NOT EXISTS casbin_rule (
    id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY COMMENT 'ID',
    ptype VARCHAR(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '策略类型',
    v0 VARCHAR(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '主体或角色',
    v1 VARCHAR(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '资源',
    v2 VARCHAR(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '操作',
    v3 VARCHAR(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '扩展字段3',
    v4 VARCHAR(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '扩展字段4',
    v5 VARCHAR(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '扩展字段5',
    UNIQUE KEY idx_casbin_rule (ptype, v0, v1, v2) USING BTREE,
    INDEX idx_v0 (v0),
    INDEX idx_v1 (v1)
) ENGINE=InnoDB CHARACTER SET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='Casbin策略规则表';

-- Casbin权限规则：参考gin-vue-admin的权限规则结构
-- p 表示权限规则: (角色ID, 资源路径, HTTP方法)
-- v0: 角色ID, v1: 资源路径, v2: HTTP方法
INSERT INTO casbin_rule (ptype, v0, v1, v2) VALUES
-- ==================== 管理员权限(角色ID: 1) ====================
-- 登录注册
('p', '1', '/api/v1/admin/login', 'POST'),
('p', '1', '/api/v1/admin/register', 'POST'),
('p', '1', '/base/login', 'POST'),
('p', '1', '/user/admin_register', 'POST'),

-- 用户管理
('p', '1', '/api/v1/user/list', 'GET'),
('p', '1', '/api/v1/user/list', 'POST'),
('p', '1', '/api/v1/user/:id', 'GET'),
('p', '1', '/api/v1/user/create', 'POST'),
('p', '1', '/api/v1/user/:id', 'PUT'),
('p', '1', '/api/v1/user/:id', 'DELETE'),
('p', '1', '/user/getUserInfo', 'GET'),
('p', '1', '/user/getUserList', 'POST'),
('p', '1', '/user/setUserInfo', 'PUT'),
('p', '1', '/user/setUserAuthority', 'POST'),
('p', '1', '/user/changePassword', 'POST'),
('p', '1', '/user/resetPassword', 'POST'),
('p', '1', '/user/deleteUser', 'DELETE'),

-- 角色管理
('p', '1', '/api/v1/role/list', 'GET'),
('p', '1', '/api/v1/role/list', 'POST'),
('p', '1', '/api/v1/role/:id', 'GET'),
('p', '1', '/api/v1/role/create', 'POST'),
('p', '1', '/api/v1/role/:id', 'PUT'),
('p', '1', '/api/v1/role/:id', 'DELETE'),
('p', '1', '/api/v1/role/:id/permissions', 'GET'),
('p', '1', '/api/v1/role/:id/permissions', 'POST'),
('p', '1', '/authority/getAuthorityList', 'POST'),
('p', '1', '/authority/createAuthority', 'POST'),
('p', '1', '/authority/updateAuthority', 'PUT'),
('p', '1', '/authority/deleteAuthority', 'POST'),
('p', '1', '/authority/setDataAuthority', 'POST'),

-- 权限管理
('p', '1', '/api/v1/permission/list', 'GET'),
('p', '1', '/api/v1/permission/list', 'POST'),
('p', '1', '/api/v1/permission/:id', 'GET'),
('p', '1', '/api/v1/permission/create', 'POST'),
('p', '1', '/api/v1/permission/:id', 'PUT'),
('p', '1', '/api/v1/permission/:id', 'DELETE'),
('p', '1', '/casbin/updateCasbin', 'POST'),
('p', '1', '/casbin/getPolicyPathByAuthorityId', 'POST'),

-- 菜单管理
('p', '1', '/api/v1/menu/list', 'GET'),
('p', '1', '/api/v1/menu/list', 'POST'),
('p', '1', '/api/v1/menu/tree', 'GET'),
('p', '1', '/api/v1/menu/:id', 'GET'),
('p', '1', '/api/v1/menu/create', 'POST'),
('p', '1', '/api/v1/menu/:id', 'PUT'),
('p', '1', '/api/v1/menu/:id', 'DELETE'),
('p', '1', '/menu/getMenu', 'POST'),
('p', '1', '/menu/addBaseMenu', 'POST'),
('p', '1', '/menu/updateBaseMenu', 'POST'),
('p', '1', '/menu/deleteBaseMenu', 'POST'),
('p', '1', '/menu/getBaseMenuTree', 'POST'),
('p', '1', '/menu/getBaseMenuById', 'POST'),
('p', '1', '/menu/getMenuAuthority', 'POST'),
('p', '1', '/menu/addMenuAuthority', 'POST'),

-- API管理
('p', '1', '/api/v1/api/list', 'GET'),
('p', '1', '/api/v1/api/list', 'POST'),
('p', '1', '/api/v1/api/:id', 'GET'),
('p', '1', '/api/v1/api/create', 'POST'),
('p', '1', '/api/v1/api/:id', 'PUT'),
('p', '1', '/api/v1/api/:id', 'DELETE'),
('p', '1', '/api/getAllApis', 'POST'),
('p', '1', '/api/createApi', 'POST'),
('p', '1', '/api/updateApi', 'POST'),
('p', '1', '/api/deleteApi', 'POST'),
('p', '1', '/api/getApiById', 'POST'),
('p', '1', '/api/getApiList', 'POST'),
('p', '1', '/api/getApiGroups', 'GET'),
('p', '1', '/api/syncApi', 'GET'),
('p', '1', '/api/enterSyncApi', 'POST'),

-- 字典管理
('p', '1', '/api/v1/dict/list', 'GET'),
('p', '1', '/api/v1/dict/list', 'POST'),
('p', '1', '/api/v1/dict/:id', 'GET'),
('p', '1', '/api/v1/dict/create', 'POST'),
('p', '1', '/api/v1/dict/:id', 'PUT'),
('p', '1', '/api/v1/dict/:id', 'DELETE'),
('p', '1', '/sysDictionary/getSysDictionaryList', 'GET'),
('p', '1', '/sysDictionary/findSysDictionary', 'GET'),
('p', '1', '/sysDictionary/createSysDictionary', 'POST'),
('p', '1', '/sysDictionary/updateSysDictionary', 'PUT'),
('p', '1', '/sysDictionary/deleteSysDictionary', 'DELETE'),

-- 文件管理
('p', '1', '/fileUploadAndDownload/upload', 'POST'),
('p', '1', '/fileUploadAndDownload/getFileList', 'POST'),
('p', '1', '/fileUploadAndDownload/editFileName', 'POST'),
('p', '1', '/fileUploadAndDownload/deleteFile', 'POST'),
('p', '1', '/fileUploadAndDownload/importURL', 'POST'),

-- 应用市场
('p', '1', '/api/v1/app-market/list', 'GET'),
('p', '1', '/api/v1/app-market/:id', 'GET'),
('p', '1', '/api/v1/app-market/category/:category', 'GET'),
('p', '1', '/api/v1/app-market/search', 'GET'),

-- 系统管理
('p', '1', '/system/getSystemConfig', 'POST'),
('p', '1', '/system/setSystemConfig', 'POST'),
('p', '1', '/system/getServerInfo', 'POST'),

-- JWT相关
('p', '1', '/jwt/jsonInBlacklist', 'POST'),

-- ==================== 编辑角色权限(角色ID: 2) ====================
-- 基础权限
('p', '2', '/base/login', 'POST'),
('p', '2', '/user/getUserInfo', 'GET'),
('p', '2', '/user/setSelfInfo', 'PUT'),
('p', '2', '/user/changePassword', 'POST'),

-- 菜单权限
('p', '2', '/menu/getMenu', 'POST'),
('p', '2', '/api/v1/menu/list', 'GET'),
('p', '2', '/api/v1/menu/tree', 'GET'),
('p', '2', '/api/v1/menu/:id', 'GET'),

-- 字典权限(只读)
('p', '2', '/api/v1/dict/list', 'GET'),
('p', '2', '/api/v1/dict/:id', 'GET'),
('p', '2', '/sysDictionary/getSysDictionaryList', 'GET'),
('p', '2', '/sysDictionary/findSysDictionary', 'GET'),
('p', '2', '/api/v1/dict/create', 'POST'),
('p', '2', '/api/v1/dict/:id', 'PUT'),

-- 文件权限
('p', '2', '/fileUploadAndDownload/upload', 'POST'),
('p', '2', '/fileUploadAndDownload/getFileList', 'POST'),
('p', '2', '/fileUploadAndDownload/deleteFile', 'POST'),

-- 应用市场
('p', '2', '/api/v1/app-market/list', 'GET'),
('p', '2', '/api/v1/app-market/:id', 'GET'),
('p', '2', '/api/v1/app-market/search', 'GET'),

-- JWT相关
('p', '2', '/jwt/jsonInBlacklist', 'POST'),

-- ==================== 普通用户权限(角色ID: 3) ====================
-- 基础权限
('p', '3', '/base/login', 'POST'),
('p', '3', '/user/getUserInfo', 'GET'),
('p', '3', '/user/setSelfInfo', 'PUT'),
('p', '3', '/user/changePassword', 'POST'),

-- 菜单权限
('p', '3', '/menu/getMenu', 'POST'),
('p', '3', '/api/v1/menu/list', 'GET'),
('p', '3', '/api/v1/menu/tree', 'GET'),

-- 应用市场(只读)
('p', '3', '/api/v1/app-market/list', 'GET'),
('p', '3', '/api/v1/app-market/:id', 'GET'),
('p', '3', '/api/v1/app-market/search', 'GET'),

-- 用户信息(只读自己的)
('p', '3', '/api/v1/user/:id', 'GET'),

-- JWT相关
('p', '3', '/jwt/jsonInBlacklist', 'POST'),

-- ==================== 角色继承关系 ====================
-- g规则: (主体, 角色ID) - 定义用户属于哪个角色
-- g2规则: (用户名, 角色ID) - 用户的角色权限继承
-- 管理员用户
('g', 'admin', '1'),
('g2', 'admin', '1')

-- 初始化数据：权限
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
('api_create', '创建 API', 'apis', 'create', 1),
('api_update', '编辑API', 'apis', 'update', 1),
('api_delete', '删除API', 'apis', 'delete', 1),
('dict_view', '查看字典', 'dicts', 'view', 1),
('dict_create', '创建字典', 'dicts', 'create', 1),
('dict_update', '编辑字典', 'dicts', 'update', 1),
('dict_delete', '删除字典', 'dicts', 'delete', 1),
('app_view', '查看应用', 'apps', 'view', 1),
('app_install', '安装应用', 'apps', 'install', 1);

-- 关联管理员用户到管理员角色
INSERT INTO user_roles (user_id, role_id) VALUES 
(1, 1),  -- admin 用户 关联管理员角色
(2, 2),  -- editor 用户 关联编辑角色
(3, 3);  -- user 用户 关联普通用户角色

-- 为管理员角色分配所有权限
INSERT INTO role_permissions (role_id, permission_id)
SELECT 1, id FROM permissions WHERE status = 1;

-- 为编辑角色分配内容管理权限
INSERT INTO role_permissions (role_id, permission_id)
SELECT 2, id FROM permissions WHERE resource IN ('menus', 'dicts') AND action IN ('view', 'create', 'update');

-- 为用户角色分配基本权限
INSERT INTO role_permissions (role_id, permission_id)
SELECT 3, id FROM permissions WHERE resource IN ('users') AND action = 'view';

-- 初始化数据：菜单
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

-- 初始化数据：字典
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

-- 初始化数据：API接口
INSERT INTO apis (api_name, api_path, api_method, description, status) VALUES
-- 认证相关API
('用户登录', '/api/admin/auth/login', 'POST', '用户登录接口', 1),
('用户登出', '/api/admin/auth/logout', 'POST', '用户登出接口', 1),
('获取用户信息', '/api/admin/auth/info', 'GET', '获取当前登录用户信息', 1),
-- 用户管理API
('获取用户列表', '/api/admin/system/users', 'GET', '分页获取用户列表', 1),
('创建用户', '/api/admin/system/users', 'POST', '创建新用户', 1),
('编辑用户', '/api/admin/system/users', 'PUT', '编辑用户信息', 1),
('删除用户', '/api/admin/system/users', 'DELETE', '删除用户', 1),
('获取用户详情', '/api/admin/system/users/:id', 'GET', '根据ID获取用户详情', 1),
('为用户分配角色', '/api/admin/system/users/:id/roles', 'POST', '为用户分配角色', 1),
('获取用户角色', '/api/admin/system/users/:id/roles', 'GET', '获取用户已分配的角色', 1),
-- 角色管理API
('获取角色列表', '/api/admin/system/roles', 'GET', '分页获取角色列表', 1),
('创建角色', '/api/admin/system/roles', 'POST', '创建新角色', 1),
('编辑角色', '/api/admin/system/roles', 'PUT', '编辑角色信息', 1),
('删除角色', '/api/admin/system/roles', 'DELETE', '删除角色', 1),
('获取角色详情', '/api/admin/system/roles/:id', 'GET', '根据ID获取角色详情', 1),
('为角色分配权限', '/api/admin/system/roles/:id/permissions', 'POST', '为角色分配权限', 1),
('获取角色权限', '/api/admin/system/roles/:id/permissions', 'GET', '获取角色已分配的权限', 1),
-- 菜单管理API
('获取菜单列表', '/api/admin/system/menus', 'GET', '分页获取菜单列表', 1),
('创建菜单', '/api/admin/system/menus', 'POST', '创建新菜单', 1),
('编辑菜单', '/api/admin/system/menus', 'PUT', '编辑菜单信息', 1),
('删除菜单', '/api/admin/system/menus', 'DELETE', '删除菜单', 1),
('获取菜单详情', '/api/admin/system/menus/:id', 'GET', '根据ID获取菜单详情', 1),
-- 权限管理API
('获取权限列表', '/api/admin/system/permissions', 'GET', '分页获取权限列表', 1),
('创建权限', '/api/admin/system/permissions', 'POST', '创建新权限', 1),
('编辑权限', '/api/admin/system/permissions', 'PUT', '编辑权限信息', 1),
('删除权限', '/api/admin/system/permissions', 'DELETE', '删除权限', 1),
('获取权限详情', '/api/admin/system/permissions/:id', 'GET', '根据ID获取权限详情', 1),
-- 字典管理API
('获取字典列表', '/api/admin/content/dicts', 'GET', '分页获取字典列表', 1),
('创建字典', '/api/admin/content/dicts', 'POST', '创建新字典项', 1),
('编辑字典', '/api/admin/content/dicts', 'PUT', '编辑字典项', 1),
('删除字典', '/api/admin/content/dicts', 'DELETE', '删除字典项', 1),
('获取字典详情', '/api/admin/content/dicts/:id', 'GET', '根据ID获取字典项详情', 1),
-- API管理API
('获取API列表', '/api/admin/system/apis', 'GET', '分页获取API列表', 1),
('创建API', '/api/admin/system/apis', 'POST', '创建新API接口记录', 1),
('编辑API', '/api/admin/system/apis', 'PUT', '编辑API接口信息', 1),
('删除API', '/api/admin/system/apis', 'DELETE', '删除API接口记录', 1),
('获取API详情', '/api/admin/system/apis/:id', 'GET', '根据ID获取API接口详情', 1);

-- 应用市场表
CREATE TABLE IF NOT EXISTS apps (
    id BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '应用ID',
    app_key VARCHAR(100) UNIQUE NOT NULL COMMENT '应用标识key',
    app_name VARCHAR(100) NOT NULL COMMENT '应用名称',
    version VARCHAR(50) COMMENT '应用版本',
    author VARCHAR(100) COMMENT '作者',
    description TEXT COMMENT '应用描述',
    icon VARCHAR(255) COMMENT '应用图标URL',
    download_url VARCHAR(255) COMMENT '下载地址',
    demo_url VARCHAR(255) COMMENT '演示地址',
    category VARCHAR(50) COMMENT '应用分类',
    tags VARCHAR(255) COMMENT '应用标签',
    rating DECIMAL(3,2) DEFAULT 0 COMMENT '评分',
    downloads BIGINT DEFAULT 0 COMMENT '下载次数',
    status INT DEFAULT 1 COMMENT '状态 1:上架 0:下架',
    published INT DEFAULT 1 COMMENT '是否已发布 1:是 0:否',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    deleted_at TIMESTAMP NULL COMMENT '删除时间',
    INDEX idx_app_key (app_key),
    INDEX idx_category (category),
    INDEX idx_created_at (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='应用市场表';

-- 应用安装记录表
CREATE TABLE IF NOT EXISTS app_installations (
    id BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '安装记录ID',
    app_key VARCHAR(100) UNIQUE NOT NULL COMMENT '应用标识',
    app_id BIGINT COMMENT '应用ID',
    app_name VARCHAR(100) COMMENT '应用名称',
    version VARCHAR(50) COMMENT '安装版本',
    status INT DEFAULT 1 COMMENT '安装状态(1:已安装,0:未安装)',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '安装时间',
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    INDEX idx_app_key (app_key),
    INDEX idx_app_id (app_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='应用安装记录表';

-- 插入应用市场数据
INSERT INTO apps (app_key, app_name, version, author, description, icon, category, tags, rating, downloads, status, published) VALUES
('cms', 'CMS内容系统', '1.0.0', 'Power Admin Team', '完整的内容管理系统，支持文章、分类、标签、评论管理，适合博客和内容运营', 'https://via.placeholder.com/200?text=CMS', 'content', '内容管理,博客', 4.8, 1250, 1, 1),
('shop', '电商系统', '2.0.0', 'Power Admin Team', '功能完整的电商平台，包含商品管理、订单、支付、物流跟踪等功能', 'https://via.placeholder.com/200?text=Shop', 'business', '电商,商城', 4.6, 980, 1, 1),
('crm', 'CRM客户管理', '1.5.0', 'Power Admin Team', '企业客户关系管理系统，支持客户跟进、商机管理、合同管理等', 'https://via.placeholder.com/200?text=CRM', 'business', '客户管理,销售', 4.5, 750, 1, 1),
('finance', '财务管理系统', '1.2.0', 'Power Admin Team', '企业财务管理工具，涵盖应收应付、报表、成本管理等功能', 'https://via.placeholder.com/200?text=Finance', 'finance', '财务,报表', 4.7, 680, 1, 1),
('hr', 'HR人力资源', '1.1.0', 'Power Admin Team', '人力资源管理平台，支持招聘、员工档案、考勤、薪资管理', 'https://via.placeholder.com/200?text=HR', 'business', '人力资源,招聘', 4.4, 520, 1, 1),
('marketing', '营销工具箱', '1.0.0', 'Power Admin Team', '集成营销工具，支持活动管理、积分管理、优惠券等营销功能', 'https://via.placeholder.com/200?text=Marketing', 'marketing', '营销,促销', 4.3, 400, 1, 1);

-- 插入已安装应用记录（仅CMS已安装）
INSERT INTO app_installations (app_key, app_id, app_name, version, status) VALUES
('cms', 1, 'CMS内容系统', '1.0.0', 1);
