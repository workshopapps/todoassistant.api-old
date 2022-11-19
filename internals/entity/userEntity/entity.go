package userEntity

type CreateUserReq struct {
	UserId        string `json:"user_id"`
	FirstName     string `json:"first_name"`
	LastName      string `json:"last_name"`
	Email         string `json:"email"`
	Phone         string `json:"phone"`
	Password      string `json:"password"`
	Gender        string `json:"gender"`
	DateOfBirth   string `json:"date_of_birth"`
	AccountStatus string `json:"account_status"`
	PaymentStatus string `json:"payment_status"`
	DateCreated   string `json:"date_created"`
}

type CreateUserRes struct {
	UserId    string `json:"user_id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
}
