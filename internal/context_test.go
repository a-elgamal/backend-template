package internal

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestUserFromGinContext(t *testing.T) {
	remoteAddr := "127.0.0.1"

	testCases := []struct {
		Name      string
		UserEmail interface{}
		UserID    interface{}
		Expected  string
	}{
		{Name: "UserEmail when email set", UserEmail: "email", UserID: "id", Expected: "email"},
		{Name: "UserId when email not set", UserEmail: nil, UserID: "id", Expected: "id"},
		{Name: "Remote IP when neither are set", UserEmail: nil, UserID: nil, Expected: remoteAddr},
	}

	for _, tt := range testCases {
		t.Run(tt.Name, func(t *testing.T) {
			r := gin.Default()
			r.Use(func(ctx *gin.Context) {
				if tt.UserEmail != nil {
					ctx.Set(UserEmailContextKey, tt.UserEmail)
				}
				if tt.UserID != nil {
					ctx.Set(UserIDContextKey, tt.UserID)
				}
			})

			r.GET("/", func(ctx *gin.Context) {
				ctx.JSON(200, gin.H{"user": UserFromGinContext(ctx)})
			})
			w := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodGet, "/", nil)
			req.RemoteAddr = remoteAddr + ":50000"
			r.ServeHTTP(w, req)

			assert.Equal(t, http.StatusOK, w.Code)
			assert.Equal(t, fmt.Sprintf("{\"user\":\"%v\"}", tt.Expected), w.Body.String())
		})
	}
}
