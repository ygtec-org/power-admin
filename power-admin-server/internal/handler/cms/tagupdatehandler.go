// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package cms

import (
	"net/http"

	"power-admin-server/internal/logic/cms"
	"power-admin-server/internal/svc"
	"power-admin-server/internal/types"

	"github.com/zeromicro/go-zero/rest/httpx"

	"power-admin-server/common/response"
)

// 更新标签
func TagUpdateHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.TagUpdateReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := cms.NewTagUpdateLogic(r.Context(), svcCtx)
		err := l.TagUpdate(&req)
		if err != nil {
			response.Error(w, 500, err.Error())
		} else {
			response.Success(w, nil)
		}
	}
}
