package models

import (
	"database/sql"
	"time"
)

// =============================================
// CMS 内容模型
// =============================================
type CmsContent struct {
	ID               int64      `gorm:"primaryKey" json:"id"`
	Title            string     `gorm:"type:varchar(255);not null" json:"title"`
	Slug             string     `gorm:"type:varchar(255);unique" json:"slug"`
	Description      string     `gorm:"type:varchar(500)" json:"description"`
	Content          string     `gorm:"type:longtext;not null" json:"content"`
	FeaturedImage    string     `gorm:"type:varchar(500)" json:"featured_image"`
	FeaturedImageAlt string     `gorm:"type:varchar(255)" json:"featured_image_alt"`
	CategoryID       *int64     `json:"category_id"`
	AuthorID         int64      `gorm:"not null" json:"author_id"`
	Status           int8       `gorm:"default:1" json:"status"` // 1:草稿 2:已发布 3:已删除
	Visibility       int8       `gorm:"default:1" json:"visibility"`
	CommentStatus    int8       `gorm:"default:1" json:"comment_status"`
	ViewCount        int        `gorm:"default:0" json:"view_count"`
	CommentCount     int        `gorm:"default:0" json:"comment_count"`
	LikeCount        int        `gorm:"default:0" json:"like_count"`
	SeoTitle         string     `gorm:"type:varchar(255)" json:"seo_title"`
	SeoKeywords      string     `gorm:"type:varchar(255)" json:"seo_keywords"`
	SeoDescription   string     `gorm:"type:varchar(500)" json:"seo_description"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"`
	PublishedAt      *time.Time `json:"published_at"`
	ScheduledAt      *time.Time `json:"scheduled_at"`
	IsFeatured       int8       `gorm:"default:0" json:"is_featured"`
	IsSticky         int8       `gorm:"default:0" json:"is_sticky"`
	RevisionCount    int        `gorm:"default:0" json:"revision_count"`
}

func (CmsContent) TableName() string {
	return "cms_content"
}

// =============================================
// CMS 分类模型
// =============================================
type CmsCategory struct {
	ID             int64          `gorm:"primaryKey" json:"id"`
	Name           string         `gorm:"type:varchar(100);not null" json:"name"`
	Slug           string         `gorm:"type:varchar(100);unique" json:"slug"`
	Description    string         `gorm:"type:text" json:"description"`
	Thumbnail      string         `gorm:"type:varchar(500)" json:"thumbnail"`
	ParentID       *int64         `json:"parent_id"`
	Sort           int            `gorm:"default:0" json:"sort"`
	Status         int8           `gorm:"default:1" json:"status"`
	ContentCount   int            `gorm:"default:0" json:"content_count"`
	SeoKeywords    string         `gorm:"type:varchar(255)" json:"seo_keywords"`
	SeoDescription string         `gorm:"type:varchar(500)" json:"seo_description"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	Children       []*CmsCategory `gorm:"-" json:"children,omitempty"`
}

func (CmsCategory) TableName() string {
	return "cms_category"
}

// =============================================
// CMS 标签模型
// =============================================
type CmsTag struct {
	ID          int64     `gorm:"primaryKey" json:"id"`
	Name        string    `gorm:"type:varchar(50);not null;unique" json:"name"`
	Slug        string    `gorm:"type:varchar(50);unique" json:"slug"`
	Description string    `gorm:"type:text" json:"description"`
	Color       string    `gorm:"type:varchar(7)" json:"color"`
	UsageCount  int       `gorm:"default:0" json:"usage_count"`
	Status      int8      `gorm:"default:1" json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (CmsTag) TableName() string {
	return "cms_tag"
}

// =============================================
// CMS 内容-标签关联模型
// =============================================
type CmsContentTag struct {
	ContentID int64     `gorm:"primaryKey;column:content_id" json:"content_id"`
	TagID     int64     `gorm:"primaryKey;column:tag_id" json:"tag_id"`
	CreatedAt time.Time `json:"created_at"`
}

func (CmsContentTag) TableName() string {
	return "cms_content_tag"
}

// =============================================
// CMS 内容版本历史模型
// =============================================
type CmsContentRevision struct {
	ID             int64     `gorm:"primaryKey" json:"id"`
	ContentID      int64     `gorm:"not null" json:"content_id"`
	Title          string    `gorm:"type:varchar(255);not null" json:"title"`
	Content        string    `gorm:"type:longtext;not null" json:"content"`
	RevisionNumber int       `json:"revision_number"`
	AuthorID       int64     `gorm:"not null" json:"author_id"`
	ChangeSummary  string    `gorm:"type:varchar(500)" json:"change_summary"`
	CreatedAt      time.Time `json:"created_at"`
}

func (CmsContentRevision) TableName() string {
	return "cms_content_revision"
}

// =============================================
// CMS 前台用户模型
// =============================================
type CmsUser struct {
	ID              int64      `gorm:"primaryKey" json:"id"`
	Username        string     `gorm:"type:varchar(100);not null;unique" json:"username"`
	Email           string     `gorm:"type:varchar(255);unique" json:"email"`
	Phone           string     `gorm:"type:varchar(20)" json:"phone"`
	Password        string     `gorm:"type:varchar(255);not null" json:"-"` // 不序列化密码
	Nickname        string     `gorm:"type:varchar(100)" json:"nickname"`
	Avatar          string     `gorm:"type:varchar(500)" json:"avatar"`
	Bio             string     `gorm:"type:text" json:"bio"`
	Gender          int8       `json:"gender"`
	Status          int8       `gorm:"default:1" json:"status"`
	EmailVerified   int8       `gorm:"default:0" json:"email_verified"`
	PhoneVerified   int8       `gorm:"default:0" json:"phone_verified"`
	EmailVerifiedAt *time.Time `json:"email_verified_at"`
	PhoneVerifiedAt *time.Time `json:"phone_verified_at"`
	LastLoginAt     *time.Time `json:"last_login_at"`
	LastLoginIP     string     `gorm:"type:varchar(45)" json:"last_login_ip"`
	LoginCount      int        `gorm:"default:0" json:"login_count"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
	DeletedAt       *time.Time `sql:"index" json:"deleted_at"`
}

func (CmsUser) TableName() string {
	return "cms_users"
}

// =============================================
// CMS 评论模型
// =============================================
type CmsComment struct {
	ID              int64      `gorm:"primaryKey" json:"id"`
	ContentID       int64      `gorm:"not null" json:"content_id"`
	UserID          *int64     `json:"user_id"`
	ParentCommentID *int64     `json:"parent_comment_id"`
	AuthorName      string     `gorm:"type:varchar(100)" json:"author_name"`
	AuthorEmail     string     `gorm:"type:varchar(255)" json:"author_email"`
	Content         string     `gorm:"type:text;not null" json:"content"`
	Status          int8       `gorm:"default:0" json:"status"`
	LikeCount       int        `gorm:"default:0" json:"like_count"`
	ReplyCount      int        `gorm:"default:0" json:"reply_count"`
	IPAddress       string     `gorm:"type:varchar(45)" json:"ip_address"`
	UserAgent       string     `gorm:"type:varchar(500)" json:"user_agent"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
	ApprovedAt      *time.Time `json:"approved_at"`
}

func (CmsComment) TableName() string {
	return "cms_comments"
}

// =============================================
// CMS 权限模型
// =============================================
type CmsPermission struct {
	ID          int64     `gorm:"primaryKey" json:"id"`
	Name        string    `gorm:"type:varchar(100);not null;unique" json:"name"`
	Description string    `gorm:"type:text" json:"description"`
	Resource    string    `gorm:"type:varchar(100)" json:"resource"`
	Action      string    `gorm:"type:varchar(100)" json:"action"`
	CreatedAt   time.Time `json:"created_at"`
}

func (CmsPermission) TableName() string {
	return "cms_permissions"
}

// =============================================
// CMS 管理员角色映射模型
// =============================================
type CmsAdminRole struct {
	ID         int64     `gorm:"primaryKey" json:"id"`
	AdminID    int64     `gorm:"not null" json:"admin_id"`
	RoleName   string    `gorm:"type:varchar(50);not null" json:"role_name"`
	AssignedBy *int64    `json:"assigned_by"`
	AssignedAt time.Time `json:"assigned_at"`
	CreatedAt  time.Time `json:"created_at"`
}

func (CmsAdminRole) TableName() string {
	return "cms_admin_roles"
}

// =============================================
// CMS 插件状态模型
// =============================================
type CmsPluginStatus struct {
	ID          int64         `gorm:"primaryKey" json:"id"`
	PluginName  string        `gorm:"type:varchar(100);not null;unique" json:"plugin_name"`
	Enabled     int8          `gorm:"default:0" json:"enabled"`
	Version     string        `gorm:"type:varchar(20)" json:"version"`
	Config      *sql.RawBytes `gorm:"type:json" json:"config"`
	InstalledAt time.Time     `json:"installed_at"`
	UpdatedAt   time.Time     `json:"updated_at"`
}

func (CmsPluginStatus) TableName() string {
	return "cms_plugin_status"
}

// =============================================
// CMS 审计日志模型
// =============================================
type CmsAuditLog struct {
	ID           int64         `gorm:"primaryKey" json:"id"`
	AdminID      int64         `gorm:"not null" json:"admin_id"`
	Action       string        `gorm:"type:varchar(50)" json:"action"`
	ResourceType string        `gorm:"type:varchar(50)" json:"resource_type"`
	ResourceID   *int64        `json:"resource_id"`
	OldValue     *sql.RawBytes `gorm:"type:json" json:"old_value"`
	NewValue     *sql.RawBytes `gorm:"type:json" json:"new_value"`
	Description  string        `gorm:"type:varchar(500)" json:"description"`
	IPAddress    string        `gorm:"type:varchar(45)" json:"ip_address"`
	UserAgent    string        `gorm:"type:varchar(500)" json:"user_agent"`
	Status       int8          `gorm:"default:1" json:"status"`
	ErrorMsg     string        `gorm:"type:varchar(500)" json:"error_msg"`
	CreatedAt    time.Time     `json:"created_at"`
}

func (CmsAuditLog) TableName() string {
	return "cms_audit_logs"
}

// =============================================
// CMS 点赞模型
// =============================================
type CmsLike struct {
	ID           int64     `gorm:"primaryKey" json:"id"`
	UserID       int64     `gorm:"not null" json:"user_id"`
	LikeableType string    `gorm:"type:varchar(50);not null" json:"likeable_type"`
	LikeableID   int64     `gorm:"not null" json:"likeable_id"`
	CreatedAt    time.Time `json:"created_at"`
}

func (CmsLike) TableName() string {
	return "cms_likes"
}

// =============================================
// CMS 草稿模型
// =============================================
type CmsDraft struct {
	ID          int64         `gorm:"primaryKey" json:"id"`
	ContentID   *int64        `json:"content_id"`
	AuthorID    int64         `gorm:"not null" json:"author_id"`
	Title       string        `gorm:"type:varchar(255)" json:"title"`
	Content     string        `gorm:"type:longtext" json:"content"`
	CategoryID  *int64        `json:"category_id"`
	Metadata    *sql.RawBytes `gorm:"type:json" json:"metadata"`
	CreatedAt   time.Time     `json:"created_at"`
	UpdatedAt   time.Time     `json:"updated_at"`
	AutoSavedAt *time.Time    `json:"auto_saved_at"`
}

func (CmsDraft) TableName() string {
	return "cms_drafts"
}
