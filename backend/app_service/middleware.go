package app_service

import "github.com/gin-gonic/gin"

// AuthenticationMiddleware is a middleware function that checks for the presence of an authentication token in the request headers.
func (a *AppHandler) AuthenticationMiddleware(c *gin.Context) {
	clientToken := c.Request.Header.Get("AuthToken")
	if clientToken == "" {
		c.JSON(401, gin.H{"error": "Token is required"})
		c.Abort()
		return
	}
	claims, err := a.tokenHelper.ValidateToken(clientToken)
	if err != nil {
		c.JSON(401, gin.H{"error": "Invalid token"})
		c.Abort()
		return
	}
	c.Set("email", claims.Email)
	c.Set("username", claims.Username)
	c.Set("createdAt", claims.CreatedAt)
	c.Next()
}
