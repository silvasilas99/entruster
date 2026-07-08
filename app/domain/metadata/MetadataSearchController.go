package metadata

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hyperledger/fabric-gateway/pkg/client"
	"github.com/silvasilas99/entruster/app/core/audit"
	"github.com/silvasilas99/entruster/app/core/elasticsearch"
	"github.com/silvasilas99/entruster/utils"
)

// GetAll handles GET /api/metadata/
//
//	@Summary		List all metadata assets
//	@Tags			metadata
//	@Produce		json
//	@Param			patient_id	query		string	false	"Filter by patient ID"
//	@Param			asset_id	query		string	false	"Filter by asset ID"
//	@Param			from		query		string	false	"Start date (created_at)"
//	@Param			to			query		string	false	"End date (created_at)"
//	@Success		200	{object}	utils.SuccessResponse{data=[]map[string]interface{}}
//	@Failure		500	{object}	utils.ErrorResponse
//	@Router			/metadata/ [get]
func GetAll(contract *client.Contract, observer *MetadataObserver, elasticSvc *elasticsearch.ElasticService) gin.HandlerFunc {
	return func(c *gin.Context) {
		patientID := c.Query("patient_id")
		assetID := c.Query("asset_id")
		from := c.Query("from")
		to := c.Query("to")

		var mustClauses []map[string]interface{}

		if patientID != "" {
			mustClauses = append(mustClauses, map[string]interface{}{
				"match": map[string]interface{}{"patient_id": patientID},
			})
		}
		if assetID != "" {
			mustClauses = append(mustClauses, map[string]interface{}{
				"match": map[string]interface{}{"asset_id": assetID},
			})
		}

		if from != "" || to != "" {
			rangeQuery := map[string]interface{}{}
			if from != "" {
				rangeQuery["gte"] = from
			}
			if to != "" {
				rangeQuery["lte"] = to
			}
			mustClauses = append(mustClauses, map[string]interface{}{
				"range": map[string]interface{}{
					"created_at": rangeQuery,
				},
			})
		}

		// Filter out soft-deleted records in ES
		mustNotClauses := []map[string]interface{}{
			{
				"exists": map[string]interface{}{
					"field": "deleted_at",
				},
			},
		}

		query := map[string]interface{}{
			"query": map[string]interface{}{
				"bool": map[string]interface{}{
					"must":     mustClauses,
					"must_not": mustNotClauses,
				},
			},
		}

		results, err := elasticSvc.SearchDocuments(c.Request.Context(), "metadata", query)
		if err != nil {
			utils.SendError(c, http.StatusInternalServerError, "Failed to fetch metadata from search engine: "+err.Error())
			return
		}

		if observer != nil {
			observer.OnList(len(results))
		}

		utils.SendSuccess(c, "Metadata list retrieved successfully", results)
	}
}

// GetMetadataByIDHandler handles GET /api/metadata/:id
//
//	@Summary		Get metadata by ID
//	@Description	Evaluates GetMetadataById on the Fabric ledger and returns the matching metadata asset.
//	@Tags			metadata
//	@Produce		json
//	@Param			id	path		string	true	"Asset ID (numeric string, e.g. \"1\")"
//	@Success		200	{object}	utils.SuccessResponse{data=MetadataModel}
//	@Failure		400	{object}	utils.ErrorResponse
//	@Failure		500	{object}	utils.ErrorResponse
//	@Router			/metadata/{id} [get]
func GetMetadataByIDHandler(contract *client.Contract) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		if id == "" {
			utils.SendError(c, http.StatusBadRequest, "ID parameter is required")
			return
		}
		service := NewMetadataService(contract, nil)
		m, err := service.GetMetadataByID(id)
		if err != nil {
			utils.SendError(c, http.StatusInternalServerError, err.Error())
			return
		}
		utils.SendSuccess(c, "Metadata retrieved successfully", m)
	}
}

// GetMetadataAuditoryByIDHandler handles GET /api/metadata/:id/auditory
//
//	@Summary		Get audit trail for a metadata asset
//	@Description	Returns the full immutable history of the asset from the audit service.
//	@Tags			metadata
//	@Produce		json
//	@Param			id	path		string	true	"Asset ID (numeric string)"
//	@Success		200	{object}	utils.SuccessResponse{data=[]audit.AuditModel}
//	@Failure		400	{object}	utils.ErrorResponse
//	@Router			/metadata/{id}/auditory [get]
func GetMetadataAuditoryByIDHandler(auditSvc *audit.AuditService) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		if id == "" {
			utils.SendError(c, http.StatusBadRequest, "ID parameter is required")
			return
		}
		history := auditSvc.GetByEntityID("Metadata", id)
		utils.SendSuccess(c, "Metadata audit trail retrieved successfully", history)
	}
}
