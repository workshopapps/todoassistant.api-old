package userService

import (
	"log"
	"net/http"
	"test-va/internals/Repository/userRepo"
	"test-va/internals/entity/ResponseEntity"
	"test-va/internals/entity/userEntity"
	"test-va/internals/service/cryptoService"
	"test-va/internals/service/timeSrv"
	"test-va/internals/service/validationService"
	"time"

	"github.com/google/uuid"
)

type UserSrv interface {
	SaveUser(req *userEntity.CreateUserReq) (*userEntity.CreateUserRes, *ResponseEntity.ResponseMessage)
	Login(req *userEntity.LoginReq) (*userEntity.LoginRes, *ResponseEntity.ResponseMessage)
}

type userSrv struct {
	repo      userRepo.UserRepository
	validator validationService.ValidationSrv
	timeSrv   timeSrv.TimeService
	cryptoSrv cryptoService.CryptoSrv
}

func (u *userSrv) Login(req *userEntity.LoginReq) (*userEntity.LoginRes, *ResponseEntity.ResponseMessage) {
	err := u.validator.Validate(req)
	if err != nil {
		return nil, ResponseEntity.NewCustomError(http.StatusBadRequest, "Bad Input Request")
	}
	// FIND BY EMAIL
	user, err := u.repo.GetByEmail(req.Email)
	if err != nil {
		return nil, ResponseEntity.NewCustomError(http.StatusNotFound, "user not found in database")
	}
	//compare password
	err = u.cryptoSrv.ComparePassword(user.Password, req.Password)
	if err != nil {
		return nil, ResponseEntity.NewCustomError(http.StatusUnauthorized, "invalid login credentials")
	}
	loggedInUser := userEntity.LoginRes{
		UserId:    user.UserId,
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Phone:     user.Phone,
		Gender:    user.Gender,
	}
	return &loggedInUser, nil
}

func (u *userSrv) SaveUser(req *userEntity.CreateUserReq) (*userEntity.CreateUserRes, *ResponseEntity.ResponseMessage) {
	// validate request
	err := u.validator.Validate(req)
	if err != nil {
		return nil, ResponseEntity.NewCustomError(http.StatusBadRequest, "Bad Input Request")
	}
	// check if user with that email exists already

	_, err = u.repo.GetByEmail(req.Email)
	if err == nil {
		return nil, ResponseEntity.NewCustomError(http.StatusBadRequest, "User Exists Already")
	}

	//hash password
	password, err := u.cryptoSrv.HashPassword(req.Password)
	if err != nil {
		return nil, ResponseEntity.NewCustomError(http.StatusInternalServerError, "Error Hashing Password")
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
		return nil, ResponseEntity.NewCustomError(http.StatusInternalServerError, "Error Saving To DB")
	}

	data := &userEntity.CreateUserRes{
		UserId:    req.UserId,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Email:     req.Email,
		Phone:     req.Phone,
	}
	// set Reminder

	return data, nil
}

func NewUserSrv(repo userRepo.UserRepository, validator validationService.ValidationSrv, timeSrv timeSrv.TimeService, cryptoSrv cryptoService.CryptoSrv) UserSrv {
	return &userSrv{repo: repo, validator: validator, timeSrv: timeSrv, cryptoSrv: cryptoSrv}
}
