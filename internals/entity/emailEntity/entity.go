package emailEntity

type SendEmailReq struct {
	EmailAddress string
	EmailSubject string
	EmailBody    string
}

type SendBatchEmail struct{
	EmailAddresses []string
	EmailSubject string
	EmailBody    string
}