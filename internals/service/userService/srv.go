package userService

import (
	"fmt"
	"math/rand"
	"test-va/internals/Repository/userRepo"
	"test-va/internals/entity/ResponseEntity"
	"test-va/internals/entity/emailEntity"
	"test-va/internals/entity/userEntity"
	"test-va/internals/service/cryptoService"
	"test-va/internals/service/emailService"
	"test-va/internals/service/timeSrv"
	tokenservice "test-va/internals/service/tokenService"
	"test-va/internals/service/validationService"
	"time"

	"github.com/google/uuid"
)

type UserSrv interface {
	SaveUser(req *userEntity.CreateUserReq) (*userEntity.CreateUserRes, *ResponseEntity.ServiceError)
	Login(req *userEntity.LoginReq) (*userEntity.LoginRes, *ResponseEntity.ServiceError)
	GetUsers(page int) ([]*userEntity.UsersRes, error)
	GetUser(user_id string) (*userEntity.GetByIdRes, error)
	UpdateUser(req *userEntity.UpdateUserReq, userId string) (*userEntity.UpdateUserRes, *ResponseEntity.ServiceError)
	ChangePassword(req *userEntity.ChangePasswordReq, userId string) *ResponseEntity.ServiceError
	ResetPassword(req *userEntity.ResetPasswordReq) (*userEntity.ResetPasswordRes, error)
	ResetPasswordWithToken(req *userEntity.ResetPasswordWithTokenReq, token, userId string) *ResponseEntity.ServiceError
	DeleteUser(user_id string) error
	AssignVAToUser(user_id, va_id string) *ResponseEntity.ServiceError
}

type userSrv struct {
	repo      userRepo.UserRepository
	validator validationService.ValidationSrv
	timeSrv   timeSrv.TimeService
	cryptoSrv cryptoService.CryptoSrv
	emailSrv  emailService.EmailService
}

func (u *userSrv) Login(req *userEntity.LoginReq) (*userEntity.LoginRes, *ResponseEntity.ServiceError) {
	err := u.validator.Validate(req)
	if err != nil {
		return nil, ResponseEntity.NewValidatingError(err)
	}
	// FIND BY EMAIL
	user, err := u.repo.GetByEmail(req.Email)
	if err != nil {
		return nil, ResponseEntity.NewInternalServiceError(fmt.Sprintf("No User Found in DB: %v", err))
	}
	//compare password
	err = u.cryptoSrv.ComparePassword(user.Password, req.Password)
	if err != nil {
		return nil, ResponseEntity.NewInternalServiceError("Passwords Don't Match")
	}

	tokenSrv := tokenservice.NewTokenSrv("tokenString")
	token, refreshToken, errToken := tokenSrv.CreateToken(user.UserId, "user", req.Email)
	if errToken != nil {
		return nil, ResponseEntity.NewInternalServiceError("Cannot create access token!")
	}

	loggedInUser := userEntity.LoginRes{
		UserId:       user.UserId,
		Email:        user.Email,
		FirstName:    user.FirstName,
		LastName:     user.LastName,
		Phone:        user.Phone,
		Gender:       user.Gender,
		Token:        token,
		RefreshToken: refreshToken,
	}
	return &loggedInUser, nil
}

func (u *userSrv) SaveUser(req *userEntity.CreateUserReq) (*userEntity.CreateUserRes, *ResponseEntity.ServiceError) {
	// validate request
	err := u.validator.Validate(req)
	if err != nil {
		return nil, ResponseEntity.NewValidatingError(err)
	}
	// check if user with that email exists already

	_, err = u.repo.GetByEmail(req.Email)
	if err == nil {
		return nil, ResponseEntity.NewInternalServiceError("Email Already Exists")
	}

	//hash password
	password, err := u.cryptoSrv.HashPassword(req.Password)
	if err != nil {
		return nil, ResponseEntity.NewInternalServiceError(err)
	}
	//set time and etc
	req.UserId = uuid.New().String()
	req.Password = password
	req.AccountStatus = "ACTIVE"
	req.DateCreated = u.timeSrv.CurrentTime().Format(time.RFC3339)

	// save to DB
	err = u.repo.Persist(req)
	if err != nil {
		return nil, ResponseEntity.NewInternalServiceError(err)
	}

	tokenSrv := tokenservice.NewTokenSrv("tokenString")
	token, refreshToken, errToken := tokenSrv.CreateToken(req.UserId, "user", req.Email)
	if errToken != nil {
		return nil, ResponseEntity.NewInternalServiceError("Cannot create access token!")
	}

	data := &userEntity.CreateUserRes{
		UserId:       req.UserId,
		FirstName:    req.FirstName,
		LastName:     req.LastName,
		Email:        req.Email,
		Phone:        req.Phone,
		Token:        token,
		RefreshToken: refreshToken,
	}
	// set Reminder

	return data, nil
}

func (u *userSrv) UpdateUser(req *userEntity.UpdateUserReq, userId string) (*userEntity.UpdateUserRes, *ResponseEntity.ServiceError) {
	err := u.validator.Validate(req)
	if err != nil {
		return nil, ResponseEntity.NewValidatingError(err)
	}

	err = u.repo.UpdateUser(req, userId)
	if err != nil {
		return nil, ResponseEntity.NewInternalServiceError(err)
	}
	data := &userEntity.UpdateUserRes{
		FirstName:   req.FirstName,
		LastName:    req.LastName,
		Email:       req.Email,
		Phone:       req.Phone,
		Gender:      req.Gender,
		DateOfBirth: req.DateOfBirth,
	}

	return data, nil
}

func (u *userSrv) ChangePassword(req *userEntity.ChangePasswordReq, userId string) *ResponseEntity.ServiceError {
	// validate request
	err := u.validator.Validate(req)
	if err != nil {
		return ResponseEntity.NewValidatingError(err)
	}

	// Get user by user id
	user, err := u.repo.GetById(userId)
	if err != nil {
		return ResponseEntity.NewInternalServiceError("Check the access token!")
	}

	// Compare password in database and password gotten from user
	err = u.cryptoSrv.ComparePassword(user.Password, req.OldPassword)
	if err != nil {
		return ResponseEntity.NewInternalServiceError("Passwords do not match!")
	}

	// Check if new password is the same as old password
	err = u.cryptoSrv.ComparePassword(user.Password, req.NewPassword)
	if err == nil {
		return ResponseEntity.NewInternalServiceError("The new password cannot be the same as your old password!")
	}

	// Create a new password hash
	newPassword, _ := u.cryptoSrv.HashPassword(req.NewPassword)
	err = u.repo.ChangePassword(userId, newPassword)
	if err != nil {
		return ResponseEntity.NewInternalServiceError("Could not change password!")
	}

	return nil
}

func (u *userSrv) GetUsers(page int) ([]*userEntity.UsersRes, error) {
	users, err := u.repo.GetUsers(page)
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (u *userSrv) GetUser(user_id string) (*userEntity.GetByIdRes, error) {
	user, err := u.repo.GetById(user_id)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (u *userSrv) DeleteUser(user_id string) error {
	_, idErr := u.repo.GetById(user_id)
	if idErr != nil {
		return idErr
	}

	delErr := u.repo.DeleteUser(user_id)
	if delErr != nil {
		return delErr
	}

	return nil
}

func (u *userSrv) ResetPassword(req *userEntity.ResetPasswordReq) (*userEntity.ResetPasswordRes, error) {
	var token userEntity.ResetPasswordRes
	var message emailEntity.SendEmailReq

	err := u.validator.Validate(req)
	if err != nil {
		return nil, err
	}

	// Check if the user exists, if he/she doesn't return error
	user, err := u.repo.GetByEmail(req.Email)
	if err != nil {
		return nil, err
	}

	// Delete old tokens from system
	err = u.repo.DeleteToken(user.UserId)
	if err != nil {
		return nil, err
	}

	// Create token, add to database and then send to user's email address
	token.UserId = user.UserId
	token.TokenId = uuid.New().String()
	token.Token = GenerateToken(4)
	token.Expiry = time.Now().Add(time.Minute * 30).Format(time.RFC3339)

	err = u.repo.AddToken(&token)
	if err != nil {
		return nil, err
	}

	// Send message to users email, if it exists
	message.EmailAddress = user.Email
	message.EmailSubject = "Reset Password Token"
	message.EmailBody = CreateMessageBody(user.FirstName, user.LastName, user.UserId, token.Token)

	err = u.emailSrv.SendMail(message)
	if err != nil {
		return nil, err
	}

	return &token, nil
}

func (u *userSrv) ResetPasswordWithToken(req *userEntity.ResetPasswordWithTokenReq, token, userId string) *ResponseEntity.ServiceError {
	err := u.validator.Validate(req)
	if err != nil {
		return ResponseEntity.NewInternalServiceError(err)
	}

	tokenDB, err := u.repo.GetTokenById(token, userId)
	fmt.Println(tokenDB)
	if err != nil {
		return ResponseEntity.NewInternalServiceError("Invalid access token!")
	}

	timeNow := time.Now().Format(time.RFC3339)
	if tokenDB.Expiry < timeNow {
		return ResponseEntity.NewInternalServiceError("Token has expired!")
	}

	user, err := u.repo.GetById(tokenDB.UserId)
	if err != nil {
		return ResponseEntity.NewInternalServiceError("Check the user!")
	}

	err = u.cryptoSrv.ComparePassword(user.Password, req.Password)
	if err == nil {
		return ResponseEntity.NewInternalServiceError("The new password cannot be the same as your old password!")
	}

	// Create a new password hash
	newPassword, _ := u.cryptoSrv.HashPassword(req.Password)
	err = u.repo.ChangePassword(tokenDB.UserId, newPassword)
	if err != nil {
		return ResponseEntity.NewInternalServiceError("Could not change password!")
	}

	return nil
}

func (u *userSrv) AssignVAToUser(user_id, va_id string) *ResponseEntity.ServiceError {
	err := u.repo.AssignVAToUser(user_id, va_id)
	if err != nil {
		fmt.Println(err)
		switch {
		case err.Error() == "user already has a VA":
			return ResponseEntity.NewCustomServiceError("user already has a VA", err)
		default:
			return ResponseEntity.NewInternalServiceError("Could Not Assign Va")
		}
	}
	return nil
}

// Auxillary Function
func GenerateToken(tokenLength int) string {
	rand.Seed(time.Now().UnixNano())
	const charset = "0123456789"
	b := make([]byte, tokenLength)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

func CreateMessageBody(firstName, lastName, email, token string) string {
	link := fmt.Sprintf("https://ticked.hng.tech/reset-password?token=%v&user_id=%v", token, email)

	subject := fmt.Sprintf("Hi %v %v, \n\n", firstName, lastName)
	mainBody := fmt.Sprintf("You have requested to reset your password, this is your otp code %v\n\nIf you have difficulty using the otp code, please copy this link and paste it in your browser:\n%v\n\nBut if you did not request for a change of password, you can forget about this email\n\nLink expires in 30 minutes!", token, link)

	message := subject + mainBody
	return string(message)
}

func NewUserSrv(repo userRepo.UserRepository, validator validationService.ValidationSrv, timeSrv timeSrv.TimeService, cryptoSrv cryptoService.CryptoSrv, emailSrv emailService.EmailService) UserSrv {
	return &userSrv{repo: repo, validator: validator, timeSrv: timeSrv, cryptoSrv: cryptoSrv, emailSrv: emailSrv}
}
