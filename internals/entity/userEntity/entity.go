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

type LoginReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginRes struct {
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Phone     string `json:"phone"`
	Gender    string `json:"gender"`
}

type GetByEmailRes struct {
	UserId    string `json:"user_id"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Phone     string `json:"phone"`
	Gender    string `json:"gender"`
}

type GetByIdRes struct {
	UserId      string `json:"user_id"`
	Password    string `json:"password"`
	Email       string `json:"email"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	Phone       string `json:"phone"`
	Gender      string `json:"gender"`
	DateOfBirth string `json:"date_of_birth"`
}

type UpdateUserReq struct {
	FirstName     string `json:"first_name"`
	LastName      string `json:"last_name"`
	Email         string `json:"email"`
	Phone         string `json:"phone"`
	Gender        string `json:"gender"`
	DateOfBirth   string `json:"date_of_birth"`
	AccountStatus string `json:"account_status"`
	PaymentStatus string `json:"payment_status"`
}

type UpdateUserRes struct {
	UserId    string `json:"user_id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
}

type UsersRes struct {
	UserId      string `json:"user_id"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	Email       string `json:"email"`
	Phone       string `json:"phone"`
	DateOfBirth string `json:"date_of_birth"`
	DateCreated string `json:"date_created"`
}

type ChangePasswordReq struct {
	UserId      string `json:"user_id"`
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}
