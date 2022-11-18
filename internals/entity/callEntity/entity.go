package callEntity

type CallRes struct {
	CallId  	string  `json:"id"`
	VaId    	string	`json:"va_id"`
	UserId  	string	`json:"user_id"`
	CallRating 	int		`json:"call_rating"`
	CallComment string 	`json:"call_comment"`
}
