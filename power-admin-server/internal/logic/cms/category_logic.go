package cms

import (
	"context"
	"errors"
	"power-admin-server/pkg/models"
	"power-admin-server/pkg/repository"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
)

// CategoryLogic 分类管理业务逻辑
type CategoryLogic struct {
	categoryRepo repository.CategoryRepository
	contentRepo  repository.ContentRepository
	logger       logx.Logger
}

// NewCategoryLogic 创建CategoryLogic实例
func NewCategoryLogic(categoryRepo repository.CategoryRepository, contentRepo repository.ContentRepository) *CategoryLogic {
	return &CategoryLogic{
		categoryRepo: categoryRepo,
		contentRepo:  contentRepo,
		logger:       logx.WithContext(context.Background()),
	}
}

// CreateCategoryRequest 创建分类请求
type CreateCategoryRequest struct {
	Name           string `json:"name"`
	Slug           string `json:"slug"`
	Description    string `json:"description"`
	Thumbnail      string `json:"thumbnail"`
	ParentID       *int64 `json:"parent_id"`
	Sort           int    `json:"sort"`
	Status         int8   `json:"status"`
	SeoKeywords    string `json:"seo_keywords"`
	SeoDescription string `json:"seo_description"`
}

// CreateCategory 创建分类
func (l *CategoryLogic) CreateCategory(ctx context.Context, req *CreateCategoryRequest) (*models.CmsCategory, error) {
	// 验证必填字段
	if req.Name == "" {
		return nil, errors.New("分类名称不能为空")
	}

	// 验证父分类存在（如果有）
	if req.ParentID != nil {
		parent, err := l.categoryRepo.Get(ctx, *req.ParentID)
		if err != nil {
			l.logger.Errorf("获取父分类失败: %v", err)
			return nil, err
		}
		if parent == nil {
			return nil, errors.New("父分类不存在")
		}
	}

	category := &models.CmsCategory{
		Name:           req.Name,
		Slug:           req.Slug,
		Description:    req.Description,
		Thumbnail:      req.Thumbnail,
		ParentID:       req.ParentID,
		Sort:           req.Sort,
		Status:         req.Status,
		SeoKeywords:    req.SeoKeywords,
		SeoDescription: req.SeoDescription,
		ContentCount:   0,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	if err := l.categoryRepo.Create(ctx, category); err != nil {
		l.logger.Errorf("创建分类失败: %v", err)
		return nil, err
	}

	return category, nil
}

// UpdateCategoryRequest 更新分类请求
type UpdateCategoryRequest struct {
	ID             int64
	Name           string
	Slug           string
	Description    string
	Thumbnail      string
	ParentID       *int64
	Sort           int
	Status         int8
	SeoKeywords    string
	SeoDescription string
}

// UpdateCategory 更新分类
func (l *CategoryLogic) UpdateCategory(ctx context.Context, req *UpdateCategoryRequest) (*models.CmsCategory, error) {
	// 获取原分类
	existing, err := l.categoryRepo.Get(ctx, req.ID)
	if err != nil {
		l.logger.Errorf("获取分类失败: %v", err)
		return nil, err
	}
	if existing == nil {
		return nil, errors.New("分类不存在")
	}

	// 验证父分类是否存在（防止自己作为自己的父分类）
	if req.ParentID != nil {
		if *req.ParentID == req.ID {
			return nil, errors.New("分类不能以自己作为父分类")
		}
		parent, err := l.categoryRepo.Get(ctx, *req.ParentID)
		if err != nil {
			l.logger.Errorf("获取父分类失败: %v", err)
			return nil, err
		}
		if parent == nil {
			return nil, errors.New("父分类不存在")
		}
	}

	updateData := &models.CmsCategory{
		Name:           req.Name,
		Slug:           req.Slug,
		Description:    req.Description,
		Thumbnail:      req.Thumbnail,
		ParentID:       req.ParentID,
		Sort:           req.Sort,
		Status:         req.Status,
		SeoKeywords:    req.SeoKeywords,
		SeoDescription: req.SeoDescription,
		UpdatedAt:      time.Now(),
	}

	if err := l.categoryRepo.Update(ctx, req.ID, updateData); err != nil {
		l.logger.Errorf("更新分类失败: %v", err)
		return nil, err
	}

	return l.categoryRepo.Get(ctx, req.ID)
}

// DeleteCategory 删除分类（包含级联处理）
func (l *CategoryLogic) DeleteCategory(ctx context.Context, categoryID int64) error {
	// 获取分类
	existing, err := l.categoryRepo.Get(ctx, categoryID)
	if err != nil {
		l.logger.Errorf("获取分类失败: %v", err)
		return err
	}
	if existing == nil {
		return errors.New("分类不存在")
	}

	// 检查是否有子分类
	children, err := l.categoryRepo.List(ctx, &categoryID)
	if err != nil {
		l.logger.Errorf("获取子分类失败: %v", err)
		return err
	}
	if len(children) > 0 {
		return errors.New("分类下存在子分类，无法删除")
	}

	// 删除分类
	if err := l.categoryRepo.Delete(ctx, categoryID); err != nil {
		l.logger.Errorf("删除分类失败: %v", err)
		return err
	}

	return nil
}

// GetCategory 获取分类详情
func (l *CategoryLogic) GetCategory(ctx context.Context, categoryID int64) (*models.CmsCategory, error) {
	category, err := l.categoryRepo.Get(ctx, categoryID)
	if err != nil {
		l.logger.Errorf("获取分类失败: %v", err)
		return nil, err
	}
	if category == nil {
		return nil, errors.New("分类不存在")
	}

	return category, nil
}

// GetCategoryBySlug 根据Slug获取分类
func (l *CategoryLogic) GetCategoryBySlug(ctx context.Context, slug string) (*models.CmsCategory, error) {
	category, err := l.categoryRepo.GetBySlug(ctx, slug)
	if err != nil {
		l.logger.Errorf("获取分类失败: %v", err)
		return nil, err
	}
	if category == nil {
		return nil, errors.New("分类不存在")
	}

	return category, nil
}

// ListCategories 获取分类列表
func (l *CategoryLogic) ListCategories(ctx context.Context, parentID *int64) ([]*models.CmsCategory, error) {
	categories, err := l.categoryRepo.List(ctx, parentID)
	if err != nil {
		l.logger.Errorf("获取分类列表失败: %v", err)
		return nil, err
	}

	return categories, nil
}

// GetCategoryTree 获取分类树形结构
func (l *CategoryLogic) GetCategoryTree(ctx context.Context) ([]*models.CmsCategory, error) {
	tree, err := l.categoryRepo.GetTree(ctx)
	if err != nil {
		l.logger.Errorf("获取分类树形结构失败: %v", err)
		return nil, err
	}

	return tree, nil
}

// UpdateCategoryContentCount 更新分类的内容数
func (l *CategoryLogic) UpdateCategoryContentCount(ctx context.Context, categoryID int64) error {
	// 这里可以根据需要计算分类下的内容数
	// 暂时使用一个简单的实现，实际可能需要从Repository查询内容数
	if err := l.categoryRepo.UpdateContentCount(ctx, categoryID, 0); err != nil {
		l.logger.Errorf("更新分类内容数失败: %v", err)
		return err
	}

	return nil
}

// MoveCategoryContent 当删除分类时，将该分类的内容移动到其他分类
func (l *CategoryLogic) MoveCategoryContent(ctx context.Context, fromCategoryID, toCategoryID int64) error {
	// 验证目标分类存在
	target, err := l.categoryRepo.Get(ctx, toCategoryID)
	if err != nil {
		l.logger.Errorf("获取目标分类失败: %v", err)
		return err
	}
	if target == nil {
		return errors.New("目标分类不存在")
	}

	// 这里需要查询并更新所有属于fromCategoryID的内容
	// 具体实现取决于Repository提供的接口
	// 目前的Repository设计中没有这个功能，可以后续扩展

	return nil
}
