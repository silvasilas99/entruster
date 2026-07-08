package metadata

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hyperledger/fabric-gateway/pkg/client"
	"github.com/silvasilas99/entruster/app/core/middleware"
	"github.com/silvasilas99/entruster/utils"
)

type MetadataCrudController struct {
	metadataService *MetadataService
}

func NewMetadataCrudController(contract *client.Contract, observer *MetadataObserver) *MetadataCrudController {
	service := NewMetadataService(contract, observer)

	return &MetadataCrudController{
		metadataService: service,
	}
}

func (c *MetadataCrudController) Store() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var metadataDTO MetadataDTO

		if err := ctx.ShouldBindJSON(&metadataDTO); err != nil {
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
				metadataDTO.CreatedBy = userInfo.Name
			}
		}

		if err := c.metadataService.RegisterMetadata(metadataDTO); err != nil {
			utils.SendError(ctx, http.StatusInternalServerError, err.Error())
			return
		}

		utils.SendSuccess(ctx, "Metadata registered on blockchain", gin.H{
			"id":         metadataDTO.ID,
			"patient_id": metadataDTO.PatientID,
			"asset_id":   metadataDTO.AssetID,
		})
	}
}

func (c *MetadataCrudController) UpdateMetadataByIDHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id := ctx.Param("id")
		if id == "" {
			utils.SendError(ctx, http.StatusBadRequest, "ID parameter is required")
			return
		}
		var req MetadataModel
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

		if err := c.metadataService.UpdateMetadataByID(id, req); err != nil {
			utils.SendError(ctx, http.StatusInternalServerError, err.Error())
			return
		}
		utils.SendSuccess(ctx, "Metadata updated successfully", gin.H{"id": id})
	}
}

func (c *MetadataCrudController) DeleteByID() gin.HandlerFunc {
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

		if err := c.metadataService.DeleteMetadataByID(id, deletedBy); err != nil {
			utils.SendError(ctx, http.StatusInternalServerError, err.Error())
			return
		}
		utils.SendSuccess(ctx, "Metadata deleted successfully", gin.H{"id": id})
	}
}
