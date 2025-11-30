package cms

import (
	"context"
	"errors"
	"power-admin-server/pkg/models"
	"power-admin-server/pkg/repository"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
)

// PublishLogic 发布和工作流业务逻辑
type PublishLogic struct {
	contentRepo  repository.ContentRepository
	categoryRepo repository.CategoryRepository
	logger       logx.Logger
}

// NewPublishLogic 创建PublishLogic实例
func NewPublishLogic(contentRepo repository.ContentRepository, categoryRepo repository.CategoryRepository) *PublishLogic {
	return &PublishLogic{
		contentRepo:  contentRepo,
		categoryRepo: categoryRepo,
		logger:       logx.WithContext(context.Background()),
	}
}

// PublishImmediateRequest 立即发布请求
type PublishImmediateRequest struct {
	ContentID int64 `json:"content_id"`
}

// PublishImmediate 立即发布内容
func (l *PublishLogic) PublishImmediate(ctx context.Context, req *PublishImmediateRequest) error {
	// 验证内容存在
	content, err := l.contentRepo.Get(ctx, req.ContentID)
	if err != nil {
		l.logger.Errorf("获取内容失败: %v", err)
		return err
	}
	if content == nil {
		return errors.New("内容不存在")
	}

	// 发布内容
	if err := l.contentRepo.Publish(ctx, req.ContentID); err != nil {
		l.logger.Errorf("发布内容失败: %v", err)
		return err
	}

	return nil
}

// PublishScheduledRequest 定时发布请求
type PublishScheduledRequest struct {
	ContentID   int64      `json:"content_id"`
	ScheduledAt *time.Time `json:"scheduled_at"`
}

// PublishScheduled 定时发布内容（暂定，等待后续实现定时任务）
func (l *PublishLogic) PublishScheduled(ctx context.Context, req *PublishScheduledRequest) error {
	// 验证内容存在
	content, err := l.contentRepo.Get(ctx, req.ContentID)
	if err != nil {
		l.logger.Errorf("获取内容失败: %v", err)
		return err
	}
	if content == nil {
		return errors.New("内容不存在")
	}

	// 验证定时时间
	if req.ScheduledAt == nil {
		return errors.New("定时发布时间不能为空")
	}

	if req.ScheduledAt.Before(time.Now()) {
		return errors.New("定时发布时间不能早于当前时间")
	}

	// 更新内容的scheduled_at字段
	updateData := &models.CmsContent{
		ScheduledAt: req.ScheduledAt,
		UpdatedAt:   time.Now(),
	}

	if err := l.contentRepo.Update(ctx, req.ContentID, updateData); err != nil {
		l.logger.Errorf("设置定时发布失败: %v", err)
		return err
	}

	// 这里应该向任务队列中添加任务，在指定时间执行发布操作
	// 可以使用消息队列（如RabbitMQ、Redis、Kafka）或定时任务框架（如cron）
	// 暂时在此记录日志，表示已创建定时任务

	l.logger.Infof("已设置内容 %d 的定时发布，计划发布时间：%s", req.ContentID, req.ScheduledAt.Format(time.RFC3339))

	return nil
}

// UnpublishRequest 取消发布请求
type UnpublishRequest struct {
	ContentID int64 `json:"content_id"`
}

// Unpublish 取消发布内容
func (l *PublishLogic) Unpublish(ctx context.Context, req *UnpublishRequest) error {
	// 验证内容存在
	content, err := l.contentRepo.Get(ctx, req.ContentID)
	if err != nil {
		l.logger.Errorf("获取内容失败: %v", err)
		return err
	}
	if content == nil {
		return errors.New("内容不存在")
	}

	// 检查是否已发布
	if content.Status != 2 {
		return errors.New("内容未发布，无法取消发布")
	}

	// 取消发布
	if err := l.contentRepo.Unpublish(ctx, req.ContentID); err != nil {
		l.logger.Errorf("取消发布失败: %v", err)
		return err
	}

	return nil
}

// CancelScheduledPublishRequest 取消定时发布请求
type CancelScheduledPublishRequest struct {
	ContentID int64 `json:"content_id"`
}

// CancelScheduledPublish 取消定时发布
func (l *PublishLogic) CancelScheduledPublish(ctx context.Context, req *CancelScheduledPublishRequest) error {
	// 验证内容存在
	content, err := l.contentRepo.Get(ctx, req.ContentID)
	if err != nil {
		l.logger.Errorf("获取内容失败: %v", err)
		return err
	}
	if content == nil {
		return errors.New("内容不存在")
	}

	// 检查是否有定时发布
	if content.ScheduledAt == nil {
		return errors.New("内容没有设置定时发布")
	}

	// 取消定时发布（清除scheduled_at字段）
	updateData := &models.CmsContent{
		ScheduledAt: nil,
		UpdatedAt:   time.Now(),
	}

	if err := l.contentRepo.Update(ctx, req.ContentID, updateData); err != nil {
		l.logger.Errorf("取消定时发布失败: %v", err)
		return err
	}

	l.logger.Infof("已取消内容 %d 的定时发布", req.ContentID)

	return nil
}

// GetPublishStatusRequest 获取发布状态请求
type GetPublishStatusRequest struct {
	ContentID int64 `json:"content_id"`
}

// PublishStatus 发布状态信息
type PublishStatus struct {
	ContentID   int64      `json:"content_id"`
	Title       string     `json:"title"`
	Status      int8       `json:"status"` // 1:草稿 2:已发布
	PublishedAt *time.Time `json:"published_at"`
	ScheduledAt *time.Time `json:"scheduled_at"`
	IsScheduled bool       `json:"is_scheduled"`
	IsDraft     bool       `json:"is_draft"`
	IsPublished bool       `json:"is_published"`
}

// GetPublishStatus 获取内容发布状态
func (l *PublishLogic) GetPublishStatus(ctx context.Context, contentID int64) (*PublishStatus, error) {
	// 验证内容存在
	content, err := l.contentRepo.Get(ctx, contentID)
	if err != nil {
		l.logger.Errorf("获取内容失败: %v", err)
		return nil, err
	}
	if content == nil {
		return nil, errors.New("内容不存在")
	}

	status := &PublishStatus{
		ContentID:   content.ID,
		Title:       content.Title,
		Status:      content.Status,
		PublishedAt: content.PublishedAt,
		ScheduledAt: content.ScheduledAt,
		IsDraft:     content.Status == 1,
		IsPublished: content.Status == 2,
		IsScheduled: content.ScheduledAt != nil && content.ScheduledAt.After(time.Now()),
	}

	return status, nil
}

// BatchPublishRequest 批量发布请求
type BatchPublishRequest struct {
	ContentIDs []int64 `json:"content_ids"`
}

// BatchPublish 批量发布内容
func (l *PublishLogic) BatchPublish(ctx context.Context, req *BatchPublishRequest) error {
	if len(req.ContentIDs) == 0 {
		return errors.New("内容ID列表不能为空")
	}

	// 批量更新状态为发布
	if err := l.contentRepo.BatchUpdateStatus(ctx, req.ContentIDs, 2); err != nil {
		l.logger.Errorf("批量发布内容失败: %v", err)
		return err
	}

	l.logger.Infof("已批量发布 %d 个内容", len(req.ContentIDs))

	return nil
}

// BatchUnpublishRequest 批量取消发布请求
type BatchUnpublishRequest struct {
	ContentIDs []int64 `json:"content_ids"`
}

// BatchUnpublish 批量取消发布内容
func (l *PublishLogic) BatchUnpublish(ctx context.Context, req *BatchUnpublishRequest) error {
	if len(req.ContentIDs) == 0 {
		return errors.New("内容ID列表不能为空")
	}

	// 批量更新状态为草稿
	if err := l.contentRepo.BatchUpdateStatus(ctx, req.ContentIDs, 1); err != nil {
		l.logger.Errorf("批量取消发布内容失败: %v", err)
		return err
	}

	l.logger.Infof("已批量取消发布 %d 个内容", len(req.ContentIDs))

	return nil
}

// GetScheduledContents 获取待发布的定时内容
func (l *PublishLogic) GetScheduledContents(ctx context.Context) ([]*models.CmsContent, error) {
	// 这里应该从数据库查询所有scheduled_at在当前时间之前的内容
	// 暂时返回空列表，等待后续实现

	l.logger.Infof("检查待发布的定时内容")

	return []*models.CmsContent{}, nil
}

// ProcessScheduledPublish 处理定时发布（应该由后台任务定期调用）
func (l *PublishLogic) ProcessScheduledPublish(ctx context.Context) error {
	// 获取所有待发布的定时内容
	scheduledContents, err := l.GetScheduledContents(ctx)
	if err != nil {
		l.logger.Errorf("获取待发布内容失败: %v", err)
		return err
	}

	// 发布到期的内容
	for _, content := range scheduledContents {
		if err := l.PublishImmediate(ctx, &PublishImmediateRequest{ContentID: content.ID}); err != nil {
			l.logger.Errorf("发布内容 %d 失败: %v", content.ID, err)
			continue
		}

		// 清除scheduled_at字段
		updateData := &models.CmsContent{
			ScheduledAt: nil,
			UpdatedAt:   time.Now(),
		}
		_ = l.contentRepo.Update(ctx, content.ID, updateData)

		l.logger.Infof("已自动发布定时内容: %d", content.ID)
	}

	return nil
}

// GetPublishStatistics 获取发布统计
type PublishStatistics struct {
	TotalContents  int64 `json:"total_contents"`
	PublishedCount int64 `json:"published_count"`
	DraftCount     int64 `json:"draft_count"`
	ScheduledCount int64 `json:"scheduled_count"`
	DeletedCount   int64 `json:"deleted_count"`
}

// GetPublishStatistics 获取发布统计信息（暂时实现为占位符）
func (l *PublishLogic) GetPublishStatisticsInfo(ctx context.Context) (*PublishStatistics, error) {
	// 这里应该从数据库查询各种状态的内容数量
	// 暂时返回空数据结构，等待后续实现

	return &PublishStatistics{}, nil
}
