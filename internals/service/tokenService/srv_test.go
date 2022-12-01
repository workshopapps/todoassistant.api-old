package tokenservice

import (
	"log"
	"testing"
)

func Test_token(t *testing.T) {
	tokensrv := NewTokenSrv("vbvkvjbkv")

	token, rtoken, err := tokensrv.CreateToken("eb@gmail.com","user","555")
	log.Println(token)
	log.Println(rtoken)
	log.Println(err)

}

func Test_voken(t *testing.T) {
	tokensrv := NewTokenSrv("vbvkvjbkv")

	token, err := tokensrv.ValidateToken("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJFbWFpbCI6ImViYkBnbWFpbC5jb20iLCJJZCI6IjU1NTUiLCJleHAiOjE2NjkyODY2MzB9.vGueXgEMpKpW28ortTI6_0lOcb9wulaBFpX08_D7lr0")

	log.Println(token, "y")
	log.Println(err, "s")

}
