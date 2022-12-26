package emailService

import (
	"fmt"
	"net/smtp"
	"test-va/internals/entity/emailEntity"
)

type EmailService interface {
	SendMail(req emailEntity.SendEmailReq) error
	SendBatchEmail(req emailEntity.SendBatchEmail) error
	SendMailToSupport(req emailEntity.SendEmailReq) error
}
type emailSrv struct {
	FromEmail string
	Password  string
	Host      string
	Port      string
}

func (e emailSrv) SendBatchEmail(req emailEntity.SendBatchEmail) error {
	auth := smtp.PlainAuth("", e.FromEmail, e.Password, e.Host)
	addr := e.Host + ":" + e.Port
	header := fmt.Sprintf("From: %v\nTo: %v\n", e.FromEmail, req.EmailAddresses)
	body := []byte(header + req.EmailSubject + "\n" + req.EmailBody)
	err := smtp.SendMail(addr, auth, e.FromEmail, req.EmailAddresses, body)
	if err != nil {
		return err
	}
	return nil
}

func (e emailSrv) SendMail(req emailEntity.SendEmailReq) error {
	auth := smtp.PlainAuth("", e.FromEmail, e.Password, e.Host)
	addr := e.Host + ":" + e.Port
	header := fmt.Sprintf("From: %v\nTo: %v\n", e.FromEmail, req.EmailAddress)
	body := []byte(header + req.EmailSubject + req.EmailBody)
	err := smtp.SendMail(addr, auth, e.FromEmail, []string{req.EmailAddress}, body)
	if err != nil {
		return err
	}
	return nil
}

func (e emailSrv) SendMailToSupport(req emailEntity.SendEmailReq) error {
	auth := smtp.PlainAuth("", e.FromEmail, e.Password, e.Host)
	addr := e.Host + ":" + e.Port
	header := fmt.Sprintf("From: %v\nTo: %v\n", req.EmailAddress, e.FromEmail)
	body := []byte(header + req.EmailSubject + req.EmailBody)
	err := smtp.SendMail(addr, auth, req.EmailAddress, []string{e.FromEmail}, body)

	if err != nil {
		return err
	}

	return nil
}

// func SendMail(req emailEntity.SendEmailReq) error {
// 	// add the email and password below, but remember to remove before pushing to giyhub
// 	from := utils.Config.FromEmailAddr
// 	password := os.Getenv("SMTPpwd")
// 	toEmail := req.EmailAddress
// 	to := []string{toEmail}
// 	host := "smtp.gmail.com"
// 	port := "587"
// 	address := host + ":" + port
// 	subject := req.EmailSubject
// 	body := req.EmailBody
// 	message := []byte(subject + "\n" + body)
// 	auth := smtp.PlainAuth("", from, password, host)
// 	err := smtp.SendMail(address, auth, from, to, message)
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }

func NewEmailSrv(fromEmail string, password string, host string, port string) EmailService {
	return &emailSrv{FromEmail: fromEmail, Password: password, Host: host, Port: port}
}
