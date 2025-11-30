package cms

import (
	"context"
	"errors"
	"power-admin-server/pkg/models"
	"power-admin-server/pkg/repository"
	"testing"
	"time"
)

// MockContentRepository 内容仓储的模拟实现
type MockContentRepository struct {
	contents map[int64]*models.CmsContent
	lastID   int64
}

func NewMockContentRepository() *MockContentRepository {
	return &MockContentRepository{
		contents: make(map[int64]*models.CmsContent),
		lastID:   0,
	}
}

func (m *MockContentRepository) Create(ctx context.Context, content *models.CmsContent) error {
	m.lastID++
	content.ID = m.lastID
	m.contents[content.ID] = content
	return nil
}

func (m *MockContentRepository) Get(ctx context.Context, id int64) (*models.CmsContent, error) {
	if content, exists := m.contents[id]; exists {
		return content, nil
	}
	return nil, nil
}

func (m *MockContentRepository) Update(ctx context.Context, id int64, content *models.CmsContent) error {
	if existing, exists := m.contents[id]; exists {
		if content.Title != "" {
			existing.Title = content.Title
		}
		if content.Content != "" {
			existing.Content = content.Content
		}
		existing.UpdatedAt = content.UpdatedAt
		return nil
	}
	return errors.New("content not found")
}

func (m *MockContentRepository) Delete(ctx context.Context, id int64) error {
	if _, exists := m.contents[id]; exists {
		delete(m.contents, id)
		return nil
	}
	return errors.New("content not found")
}

func (m *MockContentRepository) SoftDelete(ctx context.Context, id int64) error {
	if existing, exists := m.contents[id]; exists {
		existing.Status = 3
		return nil
	}
	return errors.New("content not found")
}

func (m *MockContentRepository) List(ctx context.Context, req *repository.ContentListRequest) (*repository.PagedResult, error) {
	var items []*models.CmsContent
	for _, content := range m.contents {
		items = append(items, content)
	}
	return &repository.PagedResult{
		Total:    int64(len(items)),
		Page:     1,
		PageSize: 10,
		Items:    items,
	}, nil
}

func (m *MockContentRepository) GetBySlug(ctx context.Context, slug string) (*models.CmsContent, error) {
	for _, content := range m.contents {
		if content.Slug == slug {
			return content, nil
		}
	}
	return nil, nil
}

func (m *MockContentRepository) Publish(ctx context.Context, id int64) error {
	if existing, exists := m.contents[id]; exists {
		existing.Status = 2
		now := time.Now()
		existing.PublishedAt = &now
		return nil
	}
	return errors.New("content not found")
}

func (m *MockContentRepository) Unpublish(ctx context.Context, id int64) error {
	if existing, exists := m.contents[id]; exists {
		existing.Status = 1
		return nil
	}
	return errors.New("content not found")
}

func (m *MockContentRepository) IncrementViewCount(ctx context.Context, id int64) error {
	if existing, exists := m.contents[id]; exists {
		existing.ViewCount++
		return nil
	}
	return errors.New("content not found")
}

func (m *MockContentRepository) BatchUpdateStatus(ctx context.Context, ids []int64, status int8) error {
	for _, id := range ids {
		if existing, exists := m.contents[id]; exists {
			existing.Status = status
		}
	}
	return nil
}

// MockCategoryRepository 分类仓储的模拟实现
type MockCategoryRepository struct {
	categories map[int64]*models.CmsCategory
	lastID     int64
}

func NewMockCategoryRepository() *MockCategoryRepository {
	return &MockCategoryRepository{
		categories: make(map[int64]*models.CmsCategory),
		lastID:     0,
	}
}

func (m *MockCategoryRepository) Create(ctx context.Context, category *models.CmsCategory) error {
	m.lastID++
	category.ID = m.lastID
	m.categories[category.ID] = category
	return nil
}

func (m *MockCategoryRepository) Get(ctx context.Context, id int64) (*models.CmsCategory, error) {
	if category, exists := m.categories[id]; exists {
		return category, nil
	}
	return nil, nil
}

func (m *MockCategoryRepository) Update(ctx context.Context, id int64, category *models.CmsCategory) error {
	if existing, exists := m.categories[id]; exists {
		if category.Name != "" {
			existing.Name = category.Name
		}
		existing.UpdatedAt = category.UpdatedAt
		return nil
	}
	return errors.New("category not found")
}

func (m *MockCategoryRepository) Delete(ctx context.Context, id int64) error {
	if _, exists := m.categories[id]; exists {
		delete(m.categories, id)
		return nil
	}
	return errors.New("category not found")
}

func (m *MockCategoryRepository) List(ctx context.Context, parentID *int64) ([]*models.CmsCategory, error) {
	var categories []*models.CmsCategory
	for _, cat := range m.categories {
		if parentID == nil && cat.ParentID == nil {
			categories = append(categories, cat)
		} else if parentID != nil && cat.ParentID != nil && *parentID == *cat.ParentID {
			categories = append(categories, cat)
		}
	}
	return categories, nil
}

func (m *MockCategoryRepository) GetTree(ctx context.Context) ([]*models.CmsCategory, error) {
	var rootCategories []*models.CmsCategory
	for _, cat := range m.categories {
		if cat.ParentID == nil {
			rootCategories = append(rootCategories, cat)
		}
	}
	return rootCategories, nil
}

func (m *MockCategoryRepository) GetBySlug(ctx context.Context, slug string) (*models.CmsCategory, error) {
	for _, cat := range m.categories {
		if cat.Slug == slug {
			return cat, nil
		}
	}
	return nil, nil
}

func (m *MockCategoryRepository) UpdateContentCount(ctx context.Context, categoryID int64, count int) error {
	if existing, exists := m.categories[categoryID]; exists {
		existing.ContentCount = count
		return nil
	}
	return errors.New("category not found")
}

// Test CreateContent
func TestCreateContent(t *testing.T) {
	contentRepo := NewMockContentRepository()
	categoryRepo := NewMockCategoryRepository()
	logic := NewContentLogic(contentRepo, categoryRepo)

	// 测试成功创建内容
	req := &CreateContentRequest{
		Title:    "Test Title",
		Content:  "Test Content",
		AuthorID: 1,
		Status:   1,
	}

	content, err := logic.CreateContent(context.Background(), req)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if content == nil {
		t.Error("Expected content, got nil")
	}
	if content.Title != "Test Title" {
		t.Errorf("Expected title 'Test Title', got '%s'", content.Title)
	}

	// 测试空标题
	req2 := &CreateContentRequest{
		Content:  "Test Content",
		AuthorID: 1,
	}
	_, err = logic.CreateContent(context.Background(), req2)
	if err == nil {
		t.Error("Expected error for empty title")
	}

	// 测试空内容
	req3 := &CreateContentRequest{
		Title:    "Test",
		AuthorID: 1,
	}
	_, err = logic.CreateContent(context.Background(), req3)
	if err == nil {
		t.Error("Expected error for empty content")
	}
}

// Test PublishContent
func TestPublishContent(t *testing.T) {
	contentRepo := NewMockContentRepository()
	categoryRepo := NewMockCategoryRepository()
	logic := NewContentLogic(contentRepo, categoryRepo)

	// 创建内容
	req := &CreateContentRequest{
		Title:    "Test",
		Content:  "Content",
		AuthorID: 1,
		Status:   1,
	}
	content, _ := logic.CreateContent(context.Background(), req)

	// 发布内容
	err := logic.PublishContent(context.Background(), content.ID)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// 验证状态
	published, _ := contentRepo.Get(context.Background(), content.ID)
	if published.Status != 2 {
		t.Errorf("Expected status 2, got %d", published.Status)
	}

	// 测试发布不存在的内容
	err = logic.PublishContent(context.Background(), 9999)
	if err == nil {
		t.Error("Expected error for non-existent content")
	}
}

// Test DeleteContent
func TestDeleteContent(t *testing.T) {
	contentRepo := NewMockContentRepository()
	categoryRepo := NewMockCategoryRepository()
	logic := NewContentLogic(contentRepo, categoryRepo)

	// 创建内容
	req := &CreateContentRequest{
		Title:    "Test",
		Content:  "Content",
		AuthorID: 1,
		Status:   1,
	}
	content, _ := logic.CreateContent(context.Background(), req)

	// 软删除
	err := logic.DeleteContent(context.Background(), content.ID)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// 验证状态（应该是3表示已删除）
	deleted, _ := contentRepo.Get(context.Background(), content.ID)
	if deleted.Status != 3 {
		t.Errorf("Expected status 3, got %d", deleted.Status)
	}

	// 测试删除不存在的内容
	err = logic.DeleteContent(context.Background(), 9999)
	if err == nil {
		t.Error("Expected error for non-existent content")
	}
}
