package tokenservice

import (
	"github.com/golang-jwt/jwt/v4"
	"time"
	"log"
)


type Token struct{
	Email string
	Id 	  string
	jwt.StandardClaims
}

type TokenSrv interface{
	createToken(email string, id string) (string, string, error)
	validateToken(token string) (*Token,error)
}

type tokenSrv struct {
	SecretKey string 
}


func (t *tokenSrv) createToken(email string, id string) (string, string, error){
	tokenDetails := &Token{
		Email: email,
		Id: id,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(24)).Unix(),
		},
	}

	refreshTokenDetails := &Token{
		Email: email,
		Id: id,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(30)).Unix(),
		},
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, tokenDetails).SignedString([]byte(t.SecretKey))
	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshTokenDetails).SignedString([]byte(t.SecretKey))

	if err != nil {
		log.Panic(err)
		return "","",err
	}

	return token, refreshToken, err
}

func (t *tokenSrv) validateToken( tokenUrl string) (*Token, error){
		token, err := jwt.ParseWithClaims(
			tokenUrl,
			&Token{},
			func(token *jwt.Token) (interface{}, error) {
				return []byte(t.SecretKey), nil
			},
		)
		
		
		claims, ok := token.Claims.(*Token)
		if !ok {
			return nil,err
		}
	
		if claims.ExpiresAt < time.Now().Local().Unix() {
			return nil, err
		}
	
		return claims, err
	
}


func NewTokenSrv(token string) TokenSrv {
	return &tokenSrv{token}
}
