// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package cms

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

// 获取用户列表
func NewUserListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserListLogic {
	return &UserListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserListLogic) UserList() (resp *types.CmsUserListResp, err error) {
	users, err := l.svcCtx.CmsUserRepo.List(l.ctx)
	if err != nil {
		l.Logger.Errorf("查询用户列表失败: %v", err)
		return nil, err
	}

	var data []types.UserData
	for _, user := range users {
		data = append(data, types.UserData{
			Id:        user.ID,
			Username:  user.Username,
			Email:     user.Email,
			Status:    int(user.Status),
			CreatedAt: user.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	return &types.CmsUserListResp{
		Total: int64(len(data)),
		Data:  data,
	}, nil
}
