package cms

import (
	"context"
	"errors"
	"power-admin-server/pkg/models"
	"power-admin-server/pkg/repository"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
	"golang.org/x/crypto/bcrypt"
)

// CmsUserLogic CMS访客用户管理业务逻辑
type CmsUserLogic struct {
	userRepo repository.CmsUserRepository
	logger   logx.Logger
}

// NewCmsUserLogic 创建CmsUserLogic实例
func NewCmsUserLogic(userRepo repository.CmsUserRepository) *CmsUserLogic {
	return &CmsUserLogic{
		userRepo: userRepo,
		logger:   logx.WithContext(context.Background()),
	}
}

// RegisterRequest 注册请求
type RegisterRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Nickname string `json:"nickname"`
}

// Register 用户注册
func (l *CmsUserLogic) Register(ctx context.Context, req *RegisterRequest) (*models.CmsUser, error) {
	// 验证必填字段
	if req.Username == "" {
		return nil, errors.New("用户名不能为空")
	}
	if req.Email == "" {
		return nil, errors.New("邮箱不能为空")
	}
	if req.Password == "" {
		return nil, errors.New("密码不能为空")
	}

	// 检查用户名是否已存在
	existing, err := l.userRepo.GetByUsername(ctx, req.Username)
	if err != nil {
		l.logger.Errorf("检查用户名失败: %v", err)
		return nil, err
	}
	if existing != nil {
		return nil, errors.New("用户名已存在")
	}

	// 检查邮箱是否已存在
	existing, err = l.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		l.logger.Errorf("检查邮箱失败: %v", err)
		return nil, err
	}
	if existing != nil {
		return nil, errors.New("邮箱已注册")
	}

	// 加密密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		l.logger.Errorf("密码加密失败: %v", err)
		return nil, err
	}

	user := &models.CmsUser{
		Username:  req.Username,
		Email:     req.Email,
		Password:  string(hashedPassword),
		Nickname:  req.Nickname,
		Status:    1, // 正常状态
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := l.userRepo.Create(ctx, user); err != nil {
		l.logger.Errorf("创建用户失败: %v", err)
		return nil, err
	}

	// 不返回密码
	user.Password = ""
	return user, nil
}

// LoginRequest 登录请求
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	IP       string `json:"ip"`
}

// Login 用户登录
func (l *CmsUserLogic) Login(ctx context.Context, req *LoginRequest) (*models.CmsUser, error) {
	// 验证必填字段
	if req.Username == "" {
		return nil, errors.New("用户名不能为空")
	}
	if req.Password == "" {
		return nil, errors.New("密码不能为空")
	}

	// 获取用户
	user, err := l.userRepo.GetByUsername(ctx, req.Username)
	if err != nil {
		l.logger.Errorf("获取用户失败: %v", err)
		return nil, err
	}
	if user == nil {
		return nil, errors.New("用户名或密码错误")
	}

	// 检查用户状态
	if user.Status != 1 {
		return nil, errors.New("用户已被禁用")
	}

	// 验证密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		l.logger.Errorf("密码验证失败: %v", err)
		return nil, errors.New("用户名或密码错误")
	}

	// 更新登录信息
	_ = l.userRepo.UpdateLastLogin(ctx, user.ID, req.IP)

	// 不返回密码
	user.Password = ""
	return user, nil
}

// UpdateUserRequest 更新用户信息请求
type UpdateUserRequest struct {
	ID       int64
	Email    string
	Nickname string
	Avatar   string
	Bio      string
	Gender   int8
	Phone    string
}

// UpdateUser 更新用户信息
func (l *CmsUserLogic) UpdateUser(ctx context.Context, req *UpdateUserRequest) (*models.CmsUser, error) {
	// 获取原用户
	existing, err := l.userRepo.Get(ctx, req.ID)
	if err != nil {
		l.logger.Errorf("获取用户失败: %v", err)
		return nil, err
	}
	if existing == nil {
		return nil, errors.New("用户不存在")
	}

	// 如果邮箱改变，检查新邮箱是否已存在
	if req.Email != "" && req.Email != existing.Email {
		emailExists, err := l.userRepo.GetByEmail(ctx, req.Email)
		if err != nil {
			l.logger.Errorf("检查邮箱失败: %v", err)
			return nil, err
		}
		if emailExists != nil {
			return nil, errors.New("邮箱已被其他用户使用")
		}
	}

	updateData := &models.CmsUser{
		Email:     req.Email,
		Nickname:  req.Nickname,
		Avatar:    req.Avatar,
		Bio:       req.Bio,
		Gender:    req.Gender,
		Phone:     req.Phone,
		UpdatedAt: time.Now(),
	}

	if err := l.userRepo.Update(ctx, req.ID, updateData); err != nil {
		l.logger.Errorf("更新用户失败: %v", err)
		return nil, err
	}

	user, _ := l.userRepo.Get(ctx, req.ID)
	user.Password = ""
	return user, nil
}

// ChangePasswordRequest 修改密码请求
type ChangePasswordRequest struct {
	UserID      int64
	OldPassword string
	NewPassword string
}

// ChangePassword 修改用户密码
func (l *CmsUserLogic) ChangePassword(ctx context.Context, req *ChangePasswordRequest) error {
	// 获取用户
	user, err := l.userRepo.Get(ctx, req.UserID)
	if err != nil {
		l.logger.Errorf("获取用户失败: %v", err)
		return err
	}
	if user == nil {
		return errors.New("用户不存在")
	}

	// 验证旧密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.OldPassword)); err != nil {
		return errors.New("旧密码错误")
	}

	// 加密新密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		l.logger.Errorf("密码加密失败: %v", err)
		return err
	}

	// 更新密码
	updateData := &models.CmsUser{
		Password:  string(hashedPassword),
		UpdatedAt: time.Now(),
	}

	if err := l.userRepo.Update(ctx, req.UserID, updateData); err != nil {
		l.logger.Errorf("更新密码失败: %v", err)
		return err
	}

	return nil
}

// GetUser 获取用户详情
func (l *CmsUserLogic) GetUser(ctx context.Context, userID int64) (*models.CmsUser, error) {
	user, err := l.userRepo.Get(ctx, userID)
	if err != nil {
		l.logger.Errorf("获取用户失败: %v", err)
		return nil, err
	}
	if user == nil {
		return nil, errors.New("用户不存在")
	}

	user.Password = ""
	return user, nil
}

// ListUsers 获取用户列表
func (l *CmsUserLogic) ListUsers(ctx context.Context) ([]*models.CmsUser, error) {
	users, err := l.userRepo.List(ctx)
	if err != nil {
		l.logger.Errorf("获取用户列表失败: %v", err)
		return nil, err
	}

	// 不返回密码
	for _, user := range users {
		user.Password = ""
	}

	return users, nil
}

// DisableUser 禁用用户
func (l *CmsUserLogic) DisableUser(ctx context.Context, userID int64) error {
	// 获取用户
	user, err := l.userRepo.Get(ctx, userID)
	if err != nil {
		l.logger.Errorf("获取用户失败: %v", err)
		return err
	}
	if user == nil {
		return errors.New("用户不存在")
	}

	updateData := &models.CmsUser{
		Status:    0, // 0:已禁用
		UpdatedAt: time.Now(),
	}

	if err := l.userRepo.Update(ctx, userID, updateData); err != nil {
		l.logger.Errorf("禁用用户失败: %v", err)
		return err
	}

	return nil
}

// EnableUser 启用用户
func (l *CmsUserLogic) EnableUser(ctx context.Context, userID int64) error {
	// 获取用户
	user, err := l.userRepo.Get(ctx, userID)
	if err != nil {
		l.logger.Errorf("获取用户失败: %v", err)
		return err
	}
	if user == nil {
		return errors.New("用户不存在")
	}

	updateData := &models.CmsUser{
		Status:    1, // 1:正常
		UpdatedAt: time.Now(),
	}

	if err := l.userRepo.Update(ctx, userID, updateData); err != nil {
		l.logger.Errorf("启用用户失败: %v", err)
		return err
	}

	return nil
}

// DeleteUser 删除用户（软删除）
func (l *CmsUserLogic) DeleteUser(ctx context.Context, userID int64) error {
	// 获取用户
	user, err := l.userRepo.Get(ctx, userID)
	if err != nil {
		l.logger.Errorf("获取用户失败: %v", err)
		return err
	}
	if user == nil {
		return errors.New("用户不存在")
	}

	if err := l.userRepo.SoftDelete(ctx, userID); err != nil {
		l.logger.Errorf("删除用户失败: %v", err)
		return err
	}

	return nil
}

// HardDeleteUser 永久删除用户
func (l *CmsUserLogic) HardDeleteUser(ctx context.Context, userID int64) error {
	// 获取用户
	user, err := l.userRepo.Get(ctx, userID)
	if err != nil {
		l.logger.Errorf("获取用户失败: %v", err)
		return err
	}
	if user == nil {
		return errors.New("用户不存在")
	}

	if err := l.userRepo.Delete(ctx, userID); err != nil {
		l.logger.Errorf("永久删除用户失败: %v", err)
		return err
	}

	return nil
}

// VerifyEmail 验证邮箱
func (l *CmsUserLogic) VerifyEmail(ctx context.Context, userID int64) error {
	// 获取用户
	user, err := l.userRepo.Get(ctx, userID)
	if err != nil {
		l.logger.Errorf("获取用户失败: %v", err)
		return err
	}
	if user == nil {
		return errors.New("用户不存在")
	}

	now := time.Now()
	updateData := &models.CmsUser{
		EmailVerified:   1,
		EmailVerifiedAt: &now,
		UpdatedAt:       time.Now(),
	}

	if err := l.userRepo.Update(ctx, userID, updateData); err != nil {
		l.logger.Errorf("验证邮箱失败: %v", err)
		return err
	}

	return nil
}

// VerifyPhone 验证手机
func (l *CmsUserLogic) VerifyPhone(ctx context.Context, userID int64) error {
	// 获取用户
	user, err := l.userRepo.Get(ctx, userID)
	if err != nil {
		l.logger.Errorf("获取用户失败: %v", err)
		return err
	}
	if user == nil {
		return errors.New("用户不存在")
	}

	now := time.Now()
	updateData := &models.CmsUser{
		PhoneVerified:   1,
		PhoneVerifiedAt: &now,
		UpdatedAt:       time.Now(),
	}

	if err := l.userRepo.Update(ctx, userID, updateData); err != nil {
		l.logger.Errorf("验证手机失败: %v", err)
		return err
	}

	return nil
}
