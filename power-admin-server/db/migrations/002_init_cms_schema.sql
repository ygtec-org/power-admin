-- =============================================
-- CMS 内容管理系统数据表初始化脚本
-- =============================================

-- =============================================
-- 1. CMS 内容表
-- =============================================
CREATE TABLE IF NOT EXISTS `cms_content` (
  `id` BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '内容ID',
  `title` VARCHAR(255) NOT NULL COMMENT '文章标题',
  `slug` VARCHAR(255) UNIQUE COMMENT 'URL别名，用于SEO',
  `content` LONGTEXT COMMENT '文章内容(HTML格式)',
  `excerpt` VARCHAR(500) COMMENT '文章摘要，自动截取或手动编辑',
  `category_id` BIGINT COMMENT '分类ID',
  `author_id` BIGINT NOT NULL COMMENT '作者ID(对应系统用户)',
  `status` TINYINT DEFAULT 1 COMMENT '1:草稿 2:已发布 3:已删除',
  `view_count` INT DEFAULT 0 COMMENT '浏览次数',
  `comment_count` INT DEFAULT 0 COMMENT '评论数',
  `featured_image` VARCHAR(500) COMMENT '特色图片URL',
  `seo_title` VARCHAR(255) COMMENT 'SEO标题',
  `seo_keywords` VARCHAR(255) COMMENT 'SEO关键词',
  `seo_description` VARCHAR(500) COMMENT 'SEO描述',
  `is_pinned` TINYINT DEFAULT 0 COMMENT '是否置顶',
  `is_recommended` TINYINT DEFAULT 0 COMMENT '是否推荐',
  `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `published_at` TIMESTAMP NULL COMMENT '发布时间',
  KEY `idx_category` (`category_id`),
  KEY `idx_author` (`author_id`),
  KEY `idx_status` (`status`),
  KEY `idx_created` (`created_at`),
  KEY `idx_published` (`published_at`),
  KEY `idx_is_pinned` (`is_pinned`),
  CONSTRAINT `fk_cms_content_author` FOREIGN KEY (`author_id`) REFERENCES `admin_users` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='CMS内容表';

-- =============================================
-- 2. CMS 分类表
-- =============================================
CREATE TABLE IF NOT EXISTS `cms_category` (
  `id` BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '分类ID',
  `name` VARCHAR(100) NOT NULL COMMENT '分类名称',
  `slug` VARCHAR(100) UNIQUE COMMENT 'URL别名',
  `description` TEXT COMMENT '分类描述',
  `parent_id` BIGINT COMMENT '父分类ID，支持多级分类',
  `sort` INT DEFAULT 0 COMMENT '排序号，数字越大越靠前',
  `status` TINYINT DEFAULT 1 COMMENT '1:启用 0:禁用',
  `content_count` INT DEFAULT 0 COMMENT '分类下的内容数量',
  `seo_keywords` VARCHAR(255) COMMENT 'SEO关键词',
  `seo_description` VARCHAR(500) COMMENT 'SEO描述',
  `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  KEY `idx_parent` (`parent_id`),
  KEY `idx_status` (`status`),
  KEY `idx_sort` (`sort`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='CMS分类表';

-- =============================================
-- 3. CMS 标签表 (可选)
-- =============================================
CREATE TABLE IF NOT EXISTS `cms_tag` (
  `id` BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '标签ID',
  `name` VARCHAR(50) NOT NULL UNIQUE COMMENT '标签名称',
  `slug` VARCHAR(50) UNIQUE COMMENT 'URL别名',
  `description` TEXT COMMENT '标签描述',
  `color` VARCHAR(7) COMMENT '标签颜色(HEX格式)',
  `usage_count` INT DEFAULT 0 COMMENT '使用次数',
  `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='CMS标签表';

-- =============================================
-- 4. CMS 内容-标签关联表
-- =============================================
CREATE TABLE IF NOT EXISTS `cms_content_tag` (
  `content_id` BIGINT NOT NULL,
  `tag_id` BIGINT NOT NULL,
  PRIMARY KEY (`content_id`, `tag_id`),
  CONSTRAINT `fk_content_tag_content` FOREIGN KEY (`content_id`) REFERENCES `cms_content` (`id`) ON DELETE CASCADE,
  CONSTRAINT `fk_content_tag_tag` FOREIGN KEY (`tag_id`) REFERENCES `cms_tag` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='CMS内容-标签关联表';

-- =============================================
-- 5. CMS 前台用户表 (访客)
-- =============================================
CREATE TABLE IF NOT EXISTS `cms_users` (
  `id` BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '用户ID',
  `username` VARCHAR(100) UNIQUE NOT NULL COMMENT '用户名',
  `email` VARCHAR(255) UNIQUE COMMENT '邮箱地址',
  `phone` VARCHAR(20) COMMENT '电话号码',
  `password` VARCHAR(255) NOT NULL COMMENT '密码(加密存储)',
  `nickname` VARCHAR(100) COMMENT '昵称',
  `avatar` VARCHAR(500) COMMENT '头像URL',
  `bio` TEXT COMMENT '个人简介',
  `gender` TINYINT COMMENT '性别: 0=未知 1=男 2=女',
  `status` TINYINT DEFAULT 1 COMMENT '1:正常 0:禁用 2:封禁',
  `email_verified` TINYINT DEFAULT 0 COMMENT '邮箱是否验证',
  `phone_verified` TINYINT DEFAULT 0 COMMENT '电话是否验证',
  `last_login_at` TIMESTAMP NULL COMMENT '上次登录时间',
  `last_login_ip` VARCHAR(45) COMMENT '上次登录IP(支持IPv6)',
  `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  KEY `idx_email` (`email`),
  KEY `idx_phone` (`phone`),
  KEY `idx_status` (`status`),
  KEY `idx_created` (`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='CMS前台用户表(访客)';

-- =============================================
-- 6. CMS 评论表 (可选)
-- =============================================
CREATE TABLE IF NOT EXISTS `cms_comments` (
  `id` BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '评论ID',
  `content_id` BIGINT NOT NULL COMMENT '内容ID',
  `user_id` BIGINT COMMENT '评论者ID(cms_users表)',
  `parent_comment_id` BIGINT COMMENT '父评论ID，支持回复',
  `author_name` VARCHAR(100) COMMENT '评论者名称(如果未登录)',
  `author_email` VARCHAR(255) COMMENT '评论者邮箱(如果未登录)',
  `content` TEXT NOT NULL COMMENT '评论内容',
  `status` TINYINT DEFAULT 0 COMMENT '0:待审核 1:已批准 2:垃圾评论',
  `like_count` INT DEFAULT 0 COMMENT '点赞数',
  `ip_address` VARCHAR(45) COMMENT 'IP地址',
  `user_agent` VARCHAR(500) COMMENT '浏览器标识',
  `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  KEY `idx_content` (`content_id`),
  KEY `idx_user` (`user_id`),
  KEY `idx_status` (`status`),
  CONSTRAINT `fk_comment_content` FOREIGN KEY (`content_id`) REFERENCES `cms_content` (`id`) ON DELETE CASCADE,
  CONSTRAINT `fk_comment_user` FOREIGN KEY (`user_id`) REFERENCES `cms_users` (`id`) ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='CMS评论表';

-- =============================================
-- 7. CMS 权限表
-- =============================================
CREATE TABLE IF NOT EXISTS `cms_permissions` (
  `id` BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '权限ID',
  `name` VARCHAR(100) UNIQUE NOT NULL COMMENT '权限名称',
  `description` TEXT COMMENT '权限描述',
  `resource` VARCHAR(100) COMMENT '资源(如:cms_content, cms_category)',
  `action` VARCHAR(100) COMMENT '操作(如:create, read, update, delete)',
  `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  KEY `idx_resource_action` (`resource`, `action`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='CMS权限表';

-- =============================================
-- 8. CMS 管理员角色映射表
-- =============================================
CREATE TABLE IF NOT EXISTS `cms_admin_roles` (
  `id` BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '记录ID',
  `admin_id` BIGINT NOT NULL COMMENT '主系统管理员ID',
  `role_name` VARCHAR(50) NOT NULL COMMENT 'CMS角色: cms_admin/cms_editor/cms_viewer',
  `assigned_by` BIGINT COMMENT '由哪个管理员分配',
  `assigned_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '分配时间',
  `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  UNIQUE KEY `uk_admin_role` (`admin_id`, `role_name`),
  KEY `idx_admin` (`admin_id`),
  CONSTRAINT `fk_cms_admin_role_admin` FOREIGN KEY (`admin_id`) REFERENCES `admin_users` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='CMS-系统用户角色映射表';

-- =============================================
-- 9. 插件状态表
-- =============================================
CREATE TABLE IF NOT EXISTS `plugin_status` (
  `id` BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '记录ID',
  `plugin_name` VARCHAR(100) UNIQUE NOT NULL COMMENT '插件名称(如:cms)',
  `enabled` TINYINT DEFAULT 0 COMMENT '0:禁用 1:启用',
  `version` VARCHAR(20) COMMENT '插件版本',
  `config` JSON COMMENT '插件配置(JSON格式)',
  `installed_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '安装时间',
  `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  KEY `idx_name` (`plugin_name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='插件启用状态表';

-- =============================================
-- 10. CMS 操作日志表 (审计)
-- =============================================
CREATE TABLE IF NOT EXISTS `cms_audit_logs` (
  `id` BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '日志ID',
  `admin_id` BIGINT NOT NULL COMMENT '操作管理员ID',
  `action` VARCHAR(50) COMMENT '操作类型(如:create_content, edit_content)',
  `resource_type` VARCHAR(50) COMMENT '资源类型(如:content, category)',
  `resource_id` BIGINT COMMENT '资源ID',
  `old_value` JSON COMMENT '修改前的值',
  `new_value` JSON COMMENT '修改后的值',
  `ip_address` VARCHAR(45) COMMENT 'IP地址',
  `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  KEY `idx_admin` (`admin_id`),
  KEY `idx_action` (`action`),
  KEY `idx_created` (`created_at`),
  CONSTRAINT `fk_audit_admin` FOREIGN KEY (`admin_id`) REFERENCES `admin_users` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='CMS操作审计日志表';

-- =============================================
-- 初始化数据
-- =============================================

-- 初始化CMS权限
INSERT INTO `cms_permissions` (`name`, `description`, `resource`, `action`) VALUES
('查看CMS内容', '查看CMS内容列表和详情', 'cms_content', 'read'),
('创建CMS内容', '创建新的CMS内容', 'cms_content', 'create'),
('编辑CMS内容', '编辑已存在的CMS内容', 'cms_content', 'update'),
('删除CMS内容', '删除CMS内容', 'cms_content', 'delete'),
('发布CMS内容', '发布CMS内容', 'cms_content', 'publish'),
('管理CMS分类', '创建/编辑/删除CMS分类', 'cms_category', 'manage'),
('管理CMS标签', '创建/编辑/删除CMS标签', 'cms_tag', 'manage'),
('管理CMS用户', '管理CMS前台用户', 'cms_users', 'manage'),
('管理CMS评论', '审核/删除CMS评论', 'cms_comments', 'manage'),
('查看CMS统计', '查看CMS统计数据', 'cms_stats', 'read')
ON DUPLICATE KEY UPDATE description=VALUES(description);

-- 初始化CMS插件状态
INSERT INTO `plugin_status` (`plugin_name`, `enabled`, `version`, `config`, `installed_at`) VALUES
('cms', 1, '1.0.0', '{"enable_comments": true, "enable_ratings": false, "posts_per_page": 10}', NOW())
ON DUPLICATE KEY UPDATE enabled=VALUES(enabled), version=VALUES(version), updated_at=NOW();

-- =============================================
-- Casbin 权限规则 (添加到现有的 casbin_rule 表)
-- =============================================
-- 注意: 这些规则应添加到现有的 casbin_rule 表中

INSERT IGNORE INTO `casbin_rule` (`ptype`, `v0`, `v1`, `v2`, `v3`, `v4`, `v5`) VALUES
-- CMS 管理员权限 (所有操作)
('p', 'cms_admin', '/api/cms/admin/contents', 'GET', '', '', ''),
('p', 'cms_admin', '/api/cms/admin/contents', 'POST', '', '', ''),
('p', 'cms_admin', '/api/cms/admin/contents', 'PUT', '', '', ''),
('p', 'cms_admin', '/api/cms/admin/contents', 'DELETE', '', '', ''),
('p', 'cms_admin', '/api/cms/admin/categories', 'GET', '', '', ''),
('p', 'cms_admin', '/api/cms/admin/categories', 'POST', '', '', ''),
('p', 'cms_admin', '/api/cms/admin/categories', 'PUT', '', '', ''),
('p', 'cms_admin', '/api/cms/admin/categories', 'DELETE', '', '', ''),
('p', 'cms_admin', '/api/cms/admin/publish', 'POST', '', '', ''),
('p', 'cms_admin', '/api/cms/admin/users', 'GET', '', '', ''),
('p', 'cms_admin', '/api/cms/admin/users', 'POST', '', '', ''),
('p', 'cms_admin', '/api/cms/admin/users', 'PUT', '', '', ''),
('p', 'cms_admin', '/api/cms/admin/users', 'DELETE', '', '', ''),

-- CMS 编辑权限 (创建/编辑/发布)
('p', 'cms_editor', '/api/cms/admin/contents', 'GET', '', '', ''),
('p', 'cms_editor', '/api/cms/admin/contents', 'POST', '', '', ''),
('p', 'cms_editor', '/api/cms/admin/contents', 'PUT', '', '', ''),
('p', 'cms_editor', '/api/cms/admin/categories', 'GET', '', '', ''),
('p', 'cms_editor', '/api/cms/admin/publish', 'POST', '', '', ''),

-- CMS 查看者权限 (仅查看)
('p', 'cms_viewer', '/api/cms/admin/contents', 'GET', '', '', ''),
('p', 'cms_viewer', '/api/cms/admin/categories', 'GET', '', '', ''),

-- 超级管理员 (user_id=1) 继承所有权限
('g', '1', 'cms_admin', '', '', '', ''),

-- 公开API (无需认证)
('p', 'public_user', '/api/cms/public/contents', 'GET', '', '', ''),
('p', 'public_user', '/api/cms/public/categories', 'GET', '', '', ''),
('p', 'public_user', '/api/cms/public/auth/login', 'POST', '', '', ''),
('p', 'public_user', '/api/cms/public/auth/register', 'POST', '', '', '');

-- =============================================
-- 完成初始化
-- =============================================

-- 验证表是否创建成功
-- SELECT TABLE_NAME FROM INFORMATION_SCHEMA.TABLES WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME LIKE 'cms_%';

