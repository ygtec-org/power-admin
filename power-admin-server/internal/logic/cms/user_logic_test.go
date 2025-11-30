package cms

import (
	"context"
	"errors"
	"power-admin-server/pkg/models"
	"testing"
	"time"
)

// MockCmsUserRepository CMS用户仓储的模拟实现
type MockCmsUserRepository struct {
	users  map[int64]*models.CmsUser
	lastID int64
}

func NewMockCmsUserRepository() *MockCmsUserRepository {
	return &MockCmsUserRepository{
		users:  make(map[int64]*models.CmsUser),
		lastID: 0,
	}
}

func (m *MockCmsUserRepository) Create(ctx context.Context, user *models.CmsUser) error {
	m.lastID++
	user.ID = m.lastID
	// 剋一事阐本 user 对象，以防止指针护会影响存储的數据
	newUser := *user
	m.users[newUser.ID] = &newUser
	return nil
}

func (m *MockCmsUserRepository) Get(ctx context.Context, id int64) (*models.CmsUser, error) {
	if user, exists := m.users[id]; exists {
		return user, nil
	}
	return nil, nil
}

func (m *MockCmsUserRepository) Update(ctx context.Context, id int64, user *models.CmsUser) error {
	if existing, exists := m.users[id]; exists {
		if user.Email != "" {
			existing.Email = user.Email
		}
		if user.Nickname != "" {
			existing.Nickname = user.Nickname
		}
		if user.Password != "" && user.Password != existing.Password {
			existing.Password = user.Password
		}
		// Status 可以是 0（禁用）或 1（正常），所以需要通过其他方式检页是否填换2
		if user.Status >= 0 && user.Status <= 1 {
			existing.Status = user.Status
		}
		existing.UpdatedAt = user.UpdatedAt
		return nil
	}
	return errors.New("user not found")
}

func (m *MockCmsUserRepository) Delete(ctx context.Context, id int64) error {
	if _, exists := m.users[id]; exists {
		delete(m.users, id)
		return nil
	}
	return errors.New("user not found")
}

func (m *MockCmsUserRepository) SoftDelete(ctx context.Context, id int64) error {
	if existing, exists := m.users[id]; exists {
		now := time.Now()
		existing.DeletedAt = &now
		return nil
	}
	return errors.New("user not found")
}

func (m *MockCmsUserRepository) GetByUsername(ctx context.Context, username string) (*models.CmsUser, error) {
	for _, user := range m.users {
		if user.Username == username && user.DeletedAt == nil {
			return user, nil
		}
	}
	return nil, nil
}

func (m *MockCmsUserRepository) GetByEmail(ctx context.Context, email string) (*models.CmsUser, error) {
	for _, user := range m.users {
		if user.Email == email && user.DeletedAt == nil {
			return user, nil
		}
	}
	return nil, nil
}

func (m *MockCmsUserRepository) List(ctx context.Context) ([]*models.CmsUser, error) {
	var users []*models.CmsUser
	for _, user := range m.users {
		if user.DeletedAt == nil {
			users = append(users, user)
		}
	}
	return users, nil
}

func (m *MockCmsUserRepository) UpdateLastLogin(ctx context.Context, id int64, ip string) error {
	if existing, exists := m.users[id]; exists {
		now := time.Now()
		existing.LastLoginAt = &now
		existing.LastLoginIP = ip
		existing.LoginCount++
		return nil
	}
	return errors.New("user not found")
}

// Test Register
func TestRegister(t *testing.T) {
	userRepo := NewMockCmsUserRepository()
	logic := NewCmsUserLogic(userRepo)

	// 测试成功注册
	req := &RegisterRequest{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "password123",
		Nickname: "Test User",
	}

	user, err := logic.Register(context.Background(), req)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if user == nil {
		t.Error("Expected user, got nil")
	}
	if user.Username != "testuser" {
		t.Errorf("Expected username 'testuser', got '%s'", user.Username)
	}
	if user.Password != "" {
		t.Error("Expected password to be empty in response")
	}

	// 测试空用户名
	req2 := &RegisterRequest{
		Email:    "test2@example.com",
		Password: "password123",
	}
	_, err = logic.Register(context.Background(), req2)
	if err == nil {
		t.Error("Expected error for empty username")
	}

	// 测试重复用户名
	req3 := &RegisterRequest{
		Username: "testuser",
		Email:    "test2@example.com",
		Password: "password123",
	}
	_, err = logic.Register(context.Background(), req3)
	if err == nil {
		t.Error("Expected error for duplicate username")
	}

	// 测试重复邮箱
	req4 := &RegisterRequest{
		Username: "testuser2",
		Email:    "test@example.com",
		Password: "password123",
	}
	_, err = logic.Register(context.Background(), req4)
	if err == nil {
		t.Error("Expected error for duplicate email")
	}
}

// Test Login
func TestLogin(t *testing.T) {
	userRepo := NewMockCmsUserRepository()
	logic := NewCmsUserLogic(userRepo)

	// 注册用户
	regReq := &RegisterRequest{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "password123",
		Nickname: "Test User",
	}
	_, _ = logic.Register(context.Background(), regReq)

	// 测试成功登录
	loginReq := &LoginRequest{
		Username: "testuser",
		Password: "password123",
		IP:       "127.0.0.1",
	}

	loggedInUser, err := logic.Login(context.Background(), loginReq)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if loggedInUser == nil {
		t.Error("Expected user, got nil")
	} else {
		if loggedInUser.Password != "" {
			t.Error("Expected password to be empty in response")
		}
	}

	// 验证登录新和IP
	if loggedInUser != nil {
		updatedUser, _ := userRepo.Get(context.Background(), loggedInUser.ID)
		if updatedUser.LoginCount != 1 {
			t.Errorf("Expected login count 1, got %d", updatedUser.LoginCount)
		}
	}

	// 测试错误密码
	loginReq2 := &LoginRequest{
		Username: "testuser",
		Password: "wrongpassword",
		IP:       "127.0.0.1",
	}
	_, err = logic.Login(context.Background(), loginReq2)
	if err == nil {
		t.Error("Expected error for wrong password")
	}

	// 测试不存在的用户
	loginReq3 := &LoginRequest{
		Username: "nonexistent",
		Password: "password123",
		IP:       "127.0.0.1",
	}
	_, err = logic.Login(context.Background(), loginReq3)
	if err == nil {
		t.Error("Expected error for non-existent user")
	}
}

// Test UpdateUser
func TestUpdateUser(t *testing.T) {
	userRepo := NewMockCmsUserRepository()
	logic := NewCmsUserLogic(userRepo)

	// 注册用户
	regReq := &RegisterRequest{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "password123",
		Nickname: "Test User",
	}
	user, _ := logic.Register(context.Background(), regReq)

	// 更新用户
	updateReq := &UpdateUserRequest{
		ID:       user.ID,
		Email:    "new@example.com",
		Nickname: "New Name",
		Avatar:   "avatar.jpg",
	}

	updated, err := logic.UpdateUser(context.Background(), updateReq)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if updated.Nickname != "New Name" {
		t.Errorf("Expected nickname 'New Name', got '%s'", updated.Nickname)
	}

	// 测试更新不存在的用户
	invalidReq := &UpdateUserRequest{
		ID: 9999,
	}
	_, err = logic.UpdateUser(context.Background(), invalidReq)
	if err == nil {
		t.Error("Expected error for non-existent user")
	}
}

// Test DisableUser
func TestDisableUser(t *testing.T) {
	userRepo := NewMockCmsUserRepository()
	logic := NewCmsUserLogic(userRepo)

	// 注册用户
	regReq := &RegisterRequest{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "password123",
	}
	user, _ := logic.Register(context.Background(), regReq)

	// 禁用用户
	err := logic.DisableUser(context.Background(), user.ID)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// 验证状态
	disabled, _ := userRepo.Get(context.Background(), user.ID)
	if disabled.Status != 0 {
		t.Errorf("Expected status 0, got %d", disabled.Status)
	}

	// 验证禁用后无法登录
	loginReq := &LoginRequest{
		Username: "testuser",
		Password: "password123",
		IP:       "127.0.0.1",
	}
	_, err = logic.Login(context.Background(), loginReq)
	if err == nil {
		t.Error("Expected error for disabled user")
	}
}

// Test DeleteUser
func TestDeleteUser(t *testing.T) {
	userRepo := NewMockCmsUserRepository()
	logic := NewCmsUserLogic(userRepo)

	// 注册用户
	regReq := &RegisterRequest{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "password123",
	}
	user, _ := logic.Register(context.Background(), regReq)

	// 软删除用户
	err := logic.DeleteUser(context.Background(), user.ID)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// 验证用户被标记为删除
	deleted, _ := userRepo.Get(context.Background(), user.ID)
	if deleted.DeletedAt == nil {
		t.Error("Expected deleted_at to be set")
	}

	// 验证删除后无法通过用户名获取
	_, err = userRepo.GetByUsername(context.Background(), "testuser")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}
