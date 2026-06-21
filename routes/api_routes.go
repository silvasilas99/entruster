package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/hyperledger/fabric-gateway/pkg/client"
	"github.com/silvasilas99/entruster/domain/metadata"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func SetupRoutes(contract *client.Contract) *gin.Engine {
	r := gin.Default()

	// Swagger UI route
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	metadataRoutes := r.Group("/api/metadata")
	{
		metadataRoutes.POST("/", metadata.CreateMetadataHandler(contract))
		metadataRoutes.GET("/", metadata.GetAllMetadataHandler(contract))
		metadataRoutes.GET("/:id", metadata.GetMetadataByIDHandler(contract))
		metadataRoutes.PUT("/:id", metadata.UpdateMetadataByIDHandler(contract))
		metadataRoutes.DELETE("/:id", metadata.DeleteMetadataByIDHandler(contract))
		metadataRoutes.GET("/:id/auditory", metadata.GetMetadataAuditoryByIDHandler(contract))
	}

	return r
}