package thirdpartyappevent

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hyperledger/fabric-gateway/pkg/client"
	"github.com/silvasilas99/entruster/app/core/middleware"
	"github.com/silvasilas99/entruster/utils"
)

type ThirdPartyAppEventCrudController struct {
	thirdPartyAppEventService *ThirdPartyAppEventService
}

func NewThirdPartyAppEventCrudController(contract *client.Contract, observer *ThirdPartyAppEventObserver) *ThirdPartyAppEventCrudController {
	service := NewThirdPartyAppEventService(contract, observer)

	return &ThirdPartyAppEventCrudController{
		thirdPartyAppEventService: service,
	}
}

func (c *ThirdPartyAppEventCrudController) StoreHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var thirdPartyAppEventDTO ThirdPartyAppEventDTO

		if err := ctx.ShouldBindJSON(&thirdPartyAppEventDTO); err != nil {
			utils.SendError(
				ctx,
				http.StatusBadRequest,
				"Invalid request body. Check the required fields and their types.",
			)
			return
		}

		// Set dynamic CreatedBy from authenticated user info
		if user, exists := ctx.Get("currentUser"); exists {
			if userInfo, ok := user.(*middleware.UserInfo); ok {
				thirdPartyAppEventDTO.CreatedBy = userInfo.Name
			}
		}

		if err := c.thirdPartyAppEventService.RegisterThirdPartyAppEvent(thirdPartyAppEventDTO); err != nil {
			utils.SendError(ctx, http.StatusInternalServerError, err.Error())
			return
		}

		utils.SendSuccess(ctx, "ThirdPartyAppEvent registered on blockchain", gin.H{
			"id":         thirdPartyAppEventDTO.ID,
			"app_id":     thirdPartyAppEventDTO.AppID,
			"event_type": thirdPartyAppEventDTO.EventType,
		})
	}
}