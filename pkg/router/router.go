package router

import (
	"github.com/gin-gonic/gin"
	email "github.com/jeffleon/email-service/pkg/email/infraestructure"
	"github.com/jeffleon/email-service/pkg/health"
	"github.com/jeffleon/email-service/pkg/swagger"
)

type Router interface {
	Run(addr ...string) error
}

func NewRouter(routes RoutesGroup) Router {
	route := gin.Default()
	public := route.Group("api/mail/v1/public")
	routes.Mail.PublicRoutes(public)
	routes.Health.RegisterRoutes(public)
	routes.Swagger.RegisterRoutes(public)
	return route
}

type RoutesGroup struct {
	Mail    *email.MailRoutes
	Health  *health.Routes
	Swagger *swagger.Routes
}
