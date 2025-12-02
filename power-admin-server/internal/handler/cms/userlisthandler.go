// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package cms

import (
	"net/http"

	"power-admin-server/internal/logic/cms"
	"power-admin-server/internal/svc"

	"power-admin-server/common/response"
)

// 获取用户列表
func UserListHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := cms.NewUserListLogic(r.Context(), svcCtx)
		resp, err := l.UserList()
		if err != nil {
			response.Error(w, 500, err.Error())
		} else {
			response.Success(w, resp)
		}
	}
}
