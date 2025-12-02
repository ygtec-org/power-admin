// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package codegen

import (
	"context"

	"power-admin-server/internal/svc"
	"power-admin-server/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListHistoryLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListHistoryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListHistoryLogic {
	return &ListHistoryLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListHistoryLogic) ListHistory(req *types.GenHistoryListReq) (resp *types.GenHistoryListResp, err error) {
	histories, total, err := l.svcCtx.CodegenRepo.ListHistory(l.ctx, req.Page, req.PageSize, req.TableName)
	if err != nil {
		return nil, err
	}

	resp = &types.GenHistoryListResp{
		Total: total,
		Data:  make([]types.GenHistoryResp, 0, len(histories)),
	}

	for _, history := range histories {
		resp.Data = append(resp.Data, types.GenHistoryResp{
			ID:          history.ID,
			GenConfigID: history.GenConfigID,
			TableName:   history.Table,
			FilePath:    history.FilePath,
			FileType:    history.FileType,
			Content:     history.Content,
			Status:      history.Status,
			ErrorMsg:    history.ErrorMsg,
			Operator:    history.Operator,
			CreatedAt:   history.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	return resp, nil
}
