package userService

import (
	"github.com/google/uuid"
	"log"
	"net/http"
	"test-va/internals/Repository/userRepo"
	"test-va/internals/entity/ResponseEntity"
	"test-va/internals/entity/userEntity"
	"test-va/internals/service/cryptoService"
	"test-va/internals/service/timeSrv"
	"test-va/internals/service/validationService"
	"time"
)

type UserSrv interface {
	SaveUser(req *userEntity.CreateUserReq) (*userEntity.CreateUserRes, *ResponseEntity.ResponseMessage)
}

type userSrv struct {
	repo      userRepo.UserRepository
	validator validationService.ValidationSrv
	timeSrv   timeSrv.TimeService
	cryptoSrv cryptoService.CryptoSrv
}

func (u *userSrv) SaveUser(req *userEntity.CreateUserReq) (*userEntity.CreateUserRes, *ResponseEntity.ResponseMessage) {
	// validate request
	err := u.validator.Validate(req)
	if err != nil {
		return nil, ResponseEntity.NewCustomError(http.StatusBadRequest, "Bad Input Request")
	}
	// check if user with that email exists already

	//------GET BY EMAIL----------

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
	return data, nil
}

func NewUserSrv(repo userRepo.UserRepository, validator validationService.ValidationSrv, timeSrv timeSrv.TimeService, cryptoSrv cryptoService.CryptoSrv) UserSrv {
	return &userSrv{repo: repo, validator: validator, timeSrv: timeSrv, cryptoSrv: cryptoSrv}
}
