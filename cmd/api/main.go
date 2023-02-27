package main

// @title     Mail API
// @version   1.0

// @host      localhost:8081
// @BasePath  /api/mail/v1

import (
	"fmt"
	"log"
	"net"
	"net/rpc"

	"github.com/jeffleon/email-service/internal/config"
	appl "github.com/jeffleon/email-service/pkg/email/aplication"
	infra "github.com/jeffleon/email-service/pkg/email/infraestructure"
	"github.com/jeffleon/email-service/pkg/health"
	"github.com/jeffleon/email-service/pkg/router"
	"github.com/jeffleon/email-service/pkg/swagger"
	"github.com/sirupsen/logrus"
)

func main() {
	config.InitEnvConfigs()
	smtpClient, email, err := config.InitMail()
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
	// Register the RPC server
	err = rpc.Register(new(infra.Mailer))
	if err != nil {
		logrus.Fatalf("Error initialization RPC server: %s", err)
	}
	go rpcListen()
	logrus.Fatal(r.Run(fmt.Sprintf(":%s", config.EnvConfigs.AppPort)))
}

func rpcListen() error {
	log.Println("Starting RPC server on port", config.EnvConfigs.RPCPort)
	listen, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%s", config.EnvConfigs.RPCPort))
	if err != nil {
		return err
	}
	defer listen.Close()
	for {
		rpcConn, err := listen.Accept()
		if err != nil {
			continue
		}
		go rpc.ServeConn(rpcConn)
	}
}
