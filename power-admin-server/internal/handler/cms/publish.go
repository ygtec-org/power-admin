package cms

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/zeromicro/go-zero/core/logx"
	"power-admin-server/internal/logic/cms"
	"power-admin-server/internal/svc"
	"power-admin-server/internal/types"
)

// PublishImmediateHandler 立即发布内容
func PublishImmediateHandler(serverCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		contentID, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
		if err != nil {
			logx.Errorf("parse path id failed: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"code":    400,
				"message": "内容ID格式错误",
				"data":    nil,
			})
			return
		}

		// 调用Logic层
		logicReq := &cms.PublishImmediateRequest{
			ContentID: contentID,
		}

		if err := serverCtx.CmsPublishLogic.PublishImmediate(r.Context(), logicReq); err != nil {
			logx.Errorf("publish immediate failed: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"code":    500,
				"message": "立即发布失败: " + err.Error(),
				"data":    nil,
			})
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"code":    0,
			"message": "success",
			"data":    nil,
		})
	}
}

// PublishScheduledHandler 定时发布内容
func PublishScheduledHandler(serverCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		contentID, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
		if err != nil {
			logx.Errorf("parse path id failed: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"code":    400,
				"message": "内容ID格式错误",
				"data":    nil,
			})
			return
		}

		var req types.PublishScheduledReq
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			logx.Errorf("decode request failed: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"code":    400,
				"message": "请求参数格式错误",
				"data":    nil,
			})
			return
		}

		// 调用Logic层
		logicReq := &cms.PublishScheduledRequest{
			ContentID:   contentID,
			ScheduledAt: req.ScheduledAt,
		}

		if err := serverCtx.CmsPublishLogic.PublishScheduled(r.Context(), logicReq); err != nil {
			logx.Errorf("publish scheduled failed: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"code":    500,
				"message": "定时发布失败: " + err.Error(),
				"data":    nil,
			})
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"code":    0,
			"message": "success",
			"data":    nil,
		})
	}
}

// GetPublishStatusHandler 获取发布状态
func GetPublishStatusHandler(serverCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		contentID, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
		if err != nil {
			logx.Errorf("parse path id failed: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"code":    400,
				"message": "内容ID格式错误",
				"data":    nil,
			})
			return
		}

		// 调用Logic层获取内容
		content, err := serverCtx.CmsContentLogic.GetContent(r.Context(), contentID)
		if err != nil {
			logx.Errorf("get content failed: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"code":    500,
				"message": "获取内容失败: " + err.Error(),
				"data":    nil,
			})
			return
		}

		// 构建发布状态响应
		isPublished := content.Status == 2
		isDraft := content.Status == 1
		isScheduled := content.ScheduledAt != ""

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"code":    0,
			"message": "success",
			"data": map[string]interface{}{
				"contentId":   content.ID,
				"title":       content.Title,
				"status":      content.Status,
				"publishedAt": content.PublishedAt,
				"scheduledAt": content.ScheduledAt,
				"isScheduled": isScheduled,
				"isDraft":     isDraft,
				"isPublished": isPublished,
			},
		})
	}
}

// CancelScheduledPublishHandler 取消定时发布
func CancelScheduledPublishHandler(serverCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		contentID, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
		if err != nil {
			logx.Errorf("parse path id failed: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"code":    400,
				"message": "内容ID格式错误",
				"data":    nil,
			})
			return
		}

		// 调用Logic层
		logicReq := &cms.CancelScheduledPublishRequest{
			ContentID: contentID,
		}

		if err := serverCtx.CmsPublishLogic.CancelScheduledPublish(r.Context(), logicReq); err != nil {
			logx.Errorf("cancel scheduled publish failed: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"code":    500,
				"message": "取消定时发布失败: " + err.Error(),
				"data":    nil,
			})
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"code":    0,
			"message": "success",
			"data":    nil,
		})
	}
}

// BatchPublishHandler 批量发布内容
func BatchPublishHandler(serverCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.BatchPublishReq
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			logx.Errorf("decode request failed: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"code":    400,
				"message": "请求参数格式错误",
				"data":    nil,
			})
			return
		}

		// 参数验证
		if len(req.ContentIds) == 0 {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"code":    400,
				"message": "内容ID列表不能为空",
				"data":    nil,
			})
			return
		}

		// 调用Logic层
		logicReq := &cms.BatchPublishRequest{
			ContentIDs: req.ContentIds,
		}

		if err := serverCtx.CmsPublishLogic.BatchPublish(r.Context(), logicReq); err != nil {
			logx.Errorf("batch publish failed: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"code":    500,
				"message": "批量发布失败: " + err.Error(),
				"data":    nil,
			})
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"code":    0,
			"message": "success",
			"data":    nil,
		})
	}
}
