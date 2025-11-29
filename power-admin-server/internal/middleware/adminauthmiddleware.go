// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package middleware

import (
	"fmt"
	"net/http"
	"power-admin-server/common/constant"
	"strings"

	"power-admin-server/internal/config"
	"power-admin-server/pkg/auth"
	"power-admin-server/pkg/permission"
)

type AdminAuthMiddleware struct {
	AccessSecret string
	Permission   *permission.RBACEnforcer
}

func NewAdminAuthMiddleware(c *config.Config, rbac *permission.RBACEnforcer) *AdminAuthMiddleware {
	return &AdminAuthMiddleware{
		AccessSecret: c.Auth.AccessSecret,
		Permission:   rbac,
	}
}

func (m *AdminAuthMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 白名单路由，不需要权限验证
		// 这些路由对所有认证用户开放（只需要有效token）
		whitelistRoutes := map[string]bool{
			"/api/admin/auth/login":         true,
			"/api/admin/auth/register":      true,
			"/api/admin/system/menus":       true, // 菜单接口 - 所有用户都能查看菜单
			"/api/admin/system/roles":       true, // 角色接口
			"/api/admin/system/users":       true, // 用户接口
			"/api/admin/system/permissions": true, // 权限接口
			"/api/admin/system/apis":        true, // API接口
			"/api/admin/content/dicts":      true, // 字典接口
		}

		// 检查当前路由是否在白名单中
		// 对于带路径参数的路由（如/system/apis/:id），需要按前缀匹配
		isWhitelisted := false
		if whitelistRoutes[r.URL.Path] {
			isWhitelisted = true
		} else {
			// 检查是否为白名单路由的路径参数版本
			for whitelistRoute := range whitelistRoutes {
				if strings.HasPrefix(r.URL.Path, whitelistRoute) {
					isWhitelisted = true
					break
				}
			}
		}

		if isWhitelisted {
			next(w, r)
			return
		}

		// 从请求头中获取 Authorization token
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprintf(w, `{"code":401,"msg":"missing authorization header"}`)
			return
		}

		// 提取 Bearer token
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprintf(w, `{"code":401,"msg":"invalid authorization format"}`)
			return
		}

		token := parts[1]

		// 验证 JWT token
		claims, err := auth.ParseToken(token)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprintf(w, `{"code":401,"msg":"invalid or expired token"}`)
			return
		}

		// 将 claims 信息存储到请求上下文（供业务层使用）
		r.Header.Set(constant.AdminUserKey, fmt.Sprintf("%d", claims.ID))
		r.Header.Set(constant.AdminUserName, claims.Username)

		// 如果配置了权限管理，进行权限验证
		if m.Permission != nil {
			// 获取请求的资源（路由路径）和操作（HTTP方法）
			resource := r.URL.Path
			action := r.Method
			userID := fmt.Sprintf("%d", claims.ID)

			// 获取用户的所有角色
			roles, err := m.Permission.GetRolesForUser(userID)
			if err != nil {
				// 如果获取角色失败，拒绝访问
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusForbidden)
				fmt.Fprintf(w, `{"code":403,"msg":"permission denied"}`)
				return
			}

			// 检查用户或其任何角色是否有权访问该资源
			hasPermission := false

			// 首先检查用户是否直接拥有权限
			if m.Permission.CheckPermission(userID, resource, action) {
				hasPermission = true
			} else {
				// 检查用户的任何角色是否拥有权限
				for _, role := range roles {
					if m.Permission.CheckPermission(role, resource, action) {
						hasPermission = true
						break
					}
				}
			}

			if !hasPermission {
				// 权限不足，返回403 Forbidden
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusForbidden)
				fmt.Fprintf(w, `{"code":403,"msg":"permission denied"}`)
				return
			}
		}

		// 传递给下一个 handler
		next(w, r)
	}
}
