package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/hyperledger/fabric-gateway/pkg/client"
	"github.com/silvasilas99/entruster/app/core/audit"
	"github.com/silvasilas99/entruster/app/domain/metadata"
	"github.com/silvasilas99/entruster/app/core/elasticsearch"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func SetupRoutes(contract *client.Contract) *gin.Engine {
	r := gin.Default()

	// Bootstrap the audit pipeline:
	//   AuditService  ← persists & queries audit entries
	//   MetadataObserver ← translates metadata events into audit.Record calls
	auditSvc := audit.NewAuditService()
	elasticSvc := elasticsearch.NewElasticService()
	metadataObserver := metadata.NewMetadataObserver(auditSvc, elasticSvc)

	crudCtrl := metadata.NewMetadataCrudController(contract, metadataObserver)

	// Swagger UI route
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	metadataRoutes := r.Group("/api/metadata")
	{
		metadataRoutes.POST("/", crudCtrl.Store())
		metadataRoutes.GET("/", metadata.GetAll(contract, metadataObserver, elasticSvc))
		metadataRoutes.GET("/:id", metadata.GetMetadataByIDHandler(contract))
		metadataRoutes.PUT("/:id", crudCtrl.UpdateMetadataByIDHandler())
		metadataRoutes.DELETE("/:id", crudCtrl.DeleteByID())
		metadataRoutes.GET("/:id/auditory", metadata.GetMetadataAuditoryByIDHandler(auditSvc))
	}

	return r
}