package cms

import (
	"context"
	"errors"
	"power-admin-server/pkg/models"
	"testing"
)

// MockTagRepository 标签仓储的模拟实现
type MockTagRepository struct {
	tags   map[int64]*models.CmsTag
	lastID int64
}

func NewMockTagRepository() *MockTagRepository {
	return &MockTagRepository{
		tags:   make(map[int64]*models.CmsTag),
		lastID: 0,
	}
}

func (m *MockTagRepository) Create(ctx context.Context, tag *models.CmsTag) error {
	m.lastID++
	tag.ID = m.lastID
	m.tags[tag.ID] = tag
	return nil
}

func (m *MockTagRepository) Get(ctx context.Context, id int64) (*models.CmsTag, error) {
	if tag, exists := m.tags[id]; exists {
		return tag, nil
	}
	return nil, nil
}

func (m *MockTagRepository) Update(ctx context.Context, id int64, tag *models.CmsTag) error {
	if existing, exists := m.tags[id]; exists {
		if tag.Name != "" {
			existing.Name = tag.Name
		}
		existing.UpdatedAt = tag.UpdatedAt
		return nil
	}
	return errors.New("tag not found")
}

func (m *MockTagRepository) Delete(ctx context.Context, id int64) error {
	if _, exists := m.tags[id]; exists {
		delete(m.tags, id)
		return nil
	}
	return errors.New("tag not found")
}

func (m *MockTagRepository) List(ctx context.Context) ([]*models.CmsTag, error) {
	var tags []*models.CmsTag
	for _, tag := range m.tags {
		tags = append(tags, tag)
	}
	return tags, nil
}

func (m *MockTagRepository) GetByName(ctx context.Context, name string) (*models.CmsTag, error) {
	for _, tag := range m.tags {
		if tag.Name == name {
			return tag, nil
		}
	}
	return nil, nil
}

func (m *MockTagRepository) GetByIDs(ctx context.Context, ids []int64) ([]*models.CmsTag, error) {
	var tags []*models.CmsTag
	for _, id := range ids {
		if tag, exists := m.tags[id]; exists {
			tags = append(tags, tag)
		}
	}
	return tags, nil
}

func (m *MockTagRepository) UpdateUsageCount(ctx context.Context, tagID int64, count int) error {
	if existing, exists := m.tags[tagID]; exists {
		existing.UsageCount = count
		return nil
	}
	return errors.New("tag not found")
}

// Test CreateTag
func TestCreateTag(t *testing.T) {
	tagRepo := NewMockTagRepository()
	logic := NewTagLogic(tagRepo)

	// 测试成功创建标签
	req := &CreateTagRequest{
		Name:   "Test Tag",
		Status: 1,
	}

	tag, err := logic.CreateTag(context.Background(), req)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if tag == nil {
		t.Error("Expected tag, got nil")
	}
	if tag.Name != "Test Tag" {
		t.Errorf("Expected name 'Test Tag', got '%s'", tag.Name)
	}

	// 测试空名称
	req2 := &CreateTagRequest{}
	_, err = logic.CreateTag(context.Background(), req2)
	if err == nil {
		t.Error("Expected error for empty name")
	}

	// 测试重复名称
	req3 := &CreateTagRequest{
		Name:   "Test Tag",
		Status: 1,
	}
	_, err = logic.CreateTag(context.Background(), req3)
	if err == nil {
		t.Error("Expected error for duplicate name")
	}
}

// Test UpdateTag
func TestUpdateTag(t *testing.T) {
	tagRepo := NewMockTagRepository()
	logic := NewTagLogic(tagRepo)

	// 创建标签
	req := &CreateTagRequest{
		Name:   "Original Tag",
		Status: 1,
	}
	tag, _ := logic.CreateTag(context.Background(), req)

	// 更新标签
	updateReq := &UpdateTagRequest{
		ID:   tag.ID,
		Name: "Updated Tag",
	}
	updated, err := logic.UpdateTag(context.Background(), updateReq)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if updated.Name != "Updated Tag" {
		t.Errorf("Expected name 'Updated Tag', got '%s'", updated.Name)
	}

	// 测试更新不存在的标签
	invalidReq := &UpdateTagRequest{
		ID:   9999,
		Name: "Test",
	}
	_, err = logic.UpdateTag(context.Background(), invalidReq)
	if err == nil {
		t.Error("Expected error for non-existent tag")
	}
}

// Test DeleteTag
func TestDeleteTag(t *testing.T) {
	tagRepo := NewMockTagRepository()
	logic := NewTagLogic(tagRepo)

	// 创建标签
	req := &CreateTagRequest{
		Name:   "Test Tag",
		Status: 1,
	}
	tag, _ := logic.CreateTag(context.Background(), req)

	// 删除未使用的标签
	err := logic.DeleteTag(context.Background(), tag.ID)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// 验证标签已删除
	deleted, _ := tagRepo.Get(context.Background(), tag.ID)
	if deleted != nil {
		t.Error("Expected tag to be deleted")
	}

	// 测试删除不存在的标签
	err = logic.DeleteTag(context.Background(), 9999)
	if err == nil {
		t.Error("Expected error for non-existent tag")
	}
}

// Test IncrementTagUsage
func TestIncrementTagUsage(t *testing.T) {
	tagRepo := NewMockTagRepository()
	logic := NewTagLogic(tagRepo)

	// 创建标签
	req := &CreateTagRequest{
		Name:   "Test Tag",
		Status: 1,
	}
	tag, _ := logic.CreateTag(context.Background(), req)

	// 增加使用数
	err := logic.IncrementTagUsage(context.Background(), tag.ID)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// 验证使用数
	updated, _ := tagRepo.Get(context.Background(), tag.ID)
	if updated.UsageCount != 1 {
		t.Errorf("Expected usage count 1, got %d", updated.UsageCount)
	}
}

// Test ListTags
func TestListTags(t *testing.T) {
	tagRepo := NewMockTagRepository()
	logic := NewTagLogic(tagRepo)

	// 创建多个标签
	for i := 1; i <= 3; i++ {
		req := &CreateTagRequest{
			Name:   "Tag " + string(rune(48+i)),
			Status: 1,
		}
		logic.CreateTag(context.Background(), req)
	}

	// 获取标签列表
	tags, err := logic.ListTags(context.Background())
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if len(tags) != 3 {
		t.Errorf("Expected 3 tags, got %d", len(tags))
	}
}
