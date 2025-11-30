package cms

import (
	"context"
	"errors"
	"power-admin-server/pkg/models"
	"testing"
)

// MockCommentRepository 评论仓储的模拟实现
type MockCommentRepository struct {
	comments map[int64]*models.CmsComment
	lastID   int64
}

func NewMockCommentRepository() *MockCommentRepository {
	return &MockCommentRepository{
		comments: make(map[int64]*models.CmsComment),
		lastID:   0,
	}
}

func (m *MockCommentRepository) Create(ctx context.Context, comment *models.CmsComment) error {
	m.lastID++
	comment.ID = m.lastID
	m.comments[comment.ID] = comment
	return nil
}

func (m *MockCommentRepository) Get(ctx context.Context, id int64) (*models.CmsComment, error) {
	if comment, exists := m.comments[id]; exists {
		return comment, nil
	}
	return nil, nil
}

func (m *MockCommentRepository) Update(ctx context.Context, id int64, comment *models.CmsComment) error {
	if existing, exists := m.comments[id]; exists {
		if comment.Content != "" {
			existing.Content = comment.Content
		}
		if comment.Status != 0 {
			existing.Status = comment.Status
		}
		existing.UpdatedAt = comment.UpdatedAt
		return nil
	}
	return errors.New("comment not found")
}

func (m *MockCommentRepository) Delete(ctx context.Context, id int64) error {
	if _, exists := m.comments[id]; exists {
		delete(m.comments, id)
		return nil
	}
	return errors.New("comment not found")
}

func (m *MockCommentRepository) List(ctx context.Context, contentID int64) ([]*models.CmsComment, error) {
	var comments []*models.CmsComment
	for _, comment := range m.comments {
		if comment.ContentID == contentID && comment.Status == 1 {
			comments = append(comments, comment)
		}
	}
	return comments, nil
}

func (m *MockCommentRepository) Approve(ctx context.Context, id int64) error {
	if existing, exists := m.comments[id]; exists {
		existing.Status = 1
		return nil
	}
	return errors.New("comment not found")
}

func (m *MockCommentRepository) UpdateLikeCount(ctx context.Context, id int64, count int) error {
	if existing, exists := m.comments[id]; exists {
		existing.LikeCount = count
		return nil
	}
	return errors.New("comment not found")
}

func TestCreateComment(t *testing.T) {
	commentRepo := NewMockCommentRepository()
	contentRepo := NewMockContentRepository()
	logic := NewCommentLogic(commentRepo, contentRepo)

	// 测试成功创建评论
	req := &CreateCommentRequest{
		ContentID:   1, // 不存在的内容
		AuthorName:  "Test User",
		AuthorEmail: "test@example.com",
		Content:     "Test Comment",
	}

	comment, err := logic.CreateComment(context.Background(), req)
	if err == nil {
		t.Error("Expected error for non-existent content")
	}

	// 创建一个真实的内容
	content := &models.CmsContent{
		Title:    "Test Content",
		Content:  "Content",
		AuthorID: 1,
		Status:   2,
	}
	contentRepo.Create(context.Background(), content)

	// 测试成功创建评论
	req2 := &CreateCommentRequest{
		ContentID:   content.ID,
		AuthorName:  "Test User",
		AuthorEmail: "test@example.com",
		Content:     "Test Comment",
	}

	comment, err = logic.CreateComment(context.Background(), req2)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if comment == nil {
		t.Error("Expected comment, got nil")
	} else {
		if comment.Content != "Test Comment" {
			t.Errorf("Expected content 'Test Comment', got '%s'", comment.Content)
		}
		if comment.Status != 0 {
			t.Errorf("Expected status 0 (pending), got %d", comment.Status)
		}
	}

	// 测试空内容
	req3 := &CreateCommentRequest{
		ContentID: content.ID,
	}
	_, err = logic.CreateComment(context.Background(), req3)
	if err == nil {
		t.Error("Expected error for empty comment")
	}
}

// Test ApproveComment
func TestApproveComment(t *testing.T) {
	commentRepo := NewMockCommentRepository()
	contentRepo := NewMockContentRepository()
	logic := NewCommentLogic(commentRepo, contentRepo)

	// 创建内容
	content := &models.CmsContent{
		Title:    "Test Content",
		Content:  "Content",
		AuthorID: 1,
		Status:   2,
	}
	contentRepo.Create(context.Background(), content)

	// 创建评论
	req := &CreateCommentRequest{
		ContentID:   content.ID,
		AuthorName:  "Test",
		AuthorEmail: "test@example.com",
		Content:     "Test",
	}
	comment, _ := logic.CreateComment(context.Background(), req)

	// 审核评论
	err := logic.ApproveComment(context.Background(), comment.ID)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// 验证状态
	approved, _ := commentRepo.Get(context.Background(), comment.ID)
	if approved.Status != 1 {
		t.Errorf("Expected status 1, got %d", approved.Status)
	}

	// 测试审核不存在的评论
	err = logic.ApproveComment(context.Background(), 9999)
	if err == nil {
		t.Error("Expected error for non-existent comment")
	}
}

// Test LikeComment
func TestLikeComment(t *testing.T) {
	commentRepo := NewMockCommentRepository()
	contentRepo := NewMockContentRepository()
	logic := NewCommentLogic(commentRepo, contentRepo)

	// 创建内容
	content := &models.CmsContent{
		Title:    "Test Content",
		Content:  "Content",
		AuthorID: 1,
		Status:   2,
	}
	contentRepo.Create(context.Background(), content)

	// 创建评论
	req := &CreateCommentRequest{
		ContentID:   content.ID,
		AuthorName:  "Test",
		AuthorEmail: "test@example.com",
		Content:     "Test",
	}
	comment, _ := logic.CreateComment(context.Background(), req)

	// 点赞
	err := logic.LikeComment(context.Background(), comment.ID)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// 验证点赞数
	liked, _ := commentRepo.Get(context.Background(), comment.ID)
	if liked.LikeCount != 1 {
		t.Errorf("Expected like count 1, got %d", liked.LikeCount)
	}
}

// Test DeleteComment
func TestDeleteComment(t *testing.T) {
	commentRepo := NewMockCommentRepository()
	contentRepo := NewMockContentRepository()
	logic := NewCommentLogic(commentRepo, contentRepo)

	// 创建内容
	content := &models.CmsContent{
		Title:    "Test Content",
		Content:  "Content",
		AuthorID: 1,
		Status:   2,
	}
	contentRepo.Create(context.Background(), content)

	// 创建评论
	req := &CreateCommentRequest{
		ContentID:   content.ID,
		AuthorName:  "Test",
		AuthorEmail: "test@example.com",
		Content:     "Test",
	}
	comment, _ := logic.CreateComment(context.Background(), req)

	// 删除评论
	err := logic.DeleteComment(context.Background(), comment.ID)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// 验证评论已删除
	deleted, _ := commentRepo.Get(context.Background(), comment.ID)
	if deleted != nil {
		t.Error("Expected comment to be deleted")
	}

	// 测试删除不存在的评论
	err = logic.DeleteComment(context.Background(), 9999)
	if err == nil {
		t.Error("Expected error for non-existent comment")
	}
}

// Test RejectComment
func TestRejectComment(t *testing.T) {
	commentRepo := NewMockCommentRepository()
	contentRepo := NewMockContentRepository()
	logic := NewCommentLogic(commentRepo, contentRepo)

	// 创建内容
	content := &models.CmsContent{
		Title:    "Test Content",
		Content:  "Content",
		AuthorID: 1,
		Status:   2,
	}
	contentRepo.Create(context.Background(), content)

	// 创建评论
	req := &CreateCommentRequest{
		ContentID:   content.ID,
		AuthorName:  "Test",
		AuthorEmail: "test@example.com",
		Content:     "Test",
	}
	comment, _ := logic.CreateComment(context.Background(), req)

	// 拒绝评论
	err := logic.RejectComment(context.Background(), comment.ID)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// 验证状态
	rejected, _ := commentRepo.Get(context.Background(), comment.ID)
	if rejected.Status != 2 {
		t.Errorf("Expected status 2 (rejected), got %d", rejected.Status)
	}
}

// Test SpamComment
func TestSpamComment(t *testing.T) {
	commentRepo := NewMockCommentRepository()
	contentRepo := NewMockContentRepository()
	logic := NewCommentLogic(commentRepo, contentRepo)

	// 创建内容
	content := &models.CmsContent{
		Title:    "Test Content",
		Content:  "Content",
		AuthorID: 1,
		Status:   2,
	}
	contentRepo.Create(context.Background(), content)

	// 创建评论
	req := &CreateCommentRequest{
		ContentID:   content.ID,
		AuthorName:  "Test",
		AuthorEmail: "test@example.com",
		Content:     "Test",
	}
	comment, _ := logic.CreateComment(context.Background(), req)

	// 标记为垃圾
	err := logic.SpamComment(context.Background(), comment.ID)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// 验证状态
	spam, _ := commentRepo.Get(context.Background(), comment.ID)
	if spam.Status != 3 {
		t.Errorf("Expected status 3 (spam), got %d", spam.Status)
	}
}
