package emailService
import (

	"net/smtp"
	"os"
)
type Emailservice interface{
	SendMail(req *email) error
}
type emailSrv struct {
}
type email struct{
	name string
	emailAddress string
	emailSubject string
	emailBody string
}

func SendMail(t email) error{
	
	from := os.Getenv("FromEmailAddr")
	password := os.Getenv("SMTPpwd")
	toEmail := t.emailAddress
	to := []string{toEmail}
	host := "smtp.gmail.com"
	port := "567"
	address := host + ":" + port
	subject := t.emailSubject
	body := t.emailBody
	message := []byte(subject + body)
	auth := smtp.PlainAuth("", from, password, host)
	err := smtp.SendMail(address, auth, from, to, message)
	if err != nil {
		return err
	}
	return nil
}
// func NewTokenSrv(secret string) TokenSrv {
// 	return &tokenSrv{secret}
// }