// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package cms

import (
	"context"

	"power-admin-server/internal/svc"
	"power-admin-server/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserGetLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取用户信息
func NewUserGetLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserGetLogic {
	return &UserGetLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserGetLogic) UserGet(req *types.UserGetReq) (resp *types.UserResp, err error) {
	user, err := l.svcCtx.CmsUserRepo.Get(l.ctx, req.Id)
	if err != nil {
		l.Logger.Errorf("获取用户详情失败: %v", err)
		return nil, err
	}
	if user == nil {
		l.Logger.Infof("用户不存在: id=%d", req.Id)
		return nil, nil
	}

	return &types.UserResp{
		Data: types.UserData{
			Id:        user.ID,
			Username:  user.Username,
			Email:     user.Email,
			Status:    int(user.Status),
			CreatedAt: user.CreatedAt.Format("2006-01-02 15:04:05"),
		},
	}, nil
}
