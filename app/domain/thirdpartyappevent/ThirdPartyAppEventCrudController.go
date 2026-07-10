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

func (c *ThirdPartyAppEventCrudController) Store() gin.HandlerFunc {
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

func (c *ThirdPartyAppEventCrudController) UpdateThirdPartyAppEventByIDHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id := ctx.Param("id")
		if id == "" {
			utils.SendError(ctx, http.StatusBadRequest, "ID parameter is required")
			return
		}
		var req ThirdPartyAppEventModel
		if err := ctx.ShouldBindJSON(&req); err != nil {
			utils.SendError(ctx, http.StatusBadRequest, "Invalid request body. Check the required fields and their types.")
			return
		}

		// Set dynamic UpdatedBy from authenticated user info
		if user, exists := ctx.Get("currentUser"); exists {
			if userInfo, ok := user.(*middleware.UserInfo); ok {
				req.UpdatedBy = userInfo.Name
			}
		}

		if err := c.thirdPartyAppEventService.UpdateThirdPartyAppEventByID(id, req); err != nil {
			utils.SendError(ctx, http.StatusInternalServerError, err.Error())
			return
		}
		utils.SendSuccess(ctx, "ThirdPartyAppEvent updated successfully", gin.H{"id": id})
	}
}

func (c *ThirdPartyAppEventCrudController) DeleteByID() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id := ctx.Param("id")
		if id == "" {
			utils.SendError(ctx, http.StatusBadRequest, "ID parameter is required")
			return
		}
		deletedBy := "system"
		if user, exists := ctx.Get("currentUser"); exists {
			if userInfo, ok := user.(*middleware.UserInfo); ok {
				deletedBy = userInfo.Name
			}
		}

		if err := c.thirdPartyAppEventService.DeleteThirdPartyAppEventByID(id, deletedBy); err != nil {
			utils.SendError(ctx, http.StatusInternalServerError, err.Error())
			return
		}
		utils.SendSuccess(ctx, "ThirdPartyAppEvent deleted successfully", gin.H{"id": id})
	}
}
