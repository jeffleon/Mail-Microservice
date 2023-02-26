package health

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Routes struct{}

func NewHealthCheckRoutes() *Routes {
	return &Routes{}
}

func (r *Routes) RegisterRoutes(group *gin.RouterGroup) {
	group.GET("/health", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, map[string]interface{}{"message": "pong"})
	})
}
