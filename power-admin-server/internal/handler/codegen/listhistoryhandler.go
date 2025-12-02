// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package codegen

import (
	"net/http"

	"power-admin-server/internal/logic/codegen"
	"power-admin-server/internal/svc"
	"power-admin-server/internal/types"

	"github.com/zeromicro/go-zero/rest/httpx"

	"power-admin-server/common/response"
)

func ListHistoryHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GenHistoryListReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := codegen.NewListHistoryLogic(r.Context(), svcCtx)
		resp, err := l.ListHistory(&req)
		if err != nil {
			response.Error(w, 500, err.Error())
		} else {
			response.Success(w, resp)
		}
	}
}
