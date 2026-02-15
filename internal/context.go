package internal

import "github.com/gin-gonic/gin"

// UserIDContextKey The key that is set in Gin's context that contains the authenticeted user ID
const UserIDContextKey = "x-user-id"

// UserEmailContextKey The key that is set in Gin's context that contains the authenticated user emailF
const UserEmailContextKey = "x-user-email"

// UserFromGinContext Fetches the user from the context
func UserFromGinContext(c *gin.Context) string {
	email, ok := c.Get(UserEmailContextKey)
	if ok {
		return email.(string)
	}
	id, ok := c.Get(UserIDContextKey)
	if ok {
		return id.(string)
	}
	return c.RemoteIP()
}
