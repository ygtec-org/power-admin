// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package user

import (
	"context"

	"power-admin-server/internal/svc"
	"power-admin-server/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserListLogic {
	return &UserListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserListLogic) UserList(req *types.UserListReq) (resp *types.UserListResp, err error) {
	offset := (req.Page - 1) * req.PageSize
	if offset < 0 {
		offset = 0
	}

	users, total, err := l.svcCtx.UserRepo.List(offset, req.PageSize)
	if err != nil {
		return nil, err
	}

	data := make([]types.UserInfo, 0, len(users))
	for _, u := range users {
		data = append(data, types.UserInfo{
			Id:        u.ID,
			Username:  u.Username,
			Phone:     u.Phone,
			Email:     u.Email,
			Nickname:  u.Nickname,
			Avatar:    u.Avatar,
			Gender:    u.Gender,
			Status:    u.Status,
			CreatedAt: u.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	resp = &types.UserListResp{
		Total: total,
		Data:  data,
	}
	return
}
