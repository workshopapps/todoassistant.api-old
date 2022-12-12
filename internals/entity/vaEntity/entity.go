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

type LoginReq struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type LoginRes struct {
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

type FindByIdRes struct {
	VaId           string `json:"va_id"`
	FirstName      string `json:"first_name"`
	LastName       string `json:"last_name"`
	Email          string `json:"email"`
	Phone          string `json:"phone"`
	ProfilePicture string `json:"profile_picture"`
	CreatedAt      string `json:"created_at"`
	AccountType    string `json:"account_type"`
}

type FindByEmailRes struct {
	VaId           string `json:"va_id"`
	FirstName      string `json:"first_name"`
	LastName       string `json:"last_name"`
	Email          string `json:"email"`
	Phone          string `json:"phone"`
	Password       string `json:"password"`
	ProfilePicture string `json:"profile_picture"`
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
	Avatar    string `json:"avatar"`
}

type VATask struct {
	TaskId      string `json:"task_id"`
	Title       string `json:"title"`
	EndTime     string `json:"end_time"`
	Description string `json:"description"`
	VaOption    string `json:"va_option"`
	VaId        string `json:"va_id"`
	Status      string `json:"status"`
	User        VAUser `json:"user"`
}

type VATaskAll struct {
	TaskId       string `json:"task_id"`
	Title        string `json:"title"`
	EndTime      string `json:"end_time"`
	Description  string `json:"description"`
	VaOption     string `json:"va_option"`
	VaId         string `json:"va_id"`
	Status       string `json:"status"`
	CommentCount string `json:"comment_count"`
	User         NewVAUser `json:"user"`
}

type VAUser struct {
	UserId string `json:"user_id"`
	Name   string `json:"name"`
	Phone  string `json:"phone"`
}

type NewVAUser struct {
	UserId string `json:"user_id"`
	Name   string `json:"name"`
	Phone  string `json:"phone"`
	Avatar  string `json:"avatar"`
}
