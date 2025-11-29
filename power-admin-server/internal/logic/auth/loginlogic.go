// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package auth

import (
	"context"
	"errors"

	"power-admin-server/internal/svc"
	"power-admin-server/internal/types"
	"power-admin-server/pkg/auth"

	"github.com/zeromicro/go-zero/core/logx"
	"golang.org/x/crypto/bcrypt"
)

type LoginLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LoginLogic) Login(req *types.LoginReq) (resp *types.LoginResp, err error) {
	logx.Infof("Login attempt with phone: %s", req.Phone)

	// 根据手机号查找用户
	user, err := l.svcCtx.UserRepo.GetByPhone(req.Phone)
	if err != nil {
		logx.Errorf("User not found: %v", err)
		return nil, errors.New("用户不存在")
	}

	logx.Infof("User found: %s, checking password...", user.Username)

	// 验证密码
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		logx.Errorf("Password mismatch: %v", err)
		return nil, errors.New("密码错误")
	}

	// 检查用户状态
	if user.Status != 1 {
		logx.Infof("User disabled: %s", user.Username)
		return nil, errors.New("用户已被禁用")
	}

	// 生成 JWT token
	// 将过期时间从秒转换为小时
	expirationHours := int(l.svcCtx.Config.Auth.AccessExpire / 3600)
	accessToken, err := auth.GenerateToken(user.ID, user.Username, user.Phone, expirationHours)
	if err != nil {
		logx.Errorf("Failed to generate token: %v", err)
		return nil, errors.New("生成token失败")
	}

	logx.Infof("Login successful for user: %s", user.Username)
	resp = &types.LoginResp{
		Token:    accessToken,
		UserId:   user.ID,
		Nickname: user.Nickname,
		Avatar:   user.Avatar,
	}
	return
}
