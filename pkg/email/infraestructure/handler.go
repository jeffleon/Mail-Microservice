package infra

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	application "github.com/jeffleon/email-service/pkg/email/aplication"
	"github.com/jeffleon/email-service/pkg/email/domain"
	"github.com/sirupsen/logrus"
)

type MailHandler struct {
	MailService application.EmailService
}

// 	SendEmail
//	@Tags			email
//	@Summary		Send Mail
//	@Description	Send Mail
//	@Accept			json
//	@Produce		json
//  @Param request body domain.Message true "query params"
//	@Success		200		{object}	domain.StandardResponse
//	@Failure		422		{object}	object{status=string,error=error}
//	@Router			/email [post]
func (h *MailHandler) SendEmail(ctx *gin.Context) {
	var requestBody domain.Message
	if err := ctx.BindJSON(&requestBody); err != nil {
		logrus.Errorf("Error making binding email message: %s", err.Error())
		ctx.JSON(http.StatusUnprocessableEntity, domain.StandardResponse{
			Status: "error",
			Error:  fmt.Sprintf("bad request %s", err),
		})
		return
	}
	if err := h.MailService.SendEmail(requestBody); err != nil {
		logrus.Errorf("Error sending mail: %s", err.Error())
		ctx.JSON(http.StatusUnprocessableEntity, domain.StandardResponse{
			Status: "error",
			Error:  fmt.Sprintf("bad request %s", err),
		})
		return
	}
	ctx.JSON(http.StatusOK, domain.StandardResponse{
		Status:   "OK",
		Error:    "",
		Data:     fmt.Sprintf("send email to %s", requestBody.To),
		DataType: "string",
	})
}
