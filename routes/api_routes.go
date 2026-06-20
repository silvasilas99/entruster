package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/hyperledger/fabric-gateway/pkg/client"
	"github.com/silvasilas99/entruster/domain/metadata"
)

func SetupRoutes(contract *client.Contract) *gin.Engine {
	r := gin.Default()
	metadataRoutes := r.Group("/api/metadata")
	{
		metadataRoutes.POST("/", metadata.CreateMetadataHandler(contract))
		metadataRoutes.GET("/", metadata.GetAllMetadataHandler(contract))
		// coming next:
		// metadataRoutes.GET("/:id",       metadata.GetMetadataByDIDHandler(contract))
		// metadataRoutes.PUT("/:id",    	metadata.UpdateMetadataHandler(contract))
		// metadataRoutes.DELETE("/:id", 	metadata.DeleteMetadataHandler(contract))
	}
	return r
}