-- =============================================
-- CMS 内容管理系统完整数据库设计
-- =============================================
-- 创建时间: 2024年
-- 版本: 1.0.0
-- 说明: 生产级别CMS系统数据库定义

-- =============================================
-- 1. CMS 内容表 (核心表)
-- =============================================
CREATE TABLE IF NOT EXISTS `cms_content` (
  `id` BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '内容ID',
  `title` VARCHAR(255) NOT NULL COMMENT '文章标题',
  `slug` VARCHAR(255) UNIQUE COMMENT 'URL别名，用于SEO',
  `description` VARCHAR(500) COMMENT '文章描述/摘要',
  `content` LONGTEXT NOT NULL COMMENT '文章内容(HTML格式)',
  `featured_image` VARCHAR(500) COMMENT '特色图片URL',
  `featured_image_alt` VARCHAR(255) COMMENT '特色图片alt文本',
  `category_id` BIGINT COMMENT '分类ID',
  `author_id` BIGINT NOT NULL COMMENT '作者ID(对应系统用户)',
  `status` TINYINT DEFAULT 1 COMMENT '1:草稿 2:已发布 3:已删除',
  `visibility` TINYINT DEFAULT 1 COMMENT '1:公开 2:受保护 3:私密',
  `comment_status` TINYINT DEFAULT 1 COMMENT '1:允许评论 0:不允许',
  `view_count` INT DEFAULT 0 COMMENT '浏览次数',
  `comment_count` INT DEFAULT 0 COMMENT '评论数',
  `like_count` INT DEFAULT 0 COMMENT '点赞数',
  
  -- SEO字段
  `seo_title` VARCHAR(255) COMMENT 'SEO标题',
  `seo_keywords` VARCHAR(255) COMMENT 'SEO关键词',
  `seo_description` VARCHAR(500) COMMENT 'SEO描述',
  
  -- 时间字段
  `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `published_at` TIMESTAMP NULL COMMENT '发布时间',
  `scheduled_at` TIMESTAMP NULL COMMENT '定时发布时间',
  
  -- 字段
  `is_featured` TINYINT DEFAULT 0 COMMENT '是否精选',
  `is_sticky` TINYINT DEFAULT 0 COMMENT '是否置顶',
  `revision_count` INT DEFAULT 0 COMMENT '版本数量',
  
  -- 索引
  KEY `idx_category` (`category_id`),
  KEY `idx_author` (`author_id`),
  KEY `idx_status` (`status`),
  KEY `idx_visibility` (`visibility`),
  KEY `idx_created` (`created_at`),
  KEY `idx_published` (`published_at`),
  KEY `idx_slug` (`slug`),
  KEY `idx_is_featured` (`is_featured`),
  KEY `idx_is_sticky` (`is_sticky`),
  
  -- 外键约束
  CONSTRAINT `fk_cms_content_author` FOREIGN KEY (`author_id`) 
    REFERENCES `admin_users` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='CMS内容表(核心)';

-- =============================================
-- 2. CMS 分类表 (支持多级分类)
-- =============================================
CREATE TABLE IF NOT EXISTS `cms_category` (
  `id` BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '分类ID',
  `name` VARCHAR(100) NOT NULL COMMENT '分类名称',
  `slug` VARCHAR(100) UNIQUE COMMENT '分类URL别名',
  `description` TEXT COMMENT '分类描述',
  `thumbnail` VARCHAR(500) COMMENT '分类缩略图',
  `parent_id` BIGINT COMMENT '父分类ID，支持多级分类',
  `sort` INT DEFAULT 0 COMMENT '排序号，数字越大越靠前',
  `status` TINYINT DEFAULT 1 COMMENT '1:启用 0:禁用',
  `content_count` INT DEFAULT 0 COMMENT '分类下的内容数量',
  
  -- SEO字段
  `seo_keywords` VARCHAR(255) COMMENT 'SEO关键词',
  `seo_description` VARCHAR(500) COMMENT 'SEO描述',
  
  `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  
  KEY `idx_parent` (`parent_id`),
  KEY `idx_status` (`status`),
  KEY `idx_sort` (`sort`),
  KEY `idx_slug` (`slug`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='CMS分类表(支持多级)';

-- =============================================
-- 3. CMS 标签表
-- =============================================
CREATE TABLE IF NOT EXISTS `cms_tag` (
  `id` BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '标签ID',
  `name` VARCHAR(50) NOT NULL UNIQUE COMMENT '标签名称',
  `slug` VARCHAR(50) UNIQUE COMMENT '标签URL别名',
  `description` TEXT COMMENT '标签描述',
  `color` VARCHAR(7) COMMENT '标签颜色(HEX格式)',
  `usage_count` INT DEFAULT 0 COMMENT '使用次数',
  `status` TINYINT DEFAULT 1 COMMENT '1:启用 0:禁用',
  `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  
  KEY `idx_status` (`status`),
  KEY `idx_usage_count` (`usage_count`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='CMS标签表';

-- =============================================
-- 4. CMS 内容-标签关联表
-- =============================================
CREATE TABLE IF NOT EXISTS `cms_content_tag` (
  `content_id` BIGINT NOT NULL,
  `tag_id` BIGINT NOT NULL,
  `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`content_id`, `tag_id`),
  CONSTRAINT `fk_content_tag_content` FOREIGN KEY (`content_id`) 
    REFERENCES `cms_content` (`id`) ON DELETE CASCADE,
  CONSTRAINT `fk_content_tag_tag` FOREIGN KEY (`tag_id`) 
    REFERENCES `cms_tag` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='CMS内容-标签关联表';

-- =============================================
-- 5. CMS 内容版本历史表
-- =============================================
CREATE TABLE IF NOT EXISTS `cms_content_revision` (
  `id` BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '版本ID',
  `content_id` BIGINT NOT NULL COMMENT '内容ID',
  `title` VARCHAR(255) NOT NULL,
  `content` LONGTEXT NOT NULL,
  `revision_number` INT NOT NULL COMMENT '版本号',
  `author_id` BIGINT NOT NULL COMMENT '编辑者ID',
  `change_summary` VARCHAR(500) COMMENT '更改摘要',
  `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  
  UNIQUE KEY `uk_content_revision` (`content_id`, `revision_number`),
  KEY `idx_content_id` (`content_id`),
  KEY `idx_author_id` (`author_id`),
  CONSTRAINT `fk_revision_content` FOREIGN KEY (`content_id`) 
    REFERENCES `cms_content` (`id`) ON DELETE CASCADE,
  CONSTRAINT `fk_revision_author` FOREIGN KEY (`author_id`) 
    REFERENCES `admin_users` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='CMS内容版本历史表';

-- =============================================
-- 6. CMS 前台用户表 (访客)
-- =============================================
CREATE TABLE IF NOT EXISTS `cms_users` (
  `id` BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '用户ID',
  `username` VARCHAR(100) UNIQUE NOT NULL COMMENT '用户名',
  `email` VARCHAR(255) UNIQUE COMMENT '邮箱地址',
  `phone` VARCHAR(20) COMMENT '电话号码',
  `password` VARCHAR(255) NOT NULL COMMENT '密码(bcrypt加密)',
  `nickname` VARCHAR(100) COMMENT '昵称',
  `avatar` VARCHAR(500) COMMENT '头像URL',
  `bio` TEXT COMMENT '个人简介',
  `gender` TINYINT COMMENT '性别: 0=未知 1=男 2=女',
  `status` TINYINT DEFAULT 1 COMMENT '1:正常 0:禁用 2:封禁',
  
  -- 验证字段
  `email_verified` TINYINT DEFAULT 0 COMMENT '邮箱是否验证',
  `phone_verified` TINYINT DEFAULT 0 COMMENT '电话是否验证',
  `email_verified_at` TIMESTAMP NULL COMMENT '邮箱验证时间',
  `phone_verified_at` TIMESTAMP NULL COMMENT '电话验证时间',
  
  -- 登录信息
  `last_login_at` TIMESTAMP NULL COMMENT '上次登录时间',
  `last_login_ip` VARCHAR(45) COMMENT '上次登录IP(支持IPv6)',
  `login_count` INT DEFAULT 0 COMMENT '登录次数',
  
  -- 时间字段
  `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `deleted_at` TIMESTAMP NULL COMMENT '删除时间(软删除)',
  
  KEY `idx_email` (`email`),
  KEY `idx_phone` (`phone`),
  KEY `idx_status` (`status`),
  KEY `idx_created` (`created_at`),
  KEY `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='CMS前台用户表(访客)';

-- =============================================
-- 7. CMS 评论表
-- =============================================
CREATE TABLE IF NOT EXISTS `cms_comments` (
  `id` BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '评论ID',
  `content_id` BIGINT NOT NULL COMMENT '内容ID',
  `user_id` BIGINT COMMENT '评论者ID(cms_users表)',
  `parent_comment_id` BIGINT COMMENT '父评论ID，支持回复',
  
  -- 访客信息(未登录时)
  `author_name` VARCHAR(100) COMMENT '评论者名称(如果未登录)',
  `author_email` VARCHAR(255) COMMENT '评论者邮箱(如果未登录)',
  
  -- 评论内容
  `content` TEXT NOT NULL COMMENT '评论内容',
  `status` TINYINT DEFAULT 0 COMMENT '0:待审核 1:已批准 2:垃圾评论 3:已删除',
  
  -- 互动数据
  `like_count` INT DEFAULT 0 COMMENT '点赞数',
  `reply_count` INT DEFAULT 0 COMMENT '回复数',
  
  -- 客户端信息
  `ip_address` VARCHAR(45) COMMENT 'IP地址',
  `user_agent` VARCHAR(500) COMMENT '浏览器标识',
  
  `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `approved_at` TIMESTAMP NULL COMMENT '审批时间',
  
  KEY `idx_content` (`content_id`),
  KEY `idx_user` (`user_id`),
  KEY `idx_status` (`status`),
  KEY `idx_parent_comment` (`parent_comment_id`),
  KEY `idx_created` (`created_at`),
  CONSTRAINT `fk_comment_content` FOREIGN KEY (`content_id`) 
    REFERENCES `cms_content` (`id`) ON DELETE CASCADE,
  CONSTRAINT `fk_comment_user` FOREIGN KEY (`user_id`) 
    REFERENCES `cms_users` (`id`) ON DELETE SET NULL,
  CONSTRAINT `fk_comment_parent` FOREIGN KEY (`parent_comment_id`) 
    REFERENCES `cms_comments` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='CMS评论表';

-- =============================================
-- 8. CMS 权限表
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
-- 9. CMS 管理员角色映射表
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
  KEY `idx_role_name` (`role_name`),
  CONSTRAINT `fk_cms_admin_role_admin` FOREIGN KEY (`admin_id`) 
    REFERENCES `admin_users` (`id`) ON DELETE CASCADE,
  CONSTRAINT `fk_cms_admin_role_assigned_by` FOREIGN KEY (`assigned_by`) 
    REFERENCES `admin_users` (`id`) ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='CMS-系统用户角色映射表';

-- =============================================
-- 10. 插件启用状态表
-- =============================================
CREATE TABLE IF NOT EXISTS `cms_plugin_status` (
  `id` BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '记录ID',
  `plugin_name` VARCHAR(100) UNIQUE NOT NULL COMMENT '插件名称(cms)',
  `enabled` TINYINT DEFAULT 0 COMMENT '0:禁用 1:启用',
  `version` VARCHAR(20) COMMENT '插件版本',
  `config` JSON COMMENT '插件配置(JSON格式)',
  `installed_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '安装时间',
  `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  
  KEY `idx_name` (`plugin_name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='插件启用状态表';

-- =============================================
-- 11. CMS 操作审计日志表
-- =============================================
CREATE TABLE IF NOT EXISTS `cms_audit_logs` (
  `id` BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '日志ID',
  `admin_id` BIGINT NOT NULL COMMENT '操作管理员ID',
  `action` VARCHAR(50) COMMENT '操作类型(create_content, edit_content等)',
  `resource_type` VARCHAR(50) COMMENT '资源类型(content, category, comment等)',
  `resource_id` BIGINT COMMENT '资源ID',
  `old_value` JSON COMMENT '修改前的值',
  `new_value` JSON COMMENT '修改后的值',
  `description` VARCHAR(500) COMMENT '操作描述',
  `ip_address` VARCHAR(45) COMMENT 'IP地址',
  `user_agent` VARCHAR(500) COMMENT '浏览器标识',
  `status` TINYINT DEFAULT 1 COMMENT '1:成功 0:失败',
  `error_msg` VARCHAR(500) COMMENT '错误信息',
  `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  
  KEY `idx_admin` (`admin_id`),
  KEY `idx_action` (`action`),
  KEY `idx_resource_type` (`resource_type`),
  KEY `idx_created` (`created_at`),
  CONSTRAINT `fk_audit_admin` FOREIGN KEY (`admin_id`) 
    REFERENCES `admin_users` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='CMS操作审计日志表';

-- =============================================
-- 12. CMS 用户喜欢表 (点赞)
-- =============================================
CREATE TABLE IF NOT EXISTS `cms_likes` (
  `id` BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '点赞ID',
  `user_id` BIGINT NOT NULL COMMENT '用户ID(cms_users)',
  `likeable_type` VARCHAR(50) NOT NULL COMMENT '被点赞对象类型(content, comment)',
  `likeable_id` BIGINT NOT NULL COMMENT '被点赞对象ID',
  `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  
  UNIQUE KEY `uk_user_like` (`user_id`, `likeable_type`, `likeable_id`),
  KEY `idx_likeable` (`likeable_type`, `likeable_id`),
  CONSTRAINT `fk_like_user` FOREIGN KEY (`user_id`) 
    REFERENCES `cms_users` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='CMS用户点赞表';

-- =============================================
-- 13. CMS 草稿箱表
-- =============================================
CREATE TABLE IF NOT EXISTS `cms_drafts` (
  `id` BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '草稿ID',
  `content_id` BIGINT COMMENT '对应的正式内容ID(如果有)',
  `author_id` BIGINT NOT NULL COMMENT '作者ID',
  `title` VARCHAR(255) COMMENT '标题',
  `content` LONGTEXT COMMENT '内容',
  `category_id` BIGINT COMMENT '分类ID',
  `metadata` JSON COMMENT '其他元数据',
  `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `auto_saved_at` TIMESTAMP NULL COMMENT '自动保存时间',
  
  KEY `idx_author_id` (`author_id`),
  KEY `idx_content_id` (`content_id`),
  KEY `idx_created` (`created_at`),
  CONSTRAINT `fk_draft_author` FOREIGN KEY (`author_id`) 
    REFERENCES `admin_users` (`id`) ON DELETE CASCADE,
  CONSTRAINT `fk_draft_content` FOREIGN KEY (`content_id`) 
    REFERENCES `cms_content` (`id`) ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='CMS草稿箱表';

-- =============================================
-- 14. 初始化权限数据
-- =============================================
INSERT IGNORE INTO `cms_permissions` (`name`, `description`, `resource`, `action`) VALUES
('查看CMS内容', '查看CMS内容列表和详情', 'cms_content', 'read'),
('创建CMS内容', '创建新的CMS内容', 'cms_content', 'create'),
('编辑CMS内容', '编辑已存在的CMS内容', 'cms_content', 'update'),
('删除CMS内容', '删除CMS内容', 'cms_content', 'delete'),
('发布CMS内容', '发布CMS内容', 'cms_content', 'publish'),
('管理CMS分类', '创建/编辑/删除CMS分类', 'cms_category', 'manage'),
('管理CMS标签', '创建/编辑/删除CMS标签', 'cms_tag', 'manage'),
('管理CMS用户', '管理CMS前台用户', 'cms_users', 'manage'),
('审核CMS评论', '审核/删除CMS评论', 'cms_comments', 'manage'),
('查看CMS统计', '查看CMS统计数据', 'cms_stats', 'read'),
('管理CMS权限', '管理CMS角色和权限', 'cms_permissions', 'manage');

-- =============================================
-- 15. 初始化插件状态
-- =============================================
INSERT IGNORE INTO `cms_plugin_status` (`plugin_name`, `enabled`, `version`, `config`) VALUES
('cms', 1, '1.0.0', '{"enable_comments": true, "enable_ratings": false, "posts_per_page": 10, "enable_drafts": true}');

-- =============================================
-- 创建必要的索引和优化
-- =============================================

-- 创建全文索引用于搜索
ALTER TABLE `cms_content` ADD FULLTEXT INDEX `ft_title_content` (`title`, `description`);

-- 为关键查询创建复合索引
CREATE INDEX `idx_category_status` ON `cms_content` (`category_id`, `status`);
CREATE INDEX `idx_author_status` ON `cms_content` (`author_id`, `status`);
CREATE INDEX `idx_status_published` ON `cms_content` (`status`, `published_at`);

-- =============================================
-- 版本和备注
-- =============================================
-- V1.0.0 - 初始版本
--   - 包含完整的内容管理、分类、标签、评论系统
--   - 支持权限管理和审计日志
--   - 支持版本控制和草稿箱
--   - 优化了索引性能

