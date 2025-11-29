// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package user

import (
	"context"
	"errors"
	"power-admin-server/internal/svc"
	"power-admin-server/internal/types"
	"power-admin-server/pkg/auth"
	"power-admin-server/pkg/models"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateUserLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateUserLogic {
	return &CreateUserLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateUserLogic) CreateUser(req *types.CreateUserReq) (err error) {
	if req.Username == "" || req.Password == "" || req.Phone == "" {
		return errors.New("username, password and phone are required")
	}

	exists, err := l.svcCtx.UserRepo.IsUsernameExist(req.Username)
	if err != nil || exists {
		return errors.New("username exist")
	}

	exists, err = l.svcCtx.UserRepo.IsPhoneExist(req.Phone)
	if err != nil || exists {
		return errors.New("phone exist")
	}

	if req.Email != "" {
		exists, err = l.svcCtx.UserRepo.IsEmailExist(req.Email)
		if err != nil || exists {
			return errors.New("email exist")
		}
	}

	hashedPwd, err := auth.HashPassword(req.Password)
	if err != nil {
		return errors.New("hash password error")
	}

	user := &models.User{
		Username: req.Username,
		Phone:    req.Phone,
		Email:    req.Email,
		Password: hashedPwd,
		Nickname: req.Nickname,
		Gender:   req.Gender,
		Status:   1,
	}

	err = l.svcCtx.UserRepo.Create(user)
	if err != nil {
		return errors.New("create user error")
	}

	return nil
}
