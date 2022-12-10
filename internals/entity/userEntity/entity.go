package userEntity

type CreateUserReq struct {
	UserId        string `json:"user_id"`
	FirstName     string `json:"first_name" validate:"required"`
	LastName      string `json:"last_name"  validate:"required"`
	Email         string `json:"email" validate:"email"`
	Phone         string `json:"phone"`
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
	Avatar       string `json:"avatar"`
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
	Avatar    string `json:"avatar"`
}

type GetByIdRes struct {
	UserId      string `json:"user_id"`
	Password    string `json:"password"`
	Email       string `json:"email"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	Phone       string `json:"phone"`
	Gender      string `json:"gender"`
	Avatar      string `json:"avatar"`
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
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	Email       string `json:"email"`
	Phone       string `json:"phone"`
	Gender      string `json:"gender"`
	DateOfBirth string `json:"date_of_birth"`
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

type ResetPasswordReq struct {
	Email string `json:"email" validate:"email,required"`
}

type ResetPasswordRes struct {
	UserId  string `json:"user_id"`
	TokenId string `json:"token_id"`
	Token   string `json:"token"`
	Expiry  string `json:"expiry"`
}

type ResetPasswordWithTokenReq struct {
	Password string `json:"password" validate:"required"`
}

type ResetPasswordWithTokenRes struct {
	UserId  string `json:"user_id"`
	TokenId string `json:"token_id"`
	Token   string `json:"token"`
	Expiry  string `json:"expiry"`
}

type GoogleLoginReq struct {
	Id        string `json:"googleId"`
	FirstName string `json:"givenName"`
	LastName  string `json:"familyName"`
	Email     string `json:"email"`
	Profile   string `json:"imageUrl"`
	Name      string `json:"name"`
}

type FacebookLoginReq struct {
	Email string `json:"email"`
	Name  string `json:"name"`
}

type ProfileImageRes struct {
	Image    string `json:"avatar"`
	Size     int64  `json:"size"`
	FileType string `json:"fileType"`
}
