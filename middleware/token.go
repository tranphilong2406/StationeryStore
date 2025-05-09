package middleware

import (
	"StoreServer/utils/jwt"
	"StoreServer/utils/response"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func CheckLogin() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Request.Header.Get("Authorization")
		if token == "" || !strings.HasPrefix(token, "Bearer ") {
			c.JSON(http.StatusUnauthorized, response.Response{
				Code:    http.StatusUnauthorized,
				Message: "Unauthorized",
			})
			c.Abort()
			return
		}

		token = strings.TrimPrefix(token, "Bearer ")

		user, err := jwt.ParseToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, response.Response{
				Code:    http.StatusUnauthorized,
				Message: "Unauthorized",
			})
			c.Abort()
			return
		}

		c.Set("user", user)
		c.Set("role", user.Role)
		c.Next()
	}
}

func CheckRole(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		roleValue, exists := c.Get("role")
		if !exists {
			c.JSON(http.StatusUnauthorized, response.Response{
				Code:    http.StatusUnauthorized,
				Message: "Unauthorized",
			})
			c.Abort()
			return
		}

		userRole, ok := roleValue.(string)
		if !ok {
			c.JSON(http.StatusInternalServerError, response.Response{
				Code:    http.StatusInternalServerError,
				Message: "Role parsing failed",
			})
			c.Abort()
			return
		}

		// Kiểm tra role có nằm trong danh sách allowedRoles không
		for _, role := range allowedRoles {
			if userRole == role {
				c.Next()
				return
			}
		}

		// Nếu không có role phù hợp
		c.JSON(http.StatusForbidden, response.Response{
			Code:    http.StatusForbidden,
			Message: "Permission denied",
		})
		c.Abort()
	}
}
