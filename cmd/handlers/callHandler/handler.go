package callHandler

import (
	"net/http"
	"test-va/internals/entity/ResponseEntity"
	"test-va/internals/service/callService"

	"github.com/gin-gonic/gin"
)

type callHandler struct {
	srv callService.CallService
}

func NewCallHandler(srv callService.CallService) *callHandler {
	return &callHandler{srv: srv}
}


func (t *callHandler) GetCalls(c *gin.Context) {

	calls, errRes := t.srv.GetCalls()
	if errRes != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, ResponseEntity.BuildErrorResponse(http.StatusBadRequest, "error getting calls", errRes, nil))
		return
	}

	length := len(calls)

	if length == 0 {
		c.AbortWithStatusJSON(http.StatusOK, ResponseEntity.BuildSuccessResponse(http.StatusOK, "no calls found",calls,nil))
		return
	}

	c.JSON(http.StatusOK,ResponseEntity.BuildSuccessResponse(http.StatusOK, "successfully fetched calls and details",calls,nil))
}
