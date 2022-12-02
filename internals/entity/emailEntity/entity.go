package emailEntity

type SendEmailReq struct {
	Name         string
	EmailAddress string
	EmailSubject string
	EmailBody    string
}

type SendBatchEmail struct {
	Name           string
	EmailAddresses []string
	EmailSubject   string
	EmailBody      string
}
