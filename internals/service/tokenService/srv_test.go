package tokenservice

import (
	"log"
	"testing"
)

func Test_token(t *testing.T) {
	tokensrv := NewTokenSrv("vbvkvjbkv")

	token, rtoken, err := tokensrv.createToken("eb@gmail.com","555")

	log.Println(token)
	log.Println(rtoken)
	log.Println(err)

	 y,s := tokensrv.validateToken(token)

	 log.Println(y,"y")
	 log.Println(s, "s")

}


func Test_voken(t *testing.T) {
	tokensrv := NewTokenSrv("vbvkvjbkv")

	 token,err := tokensrv.validateToken("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJFbWFpbCI6ImViYkBnbWFpbC5jb20iLCJJZCI6IjU1NTUiLCJleHAiOjE2NjkyODY2MzB9.vGueXgEMpKpW28ortTI6_0lOcb9wulaBFpX08_D7lr0")

	 log.Println(token,"y")
	 log.Println(err, "s")

}


