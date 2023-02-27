package config

import (
	"time"

	"github.com/jeffleon/email-service/pkg/email/domain"
	mail "github.com/xhit/go-simple-mail/v2"
)

func InitMail() (*mail.SMTPClient, *domain.EMail, error) {
	server := mail.NewSMTPClient()
	server.Port = EnvConfigs.EmailPort
	server.Host = EnvConfigs.EmailHost
	server.Password = EnvConfigs.EmailPassword
	server.Username = EnvConfigs.EmailUserName
	server.Encryption = mail.EncryptionNone
	server.ConnectTimeout = 10 * time.Second
	server.SendTimeout = 10 * time.Second
	server.KeepAlive = true

	email := &domain.EMail{
		Domain:      EnvConfigs.EmailDomain,
		FromAddress: EnvConfigs.EmailFromAddress,
		FromName:    EnvConfigs.EmailFromName,
	}

	smtpClient, err := server.Connect()
	if err != nil {
		return nil, nil, err
	}
	return smtpClient, email, nil
}
