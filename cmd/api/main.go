package main

import (
	"fmt"
	"time"

	"github.com/jeffleon/email-service/internal/config"
	appl "github.com/jeffleon/email-service/pkg/email/aplication"
	"github.com/jeffleon/email-service/pkg/email/domain"
	infra "github.com/jeffleon/email-service/pkg/email/infraestructure"
	"github.com/jeffleon/email-service/pkg/health"
	"github.com/jeffleon/email-service/pkg/router"
	"github.com/jeffleon/email-service/pkg/swagger"
	"github.com/sirupsen/logrus"
	mail "github.com/xhit/go-simple-mail/v2"
)

func main() {
	config.InitEnvConfigs()
	smtpClient, email, err := InitMail()
	if err != nil {
		logrus.Fatalf("Error initialization the mailer error: %s", err)
	}
	emailRepo := infra.NewMailer(smtpClient, email)
	emailService := appl.NewEmailService(emailRepo)
	emailHandler := infra.MailHandler{MailService: emailService}
	emailRoutes := infra.NewRoutes(emailHandler)
	r := router.NewRouter(router.RoutesGroup{
		Mail:    emailRoutes,
		Health:  health.NewHealthCheckRoutes(),
		Swagger: swagger.NewSwaggerDocsRoutes(),
	})
	logrus.Fatal(r.Run(fmt.Sprintf(":%s", config.EnvConfigs.AppPort)))
}

func InitMail() (*mail.SMTPClient, *domain.EMail, error) {
	server := mail.NewSMTPClient()
	server.Port = config.EnvConfigs.EmailPort
	server.Host = config.EnvConfigs.EmailHost
	server.Password = config.EnvConfigs.EmailPassword
	server.Username = config.EnvConfigs.EmailUserName
	server.Encryption = mail.EncryptionNone
	server.ConnectTimeout = 10 * time.Second
	server.SendTimeout = 10 * time.Second
	server.KeepAlive = true

	email := &domain.EMail{
		Domain:      config.EnvConfigs.EmailDomain,
		FromAddress: config.EnvConfigs.EmailFromAddress,
		FromName:    config.EnvConfigs.EmailFromName,
	}

	smtpClient, err := server.Connect()
	if err != nil {
		return nil, nil, err
	}
	return smtpClient, email, nil
}
