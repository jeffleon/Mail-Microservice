package infra

import (
	"github.com/gin-gonic/gin"
)

type MailRoutes struct {
	handler MailHandler
}

func (ro *MailRoutes) PublicRoutes(public *gin.RouterGroup) {
	public.POST("/email", ro.handler.SendEmail)
}

func NewRoutes(handler MailHandler) *MailRoutes {
	return &MailRoutes{
		handler: handler,
	}
}
