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
        metadataRoutes.POST("/", metadata.createMetadata(contract))
        metadataRoutes.GET("/", metadata.getAllMetadata(contract))
        // coming next:
        // metadataRoutes.GET("/:id",       metadata.getMetadataByDID(contract))
        // metadataRoutes.PUT("/:id",    	metadata.updateMetadata(contract))
        // metadataRoutes.DELETE("/:id", 	metadata.deleteMetadata(contract))
    }
    return r
}