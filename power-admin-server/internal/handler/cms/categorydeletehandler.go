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

// 删除分类
func CategoryDeleteHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.CategoryDeleteReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := cms.NewCategoryDeleteLogic(r.Context(), svcCtx)
		err := l.CategoryDelete(&req)
		if err != nil {
			response.Error(w, 500, err.Error())
		} else {
			response.Success(w, nil)
		}
	}
}
