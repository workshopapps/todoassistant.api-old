package userService

import (
	"fmt"
	"log"
	"math/rand"
	"mime/multipart"
	"strings"
	"test-va/internals/Repository/userRepo"
	"test-va/internals/entity/ResponseEntity"
	"test-va/internals/entity/emailEntity"
	"test-va/internals/entity/userEntity"
	"test-va/internals/service/awsService"
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
	UploadImage(file *multipart.FileHeader, userId string) (*userEntity.ProfileImageRes, error)
	ChangePassword(req *userEntity.ChangePasswordReq) *ResponseEntity.ServiceError
	ResetPassword(req *userEntity.ResetPasswordReq) (*userEntity.ResetPasswordRes, *ResponseEntity.ServiceError)
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
	awsSrv    awsService.AWSService
	tokenSrv  tokenservice.TokenSrv
}

// Login User godoc
// @Summary	Provide email and password to be logged in
// @Description	Login to the server
// @Tags	Users
// @Accept	json
// @Produce	json
// @Param	request	body	userEntity.LoginReq	true "Login Details"
// @Success	200  {object}  userEntity.LoginRes
// @Failure	400  {object}  ResponseEntity.ServiceError
// @Failure	404  {object}  ResponseEntity.ServiceError
// @Failure	500  {object}  ResponseEntity.ServiceError
// @Router	/user/login [post]
func (u *userSrv) Login(req *userEntity.LoginReq) (*userEntity.LoginRes, *ResponseEntity.ServiceError) {
	err := u.validator.Validate(req)
	if err != nil {
		return nil, ResponseEntity.NewValidatingError(err)
	}
	// FIND BY EMAIL
	user, err := u.repo.GetByEmail(req.Email)
	if err != nil {
		return nil, ResponseEntity.NewInternalServiceError("Invalid Login Credentials")
	}
	//compare password
	err = u.cryptoSrv.ComparePassword(user.Password, req.Password)
	if err != nil {
		return nil, ResponseEntity.NewInternalServiceError("Passwords Don't Match")
	}

	token, refreshToken, errToken := u.tokenSrv.CreateToken(user.UserId, "user", req.Email)
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
		Avatar:       user.Avatar,
		Token:        token,
		RefreshToken: refreshToken,
	}
	return &loggedInUser, nil
}

// Register User godoc
// @Summary	Register route
// @Description	Register route
// @Tags	Users
// @Accept	json
// @Produce	json
// @Param	request	body	userEntity.CreateUserReq	true "Signup Details"
// @Success	200  {object}  userEntity.CreateUserRes
// @Failure	400  {object}  ResponseEntity.ServiceError
// @Failure	404  {object}  ResponseEntity.ServiceError
// @Failure	500  {object}  ResponseEntity.ServiceError
// @Router	/user [post]
func (u *userSrv) SaveUser(req *userEntity.CreateUserReq) (*userEntity.CreateUserRes, *ResponseEntity.ServiceError) {
	err := u.validator.Validate(req)
	if err != nil {
		log.Println(err)
		return nil, ResponseEntity.NewValidatingError(err)
	}
	// check if user with that email exists already

	_, err = u.repo.GetByEmail(req.Email)
	if err == nil {
		log.Println(err)
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
		log.Println(err)
		return nil, ResponseEntity.NewInternalServiceError(err)
	}

	token, refreshToken, errToken := u.tokenSrv.CreateToken(req.UserId, "user", req.Email)
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

// Update User godoc
// @Summary	Update a user profile
// @Description	Register route
// @Tags	Users
// @Accept	json
// @Produce	json
// @Param	userId	path	string	true	"User Id"
// @Param	request	body	userEntity.UpdateUserReq	true "Update User Details"
// @Success	200  {object}  userEntity.UpdateUserRes
// @Failure	400  {object}  ResponseEntity.ServiceError
// @Failure	404  {object}  ResponseEntity.ServiceError
// @Failure	500  {object}  ResponseEntity.ServiceError
// @Security ApiKeyAuth
// @Router	/user/{userId} [put]
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

// Update Profile Picture godoc
// @Summary	Update the current user profile image
// @Description	Upload image route
// @Tags	Users
// @Accept	mpfd
// @Produce	json
// @Param	Upload-Image	formData	file	true	"Update profile picture"
// @Success	200  {object}	userEntity.ProfileImageRes
// @Failure	400  {object}  ResponseEntity.ServiceError
// @Failure	404  {object}  ResponseEntity.ServiceError
// @Failure	500  {object}  ResponseEntity.ServiceError
// @Security BasicAuth
// @Router	/user/upload [post]
func (u *userSrv) UploadImage(file *multipart.FileHeader, userId string) (*userEntity.ProfileImageRes, error) {
	var res userEntity.ProfileImageRes
	fileType := strings.Split(file.Header.Get("Content-Type"), "/")[1]

	fileName := fmt.Sprintf("%s/%s.%s", userId, uuid.New().String(), fileType)

	err := u.awsSrv.UploadImage(file, fileName)

	if err != nil {
		return nil, err
	}

	res.Image = fmt.Sprintf("https://ticked-v1-backend-bucket.s3.amazonaws.com/%v", fileName)
	err = u.repo.UpdateImage(userId, res.Image)

	if err != nil {
		return nil, err
	}

	res.Size = file.Size
	res.FileType = fileType
	return &res, nil
}

// Change Password godoc
// @Summary	Change a user password
// @Description	Change password route
// @Tags	Users
// @Accept	json
// @Produce	json
// @Param	userId	path	string	true	"User Id"
// @Param	request	body	userEntity.ChangePasswordReq	true	"New password"
// @Success	200  {string}  string    "ok"
// @Failure	400  {object}  ResponseEntity.ServiceError
// @Failure	404  {object}  ResponseEntity.ServiceError
// @Failure	500  {object}  ResponseEntity.ServiceError
// @Security ApiKeyAuth
// @Router	/user/{userId}/change-password [put]
func (u *userSrv) ChangePassword(req *userEntity.ChangePasswordReq) *ResponseEntity.ServiceError {
	err := u.validator.Validate(req)
	if err != nil {
		return ResponseEntity.NewValidatingError(err)
	}

	// Get user by user id
	user, err := u.repo.GetById(req.UserId)
	if err != nil {
		return ResponseEntity.NewInternalServiceError("Check the access token!")
	}

	// Compare password in database and password gotten from user
	err = u.cryptoSrv.ComparePassword(user.Password, req.OldPassword)
	if err != nil {
		log.Println("request", req)
		log.Println(err)
		return ResponseEntity.NewInternalServiceError("Passwords do not match!")
	}

	// Check if new password is the same as old password
	//err = u.cryptoSrv.ComparePassword(user.Password, req.NewPassword)

	// if err == nil {
	// 	return ResponseEntity.NewInternalServiceError("The new password cannot be the same as your old password!")
	// }

	// Create a new password hash
	newPassword, _ := u.cryptoSrv.HashPassword(req.NewPassword)
	err = u.repo.ChangePassword(req.UserId, newPassword)
	if err != nil {
		return ResponseEntity.NewInternalServiceError("Could not change password!")
	}

	return nil
}

// Get All Users godoc
// @Summary	Get all users in the database
// @Description	Get all users route
// @Tags	Users
// @Accept	json
// @Produce	json
// @Param	page	query	string	false	"page"
// @Success	200  {object}  []userEntity.UsersRes
// @Failure	400  {object}  ResponseEntity.ServiceError
// @Failure	404  {object}  ResponseEntity.ServiceError
// @Failure	500  {object}  ResponseEntity.ServiceError
// @Security ApiKeyAuth
// @Router	/user [get]
func (u *userSrv) GetUsers(page int) ([]*userEntity.UsersRes, error) {
	users, err := u.repo.GetUsers(page)
	if err != nil {
		return nil, err
	}

	return users, nil
}

// Get User godoc
// @Summary	Get a specific user
// @Description	Get user route
// @Tags	Users
// @Accept	json
// @Produce	json
// @Param	userId	path	string	true	"User Id"
// @Success	200  {object}  userEntity.GetByIdRes
// @Failure	400  {object}  ResponseEntity.ServiceError
// @Failure	404  {object}  ResponseEntity.ServiceError
// @Failure	500  {object}  ResponseEntity.ServiceError
// @Security ApiKeyAuth
// @Router	/user/{userId} [get]

// Get User godoc
// @Summary	Get a specific user
// @Description	Get user route
// @Tags	VA
// @Accept	json
// @Produce	json
// @Param	userId	path	string	true	"User Id"
// @Success	200  {object}  userEntity.GetByIdRes
// @Failure	400  {object}  ResponseEntity.ServiceError
// @Failure	404  {object}  ResponseEntity.ServiceError
// @Failure	500  {object}  ResponseEntity.ServiceError
// @Security ApiKeyAuth
// @Router	/user/profile/{userId} [get]
func (u *userSrv) GetUser(user_id string) (*userEntity.GetByIdRes, error) {
	user, err := u.repo.GetById(user_id)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// Delete User godoc
// @Summary	Delete a user from the database
// @Description	Delete route
// @Tags	Users
// @Accept	json
// @Produce	json
// @Param	userId	path	string	true	"User Id"
// @Success	200  {string}  string    "ok"
// @Failure	400  {object}  ResponseEntity.ServiceError
// @Failure	404  {object}  ResponseEntity.ServiceError
// @Failure	500  {object}  ResponseEntity.ServiceError
// @Security ApiKeyAuth
// @Router	/user/{userId} [delete]
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

// Reset password godoc
// @Summary	Generate a token to reset users password
// @Description	Generate token
// @Tags	Users
// @Accept	json
// @Produce	json
// @Param	request	body userEntity.ResetPasswordReq	true "Input your email"
// @Success	200  {object}  userEntity.ResetPasswordRes
// @Failure	400  {object}  ResponseEntity.ServiceError
// @Failure	404  {object}  ResponseEntity.ServiceError
// @Failure	500  {object}  ResponseEntity.ServiceError
// @Router	/user/reset-password [post]
func (u *userSrv) ResetPassword(req *userEntity.ResetPasswordReq) (*userEntity.ResetPasswordRes, *ResponseEntity.ServiceError) {
	var token userEntity.ResetPasswordRes
	var message emailEntity.SendEmailReq

	err := u.validator.Validate(req)
	if err != nil {
		return nil, ResponseEntity.NewInternalServiceError(err)
	}

	// Check if the user exists, if he/she doesn't return error
	user, err := u.repo.GetByEmail(req.Email)
	if err != nil {
		return nil, ResponseEntity.NewInternalServiceError(err)
	}

	// Delete old tokens from system
	err = u.repo.DeleteToken(user.UserId)
	if err != nil {
		return nil, ResponseEntity.NewInternalServiceError(err)
	}

	// Create token, add to database and then send to user's email address
	token.UserId = user.UserId
	token.TokenId = uuid.New().String()
	token.Token = GenerateToken(4)
	token.Expiry = time.Now().Add(time.Minute * 30).Format(time.RFC3339)

	err = u.repo.AddToken(&token)
	if err != nil {
		return nil, ResponseEntity.NewInternalServiceError(err)
	}

	// Send message to users email, if it exists
	message.EmailAddress = user.Email
	message.EmailSubject = "Subject: Reset Password Token\n"
	message.EmailBody = CreateMessageBody(user.FirstName, user.LastName, token.Token)

	err = u.emailSrv.SendMail(message)
	if err != nil {
		return nil, ResponseEntity.NewInternalServiceError(err)
	}

	return &token, nil
}

// Reset password with token godoc
// @Summary	Check the provided token and reset the user's password
// @Description	Reset password
// @Tags	Users
// @Accept	json
// @Produce	json
// @Param	token	query	string	true	"Token"
// @Param	user_id	query	string	true	"User Id"
// @Success	200  {string}  string    "ok"
// @Failure	400  {object}  ResponseEntity.ServiceError
// @Failure	404  {object}  ResponseEntity.ServiceError
// @Failure	500  {object}  ResponseEntity.ServiceError
// @Router	/reset-password-token [post]
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

// Assign VA To User godoc
// @Summary	Assign VA to a User
// @Description	Assing VA to User route
// @Tags	Users
// @Accept	json
// @Produce	json
// @Param	vaId	path	string	true	"VA Id"
// @Success	200  {string}	string	"Ok"
// @Failure	400  {object}  ResponseEntity.ServiceError
// @Failure	404  {object}  ResponseEntity.ServiceError
// @Failure	500  {object}  ResponseEntity.ServiceError
// @Security ApiKeyAuth
// @Router	/assign-va/{vaId} [post]
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

func CreateMessageBody(firstName, lastName, token string) string {
	subject := fmt.Sprintf("Hi %v %v, \n\n", firstName, lastName)
	mainBody := fmt.Sprintf("You have requested to reset your password, this is your otp code %v\nBut if you did not request for a change of password, you can ignore this email.\n\nLink expires in 30 minutes!", token)

	message := subject + mainBody
	return string(message)
}

func NewUserSrv(repo userRepo.UserRepository, validator validationService.ValidationSrv, timeSrv timeSrv.TimeService,
	cryptoSrv cryptoService.CryptoSrv, emailSrv emailService.EmailService, awsSrv awsService.AWSService,
	tokenSrv tokenservice.TokenSrv) UserSrv {
	return &userSrv{repo: repo, validator: validator, timeSrv: timeSrv,
		cryptoSrv: cryptoSrv, emailSrv: emailSrv, awsSrv: awsSrv, tokenSrv: tokenSrv}
}
