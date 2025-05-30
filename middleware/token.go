package middleware

import (
	"StoreServer/utils/jwt"
	"StoreServer/utils/response"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func CheckLogin() gin.HandlerFunc {
	return func(c *gin.Context) {
		fmt.Println("===> CheckLogin middleware called")
		token := c.Request.Header.Get("Token")
		if token == "" || !strings.HasPrefix(token, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, response.Response{
				Code:    http.StatusUnauthorized,
				Message: "Unauthorized",
			})
			return
		}

		token = strings.TrimPrefix(token, "Bearer ")
		fmt.Println("token:", token)

		user, err := jwt.ParseToken(token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, response.Response{
				Code:    http.StatusUnauthorized,
				Message: "Unauthorized",
			})
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
			c.AbortWithStatusJSON(http.StatusUnauthorized, response.Response{
				Code:    http.StatusUnauthorized,
				Message: "Unauthorized",
			})
			return
		}

		userRole, ok := roleValue.(string)
		if !ok {
			c.AbortWithStatusJSON(http.StatusInternalServerError, response.Response{
				Code:    http.StatusInternalServerError,
				Message: "Role parsing failed",
			})
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
		c.AbortWithStatusJSON(http.StatusForbidden, response.Response{
			Code:    http.StatusForbidden,
			Message: "Permission denied",
		})
	}
}
