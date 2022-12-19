package routes

import (
	"test-va/cmd/handlers/callHandler"
	"test-va/internals/service/callService"

	"github.com/gin-gonic/gin"
)

func CallRoute(v1 *gin.RouterGroup, srv callService.CallService) {

	callHandler := callHandler.NewCallHandler(srv)

	v1.GET("/calls", callHandler.GetCalls)

}
