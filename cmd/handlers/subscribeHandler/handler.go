package subscribeHandler

import (
	"net/http"
	"test-va/internals/entity/ResponseEntity"
	"test-va/internals/entity/subscribeEntity"
	"test-va/internals/service/subscribeService"

	"github.com/gin-gonic/gin"
)

type subscribeHandler struct {
	srv subscribeService.SubscribeService
}

func NewSubscribeHandler(srv subscribeService.SubscribeService) *subscribeHandler {
	return &subscribeHandler{srv: srv}
}

func (t *subscribeHandler) AddSubscriber(c *gin.Context) {
	// create a request of subsctibeEntity type
	var req subscribeEntity.SubscribeReq

	// copy data from gin context to req
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest,
			ResponseEntity.BuildErrorResponse(http.StatusBadRequest, "email field required", err, nil))
		return
	}
	// call function from service that saves email to DB
	response, errRes := t.srv.PersistEmail(&req)
	if errRes != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest,
			ResponseEntity.BuildErrorResponse(http.StatusBadRequest, "error adding email to list", errRes, nil))
		return
	}
	c.JSON(http.StatusOK, response)
}

func (t *subscribeHandler) ContactUs(c *gin.Context) {
	var req subscribeEntity.ContactUsReq

	// copy data from gin context to req
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest,
			ResponseEntity.BuildErrorResponse(http.StatusBadRequest, "Invalid request", err, nil))
		return
	}
	// call function from service that saves email to DB
	response, errRes := t.srv.Contact(&req)
	if errRes != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError,
			ResponseEntity.BuildErrorResponse(http.StatusInternalServerError, "Error sending email", errRes, nil))
		return
	}
	c.JSON(http.StatusOK, response)
}
