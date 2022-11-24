package userEntity

type CreateUserReq struct {
	UserId        string `json:"user_id"`
	FirstName     string `json:"first_name" validate:"required"`
	LastName      string `json:"last_name"  validate:"required"`
	Email         string `json:"email" validate:"email"`
	Phone         string `json:"phone" validate:"required"`
	Password      string `json:"password" validate:"required,min=6"`
	Gender        string `json:"gender"  validate:"required,oneof='Male' 'Female'"`
	DateOfBirth   string `json:"date_of_birth" validate:"required"`
	AccountStatus string `json:"account_status"`
	PaymentStatus string `json:"payment_status"`
	DateCreated   string `json:"date_created"`
}

type CreateUserRes struct {
	UserId       string `json:"user_id"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	Email        string `json:"email"`
	Phone        string `json:"phone"`
	Token        string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type LoginReq struct {
	Email    string `json:"email" validate:"email"`
	Password string `json:"password" validate:"required"`
}

type LoginRes struct {
	UserId       string `json:"user_id"`
	Email        string `json:"email"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	Phone        string `json:"phone"`
	Gender       string `json:"gender"`
	Token        string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
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
