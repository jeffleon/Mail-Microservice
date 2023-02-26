package infra

import (
	"bytes"
	"html/template"
	"os"

	"github.com/jeffleon/email-service/pkg/email/domain"
	"github.com/sirupsen/logrus"
	"github.com/vanng822/go-premailer/premailer"
	mail "github.com/xhit/go-simple-mail/v2"
)

type Mailer struct {
	smtpClient *mail.SMTPClient
	email      *domain.EMail
}

func NewMailer(smtpClient *mail.SMTPClient, email *domain.EMail) domain.MailerRepository {
	return &Mailer{
		smtpClient,
		email,
	}
}

func (m *Mailer) buildHTMLMessage(msg domain.Message) (string, error) {
	templateToRender := "./templates/mail.html.gohtml"
	// Parse the go html to template
	t, err := template.New("email-html").ParseFiles(templateToRender)
	if err != nil {
		return "", nil
	}
	var tpl bytes.Buffer

	// Excute template pass the variables to template and transform in bytes buffer
	if err = t.ExecuteTemplate(&tpl, "body", msg); err != nil {
		return "", nil
	}

	// Parse bytes buffer into string
	formattedMessage := tpl.String()

	// Use premailer to transform css for the email
	formattedMessage, err = m.inlineCSS(formattedMessage)
	if err != nil {
		return "", err
	}

	return formattedMessage, nil
}

func (m *Mailer) inlineCSS(s string) (string, error) {
	options := premailer.Options{
		RemoveClasses:     false,
		CssToAttributes:   false,
		KeepBangImportant: true,
	}

	prem, err := premailer.NewPremailerFromString(s, &options)
	if err != nil {
		return "", err
	}

	html, err := prem.Transform()
	if err != nil {
		return "", err
	}

	return html, nil
}

func (m *Mailer) SendEmail(msg domain.Message) error {

	htmlBody, err := m.buildHTMLMessage(msg)
	if err != nil {
		logrus.Errorf("Error whith creation mail body : %s", err)
		return err
	}
	textbody, err := m.buildPlainTextMessage(msg)
	if err != nil {
		logrus.Errorf("Error whith creation mail body : %s", err)
		return err
	}

	email := mail.NewMSG()
	email.SetFrom(msg.From).
		AddTo(msg.To).
		SetSubject(msg.Subject)

	newImage, _ := os.ReadFile("./attachment/example.pdf")
	email.SetBody(mail.TextHTML, htmlBody)
	email.AddAlternative(mail.TextPlain, textbody)
	email.AddAttachmentData(newImage, "example.pdf", "pdf")
	err = email.Send(m.smtpClient)
	if err != nil {
		logrus.Errorf("Error sending mail : %s", err)
		return err
	}

	logrus.Infof("Sent Email to %s", msg.To)
	return nil
}

func (m *Mailer) buildPlainTextMessage(msg domain.Message) (string, error) {
	templateToRender := "./templates/mail.plain.gohtml"

	t, err := template.New("email-plain").ParseFiles(templateToRender)
	if err != nil {
		return "", err
	}

	var tpl bytes.Buffer
	if err = t.ExecuteTemplate(&tpl, "body", msg); err != nil {
		return "", err
	}

	plainMessage := tpl.String()

	return plainMessage, nil
}
