package models

import "time"

// User 用户表
type User struct {
	ID          int64      `gorm:"primaryKey;comment:用户ID" json:"id"`
	Username    string     `gorm:"uniqueIndex;size:50;comment:用户名" json:"username"`
	Phone       string     `gorm:"uniqueIndex;size:20;comment:手机号" json:"phone"`
	Email       string     `gorm:"uniqueIndex;size:100;comment:邮箱" json:"email"`
	Password    string     `gorm:"size:255;comment:密码（加密存储）" json:"-"`
	Nickname    string     `gorm:"size:100;comment:昵称" json:"nickname"`
	Avatar      string     `gorm:"size:255;comment:头像URL" json:"avatar"`
	Gender      int        `gorm:"comment:性别 1:男 2:女 0:未知" json:"gender"`
	Status      int        `gorm:"default:1;comment:状态 1:激活 0:禁用" json:"status"`
	LastLoginAt *time.Time `gorm:"comment:最后登录时间" json:"lastLoginAt"`
	CreatedAt   time.Time  `gorm:"autoCreateTime;comment:创建时间" json:"createdAt"`
	UpdatedAt   time.Time  `gorm:"autoUpdateTime;comment:更新时间" json:"updatedAt"`
	DeletedAt   *time.Time `gorm:"comment:删除时间" json:"deletedAt"`

	Roles []*Role `gorm:"many2many:user_roles;" json:"roles,omitempty"`
}

// TableName 指定表名
func (User) TableName() string {
	return "users"
}

// Role 角色表
type Role struct {
	ID          int64      `gorm:"primaryKey;comment:角色ID" json:"id"`
	Name        string     `gorm:"uniqueIndex;size:50;comment:角色名称" json:"name"`
	Description string     `gorm:"size:255;comment:角色描述" json:"description"`
	Status      int        `gorm:"default:1;comment:状态 1:启用 0:禁用" json:"status"`
	Remark      string     `gorm:"size:255;comment:备注" json:"remark"`
	CreatedAt   time.Time  `gorm:"autoCreateTime;comment:创建时间" json:"createdAt"`
	UpdatedAt   time.Time  `gorm:"autoUpdateTime;comment:更新时间" json:"updatedAt"`
	DeletedAt   *time.Time `gorm:"comment:删除时间" json:"deletedAt"`

	Permissions []*Permission `gorm:"many2many:role_permissions;" json:"permissions,omitempty"`
	Menus       []*Menu       `gorm:"many2many:role_menus;" json:"menus,omitempty"`
	Users       []*User       `gorm:"many2many:user_roles;" json:"users,omitempty"`
}

// TableName 指定表名
func (Role) TableName() string {
	return "roles"
}

// Permission 权限表
type Permission struct {
	ID          int64      `gorm:"primaryKey;comment:权限ID" json:"id"`
	Name        string     `gorm:"uniqueIndex;size:100;comment:权限名称" json:"name"`
	Description string     `gorm:"size:255;comment:权限描述" json:"description"`
	Resource    string     `gorm:"size:100;comment:资源 如 users:view" json:"resource"`
	Action      string     `gorm:"size:50;comment:操作 如 view, create, update, delete" json:"action"`
	Status      int        `gorm:"default:1;comment:状态 1:启用 0:禁用" json:"status"`
	CreatedAt   time.Time  `gorm:"autoCreateTime;comment:创建时间" json:"createdAt"`
	UpdatedAt   time.Time  `gorm:"autoUpdateTime;comment:更新时间" json:"updatedAt"`
	DeletedAt   *time.Time `gorm:"comment:删除时间" json:"deletedAt"`

	Roles []*Role `gorm:"many2many:role_permissions;" json:"roles,omitempty"`
}

// TableName 指定表名
func (Permission) TableName() string {
	return "permissions"
}

// Menu 菜单表
type Menu struct {
	ID           int64      `gorm:"primaryKey;comment:菜单ID" json:"id"`
	ParentID     int64      `gorm:"default:0;comment:父菜单ID" json:"parentId"`
	MenuName     string     `gorm:"size:100;comment:菜单名称" json:"menuName"`
	MenuPath     string     `gorm:"size:200;comment:菜单路径" json:"menuPath"`
	Component    string     `gorm:"size:255;comment:组件路径" json:"component"`
	Icon         string     `gorm:"size:100;comment:菜单图标" json:"icon"`
	Sort         int        `gorm:"default:0;comment:排序号" json:"sort"`
	Status       int        `gorm:"default:1;comment:状态 1:显示 0:隐藏" json:"status"`
	MenuType     int        `gorm:"default:1;comment:菜单类型 1:菜单 2:按钮" json:"menuType"`
	PermissionID *int64     `gorm:"comment:关联权限ID" json:"permissionId"`
	Remark       string     `gorm:"size:255;comment:备注" json:"remark"`
	CreatedAt    time.Time  `gorm:"autoCreateTime;comment:创建时间" json:"createdAt"`
	UpdatedAt    time.Time  `gorm:"autoUpdateTime;comment:更新时间" json:"updatedAt"`
	DeletedAt    *time.Time `gorm:"comment:删除时间" json:"deletedAt"`

	Children   []*Menu     `gorm:"-" json:"children,omitempty"`
	Permission *Permission `gorm:"foreignKey:PermissionID" json:"permission,omitempty"`
	Roles      []*Role     `gorm:"many2many:role_menus;" json:"roles,omitempty"`
}

// TableName 指定表名
func (Menu) TableName() string {
	return "menus"
}

// Dictionary 字典表
type Dictionary struct {
	ID          int64      `gorm:"primaryKey;comment:字典ID" json:"id"`
	DictType    string     `gorm:"size:100;index;comment:字典类型" json:"dictType"`
	DictLabel   string     `gorm:"size:100;comment:字典标签" json:"dictLabel"`
	DictValue   string     `gorm:"size:255;comment:字典值" json:"dictValue"`
	Description string     `gorm:"size:255;comment:字典描述" json:"description"`
	Sort        int        `gorm:"default:0;comment:排序号" json:"sort"`
	Status      int        `gorm:"default:1;comment:状态 1:启用 0:禁用" json:"status"`
	Remark      string     `gorm:"size:255;comment:备注" json:"remark"`
	CreatedAt   time.Time  `gorm:"autoCreateTime;comment:创建时间" json:"createdAt"`
	UpdatedAt   time.Time  `gorm:"autoUpdateTime;comment:更新时间" json:"updatedAt"`
	DeletedAt   *time.Time `gorm:"comment:删除时间" json:"deletedAt"`
}

// TableName 指定表名
func (Dictionary) TableName() string {
	return "dictionaries"
}

// API API管理表
type API struct {
	ID           int64      `gorm:"primaryKey;comment:API ID" json:"id"`
	APIName      string     `gorm:"size:100;comment:API名称" json:"apiName"`
	APIPath      string     `gorm:"size:255;comment:API路径" json:"apiPath"`
	APIMethod    string     `gorm:"size:10;comment:HTTP方法 GET, POST, PUT, DELETE等" json:"apiMethod"`
	Description  string     `gorm:"size:255;comment:API描述" json:"description"`
	Group        string     `gorm:"size:100;comment:API分组" json:"group"`
	PermissionID *int64     `gorm:"comment:关联权限ID" json:"permissionId"`
	Status       int        `gorm:"default:1;comment:状态 1:启用 0:禁用" json:"status"`
	Remark       string     `gorm:"size:255;comment:备注" json:"remark"`
	CreatedAt    time.Time  `gorm:"autoCreateTime;comment:创建时间" json:"createdAt"`
	UpdatedAt    time.Time  `gorm:"autoUpdateTime;comment:更新时间" json:"updatedAt"`
	DeletedAt    *time.Time `gorm:"comment:删除时间" json:"deletedAt"`

	Permission *Permission `gorm:"foreignKey:PermissionID" json:"permission,omitempty"`
}

// TableName 指定表名
func (API) TableName() string {
	return "apis"
}

// Plugin 插件表
type Plugin struct {
	ID          int64      `gorm:"primaryKey;comment:插件ID" json:"id"`
	PluginName  string     `gorm:"uniqueIndex;size:100;comment:插件名称" json:"pluginName"`
	PluginKey   string     `gorm:"uniqueIndex;size:100;comment:插件标识" json:"pluginKey"`
	Description string     `gorm:"size:255;comment:插件描述" json:"description"`
	Version     string     `gorm:"size:50;comment:插件版本" json:"version"`
	Author      string     `gorm:"size:100;comment:作者" json:"author"`
	Status      int        `gorm:"default:1;comment:状态 1:启用 0:禁用" json:"status"`
	Config      string     `gorm:"type:json;comment:插件配置JSON" json:"config"`
	CreatedAt   time.Time  `gorm:"autoCreateTime;comment:创建时间" json:"createdAt"`
	UpdatedAt   time.Time  `gorm:"autoUpdateTime;comment:更新时间" json:"updatedAt"`
	DeletedAt   *time.Time `gorm:"comment:删除时间" json:"deletedAt"`
}

// TableName 指定表名
func (Plugin) TableName() string {
	return "plugins"
}

// Log 系统日志表
type Log struct {
	ID           int64     `gorm:"primaryKey;comment:日志ID" json:"id"`
	UserID       *int64    `gorm:"comment:用户ID" json:"userId"`
	Username     string    `gorm:"size:50;comment:用户名" json:"username"`
	Operation    string    `gorm:"size:100;comment:操作名称" json:"operation"`
	Method       string    `gorm:"size:10;comment:请求方法" json:"method"`
	Path         string    `gorm:"size:255;comment:请求路径" json:"path"`
	IP           string    `gorm:"size:50;comment:请求IP" json:"ip"`
	Status       int       `gorm:"comment:响应状态码" json:"status"`
	ErrorMsg     string    `gorm:"type:longtext;comment:错误信息" json:"errorMsg"`
	RequestBody  string    `gorm:"type:longtext;comment:请求体" json:"requestBody"`
	ResponseBody string    `gorm:"type:longtext;comment:响应体" json:"responseBody"`
	Duration     int64     `gorm:"comment:耗时(ms)" json:"duration"`
	CreatedAt    time.Time `gorm:"autoCreateTime;comment:创建时间" json:"createdAt"`
}

// TableName 指定表名
func (Log) TableName() string {
	return "logs"
}

// App 应用市场应用表
type App struct {
	ID          int64      `gorm:"primaryKey;comment:应用ID" json:"id"`
	AppName     string     `gorm:"size:100;comment:应用名称" json:"appName"`
	AppKey      string     `gorm:"uniqueIndex;size:100;comment:应用标识" json:"appKey"`
	Version     string     `gorm:"size:50;comment:应用版本" json:"version"`
	Author      string     `gorm:"size:100;comment:作者" json:"author"`
	Description string     `gorm:"type:text;comment:应用描述" json:"description"`
	Icon        string     `gorm:"size:255;comment:应用图标" json:"icon"`
	DownloadUrl string     `gorm:"size:255;comment:下载地址" json:"downloadUrl"`
	DemoUrl     string     `gorm:"size:255;comment:演示地址" json:"demoUrl"`
	Category    string     `gorm:"size:50;comment:应用分类" json:"category"`
	Tags        string     `gorm:"size:255;comment:应用标签" json:"tags"`
	Rating      float64    `gorm:"comment:应用评分" json:"rating"`
	Downloads   int64      `gorm:"comment:下载次数" json:"downloads"`
	Status      int        `gorm:"comment:应用状态" json:"status"`
	Published   int        `gorm:"comment:发布状态" json:"published"`
	CreatedAt   time.Time  `gorm:"autoCreateTime;comment:创建时间" json:"createdAt"`
	UpdatedAt   time.Time  `gorm:"autoUpdateTime;comment:更新时间" json:"updatedAt"`
	DeletedAt   *time.Time `gorm:"comment:删除时间" json:"deletedAt"`
}

// TableName 指定表名
func (App) TableName() string {
	return "apps"
}

// Review 应用评价表
type Review struct {
	ID        int64      `gorm:"primaryKey;comment:评价ID" json:"id"`
	AppID     int64      `gorm:"index;comment:应用ID" json:"appId"`
	UserID    int64      `gorm:"comment:用户ID" json:"userId"`
	Rating    int        `gorm:"comment:评分" json:"rating"`
	Comment   string     `gorm:"type:text;comment:评价内容" json:"comment"`
	CreatedAt time.Time  `gorm:"autoCreateTime;comment:创建时间" json:"createdAt"`
	UpdatedAt time.Time  `gorm:"autoUpdateTime;comment:更新时间" json:"updatedAt"`
	DeletedAt *time.Time `gorm:"comment:删除时间" json:"deletedAt"`
}

// TableName 指定表名
func (Review) TableName() string {
	return "reviews"
}

// AppInstallation 应用安装记录表
type AppInstallation struct {
	ID        int64     `gorm:"primaryKey;comment:安装记录ID" json:"id"`
	AppKey    string    `gorm:"uniqueIndex;size:100;comment:应用标识" json:"appKey"`
	AppID     int64     `gorm:"index;comment:应用ID" json:"appId"`
	AppName   string    `gorm:"size:100;comment:应用名称" json:"appName"`
	Version   string    `gorm:"size:50;comment:安装版本" json:"version"`
	Status    int       `gorm:"comment:安装状态(1:已安装,0:未安装)" json:"status"`
	CreatedAt time.Time `gorm:"autoCreateTime;comment:安装时间" json:"createdAt"`
	UpdatedAt time.Time `gorm:"autoUpdateTime;comment:更新时间" json:"updatedAt"`
}

// TableName 指定表名
func (AppInstallation) TableName() string {
	return "app_installations"
}

// RoleMenu 角色菜单关联表
type RoleMenu struct {
	ID        int64     `gorm:"primaryKey;autoIncrement;comment:关联ID" json:"id"`
	RoleID    int64     `gorm:"index;comment:角色ID" json:"roleId"`
	MenuID    int64     `gorm:"index;comment:菜单ID" json:"menuId"`
	CreatedAt time.Time `gorm:"autoCreateTime;comment:创建时间" json:"createdAt"`
}

// TableName 指定表名
func (RoleMenu) TableName() string {
	return "role_menus"
}

type CasbinRule struct {
	ID    int64  `gorm:"primaryKey;autoIncrement;comment:规则ID" json:"id"`
	PType string `gorm:"column:ptype; size:10;index;comment:策略类型" json:"ptype"`
	V0    string `gorm:"column:v0; size:50;index;comment:字段V0" json:"v0"`
	V1    string `gorm:"column:v1; size:50;index;comment:字段V1" json:"v1"`
	V2    string `gorm:"column:v2; size:50;comment:字段V2" json:"v2"`
	V3    string `gorm:"column:v3; size:50;comment:字段V3" json:"v3"`
	V4    string `gorm:"column:v4;	 size:50;comment:字段V4" json:"v4"`
	V5    string `gorm:"column:v5; size:50;comment:字段V5" json:"v5"`
}

// TableName 指定表名
func (CasbinRule) TableName() string {
	return "casbin_rule"
}
