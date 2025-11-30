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

// CmsUserListHandler 用户列表
func CmsUserListHandler(serverCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 调用Logic层
		users, err := serverCtx.CmsUserLogic.ListUsers(r.Context())
		if err != nil {
			logx.Errorf("list users failed: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"code":    500,
				"message": "获取用户列表失败: " + err.Error(),
				"data":    nil,
			})
			return
		}

		// 构建响应数据
		userInfos := make([]map[string]interface{}, 0)
		for _, user := range users {
			userInfos = append(userInfos, map[string]interface{}{
				"id":            user.ID,
				"username":      user.Username,
				"email":         user.Email,
				"phone":         user.Phone,
				"nickname":      user.Nickname,
				"avatar":        user.Avatar,
				"bio":           user.Bio,
				"gender":        user.Gender,
				"status":        user.Status,
				"emailVerified": user.EmailVerified,
				"phoneVerified": user.PhoneVerified,
				"lastLoginAt":   user.LastLoginAt,
				"lastLoginIp":   user.LastLoginIP,
				"loginCount":    user.LoginCount,
				"createdAt":     user.CreatedAt.Format("2006-01-02 15:04:05"),
				"updatedAt":     user.UpdatedAt.Format("2006-01-02 15:04:05"),
			})
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"code":    0,
			"message": "success",
			"data": map[string]interface{}{
				"total": len(users),
				"list":  userInfos,
			},
		})
	}
}

// CmsRegisterHandler 用户注册
func CmsRegisterHandler(serverCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.CmsRegisterReq
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
		logicReq := &cms.RegisterRequest{
			Username: req.Username,
			Email:    req.Email,
			Password: req.Password,
			Nickname: req.Nickname,
		}

		user, err := serverCtx.CmsUserLogic.Register(r.Context(), logicReq)
		if err != nil {
			logx.Errorf("register failed: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"code":    500,
				"message": "注册失败: " + err.Error(),
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
				"id":       user.ID,
				"username": user.Username,
				"email":    user.Email,
				"nickname": user.Nickname,
				"avatar":   user.Avatar,
				"status":   user.Status,
			},
		})
	}
}

// CmsLoginHandler 用户登录
func CmsLoginHandler(serverCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.CmsLoginReq
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
		logicReq := &cms.LoginRequest{
			Username: req.Username,
			Password: req.Password,
			IP:       req.Ip,
		}

		user, err := serverCtx.CmsUserLogic.Login(r.Context(), logicReq)
		if err != nil {
			logx.Errorf("login failed: %v", err)
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"code":    401,
				"message": "登录失败: " + err.Error(),
				"data":    nil,
			})
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"code":    0,
			"message": "success",
			"data": map[string]interface{}{
				"id":       user.ID,
				"username": user.Username,
				"email":    user.Email,
				"nickname": user.Nickname,
				"avatar":   user.Avatar,
			},
		})
	}
}

// GetCmsUserHandler 获取用户详情
func GetCmsUserHandler(serverCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
		if err != nil {
			logx.Errorf("parse path id failed: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"code":    400,
				"message": "用户ID格式错误",
				"data":    nil,
			})
			return
		}

		// 调用Logic层
		user, err := serverCtx.CmsUserLogic.GetUser(r.Context(), userID)
		if err != nil {
			logx.Errorf("get user failed: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"code":    500,
				"message": "获取用户失败: " + err.Error(),
				"data":    nil,
			})
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"code":    0,
			"message": "success",
			"data": map[string]interface{}{
				"id":            user.ID,
				"username":      user.Username,
				"email":         user.Email,
				"phone":         user.Phone,
				"nickname":      user.Nickname,
				"avatar":        user.Avatar,
				"bio":           user.Bio,
				"gender":        user.Gender,
				"status":        user.Status,
				"emailVerified": user.EmailVerified,
				"phoneVerified": user.PhoneVerified,
				"lastLoginAt":   user.LastLoginAt,
				"lastLoginIp":   user.LastLoginIP,
				"loginCount":    user.LoginCount,
				"createdAt":     user.CreatedAt.Format("2006-01-02 15:04:05"),
				"updatedAt":     user.UpdatedAt.Format("2006-01-02 15:04:05"),
			},
		})
	}
}

// UpdateCmsUserHandler 更新用户信息
func UpdateCmsUserHandler(serverCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
		if err != nil {
			logx.Errorf("parse path id failed: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"code":    400,
				"message": "用户ID格式错误",
				"data":    nil,
			})
			return
		}

		var req types.UpdateCmsUserReq
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
		logicReq := &cms.UpdateUserRequest{
			ID:       userID,
			Email:    req.Email,
			Nickname: req.Nickname,
			Avatar:   req.Avatar,
			Bio:      req.Bio,
			Gender:   int8(req.Gender),
			Phone:    req.Phone,
		}

		user, err := serverCtx.CmsUserLogic.UpdateUser(r.Context(), logicReq)
		if err != nil {
			logx.Errorf("update user failed: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"code":    500,
				"message": "更新用户失败: " + err.Error(),
				"data":    nil,
			})
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"code":    0,
			"message": "success",
			"data": map[string]interface{}{
				"id":            user.ID,
				"username":      user.Username,
				"email":         user.Email,
				"phone":         user.Phone,
				"nickname":      user.Nickname,
				"avatar":        user.Avatar,
				"bio":           user.Bio,
				"gender":        user.Gender,
				"status":        user.Status,
				"emailVerified": user.EmailVerified,
				"phoneVerified": user.PhoneVerified,
				"createdAt":     user.CreatedAt.Format("2006-01-02 15:04:05"),
				"updatedAt":     user.UpdatedAt.Format("2006-01-02 15:04:05"),
			},
		})
	}
}

// ChangePasswordHandler 修改密码
func ChangePasswordHandler(serverCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ChangePasswordReq
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
		logicReq := &cms.ChangePasswordRequest{
			UserID:      req.UserId,
			OldPassword: req.OldPassword,
			NewPassword: req.NewPassword,
		}

		if err := serverCtx.CmsUserLogic.ChangePassword(r.Context(), logicReq); err != nil {
			logx.Errorf("change password failed: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"code":    500,
				"message": "修改密码失败: " + err.Error(),
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

// DisableCmsUserHandler 禁用用户
func DisableCmsUserHandler(serverCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
		if err != nil {
			logx.Errorf("parse path id failed: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"code":    400,
				"message": "用户ID格式错误",
				"data":    nil,
			})
			return
		}

		// 调用Logic层
		if err := serverCtx.CmsUserLogic.DisableUser(r.Context(), userID); err != nil {
			logx.Errorf("disable user failed: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"code":    500,
				"message": "禁用用户失败: " + err.Error(),
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

// DeleteCmsUserHandler 删除用户
func DeleteCmsUserHandler(serverCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
		if err != nil {
			logx.Errorf("parse path id failed: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"code":    400,
				"message": "用户ID格式错误",
				"data":    nil,
			})
			return
		}

		// 调用Logic层
		if err := serverCtx.CmsUserLogic.DeleteUser(r.Context(), userID); err != nil {
			logx.Errorf("delete user failed: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"code":    500,
				"message": "删除用户失败: " + err.Error(),
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
