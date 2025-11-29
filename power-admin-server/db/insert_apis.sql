-- 插入所有API接口数据到apis表
-- 注意：这个脚本应该在初始化脚本后执行

-- 认证相关API
INSERT INTO apis (api_name, api_path, api_method, description, status) VALUES
('用户登录', '/api/admin/auth/login', 'POST', '用户登录接口', 1),
('用户登出', '/api/admin/auth/logout', 'POST', '用户登出接口', 1),
('获取用户信息', '/api/admin/auth/info', 'GET', '获取当前登录用户信息', 1);

-- 用户管理API
INSERT INTO apis (api_name, api_path, api_method, description, status) VALUES
('获取用户列表', '/api/admin/system/users', 'GET', '分页获取用户列表', 1),
('创建用户', '/api/admin/system/users', 'POST', '创建新用户', 1),
('编辑用户', '/api/admin/system/users', 'PUT', '编辑用户信息', 1),
('删除用户', '/api/admin/system/users', 'DELETE', '删除用户', 1),
('获取用户详情', '/api/admin/system/users/:id', 'GET', '根据ID获取用户详情', 1),
('为用户分配角色', '/api/admin/system/users/:id/roles', 'POST', '为用户分配角色', 1),
('获取用户角色', '/api/admin/system/users/:id/roles', 'GET', '获取用户已分配的角色', 1);

-- 角色管理API
INSERT INTO apis (api_name, api_path, api_method, description, status) VALUES
('获取角色列表', '/api/admin/system/roles', 'GET', '分页获取角色列表', 1),
('创建角色', '/api/admin/system/roles', 'POST', '创建新角色', 1),
('编辑角色', '/api/admin/system/roles', 'PUT', '编辑角色信息', 1),
('删除角色', '/api/admin/system/roles', 'DELETE', '删除角色', 1),
('获取角色详情', '/api/admin/system/roles/:id', 'GET', '根据ID获取角色详情', 1),
('为角色分配权限', '/api/admin/system/roles/:id/permissions', 'POST', '为角色分配权限', 1),
('获取角色权限', '/api/admin/system/roles/:id/permissions', 'GET', '获取角色已分配的权限', 1);

-- 菜单管理API
INSERT INTO apis (api_name, api_path, api_method, description, status) VALUES
('获取菜单列表', '/api/admin/system/menus', 'GET', '分页获取菜单列表', 1),
('创建菜单', '/api/admin/system/menus', 'POST', '创建新菜单', 1),
('编辑菜单', '/api/admin/system/menus', 'PUT', '编辑菜单信息', 1),
('删除菜单', '/api/admin/system/menus', 'DELETE', '删除菜单', 1),
('获取菜单详情', '/api/admin/system/menus/:id', 'GET', '根据ID获取菜单详情', 1);

-- 权限管理API
INSERT INTO apis (api_name, api_path, api_method, description, status) VALUES
('获取权限列表', '/api/admin/system/permissions', 'GET', '分页获取权限列表', 1),
('创建权限', '/api/admin/system/permissions', 'POST', '创建新权限', 1),
('编辑权限', '/api/admin/system/permissions', 'PUT', '编辑权限信息', 1),
('删除权限', '/api/admin/system/permissions', 'DELETE', '删除权限', 1),
('获取权限详情', '/api/admin/system/permissions/:id', 'GET', '根据ID获取权限详情', 1);

-- 字典管理API
INSERT INTO apis (api_name, api_path, api_method, description, status) VALUES
('获取字典列表', '/api/admin/content/dicts', 'GET', '分页获取字典列表', 1),
('创建字典', '/api/admin/content/dicts', 'POST', '创建新字典项', 1),
('编辑字典', '/api/admin/content/dicts', 'PUT', '编辑字典项', 1),
('删除字典', '/api/admin/content/dicts', 'DELETE', '删除字典项', 1),
('获取字典详情', '/api/admin/content/dicts/:id', 'GET', '根据ID获取字典项详情', 1);

-- API管理API
INSERT INTO apis (api_name, api_path, api_method, description, status) VALUES
('获取API列表', '/api/admin/system/apis', 'GET', '分页获取API列表', 1),
('创建API', '/api/admin/system/apis', 'POST', '创建新API接口记录', 1),
('编辑API', '/api/admin/system/apis', 'PUT', '编辑API接口信息', 1),
('删除API', '/api/admin/system/apis', 'DELETE', '删除API接口记录', 1),
('获取API详情', '/api/admin/system/apis/:id', 'GET', '根据ID获取API接口详情', 1);
