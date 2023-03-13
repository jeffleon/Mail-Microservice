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

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/jeffleon/email-service/internal/config"
	appl "github.com/jeffleon/email-service/pkg/email/aplication"
	"github.com/jeffleon/email-service/pkg/email/domain"
	infra "github.com/jeffleon/email-service/pkg/email/infraestructure"
	"github.com/jeffleon/email-service/pkg/email/infraestructure/repository"
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

	emailRepo := repository.NewMailer(smtpClient, email)
	emailService := appl.NewEmailService(emailRepo)

	InitKafkaConsumer(emailService.SendEmail)
	logrus.Info("Kafka connected")

	emailHandler := infra.MailHandler{MailService: emailService}
	emailRoutes := infra.NewRoutes(emailHandler)
	r := router.NewRouter(router.RoutesGroup{
		Mail:    emailRoutes,
		Health:  health.NewHealthCheckRoutes(),
		Swagger: swagger.NewSwaggerDocsRoutes(),
	})
	// Register the RPC server
	err = rpc.Register(new(repository.Mailer))
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

func InitKafkaConsumer(sendMail domain.SendMailer) {
	fmt.Println(config.EnvConfigs)
	fmt.Printf("%s %s", config.EnvConfigs.KafkaUsername, config.EnvConfigs.KafkaPassword)
	fmt.Printf("%s %s", config.EnvConfigs.KafkaHost, config.EnvConfigs.KafkaPort)
	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": fmt.Sprintf("%s:%s", config.EnvConfigs.KafkaHost, config.EnvConfigs.KafkaPort),
		"security.protocol": "SASL_SSL",
		"sasl.username":     config.EnvConfigs.KafkaUsername,
		"sasl.password":     config.EnvConfigs.KafkaPassword,
		"sasl.mechanism":    "PLAIN",
		"group.id":          "myGroup",
		"auto.offset.reset": "earliest",
	})

	if err != nil {
		panic(err)
	}

	kafkaRepository := repository.NewKafkaRepository(consumer, sendMail)
	err = kafkaRepository.TopicConsume()

	if err != nil {
		panic(err)
	}

}
