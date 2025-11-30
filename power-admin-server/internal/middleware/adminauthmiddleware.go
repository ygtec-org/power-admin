// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package middleware

import (
	"context"
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
		// 格式: 路由路径 -> {HTTP方法 -> true}
		whitelistRoutes := map[string]map[string]bool{
			"/api/admin/auth/login":    {"POST": true},
			"/api/admin/auth/register": {"POST": true},
			// 菜单查询接口 - 所有用户都能查看菜单，但修改需要权限
			"/api/admin/system/menus": {"GET": true},
			// 其他接口 - 仅查询开放，修改需要权限
			"/api/admin/system/roles":       {"GET": true},
			"/api/admin/system/users":       {"GET": true},
			"/api/admin/system/permissions": {"GET": true},
			"/api/admin/system/apis":        {"GET": true},
			"/api/admin/content/dicts":      {"GET": true},
		}

		// 检查当前路由和HTTP方法是否在白名单中
		isWhitelisted := false
		currentMethod := r.Method
		currentPath := r.URL.Path

		// 精确匹配：路由 + HTTP方法都匹配
		if methods, exists := whitelistRoutes[currentPath]; exists {
			if methods[currentMethod] {
				isWhitelisted = true
			}
		}

		// 如果精确匹配失败，检查是否为白名单路由的路径参数版本
		if !isWhitelisted {
			for whitelistRoute, methods := range whitelistRoutes {
				if strings.HasPrefix(currentPath, whitelistRoute) && whitelistRoute != currentPath {
					// 路径参数版本（如 /system/menus/123）也需要检查方法
					if methods[currentMethod] {
						isWhitelisted = true
						break
					}
				}
			}
		}

		if isWhitelisted {
			// 即使是白名单路由，也要验证JWT并存储用户信息到context
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

			// 将用户ID存储到context中（供logic层使用）
			ctx := context.WithValue(r.Context(), constant.AdminUserKey, fmt.Sprintf("%d", claims.ID))
			ctx = context.WithValue(ctx, constant.AdminUserName, claims.Username)
			r = r.WithContext(ctx)

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

		// 也将用户ID存储到context中（供logic层使用）
		ctx := context.WithValue(r.Context(), constant.AdminUserKey, fmt.Sprintf("%d", claims.ID))
		ctx = context.WithValue(ctx, constant.AdminUserName, claims.Username)
		r = r.WithContext(ctx)

		// 获取请求的资源（路由路径）和操作（HTTP方法）
		resource := r.URL.Path
		action := r.Method
		userID := fmt.Sprintf("%d", claims.ID)

		// 首先检查是否是超级管理员（role_id = 1 或 user_id = 1）
		// 如果是超级管理员，跳过所有权限检查，直接放行
		isSuperAdmin := false

		// 方式1：如果user_id = 1，则为超级管理员
		if claims.ID == 1 {
			isSuperAdmin = true
		}

		// 方式2：也检查用户的角色中是否包含role_id=1
		if !isSuperAdmin && m.Permission != nil {
			roles, err := m.Permission.GetRolesForUser(userID)
			if err == nil {
				for _, role := range roles {
					if role == "1" { // 超级管理员ID为1
						isSuperAdmin = true
						break
					}
				}
			}
		}

		// 如果是超级管理员，直接放行所有接口
		if isSuperAdmin {
			next(w, r)
			return
		}

		// 如果配置了权限管理且不是超级管理员，进行权限验证
		if m.Permission != nil {
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
		} else {
			// 如果权限管理未初始化，默认拒绝非白名单的路由
			// （超级管理员已经在前面被放行了）
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusForbidden)
			fmt.Fprintf(w, `{"code":403,"msg":"permission manager not initialized"}`)
			return
		}

		// 传递给下一个 handler
		next(w, r)
	}
}
