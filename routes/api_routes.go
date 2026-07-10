package routes

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/hyperledger/fabric-gateway/pkg/client"
	"github.com/silvasilas99/entruster/app/core/audit"
	"github.com/silvasilas99/entruster/app/core/chaincode"
	"github.com/silvasilas99/entruster/app/core/middleware"
	"github.com/silvasilas99/entruster/app/domain/thirdpartyappevent"
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
	elasticSvc := elasticsearch.NewElasticService()
	chaincodeQuery := chaincode.NewChaincodeQuery(contract)
	auditSvc := audit.NewAuditService(elasticSvc, chaincodeQuery)
	metadataObserver := metadata.NewMetadataObserver(auditSvc, elasticSvc)

	crudCtrl := metadata.NewMetadataCrudController(contract, metadataObserver)
	thirdPartyAppEventObserver := thirdpartyappevent.NewThirdPartyAppEventObserver(auditSvc, elasticSvc)
	thirdPartyAppEventCrudCtrl := thirdpartyappevent.NewThirdPartyAppEventCrudController(contract, thirdPartyAppEventObserver)


	// Swagger UI route
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Mock User API (Fictitious API called by JWTAuth middleware)
	r.GET("/api/mock/user", MockUserHandler)

	metadataRoutes := r.Group("/api/metadata")
	metadataRoutes.Use(middleware.JWTAuth())
	{
		metadataRoutes.POST("/", crudCtrl.Store())
		metadataRoutes.GET("/", metadata.GetAll(contract, metadataObserver, elasticSvc))
		metadataRoutes.GET("/:id", metadata.GetMetadataByIDHandler(contract))
		metadataRoutes.PUT("/:id", crudCtrl.UpdateMetadataByIDHandler())
		metadataRoutes.DELETE("/:id", crudCtrl.DeleteByID())
		metadataRoutes.GET("/:id/auditory", metadata.GetMetadataAuditoryByIDHandler(auditSvc))
		metadataRoutes.GET("/:id/history", metadata.GetMetadataNativeHistoryByIDHandler(auditSvc))
	}

	thirdpartyappeventRoutes := r.Group("/api/thirdpartyappevent")
	thirdpartyappeventRoutes.Use(middleware.JWTAuth())
	{
		thirdpartyappeventRoutes.POST("/", thirdPartyAppEventCrudCtrl.Store())
		thirdpartyappeventRoutes.GET("/", thirdpartyappevent.GetAll(contract, thirdPartyAppEventObserver, elasticSvc))
		thirdpartyappeventRoutes.GET("/:id", thirdpartyappevent.GetThirdPartyAppEventByIDHandler(contract))
		thirdpartyappeventRoutes.PUT("/:id", thirdPartyAppEventCrudCtrl.UpdateThirdPartyAppEventByIDHandler())
		thirdpartyappeventRoutes.DELETE("/:id", thirdPartyAppEventCrudCtrl.DeleteByID())
		thirdpartyappeventRoutes.GET("/:id/auditory", thirdpartyappevent.GetThirdPartyAppEventAuditoryByIDHandler(auditSvc))
		thirdpartyappeventRoutes.GET("/:id/history", thirdpartyappevent.GetThirdPartyAppEventNativeHistoryByIDHandler(auditSvc))
	}

	// Health check routes
	r.GET("/api/health/elasticsearch", metadata.GetHealthHandler(elasticSvc))

	return r
}

// MockUserHandler provides a fictitious API to return user details based on the token.
func MockUserHandler(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	token := ""
	if strings.HasPrefix(authHeader, "Bearer ") {
		token = strings.TrimPrefix(authHeader, "Bearer ")
	}

	// Default mock user returned
	user := gin.H{
		"id":    "usr_777",
		"name":  "Dr. Silas Silva",
		"email": "silas.silva@hospital.com",
		"role":  "Chief Clinician",
	}

	// We can change response based on different tokens if needed
	if token == "some-other-token-example" {
		user = gin.H{
			"id":    "usr_888",
			"name":  "Jane Doe",
			"email": "jane.doe@hospital.com",
			"role":  "Researcher",
		}
	}

	c.JSON(200, user)
}