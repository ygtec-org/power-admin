// 获取用户角色逻辑
package user

import (
	"context"
	"power-admin-server/common/constant"

	"power-admin-server/internal/svc"
	"power-admin-server/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserRolesLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetUserRolesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserRolesLogic {
	return &GetUserRolesLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetUserRolesLogic) GetUserRoles(req *types.GetUserRolesReq) (resp *types.GetUserRolesResp, err error) {
	userID := req.UserID

	if userID == 0 {
		userID = l.ctx.Value(constant.AdminUserKey).(int64)
	}
	// 获取该用户的所有角色
	roles, err := l.svcCtx.UserRepo.GetRoles(userID)
	if err != nil {
		l.Errorf("Failed to get roles for user %d: %v", userID, err)
		return nil, err
	}

	data := make([]types.RoleInfo, 0)
	if roles != nil {
		for _, r := range roles {
			data = append(data, types.RoleInfo{
				Id:          r.ID,
				Name:        r.Name,
				Description: r.Description,
				Status:      r.Status,
				CreatedAt:   r.CreatedAt.Format("2006-01-02 15:04:05"),
			})
		}
	}

	resp = &types.GetUserRolesResp{
		Data: data,
	}
	return
}
