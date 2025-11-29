package role

import (
	"net/http"

	"power-admin-server/common/response"
	"power-admin-server/internal/logic/role"
	"power-admin-server/internal/svc"
	"power-admin-server/internal/types"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func AssignRoleMenuAndApiHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.AssignRoleMenuAndApiReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := role.NewAssignRoleMenuAndApiLogic(r.Context(), svcCtx)
		resp, err := l.AssignRoleMenuAndApi(&req)
		if err != nil {
			response.Error(w, 500, err.Error())
		} else {
			response.Success(w, resp)
		}
	}
}
