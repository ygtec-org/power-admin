package cms

import (
	"context"
	"errors"
	"power-admin-server/pkg/models"
	"power-admin-server/pkg/repository"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
)

// TagLogic 标签管理业务逻辑
type TagLogic struct {
	tagRepo repository.TagRepository
	logger  logx.Logger
}

// NewTagLogic 创建TagLogic实例
func NewTagLogic(tagRepo repository.TagRepository) *TagLogic {
	return &TagLogic{
		tagRepo: tagRepo,
		logger:  logx.WithContext(context.Background()),
	}
}

// CreateTagRequest 创建标签请求
type CreateTagRequest struct {
	Name        string `json:"name"`
	Slug        string `json:"slug"`
	Description string `json:"description"`
	Color       string `json:"color"`
	Status      int8   `json:"status"`
}

// CreateTag 创建标签
func (l *TagLogic) CreateTag(ctx context.Context, req *CreateTagRequest) (*models.CmsTag, error) {
	// 验证必填字段
	if req.Name == "" {
		return nil, errors.New("标签名称不能为空")
	}

	// 检查标签是否已存在
	existing, err := l.tagRepo.GetByName(ctx, req.Name)
	if err != nil {
		l.logger.Errorf("检查标签失败: %v", err)
		return nil, err
	}
	if existing != nil {
		return nil, errors.New("标签已存在")
	}

	tag := &models.CmsTag{
		Name:        req.Name,
		Slug:        req.Slug,
		Description: req.Description,
		Color:       req.Color,
		Status:      req.Status,
		UsageCount:  0,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := l.tagRepo.Create(ctx, tag); err != nil {
		l.logger.Errorf("创建标签失败: %v", err)
		return nil, err
	}

	return tag, nil
}

// UpdateTagRequest 更新标签请求
type UpdateTagRequest struct {
	ID          int64
	Name        string
	Slug        string
	Description string
	Color       string
	Status      int8
}

// UpdateTag 更新标签
func (l *TagLogic) UpdateTag(ctx context.Context, req *UpdateTagRequest) (*models.CmsTag, error) {
	// 获取原标签
	existing, err := l.tagRepo.Get(ctx, req.ID)
	if err != nil {
		l.logger.Errorf("获取标签失败: %v", err)
		return nil, err
	}
	if existing == nil {
		return nil, errors.New("标签不存在")
	}

	// 如果名称改变，检查新名称是否已存在
	if req.Name != existing.Name {
		nameExists, err := l.tagRepo.GetByName(ctx, req.Name)
		if err != nil {
			l.logger.Errorf("检查标签失败: %v", err)
			return nil, err
		}
		if nameExists != nil {
			return nil, errors.New("标签名称已存在")
		}
	}

	updateData := &models.CmsTag{
		Name:        req.Name,
		Slug:        req.Slug,
		Description: req.Description,
		Color:       req.Color,
		Status:      req.Status,
		UpdatedAt:   time.Now(),
	}

	if err := l.tagRepo.Update(ctx, req.ID, updateData); err != nil {
		l.logger.Errorf("更新标签失败: %v", err)
		return nil, err
	}

	return l.tagRepo.Get(ctx, req.ID)
}

// DeleteTag 删除标签
func (l *TagLogic) DeleteTag(ctx context.Context, tagID int64) error {
	// 获取标签
	existing, err := l.tagRepo.Get(ctx, tagID)
	if err != nil {
		l.logger.Errorf("获取标签失败: %v", err)
		return err
	}
	if existing == nil {
		return errors.New("标签不存在")
	}

	// 检查标签的使用数
	if existing.UsageCount > 0 {
		return errors.New("标签正在使用中，无法删除")
	}

	if err := l.tagRepo.Delete(ctx, tagID); err != nil {
		l.logger.Errorf("删除标签失败: %v", err)
		return err
	}

	return nil
}

// GetTag 获取标签详情
func (l *TagLogic) GetTag(ctx context.Context, tagID int64) (*models.CmsTag, error) {
	tag, err := l.tagRepo.Get(ctx, tagID)
	if err != nil {
		l.logger.Errorf("获取标签失败: %v", err)
		return nil, err
	}
	if tag == nil {
		return nil, errors.New("标签不存在")
	}

	return tag, nil
}

// GetTagByName 根据名称获取标签
func (l *TagLogic) GetTagByName(ctx context.Context, name string) (*models.CmsTag, error) {
	tag, err := l.tagRepo.GetByName(ctx, name)
	if err != nil {
		l.logger.Errorf("获取标签失败: %v", err)
		return nil, err
	}
	if tag == nil {
		return nil, errors.New("标签不存在")
	}

	return tag, nil
}

// ListTags 获取所有标签
func (l *TagLogic) ListTags(ctx context.Context) ([]*models.CmsTag, error) {
	tags, err := l.tagRepo.List(ctx)
	if err != nil {
		l.logger.Errorf("获取标签列表失败: %v", err)
		return nil, err
	}

	return tags, nil
}

// GetTagsByIDs 根据IDs获取多个标签
func (l *TagLogic) GetTagsByIDs(ctx context.Context, tagIDs []int64) ([]*models.CmsTag, error) {
	if len(tagIDs) == 0 {
		return []*models.CmsTag{}, nil
	}

	tags, err := l.tagRepo.GetByIDs(ctx, tagIDs)
	if err != nil {
		l.logger.Errorf("获取标签列表失败: %v", err)
		return nil, err
	}

	return tags, nil
}

// IncrementTagUsage 增加标签使用数
func (l *TagLogic) IncrementTagUsage(ctx context.Context, tagID int64) error {
	tag, err := l.tagRepo.Get(ctx, tagID)
	if err != nil {
		l.logger.Errorf("获取标签失败: %v", err)
		return err
	}
	if tag == nil {
		return errors.New("标签不存在")
	}

	newCount := tag.UsageCount + 1
	if err := l.tagRepo.UpdateUsageCount(ctx, tagID, newCount); err != nil {
		l.logger.Errorf("更新标签使用数失败: %v", err)
		return err
	}

	return nil
}

// DecrementTagUsage 减少标签使用数
func (l *TagLogic) DecrementTagUsage(ctx context.Context, tagID int64) error {
	tag, err := l.tagRepo.Get(ctx, tagID)
	if err != nil {
		l.logger.Errorf("获取标签失败: %v", err)
		return err
	}
	if tag == nil {
		return errors.New("标签不存在")
	}

	newCount := tag.UsageCount - 1
	if newCount < 0 {
		newCount = 0
	}

	if err := l.tagRepo.UpdateUsageCount(ctx, tagID, newCount); err != nil {
		l.logger.Errorf("更新标签使用数失败: %v", err)
		return err
	}

	return nil
}

// BatchGetOrCreateTags 批量获取或创建标签
func (l *TagLogic) BatchGetOrCreateTags(ctx context.Context, tagNames []string) ([]*models.CmsTag, error) {
	if len(tagNames) == 0 {
		return []*models.CmsTag{}, nil
	}

	var tags []*models.CmsTag

	for _, name := range tagNames {
		tag, err := l.tagRepo.GetByName(ctx, name)
		if err != nil {
			l.logger.Errorf("获取标签失败: %v", err)
			return nil, err
		}

		// 如果标签不存在则创建
		if tag == nil {
			req := &CreateTagRequest{
				Name:   name,
				Status: 1,
			}
			newTag, err := l.CreateTag(ctx, req)
			if err != nil {
				l.logger.Errorf("创建标签失败: %v", err)
				return nil, err
			}
			tags = append(tags, newTag)
		} else {
			tags = append(tags, tag)
		}
	}

	return tags, nil
}
