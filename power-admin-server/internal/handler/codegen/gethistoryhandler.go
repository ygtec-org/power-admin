// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package codegen

import (
	"net/http"

	"power-admin-server/internal/logic/codegen"
	"power-admin-server/internal/svc"

	"power-admin-server/common/response"
)

func GetHistoryHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := codegen.NewGetHistoryLogic(r.Context(), svcCtx)
		resp, err := l.GetHistory()
		if err != nil {
			response.Error(w, 500, err.Error())
		} else {
			response.Success(w, resp)
		}
	}
}
