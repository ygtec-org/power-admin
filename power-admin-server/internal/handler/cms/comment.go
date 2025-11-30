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

// CommentListHandler 评论列表
func CommentListHandler(serverCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.CommentListReq

		if contentIdStr := r.URL.Query().Get("contentId"); contentIdStr != "" {
			if cid, err := strconv.ParseInt(contentIdStr, 10, 64); err == nil && cid > 0 {
				req.ContentId = cid
			}
		}

		if pageStr := r.URL.Query().Get("page"); pageStr != "" {
			if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
				req.Page = p
			}
		}
		if req.Page == 0 {
			req.Page = 1
		}

		if pageSizeStr := r.URL.Query().Get("pageSize"); pageSizeStr != "" {
			if ps, err := strconv.Atoi(pageSizeStr); err == nil && ps > 0 && ps <= 100 {
				req.PageSize = ps
			}
		}
		if req.PageSize == 0 {
			req.PageSize = 10
		}

		// 调用Logic层
		comments, err := serverCtx.CmsCommentLogic.ListComments(r.Context(), req.ContentId)
		if err != nil {
			logx.Errorf("list comments failed: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"code":    500,
				"message": "获取评论列表失败: " + err.Error(),
				"data":    nil,
			})
			return
		}

		// 构建响应数据
		commentInfos := make([]map[string]interface{}, 0)
		for _, comment := range comments {
			commentInfos = append(commentInfos, map[string]interface{}{
				"id":              comment.ID,
				"contentId":       comment.ContentID,
				"userId":          comment.UserID,
				"parentCommentId": comment.ParentCommentID,
				"authorName":      comment.AuthorName,
				"authorEmail":     comment.AuthorEmail,
				"content":         comment.Content,
				"status":          comment.Status,
				"likeCount":       comment.LikeCount,
				"replyCount":      comment.ReplyCount,
				"ipAddress":       comment.IPAddress,
				"userAgent":       comment.UserAgent,
				"createdAt":       comment.CreatedAt.Format("2006-01-02 15:04:05"),
				"updatedAt":       comment.UpdatedAt.Format("2006-01-02 15:04:05"),
				"approvedAt":      comment.ApprovedAt,
			})
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"code":    0,
			"message": "success",
			"data": map[string]interface{}{
				"total": len(comments),
				"list":  commentInfos,
			},
		})
	}
}

// CreateCommentHandler 创建评论
func CreateCommentHandler(serverCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.CreateCommentReq
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
		logicReq := &cms.CreateCommentRequest{
			ContentID:       req.ContentId,
			UserID:          req.UserId,
			ParentCommentID: req.ParentCommentId,
			AuthorName:      req.AuthorName,
			AuthorEmail:     req.AuthorEmail,
			Content:         req.Content,
			IPAddress:       req.IpAddress,
			UserAgent:       req.UserAgent,
		}

		comment, err := serverCtx.CmsCommentLogic.CreateComment(r.Context(), logicReq)
		if err != nil {
			logx.Errorf("create comment failed: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"code":    500,
				"message": "创建评论失败: " + err.Error(),
				"data":    nil,
			})
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"code":    0,
			"message": "success",
			"data": map[string]interface{}{
				"id":              comment.ID,
				"contentId":       comment.ContentID,
				"userId":          comment.UserID,
				"parentCommentId": comment.ParentCommentID,
				"authorName":      comment.AuthorName,
				"authorEmail":     comment.AuthorEmail,
				"content":         comment.Content,
				"status":          comment.Status,
				"likeCount":       comment.LikeCount,
				"replyCount":      comment.ReplyCount,
				"createdAt":       comment.CreatedAt.Format("2006-01-02 15:04:05"),
				"updatedAt":       comment.UpdatedAt.Format("2006-01-02 15:04:05"),
			},
		})
	}
}

// GetCommentHandler 获取评论详情
func GetCommentHandler(serverCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		commentID, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
		if err != nil {
			logx.Errorf("parse path id failed: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"code":    400,
				"message": "评论ID格式错误",
				"data":    nil,
			})
			return
		}

		// 调用Logic层
		comment, err := serverCtx.CmsCommentLogic.GetComment(r.Context(), commentID)
		if err != nil {
			logx.Errorf("get comment failed: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"code":    500,
				"message": "获取评论失败: " + err.Error(),
				"data":    nil,
			})
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"code":    0,
			"message": "success",
			"data": map[string]interface{}{
				"id":              comment.ID,
				"contentId":       comment.ContentID,
				"userId":          comment.UserID,
				"parentCommentId": comment.ParentCommentID,
				"authorName":      comment.AuthorName,
				"authorEmail":     comment.AuthorEmail,
				"content":         comment.Content,
				"status":          comment.Status,
				"likeCount":       comment.LikeCount,
				"replyCount":      comment.ReplyCount,
				"createdAt":       comment.CreatedAt.Format("2006-01-02 15:04:05"),
				"updatedAt":       comment.UpdatedAt.Format("2006-01-02 15:04:05"),
			},
		})
	}
}

// UpdateCommentHandler 更新评论
func UpdateCommentHandler(serverCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		commentID, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
		if err != nil {
			logx.Errorf("parse path id failed: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"code":    400,
				"message": "评论ID格式错误",
				"data":    nil,
			})
			return
		}

		var req types.UpdateCommentReq
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
		comment, err := serverCtx.CmsCommentLogic.UpdateComment(r.Context(), commentID, req.Content, int8(req.Status))
		if err != nil {
			logx.Errorf("update comment failed: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"code":    500,
				"message": "更新评论失败: " + err.Error(),
				"data":    nil,
			})
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"code":    0,
			"message": "success",
			"data": map[string]interface{}{
				"id":              comment.ID,
				"contentId":       comment.ContentID,
				"userId":          comment.UserID,
				"parentCommentId": comment.ParentCommentID,
				"authorName":      comment.AuthorName,
				"authorEmail":     comment.AuthorEmail,
				"content":         comment.Content,
				"status":          comment.Status,
				"likeCount":       comment.LikeCount,
				"replyCount":      comment.ReplyCount,
				"createdAt":       comment.CreatedAt.Format("2006-01-02 15:04:05"),
				"updatedAt":       comment.UpdatedAt.Format("2006-01-02 15:04:05"),
			},
		})
	}
}

// DeleteCommentHandler 删除评论
func DeleteCommentHandler(serverCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		commentID, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
		if err != nil {
			logx.Errorf("parse path id failed: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"code":    400,
				"message": "评论ID格式错误",
				"data":    nil,
			})
			return
		}

		// 调用Logic层
		if err := serverCtx.CmsCommentLogic.DeleteComment(r.Context(), commentID); err != nil {
			logx.Errorf("delete comment failed: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"code":    500,
				"message": "删除评论失败: " + err.Error(),
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

// ApproveCommentHandler 审核通过评论
func ApproveCommentHandler(serverCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		commentID, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
		if err != nil {
			logx.Errorf("parse path id failed: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"code":    400,
				"message": "评论ID格式错误",
				"data":    nil,
			})
			return
		}

		// 调用Logic层
		if err := serverCtx.CmsCommentLogic.ApproveComment(r.Context(), commentID); err != nil {
			logx.Errorf("approve comment failed: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"code":    500,
				"message": "审核评论失败: " + err.Error(),
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

// RejectCommentHandler 拒绝评论
func RejectCommentHandler(serverCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		commentID, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
		if err != nil {
			logx.Errorf("parse path id failed: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"code":    400,
				"message": "评论ID格式错误",
				"data":    nil,
			})
			return
		}

		// 调用Logic层
		if err := serverCtx.CmsCommentLogic.RejectComment(r.Context(), commentID); err != nil {
			logx.Errorf("reject comment failed: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"code":    500,
				"message": "拒绝评论失败: " + err.Error(),
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

// LikeCommentHandler 点赞评论
func LikeCommentHandler(serverCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		commentID, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
		if err != nil {
			logx.Errorf("parse path id failed: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"code":    400,
				"message": "评论ID格式错误",
				"data":    nil,
			})
			return
		}

		// 调用Logic层
		if err := serverCtx.CmsCommentLogic.LikeComment(r.Context(), commentID); err != nil {
			logx.Errorf("like comment failed: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"code":    500,
				"message": "点赞评论失败: " + err.Error(),
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
