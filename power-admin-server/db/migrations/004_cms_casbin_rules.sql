-- =============================================
-- CMS权限规则初始化 (Casbin RBAC规则)
-- =============================================
-- 说明: 这个脚本添加CMS相关的Casbin规则到现有的权限系统

-- =============================================
-- CMS角色权限定义 (p规则)
-- =============================================

-- CMS管理员权限 (所有操作)
INSERT IGNORE INTO `casbin_rule` (`ptype`, `v0`, `v1`, `v2`, `v3`) VALUES
('p', 'cms_admin', '/api/cms/admin/contents', 'GET', ''),
('p', 'cms_admin', '/api/cms/admin/contents', 'POST', ''),
('p', 'cms_admin', '/api/cms/admin/contents', 'PUT', ''),
('p', 'cms_admin', '/api/cms/admin/contents/:id', 'DELETE', ''),
('p', 'cms_admin', '/api/cms/admin/categories', 'GET', ''),
('p', 'cms_admin', '/api/cms/admin/categories', 'POST', ''),
('p', 'cms_admin', '/api/cms/admin/categories', 'PUT', ''),
('p', 'cms_admin', '/api/cms/admin/categories/:id', 'DELETE', ''),
('p', 'cms_admin', '/api/cms/admin/tags', 'GET', ''),
('p', 'cms_admin', '/api/cms/admin/tags', 'POST', ''),
('p', 'cms_admin', '/api/cms/admin/tags', 'PUT', ''),
('p', 'cms_admin', '/api/cms/admin/tags/:id', 'DELETE', ''),
('p', 'cms_admin', '/api/cms/admin/comments', 'GET', ''),
('p', 'cms_admin', '/api/cms/admin/comments/:id/approve', 'POST', ''),
('p', 'cms_admin', '/api/cms/admin/comments/:id', 'DELETE', ''),
('p', 'cms_admin', '/api/cms/admin/users', 'GET', ''),
('p', 'cms_admin', '/api/cms/admin/users/:id', 'PUT', ''),
('p', 'cms_admin', '/api/cms/admin/users/:id', 'DELETE', ''),
('p', 'cms_admin', '/api/cms/admin/revisions', 'GET', ''),
('p', 'cms_admin', '/api/cms/admin/drafts', 'GET', ''),
('p', 'cms_admin', '/api/cms/admin/permissions', 'GET', ''),
('p', 'cms_admin', '/api/cms/admin/permissions', 'POST', ''),
('p', 'cms_admin', '/api/cms/admin/stats', 'GET', ''),

-- CMS编辑权限 (创建/编辑/发布，不能删除或管理用户)
('p', 'cms_editor', '/api/cms/admin/contents', 'GET', ''),
('p', 'cms_editor', '/api/cms/admin/contents', 'POST', ''),
('p', 'cms_editor', '/api/cms/admin/contents/:id', 'PUT', ''),
('p', 'cms_editor', '/api/cms/admin/categories', 'GET', ''),
('p', 'cms_editor', '/api/cms/admin/tags', 'GET', ''),
('p', 'cms_editor', '/api/cms/admin/tags', 'POST', ''),
('p', 'cms_editor', '/api/cms/admin/comments', 'GET', ''),
('p', 'cms_editor', '/api/cms/admin/revisions', 'GET', ''),
('p', 'cms_editor', '/api/cms/admin/drafts', 'GET', ''),
('p', 'cms_editor', '/api/cms/admin/drafts', 'POST', ''),
('p', 'cms_editor', '/api/cms/admin/stats', 'GET', ''),

-- CMS查看者权限 (仅查看)
('p', 'cms_viewer', '/api/cms/admin/contents', 'GET', ''),
('p', 'cms_viewer', '/api/cms/admin/categories', 'GET', ''),
('p', 'cms_viewer', '/api/cms/admin/tags', 'GET', ''),
('p', 'cms_viewer', '/api/cms/admin/comments', 'GET', ''),
('p', 'cms_viewer', '/api/cms/admin/stats', 'GET', ''),

-- =============================================
-- 用户角色映射 (g规则)
-- =============================================

-- 注意: 实际的用户-CMS角色映射应该通过cms_admin_roles表管理
-- 这里仅作为参考示例，生产环境中应该通过系统动态管理

-- 示例: 超级管理员映射到CMS管理员角色
-- INSERT IGNORE INTO `casbin_rule` (`ptype`, `v0`, `v1`) VALUES ('g', '1', 'cms_admin');

-- =============================================
-- 资源权限映射 (可选，用于细粒度控制)
-- =============================================

-- 这部分根据具体的业务需求添加
-- 例如: 特定分类只能由特定用户管理

-- =============================================
-- 创建CMS权限规则支持的p规则模式
-- =============================================
-- 支持的规则格式:
-- ('p', role, resource, action)
-- 例如: ('p', 'cms_editor', '/api/cms/admin/contents', 'GET')

