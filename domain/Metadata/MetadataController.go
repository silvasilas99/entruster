package metadata

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hyperledger/fabric-gateway/pkg/client"
	"github.com/silvasilas99/entruster/utils"
)

func CreateMetadataHandler(contract *client.Contract) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req MetadataModel
		if err := c.ShouldBindJSON(&req); err != nil {
			utils.SendError(c, http.StatusBadRequest, "Invalid request body. Check the required fields and their types.")
			return
		}
		// Step 2 — call service layer
		if err := RegisterMetadata(contract, req); err != nil {
			utils.SendError(c, http.StatusInternalServerError, err.Error())
			return
		}
		// Step 3 — send success response
		utils.SendSuccess(c, "Metadata registered on blockchain", gin.H{
			"id": req.ID,
		})
	}
}

func GetAllMetadataHandler(contract *client.Contract) gin.HandlerFunc {
	return func(c *gin.Context) {
		metadataList, err := GetAllMetadata(contract)
		if err != nil {
			utils.SendError(c, http.StatusInternalServerError, err.Error())
			return
		}
		utils.SendSuccess(c, "Metadata list retrieved successfully", metadataList)
	}
}

func GetMetadataByDIDHandler(contract *client.Contract) gin.HandlerFunc {
	return func(c *gin.Context) {
		utils.SendError(c, http.StatusNotImplemented, "GetMetadataByDID is not implemented yet")
		// id := c.Param("id")
		// if id == "" {
		// 	utils.SendError(c, http.StatusBadRequest, "ID parameter is required")
		// 	return
		// }
		// metadata, err := GetMetadataByDID(contract, id)
		// if err != nil {
		// 	utils.SendError(c, http.StatusInternalServerError, err.Error())
		// 	return
		// }
		// utils.SendSuccess(c, "Metadata retrieved successfully", metadata)
	}
}

func UpdateMetadataHandler(contract *client.Contract) gin.HandlerFunc {
	return func(c *gin.Context) {
		utils.SendError(c, http.StatusNotImplemented, "UpdateMetadata is not implemented yet")
		// id := c.Param("id")
		// if id == "" {
		// 	utils.SendError(c, http.StatusBadRequest, "ID parameter is required")
		// 	return
		// }
		// metadata, err := GetMetadataByDID(contract, id)
		// if err != nil {
		// 	utils.SendError(c, http.StatusInternalServerError, err.Error())
		// 	return
		// }
		// utils.SendSuccess(c, "Metadata updated successfully", metadata)
	}
}

func DeleteMetadataHandler(contract *client.Contract) gin.HandlerFunc {
	return func(c *gin.Context) {
		utils.SendError(c, http.StatusNotImplemented, "DeleteMetadata is not implemented yet")
		// id := c.Param("id")
		// if id == "" {
		// 	utils.SendError(c, http.StatusBadRequest, "ID parameter is required")
		// 	return
		// }
		// err := DeleteMetadata(contract, id)
		// if err != nil {
		// 	utils.SendError(c, http.StatusInternalServerError, err.Error())
		// 	return
		// }
		// utils.SendSuccess(c, "Metadata deleted successfully", nil)
	}
}

func ExportMetadataAsCsvHandler(contract *client.Contract) gin.HandlerFunc {
	return func(c *gin.Context) {
		utils.SendError(c, http.StatusNotImplemented, "ExportMetadataAsCsv is not implemented yet")
	}
}
