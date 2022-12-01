package socialLoginService

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"test-va/internals/Repository/userRepo"
	"test-va/internals/entity/userEntity"
	tokenservice "test-va/internals/service/tokenService"
	"test-va/utils"
	// "golang.org/x/oauth2"
)



type LoginSrv interface{
	GetAuthUrl() string
	LoginResponse(code string) *userEntity.LoginRes
}

type loginSrv struct {
	repo      userRepo.UserRepository
}


func (t *loginSrv)GetAuthUrl() string {

	
	u := utils.LoginConfig.GoogleLoginConfig.AuthCodeURL("state")
	return u
}

func (t *loginSrv)LoginResponse(code string) *userEntity.LoginRes {
	token, err := utils.LoginConfig.GoogleLoginConfig.Exchange(
		context.Background(), code)

		if err != nil{
			log.Panic(err)
		}
         
		response, err := http.Get(utils.GoogleAPI + token.AccessToken)

		if err != nil {
			log.Panic(err)
		}

		data , err:= io.ReadAll(response.Body)


		if err != nil {
			log.Panic(err)
		}



	if err != nil {
		log.Panic(err)
	}

	

	defer response.Body.Close() 


	var result userEntity.GoogleLoginRes

	

	if err := json.Unmarshal(data, &result); err != nil { 
		if err != nil {
			log.Panic(err)
		}
	
    }

	user, err := t.repo.GetByEmail(result.Email)

	if err != nil {
		log.Panic(err)
	}




	tokenSrv := tokenservice.NewTokenSrv("fvmvmvmvf")

	accessToken, refreshToken, err:= tokenSrv.CreateToken(user.Email,"user",user.UserId)

	if err != nil {
		log.Panic(err)
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

	return loginUser
}


func NewLoginSrv(repo userRepo.UserRepository) LoginSrv {
	return &loginSrv{repo}
}
