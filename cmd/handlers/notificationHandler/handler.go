package notificationHandler

import (
	"net/http"
	"test-va/internals/entity/ResponseEntity"
	"test-va/internals/entity/notificationEntity"
	"test-va/internals/service/notificationService"

	"github.com/gin-gonic/gin"
)

type notificationHandler struct {
	srv notificationService.NotificationSrv
}

func NewNotificationHandler(srv notificationService.NotificationSrv) *notificationHandler {
	return &notificationHandler{srv: srv}
}

func (n *notificationHandler) RegisterForNotifications(c *gin.Context) {
	var req notificationEntity.CreateNotification

	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest,
			ResponseEntity.BuildErrorResponse(http.StatusBadRequest, "Bad Request", err, nil))
		return
	}

	errorRes := n.srv.RegisterForNotifications(&req)
	if errorRes != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized,
			ResponseEntity.BuildErrorResponse(http.StatusUnauthorized, "Authorization Error", errorRes, nil))
		return
	}

	c.JSON(http.StatusOK, ResponseEntity.BuildSuccessResponse(200, "Notification Registration Successful", "", nil))
}
