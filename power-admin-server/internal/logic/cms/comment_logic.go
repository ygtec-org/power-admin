package cms

import (
	"context"
	"errors"
	"power-admin-server/pkg/models"
	"power-admin-server/pkg/repository"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
)

// CommentLogic 评论管理业务逻辑
type CommentLogic struct {
	commentRepo repository.CommentRepository
	contentRepo repository.ContentRepository
	logger      logx.Logger
}

// NewCommentLogic 创建CommentLogic实例
func NewCommentLogic(commentRepo repository.CommentRepository, contentRepo repository.ContentRepository) *CommentLogic {
	return &CommentLogic{
		commentRepo: commentRepo,
		contentRepo: contentRepo,
		logger:      logx.WithContext(context.Background()),
	}
}

// CreateCommentRequest 创建评论请求
type CreateCommentRequest struct {
	ContentID       int64  `json:"content_id"`
	UserID          *int64 `json:"user_id"`
	ParentCommentID *int64 `json:"parent_comment_id"`
	AuthorName      string `json:"author_name"`
	AuthorEmail     string `json:"author_email"`
	Content         string `json:"content"`
	IPAddress       string `json:"ip_address"`
	UserAgent       string `json:"user_agent"`
}

// CreateComment 创建评论
func (l *CommentLogic) CreateComment(ctx context.Context, req *CreateCommentRequest) (*models.CmsComment, error) {
	// 验证必填字段
	if req.ContentID == 0 {
		return nil, errors.New("内容ID不能为空")
	}
	if req.Content == "" {
		return nil, errors.New("评论内容不能为空")
	}

	// 验证内容存在
	content, err := l.contentRepo.Get(ctx, req.ContentID)
	if err != nil {
		l.logger.Errorf("获取内容失败: %v", err)
		return nil, err
	}
	if content == nil {
		return nil, errors.New("内容不存在")
	}

	// 验证父评论（如果有）
	if req.ParentCommentID != nil {
		parentComment, err := l.commentRepo.Get(ctx, *req.ParentCommentID)
		if err != nil {
			l.logger.Errorf("获取父评论失败: %v", err)
			return nil, err
		}
		if parentComment == nil {
			return nil, errors.New("父评论不存在")
		}
	}

	// 默认状态为待审核（0）
	comment := &models.CmsComment{
		ContentID:       req.ContentID,
		UserID:          req.UserID,
		ParentCommentID: req.ParentCommentID,
		AuthorName:      req.AuthorName,
		AuthorEmail:     req.AuthorEmail,
		Content:         req.Content,
		Status:          0, // 待审核
		LikeCount:       0,
		ReplyCount:      0,
		IPAddress:       req.IPAddress,
		UserAgent:       req.UserAgent,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	if err := l.commentRepo.Create(ctx, comment); err != nil {
		l.logger.Errorf("创建评论失败: %v", err)
		return nil, err
	}

	return comment, nil
}

// UpdateCommentRequest 更新评论请求
type UpdateCommentRequest struct {
	ID              int64
	Content         string
	Status          int8
	ParentCommentID *int64
}

// UpdateComment 更新评论
func (l *CommentLogic) UpdateComment(ctx context.Context, req *UpdateCommentRequest) (*models.CmsComment, error) {
	// 获取原评论
	existing, err := l.commentRepo.Get(ctx, req.ID)
	if err != nil {
		l.logger.Errorf("获取评论失败: %v", err)
		return nil, err
	}
	if existing == nil {
		return nil, errors.New("评论不存在")
	}

	updateData := &models.CmsComment{
		Content:         req.Content,
		Status:          req.Status,
		ParentCommentID: req.ParentCommentID,
		UpdatedAt:       time.Now(),
	}

	if err := l.commentRepo.Update(ctx, req.ID, updateData); err != nil {
		l.logger.Errorf("更新评论失败: %v", err)
		return nil, err
	}

	return l.commentRepo.Get(ctx, req.ID)
}

// DeleteComment 删除评论
func (l *CommentLogic) DeleteComment(ctx context.Context, commentID int64) error {
	// 获取评论
	existing, err := l.commentRepo.Get(ctx, commentID)
	if err != nil {
		l.logger.Errorf("获取评论失败: %v", err)
		return err
	}
	if existing == nil {
		return errors.New("评论不存在")
	}

	if err := l.commentRepo.Delete(ctx, commentID); err != nil {
		l.logger.Errorf("删除评论失败: %v", err)
		return err
	}

	return nil
}

// GetComment 获取评论详情
func (l *CommentLogic) GetComment(ctx context.Context, commentID int64) (*models.CmsComment, error) {
	comment, err := l.commentRepo.Get(ctx, commentID)
	if err != nil {
		l.logger.Errorf("获取评论失败: %v", err)
		return nil, err
	}
	if comment == nil {
		return nil, errors.New("评论不存在")
	}

	return comment, nil
}

// ListComments 获取内容的评论列表
func (l *CommentLogic) ListComments(ctx context.Context, contentID int64) ([]*models.CmsComment, error) {
	// 验证内容存在
	content, err := l.contentRepo.Get(ctx, contentID)
	if err != nil {
		l.logger.Errorf("获取内容失败: %v", err)
		return nil, err
	}
	if content == nil {
		return nil, errors.New("内容不存在")
	}

	comments, err := l.commentRepo.List(ctx, contentID)
	if err != nil {
		l.logger.Errorf("获取评论列表失败: %v", err)
		return nil, err
	}

	return comments, nil
}

// ApproveComment 审核通过评论
func (l *CommentLogic) ApproveComment(ctx context.Context, commentID int64) error {
	// 获取评论
	existing, err := l.commentRepo.Get(ctx, commentID)
	if err != nil {
		l.logger.Errorf("获取评论失败: %v", err)
		return err
	}
	if existing == nil {
		return errors.New("评论不存在")
	}

	if err := l.commentRepo.Approve(ctx, commentID); err != nil {
		l.logger.Errorf("审核评论失败: %v", err)
		return err
	}

	return nil
}

// RejectComment 拒绝评论
func (l *CommentLogic) RejectComment(ctx context.Context, commentID int64) error {
	// 获取评论
	existing, err := l.commentRepo.Get(ctx, commentID)
	if err != nil {
		l.logger.Errorf("获取评论失败: %v", err)
		return err
	}
	if existing == nil {
		return errors.New("评论不存在")
	}

	// 更新评论状态为已拒绝
	updateData := &models.CmsComment{
		Status:    2, // 2:已拒绝
		UpdatedAt: time.Now(),
	}

	if err := l.commentRepo.Update(ctx, commentID, updateData); err != nil {
		l.logger.Errorf("拒绝评论失败: %v", err)
		return err
	}

	return nil
}

// LikeComment 点赞评论
func (l *CommentLogic) LikeComment(ctx context.Context, commentID int64) error {
	// 获取评论
	existing, err := l.commentRepo.Get(ctx, commentID)
	if err != nil {
		l.logger.Errorf("获取评论失败: %v", err)
		return err
	}
	if existing == nil {
		return errors.New("评论不存在")
	}

	newLikeCount := existing.LikeCount + 1
	if err := l.commentRepo.UpdateLikeCount(ctx, commentID, newLikeCount); err != nil {
		l.logger.Errorf("更新评论点赞数失败: %v", err)
		return err
	}

	return nil
}

// UnlikeComment 取消点赞评论
func (l *CommentLogic) UnlikeComment(ctx context.Context, commentID int64) error {
	// 获取评论
	existing, err := l.commentRepo.Get(ctx, commentID)
	if err != nil {
		l.logger.Errorf("获取评论失败: %v", err)
		return err
	}
	if existing == nil {
		return errors.New("评论不存在")
	}

	newLikeCount := existing.LikeCount - 1
	if newLikeCount < 0 {
		newLikeCount = 0
	}

	if err := l.commentRepo.UpdateLikeCount(ctx, commentID, newLikeCount); err != nil {
		l.logger.Errorf("更新评论点赞数失败: %v", err)
		return err
	}

	return nil
}

// ReplyComment 回复评论
func (l *CommentLogic) ReplyComment(ctx context.Context, parentCommentID int64, req *CreateCommentRequest) (*models.CmsComment, error) {
	// 获取父评论
	parentComment, err := l.commentRepo.Get(ctx, parentCommentID)
	if err != nil {
		l.logger.Errorf("获取父评论失败: %v", err)
		return nil, err
	}
	if parentComment == nil {
		return nil, errors.New("父评论不存在")
	}

	// 设置父评论ID
	req.ParentCommentID = &parentCommentID
	req.ContentID = parentComment.ContentID

	// 创建回复
	reply, err := l.CreateComment(ctx, req)
	if err != nil {
		return nil, err
	}

	// 更新父评论的回复数
	parentComment.ReplyCount++
	updateParent := &models.CmsComment{
		ReplyCount: parentComment.ReplyCount,
		UpdatedAt:  time.Now(),
	}
	_ = l.commentRepo.Update(ctx, parentCommentID, updateParent)

	return reply, nil
}

// SpamComment 标记评论为垃圾
func (l *CommentLogic) SpamComment(ctx context.Context, commentID int64) error {
	// 获取评论
	existing, err := l.commentRepo.Get(ctx, commentID)
	if err != nil {
		l.logger.Errorf("获取评论失败: %v", err)
		return err
	}
	if existing == nil {
		return errors.New("评论不存在")
	}

	// 更新评论状态为垃圾（3）
	updateData := &models.CmsComment{
		Status:    3,
		UpdatedAt: time.Now(),
	}

	if err := l.commentRepo.Update(ctx, commentID, updateData); err != nil {
		l.logger.Errorf("标记评论为垃圾失败: %v", err)
		return err
	}

	return nil
}
