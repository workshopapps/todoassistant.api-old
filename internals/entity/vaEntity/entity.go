package vaEntity

type CreateVAReq struct {
	VaId           string `json:"va_id"`
	FirstName      string `json:"first_name"`
	LastName       string `json:"last_name"`
	Email          string `json:"email"`
	Phone          string `json:"phone"`
	ProfilePicture string `json:"profile_picture"`
	Password       string `json:"password"`
	AccountType    string `json:"account_type"`
	CreatedAt      string `json:"created_at"`
}

type CreateVARes struct {
	VaId           string `json:"va_id"`
	FirstName      string `json:"first_name"`
	LastName       string `json:"last_name"`
	Email          string `json:"email"`
	Phone          string `json:"phone"`
	ProfilePicture string `json:"profile_picture"`
	AccountType    string `json:"account_type"`
}

type FindByIdRes struct {
	VaId           string `json:"va_id"`
	FirstName      string `json:"first_name"`
	LastName       string `json:"last_name"`
	Email          string `json:"email"`
	Phone          string `json:"phone"`
	ProfilePicture string `json:"profile_picture"`
	Password       string `json:"password"`
	CreatedAt      string `json:"created_at"`
	AccountType    string `json:"account_type"`
}

type EditVaReq struct {
	FirstName      string `json:"first_name"`
	LastName       string `json:"last_name"`
	Email          string `json:"email"`
	Phone          string `json:"phone"`
	ProfilePicture string `json:"profile_picture"`
}

type EditVARes struct {
	FirstName      string `json:"first_name"`
	LastName       string `json:"last_name"`
	Email          string `json:"email"`
	Phone          string `json:"phone"`
	ProfilePicture string `json:"profile_picture"`
}

type ChangeVAPassword struct {
	VaId        string `json:"va_id"`
	NewPassword string `json:"new_password"`
}

type VAStruct struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	UserId    string `json:"user_id"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
	Status    string `json:"status"`
}
