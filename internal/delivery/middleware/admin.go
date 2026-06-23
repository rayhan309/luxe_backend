package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/luxe/backend/internal/domain"
	"github.com/luxe/backend/internal/pkg/response"
)

// AdminOnly restricts access to admin role users only
func AdminOnly() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("user_role")
		if !exists {
			response.Unauthorized(c, "authentication required")
			c.Abort()
			return
		}

		if role.(string) != domain.RoleAdmin {
			response.Forbidden(c, "admin access required")
			c.Abort()
			return
		}

		c.Next()
	}
}
