package socialLoginService

import (
	"test-va/internals/Repository/userRepo"
	"test-va/internals/entity/ResponseEntity"
	"test-va/internals/entity/userEntity"
	"test-va/internals/service/timeSrv"
	tokenservice "test-va/internals/service/tokenService"
	"time"

	"github.com/google/uuid"
)


type LoginSrv interface{
	LoginResponse(req *userEntity.GoogleLoginReq) (*userEntity.LoginRes, *ResponseEntity.ServiceError)
}

type loginSrv struct {
	repo      userRepo.UserRepository
	timeSrv   timeSrv.TimeService
}




func (t *loginSrv)LoginResponse(req *userEntity.GoogleLoginReq) (*userEntity.LoginRes, *ResponseEntity.ServiceError) {



	user, _ := t.repo.GetByEmail(req.Email)

	if user == nil{

		resData := &userEntity.CreateUserReq{
			UserId: uuid.New().String(),      
    		FirstName: req.FirstName,     
    		LastName: req.LastName,      
    		Email: req.Email,               
    		DateOfBirth: "2022-11-19T12:56:04+01:00",   
    		AccountStatus: "ACTIVE", 
    		DateCreated: t.timeSrv.CurrentTime().Format(time.RFC3339),   
		}
		
		err := t.repo.Persist(resData)

		if err != nil {
			return nil, ResponseEntity.NewInternalServiceError(err)
		}
	}


	user, _ = t.repo.GetByEmail(req.Email)
	tokenSrv := tokenservice.NewTokenSrv("fvmvmvmvf")

	accessToken, refreshToken, err:= tokenSrv.CreateToken(user.Email,"user",user.UserId)

	if err != nil {
		return nil, ResponseEntity.NewInternalServiceError(err)
	}




	loginUser := &userEntity.LoginRes{
		UserId: user.UserId,
		Email: user.Email,
		FirstName: user.FirstName,
		LastName: user.LastName,
		Phone: user.Phone,
		Gender: user.Gender,
		Token: accessToken,
		RefreshToken: refreshToken,
	}

	return loginUser, nil
}


func NewLoginSrv(repo userRepo.UserRepository, timeSrv timeSrv.TimeService) LoginSrv {
	return &loginSrv{repo, timeSrv}
}
