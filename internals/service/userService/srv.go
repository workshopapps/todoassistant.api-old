package userService

import (
	"fmt"
	"test-va/internals/Repository/userRepo"
	"test-va/internals/entity/ResponseEntity"
	"test-va/internals/entity/userEntity"
	"test-va/internals/service/cryptoService"
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
	UpdateUser(req *userEntity.UpdateUserReq, userId string) (*userEntity.GetByIdRes, *ResponseEntity.ServiceError)
	ChangePassword(req *userEntity.ChangePasswordReq, userId string) *ResponseEntity.ServiceError
	DeleteUser(user_id string) error
}

type userSrv struct {
	repo      userRepo.UserRepository
	validator validationService.ValidationSrv
	timeSrv   timeSrv.TimeService
	cryptoSrv cryptoService.CryptoSrv
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
	token, refreshToken, errToken := tokenSrv.CreateToken(user.UserId, req.Email)
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
	token, refreshToken, errToken := tokenSrv.CreateToken(req.UserId, req.Email)
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

func (u *userSrv) UpdateUser(req *userEntity.UpdateUserReq, userId string) (*userEntity.GetByIdRes, *ResponseEntity.ServiceError) {
	err := u.validator.Validate(req)
	if err != nil {
		return nil, ResponseEntity.NewValidatingError(err)
	}

	result, err := u.repo.UpdateUser(req, userId)
	if err != nil {
		return nil, ResponseEntity.NewInternalServiceError(err)
	}

	return result, nil
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
		return ResponseEntity.NewInternalServiceError("The password cannot be the same as your old password!")
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

func NewUserSrv(repo userRepo.UserRepository, validator validationService.ValidationSrv, timeSrv timeSrv.TimeService, cryptoSrv cryptoService.CryptoSrv) UserSrv {
	return &userSrv{repo: repo, validator: validator, timeSrv: timeSrv, cryptoSrv: cryptoSrv}
}
