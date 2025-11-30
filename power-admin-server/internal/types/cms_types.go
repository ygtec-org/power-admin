package types

// =====================================================
// 内容管理 (Content)
// =====================================================

type ContentListReq struct {
	Page       int    `json:"page"`
	PageSize   int    `json:"pageSize"`
	CategoryId int64  `json:"categoryId,optional"`
	Status     int    `json:"status,optional"`
	Search     string `json:"search,optional"`
	SortBy     string `json:"sortBy,optional"`
	SortOrder  string `json:"sortOrder,optional"`
}

type ContentInfo struct {
	Id               int64  `json:"id"`
	Title            string `json:"title"`
	Slug             string `json:"slug"`
	Description      string `json:"description"`
	Content          string `json:"content"`
	FeaturedImage    string `json:"featuredImage,optional"`
	FeaturedImageAlt string `json:"featuredImageAlt,optional"`
	CategoryId       int64  `json:"categoryId,optional"`
	AuthorId         int64  `json:"authorId"`
	Status           int    `json:"status"`
	Visibility       int    `json:"visibility"`
	CommentStatus    int    `json:"commentStatus"`
	ViewCount        int    `json:"viewCount"`
	CommentCount     int    `json:"commentCount"`
	LikeCount        int    `json:"likeCount"`
	SeoTitle         string `json:"seoTitle,optional"`
	SeoKeywords      string `json:"seoKeywords,optional"`
	SeoDescription   string `json:"seoDescription,optional"`
	PublishedAt      string `json:"publishedAt,optional"`
	ScheduledAt      string `json:"scheduledAt,optional"`
	IsFeatured       int    `json:"isFeatured"`
	IsSticky         int    `json:"isSticky"`
	CreatedAt        string `json:"createdAt"`
	UpdatedAt        string `json:"updatedAt"`
}

type ContentListResp struct {
	Total int64         `json:"total"`
	Page  int64         `json:"page"`
	Size  int64         `json:"pageSize"`
	Data  []ContentInfo `json:"list"`
}

type CreateContentReq struct {
	Title            string `json:"title"`
	Slug             string `json:"slug,optional"`
	Description      string `json:"description,optional"`
	Content          string `json:"content"`
	FeaturedImage    string `json:"featuredImage,optional"`
	FeaturedImageAlt string `json:"featuredImageAlt,optional"`
	CategoryId       int64  `json:"categoryId,optional"`
	AuthorId         int64  `json:"authorId"`
	Status           int    `json:"status"`
	Visibility       int    `json:"visibility"`
	CommentStatus    int    `json:"commentStatus"`
	SeoTitle         string `json:"seoTitle,optional"`
	SeoKeywords      string `json:"seoKeywords,optional"`
	SeoDescription   string `json:"seoDescription,optional"`
	IsFeatured       int    `json:"isFeatured"`
	IsSticky         int    `json:"isSticky"`
	ScheduledAt      string `json:"scheduledAt,optional"`
}

type UpdateContentReq struct {
	Id               int64  `json:"id"`
	Title            string `json:"title,optional"`
	Slug             string `json:"slug,optional"`
	Description      string `json:"description,optional"`
	Content          string `json:"content,optional"`
	FeaturedImage    string `json:"featuredImage,optional"`
	FeaturedImageAlt string `json:"featuredImageAlt,optional"`
	CategoryId       int64  `json:"categoryId,optional"`
	Status           int    `json:"status,optional"`
	Visibility       int    `json:"visibility,optional"`
	CommentStatus    int    `json:"commentStatus,optional"`
	SeoTitle         string `json:"seoTitle,optional"`
	SeoKeywords      string `json:"seoKeywords,optional"`
	SeoDescription   string `json:"seoDescription,optional"`
	IsFeatured       int    `json:"isFeatured,optional"`
	IsSticky         int    `json:"isSticky,optional"`
	ScheduledAt      string `json:"scheduledAt,optional"`
}

type GetContentReq struct {
	Id int64 `json:"id"`
}

type DeleteContentReq struct {
	Id int64 `json:"id"`
}

type PublishContentReq struct {
	Id int64 `json:"id"`
}

type UnpublishContentReq struct {
	Id int64 `json:"id"`
}

type BatchUpdateStatusReq struct {
	Ids    []int64 `json:"ids"`
	Status int     `json:"status"`
}

// =====================================================
// 分类管理 (Category)
// =====================================================

type CategoryListReq struct {
	Page     int   `json:"page"`
	PageSize int   `json:"pageSize"`
	ParentId int64 `json:"parentId,optional"`
}

type CategoryInfo struct {
	Id             int64          `json:"id"`
	Name           string         `json:"name"`
	Slug           string         `json:"slug"`
	Description    string         `json:"description,optional"`
	Thumbnail      string         `json:"thumbnail,optional"`
	ParentId       int64          `json:"parentId,optional"`
	Sort           int            `json:"sort"`
	Status         int            `json:"status"`
	ContentCount   int            `json:"contentCount"`
	SeoKeywords    string         `json:"seoKeywords,optional"`
	SeoDescription string         `json:"seoDescription,optional"`
	CreatedAt      string         `json:"createdAt"`
	UpdatedAt      string         `json:"updatedAt"`
	Children       []CategoryInfo `json:"children,optional"`
}

type CategoryListResp struct {
	Total int64          `json:"total"`
	Data  []CategoryInfo `json:"list"`
}

type CategoryTreeReq struct {
}

type CategoryTreeResp struct {
	Data []CategoryInfo `json:"data"`
}

type CreateCategoryReq struct {
	Name           string `json:"name"`
	Slug           string `json:"slug,optional"`
	Description    string `json:"description,optional"`
	Thumbnail      string `json:"thumbnail,optional"`
	ParentId       *int64 `json:"parentId,optional"`
	Sort           int    `json:"sort"`
	Status         int    `json:"status"`
	SeoKeywords    string `json:"seoKeywords,optional"`
	SeoDescription string `json:"seoDescription,optional"`
}

type UpdateCategoryReq struct {
	Id             int64  `json:"id"`
	Name           string `json:"name,optional"`
	Slug           string `json:"slug,optional"`
	Description    string `json:"description,optional"`
	Thumbnail      string `json:"thumbnail,optional"`
	ParentId       *int64 `json:"parentId,optional"`
	Sort           int    `json:"sort,optional"`
	Status         int    `json:"status,optional"`
	SeoKeywords    string `json:"seoKeywords,optional"`
	SeoDescription string `json:"seoDescription,optional"`
}

type GetCategoryReq struct {
	Id int64 `json:"id"`
}

type DeleteCategoryReq struct {
	Id int64 `json:"id"`
}

// =====================================================
// 标签管理 (Tag)
// =====================================================

type TagListReq struct {
	Page     int `json:"page"`
	PageSize int `json:"pageSize"`
}

type TagInfo struct {
	Id          int64  `json:"id"`
	Name        string `json:"name"`
	Slug        string `json:"slug"`
	Description string `json:"description,optional"`
	Color       string `json:"color,optional"`
	UsageCount  int    `json:"usageCount"`
	Status      int    `json:"status"`
	CreatedAt   string `json:"createdAt"`
	UpdatedAt   string `json:"updatedAt"`
}

type TagListResp struct {
	Total int64     `json:"total"`
	Data  []TagInfo `json:"list"`
}

type CreateTagReq struct {
	Name        string `json:"name"`
	Slug        string `json:"slug,optional"`
	Description string `json:"description,optional"`
	Color       string `json:"color,optional"`
	Status      int    `json:"status"`
}

type UpdateTagReq struct {
	Id          int64  `json:"id"`
	Name        string `json:"name,optional"`
	Slug        string `json:"slug,optional"`
	Description string `json:"description,optional"`
	Color       string `json:"color,optional"`
	Status      int    `json:"status,optional"`
}

type GetTagReq struct {
	Id int64 `json:"id"`
}

type DeleteTagReq struct {
	Id int64 `json:"id"`
}

// =====================================================
// 评论管理 (Comment)
// =====================================================

type CommentListReq struct {
	Page      int   `json:"page"`
	PageSize  int   `json:"pageSize"`
	ContentId int64 `json:"contentId"`
}

type CommentInfo struct {
	Id              int64  `json:"id"`
	ContentId       int64  `json:"contentId"`
	UserId          int64  `json:"userId,optional"`
	ParentCommentId int64  `json:"parentCommentId,optional"`
	AuthorName      string `json:"authorName"`
	AuthorEmail     string `json:"authorEmail,optional"`
	Content         string `json:"content"`
	Status          int    `json:"status"`
	LikeCount       int    `json:"likeCount"`
	ReplyCount      int    `json:"replyCount"`
	IpAddress       string `json:"ipAddress,optional"`
	UserAgent       string `json:"userAgent,optional"`
	CreatedAt       string `json:"createdAt"`
	UpdatedAt       string `json:"updatedAt"`
	ApprovedAt      string `json:"approvedAt,optional"`
}

type CommentListResp struct {
	Total int64         `json:"total"`
	Data  []CommentInfo `json:"list"`
}

type CreateCommentReq struct {
	ContentId       int64  `json:"contentId"`
	UserId          int64  `json:"userId,optional"`
	ParentCommentId int64  `json:"parentCommentId,optional"`
	AuthorName      string `json:"authorName"`
	AuthorEmail     string `json:"authorEmail,optional"`
	Content         string `json:"content"`
	IpAddress       string `json:"ipAddress,optional"`
	UserAgent       string `json:"userAgent,optional"`
}

type UpdateCommentReq struct {
	Id      int64  `json:"id"`
	Content string `json:"content,optional"`
	Status  int    `json:"status,optional"`
}

type GetCommentReq struct {
	Id int64 `json:"id"`
}

type DeleteCommentReq struct {
	Id int64 `json:"id"`
}

type ApproveCommentReq struct {
	Id int64 `json:"id"`
}

type RejectCommentReq struct {
	Id int64 `json:"id"`
}

type LikeCommentReq struct {
	Id int64 `json:"id"`
}

// =====================================================
// CMS用户管理 (CMS User)
// =====================================================

type CmsUserListReq struct {
	Page     int `json:"page"`
	PageSize int `json:"pageSize"`
}

type CmsUserInfo struct {
	Id            int64  `json:"id"`
	Username      string `json:"username"`
	Email         string `json:"email"`
	Phone         string `json:"phone,optional"`
	Nickname      string `json:"nickname,optional"`
	Avatar        string `json:"avatar,optional"`
	Bio           string `json:"bio,optional"`
	Gender        int    `json:"gender,optional"`
	Status        int    `json:"status"`
	EmailVerified int    `json:"emailVerified"`
	PhoneVerified int    `json:"phoneVerified"`
	LastLoginAt   string `json:"lastLoginAt,optional"`
	LastLoginIp   string `json:"lastLoginIp,optional"`
	LoginCount    int    `json:"loginCount"`
	CreatedAt     string `json:"createdAt"`
	UpdatedAt     string `json:"updatedAt"`
}

type CmsUserListResp struct {
	Total int64         `json:"total"`
	Data  []CmsUserInfo `json:"list"`
}

type CmsRegisterReq struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Nickname string `json:"nickname,optional"`
}

type CmsLoginReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Ip       string `json:"ip,optional"`
}

type CmsLoginResp struct {
	Id       int64  `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Nickname string `json:"nickname"`
	Avatar   string `json:"avatar,optional"`
}

type UpdateCmsUserReq struct {
	Id       int64  `json:"id"`
	Email    string `json:"email,optional"`
	Nickname string `json:"nickname,optional"`
	Avatar   string `json:"avatar,optional"`
	Bio      string `json:"bio,optional"`
	Gender   int    `json:"gender,optional"`
	Phone    string `json:"phone,optional"`
}

type ChangePasswordReq struct {
	UserId      int64  `json:"userId"`
	OldPassword string `json:"oldPassword"`
	NewPassword string `json:"newPassword"`
}

type GetCmsUserReq struct {
	Id int64 `json:"id"`
}

type DeleteCmsUserReq struct {
	Id int64 `json:"id"`
}

type DisableCmsUserReq struct {
	Id int64 `json:"id"`
}

// =====================================================
// 发布和工作流 (Publish)
// =====================================================

type PublishInfoResp struct {
	ContentId   int64  `json:"contentId"`
	Title       string `json:"title"`
	Status      int    `json:"status"`
	PublishedAt string `json:"publishedAt,optional"`
	ScheduledAt string `json:"scheduledAt,optional"`
	IsScheduled bool   `json:"isScheduled"`
	IsDraft     bool   `json:"isDraft"`
	IsPublished bool   `json:"isPublished"`
}

type PublishImmediateReq struct {
	ContentId int64 `json:"contentId"`
}

type PublishScheduledReq struct {
	ContentId   int64  `json:"contentId"`
	ScheduledAt string `json:"scheduledAt"`
}

type UnpublishReq struct {
	ContentId int64 `json:"contentId"`
}

type CancelScheduledPublishReq struct {
	ContentId int64 `json:"contentId"`
}

type GetPublishStatusReq struct {
	ContentId int64 `json:"contentId"`
}

type BatchPublishReq struct {
	ContentIds []int64 `json:"contentIds"`
}
