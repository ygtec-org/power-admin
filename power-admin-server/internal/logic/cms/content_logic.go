package cms

import (
	"context"
	"errors"
	"power-admin-server/pkg/models"
	"power-admin-server/pkg/repository"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
)

// ContentLogic 内容管理业务逻辑
type ContentLogic struct {
	contentRepo  repository.ContentRepository
	categoryRepo repository.CategoryRepository
	logger       logx.Logger
}

// NewContentLogic 创建ContentLogic实例
func NewContentLogic(contentRepo repository.ContentRepository, categoryRepo repository.CategoryRepository) *ContentLogic {
	return &ContentLogic{
		contentRepo:  contentRepo,
		categoryRepo: categoryRepo,
		logger:       logx.WithContext(context.Background()),
	}
}

// CreateContentRequest 创建内容请求
type CreateContentRequest struct {
	Title            string     `json:"title"`
	Slug             string     `json:"slug"`
	Description      string     `json:"description"`
	Content          string     `json:"content"`
	FeaturedImage    string     `json:"featured_image"`
	FeaturedImageAlt string     `json:"featured_image_alt"`
	CategoryID       *int64     `json:"category_id"`
	AuthorID         int64      `json:"author_id"`
	Status           int8       `json:"status"` // 1:草稿 2:已发布
	Visibility       int8       `json:"visibility"`
	CommentStatus    int8       `json:"comment_status"`
	SeoTitle         string     `json:"seo_title"`
	SeoKeywords      string     `json:"seo_keywords"`
	SeoDescription   string     `json:"seo_description"`
	IsFeatured       int8       `json:"is_featured"`
	IsSticky         int8       `json:"is_sticky"`
	ScheduledAt      *time.Time `json:"scheduled_at"`
}

// CreateContent 创建内容
func (l *ContentLogic) CreateContent(ctx context.Context, req *CreateContentRequest) (*models.CmsContent, error) {
	// 验证必填字段
	if req.Title == "" {
		return nil, errors.New("标题不能为空")
	}
	if req.Content == "" {
		return nil, errors.New("内容不能为空")
	}
	if req.AuthorID == 0 {
		return nil, errors.New("作者ID不能为空")
	}

	// 验证分类存在
	if req.CategoryID != nil {
		category, err := l.categoryRepo.Get(ctx, *req.CategoryID)
		if err != nil {
			l.logger.Errorf("获取分类失败: %v", err)
			return nil, err
		}
		if category == nil {
			return nil, errors.New("分类不存在")
		}
	}

	// 创建内容
	content := &models.CmsContent{
		Title:            req.Title,
		Slug:             req.Slug,
		Description:      req.Description,
		Content:          req.Content,
		FeaturedImage:    req.FeaturedImage,
		FeaturedImageAlt: req.FeaturedImageAlt,
		CategoryID:       req.CategoryID,
		AuthorID:         req.AuthorID,
		Status:           req.Status,
		Visibility:       req.Visibility,
		CommentStatus:    req.CommentStatus,
		SeoTitle:         req.SeoTitle,
		SeoKeywords:      req.SeoKeywords,
		SeoDescription:   req.SeoDescription,
		IsFeatured:       req.IsFeatured,
		IsSticky:         req.IsSticky,
		ScheduledAt:      req.ScheduledAt,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}

	if err := l.contentRepo.Create(ctx, content); err != nil {
		l.logger.Errorf("创建内容失败: %v", err)
		return nil, err
	}

	return content, nil
}

// UpdateContentRequest 更新内容请求
type UpdateContentRequest struct {
	ID               int64
	Title            string
	Slug             string
	Description      string
	Content          string
	FeaturedImage    string
	FeaturedImageAlt string
	CategoryID       *int64
	Status           int8
	Visibility       int8
	CommentStatus    int8
	SeoTitle         string
	SeoKeywords      string
	SeoDescription   string
	IsFeatured       int8
	IsSticky         int8
	ScheduledAt      *time.Time
}

// UpdateContent 更新内容
func (l *ContentLogic) UpdateContent(ctx context.Context, req *UpdateContentRequest) (*models.CmsContent, error) {
	// 获取原内容
	existing, err := l.contentRepo.Get(ctx, req.ID)
	if err != nil {
		l.logger.Errorf("获取内容失败: %v", err)
		return nil, err
	}
	if existing == nil {
		return nil, errors.New("内容不存在")
	}

	// 验证新分类
	if req.CategoryID != nil && (existing.CategoryID == nil || *existing.CategoryID != *req.CategoryID) {
		category, err := l.categoryRepo.Get(ctx, *req.CategoryID)
		if err != nil {
			l.logger.Errorf("获取分类失败: %v", err)
			return nil, err
		}
		if category == nil {
			return nil, errors.New("分类不存在")
		}
	}

	// 更新内容
	updateData := &models.CmsContent{
		Title:            req.Title,
		Slug:             req.Slug,
		Description:      req.Description,
		Content:          req.Content,
		FeaturedImage:    req.FeaturedImage,
		FeaturedImageAlt: req.FeaturedImageAlt,
		CategoryID:       req.CategoryID,
		Status:           req.Status,
		Visibility:       req.Visibility,
		CommentStatus:    req.CommentStatus,
		SeoTitle:         req.SeoTitle,
		SeoKeywords:      req.SeoKeywords,
		SeoDescription:   req.SeoDescription,
		IsFeatured:       req.IsFeatured,
		IsSticky:         req.IsSticky,
		ScheduledAt:      req.ScheduledAt,
		UpdatedAt:        time.Now(),
	}

	if err := l.contentRepo.Update(ctx, req.ID, updateData); err != nil {
		l.logger.Errorf("更新内容失败: %v", err)
		return nil, err
	}

	// 重新获取更新后的内容
	return l.contentRepo.Get(ctx, req.ID)
}

// DeleteContent 删除内容（软删除）
func (l *ContentLogic) DeleteContent(ctx context.Context, contentID int64) error {
	// 检查内容是否存在
	existing, err := l.contentRepo.Get(ctx, contentID)
	if err != nil {
		l.logger.Errorf("获取内容失败: %v", err)
		return err
	}
	if existing == nil {
		return errors.New("内容不存在")
	}

	if err := l.contentRepo.SoftDelete(ctx, contentID); err != nil {
		l.logger.Errorf("删除内容失败: %v", err)
		return err
	}

	return nil
}

// HardDeleteContent 永久删除内容
func (l *ContentLogic) HardDeleteContent(ctx context.Context, contentID int64) error {
	// 检查内容是否存在
	existing, err := l.contentRepo.Get(ctx, contentID)
	if err != nil {
		l.logger.Errorf("获取内容失败: %v", err)
		return err
	}
	if existing == nil {
		return errors.New("内容不存在")
	}

	if err := l.contentRepo.Delete(ctx, contentID); err != nil {
		l.logger.Errorf("永久删除内容失败: %v", err)
		return err
	}

	return nil
}

// GetContent 获取内容详情
func (l *ContentLogic) GetContent(ctx context.Context, contentID int64) (*models.CmsContent, error) {
	content, err := l.contentRepo.Get(ctx, contentID)
	if err != nil {
		l.logger.Errorf("获取内容失败: %v", err)
		return nil, err
	}
	if content == nil {
		return nil, errors.New("内容不存在")
	}

	// 增加浏览数
	_ = l.contentRepo.IncrementViewCount(ctx, contentID)

	return content, nil
}

// GetContentBySlug 根据Slug获取内容
func (l *ContentLogic) GetContentBySlug(ctx context.Context, slug string) (*models.CmsContent, error) {
	content, err := l.contentRepo.GetBySlug(ctx, slug)
	if err != nil {
		l.logger.Errorf("获取内容失败: %v", err)
		return nil, err
	}
	if content == nil {
		return nil, errors.New("内容不存在")
	}

	// 增加浏览数
	_ = l.contentRepo.IncrementViewCount(ctx, content.ID)

	return content, nil
}

// ListContentRequest 列表查询请求
type ListContentRequest struct {
	Page       int
	PageSize   int
	CategoryID *int64
	Status     *int8
	Visibility *int8
	Search     string
	SortBy     string
	SortOrder  string
	IsSticky   *int8
	IsFeatured *int8
}

// ListContent 查询内容列表
func (l *ContentLogic) ListContent(ctx context.Context, req *ListContentRequest) (*repository.PagedResult, error) {
	listReq := &repository.ContentListRequest{
		Page:       req.Page,
		PageSize:   req.PageSize,
		CategoryID: req.CategoryID,
		Status:     req.Status,
		Visibility: req.Visibility,
		Search:     req.Search,
		SortBy:     req.SortBy,
		SortOrder:  req.SortOrder,
		IsSticky:   req.IsSticky,
		IsFeatured: req.IsFeatured,
	}

	result, err := l.contentRepo.List(ctx, listReq)
	if err != nil {
		l.logger.Errorf("查询内容列表失败: %v", err)
		return nil, err
	}

	return result, nil
}

// PublishContent 发布内容
func (l *ContentLogic) PublishContent(ctx context.Context, contentID int64) error {
	// 检查内容是否存在
	existing, err := l.contentRepo.Get(ctx, contentID)
	if err != nil {
		l.logger.Errorf("获取内容失败: %v", err)
		return err
	}
	if existing == nil {
		return errors.New("内容不存在")
	}

	if err := l.contentRepo.Publish(ctx, contentID); err != nil {
		l.logger.Errorf("发布内容失败: %v", err)
		return err
	}

	return nil
}

// UnpublishContent 取消发布内容
func (l *ContentLogic) UnpublishContent(ctx context.Context, contentID int64) error {
	// 检查内容是否存在
	existing, err := l.contentRepo.Get(ctx, contentID)
	if err != nil {
		l.logger.Errorf("获取内容失败: %v", err)
		return err
	}
	if existing == nil {
		return errors.New("内容不存在")
	}

	if err := l.contentRepo.Unpublish(ctx, contentID); err != nil {
		l.logger.Errorf("取消发布内容失败: %v", err)
		return err
	}

	return nil
}

// BatchUpdateContentStatus 批量更新内容状态
func (l *ContentLogic) BatchUpdateContentStatus(ctx context.Context, ids []int64, status int8) error {
	if len(ids) == 0 {
		return errors.New("内容ID列表不能为空")
	}

	if status < 1 || status > 3 {
		return errors.New("无效的状态值")
	}

	if err := l.contentRepo.BatchUpdateStatus(ctx, ids, status); err != nil {
		l.logger.Errorf("批量更新内容状态失败: %v", err)
		return err
	}

	return nil
}
