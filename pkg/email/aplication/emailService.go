package application

import "github.com/jeffleon/email-service/pkg/email/domain"

type emailService struct {
	emailRepo domain.MailerRepository
}

type EmailService interface {
	SendEmail(msg domain.Message) error
}

func NewEmailService(emailRepo domain.MailerRepository) EmailService {
	return &emailService{
		emailRepo,
	}
}

func (e *emailService) SendEmail(msg domain.Message) error {
	return e.emailRepo.SendEmail(msg)
}
