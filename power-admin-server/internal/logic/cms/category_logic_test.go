package cms

import (
	"context"
	"testing"
)

// Test CreateCategory
func TestCreateCategory(t *testing.T) {
	categoryRepo := NewMockCategoryRepository()
	contentRepo := NewMockContentRepository()
	logic := NewCategoryLogic(categoryRepo, contentRepo)

	// 测试成功创建分类
	req := &CreateCategoryRequest{
		Name:   "Test Category",
		Slug:   "test-category",
		Status: 1,
	}

	category, err := logic.CreateCategory(context.Background(), req)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if category == nil {
		t.Error("Expected category, got nil")
	}
	if category.Name != "Test Category" {
		t.Errorf("Expected name 'Test Category', got '%s'", category.Name)
	}

	// 测试空名称
	req2 := &CreateCategoryRequest{
		Slug: "test",
	}
	_, err = logic.CreateCategory(context.Background(), req2)
	if err == nil {
		t.Error("Expected error for empty name")
	}
}

// Test UpdateCategory
func TestUpdateCategory(t *testing.T) {
	categoryRepo := NewMockCategoryRepository()
	contentRepo := NewMockContentRepository()
	logic := NewCategoryLogic(categoryRepo, contentRepo)

	// 创建分类
	req := &CreateCategoryRequest{
		Name:   "Original Name",
		Slug:   "original-slug",
		Status: 1,
	}
	category, _ := logic.CreateCategory(context.Background(), req)

	// 更新分类
	updateReq := &UpdateCategoryRequest{
		ID:   category.ID,
		Name: "Updated Name",
		Slug: "updated-slug",
	}
	updated, err := logic.UpdateCategory(context.Background(), updateReq)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if updated.Name != "Updated Name" {
		t.Errorf("Expected name 'Updated Name', got '%s'", updated.Name)
	}

	// 测试更新不存在的分类
	invalidReq := &UpdateCategoryRequest{
		ID:   9999,
		Name: "Test",
	}
	_, err = logic.UpdateCategory(context.Background(), invalidReq)
	if err == nil {
		t.Error("Expected error for non-existent category")
	}

	// 测试自己作为自己的父分类
	circularReq := &UpdateCategoryRequest{
		ID:       category.ID,
		Name:     "Test",
		ParentID: &category.ID,
	}
	_, err = logic.UpdateCategory(context.Background(), circularReq)
	if err == nil {
		t.Error("Expected error for circular parent reference")
	}
}

// Test DeleteCategory
func TestDeleteCategory(t *testing.T) {
	categoryRepo := NewMockCategoryRepository()
	contentRepo := NewMockContentRepository()
	logic := NewCategoryLogic(categoryRepo, contentRepo)

	// 创建分类
	req := &CreateCategoryRequest{
		Name:   "Parent Category",
		Status: 1,
	}
	parentCat, _ := logic.CreateCategory(context.Background(), req)

	// 删除不存在的分类应该出错
	err := logic.DeleteCategory(context.Background(), 9999)
	if err == nil {
		t.Error("Expected error for non-existent category")
	}

	// 删除空分类应该成功
	err = logic.DeleteCategory(context.Background(), parentCat.ID)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// 验证分类已删除
	deleted, _ := categoryRepo.Get(context.Background(), parentCat.ID)
	if deleted != nil {
		t.Error("Expected category to be deleted")
	}
}

// Test GetCategoryTree
func TestGetCategoryTree(t *testing.T) {
	categoryRepo := NewMockCategoryRepository()
	contentRepo := NewMockContentRepository()
	logic := NewCategoryLogic(categoryRepo, contentRepo)

	// 创建多级分类
	req1 := &CreateCategoryRequest{
		Name:   "Root Category",
		Status: 1,
	}
	root, _ := logic.CreateCategory(context.Background(), req1)

	// 获取树形结构
	tree, err := logic.GetCategoryTree(context.Background())
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if len(tree) == 0 {
		t.Error("Expected tree with categories")
	}

	// 检查根分类是否在树中
	found := false
	for _, cat := range tree {
		if cat.ID == root.ID {
			found = true
			break
		}
	}
	if !found {
		t.Error("Root category not found in tree")
	}
}
