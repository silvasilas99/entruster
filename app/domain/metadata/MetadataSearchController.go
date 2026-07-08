package metadata

import (
	"net/http"
	"strconv"
	"time"

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
//	@Param			q			query		string	false	"Full-text search query"
//	@Param			from		query		string	false	"Start date (created_at, RFC3339)"
//	@Param			to			query		string	false	"End date (created_at, RFC3339)"
//	@Param			limit		query		int		false	"Max results (default 20, max 100)"
//	@Param			offset		query		int		false	"Offset for pagination"
//	@Success		200	{object}	utils.SuccessResponse{data=[]map[string]interface{}}
//	@Failure		500	{object}	utils.ErrorResponse
//	@Router			/metadata/ [get]
func GetAll(contract *client.Contract, observer *MetadataObserver, elasticSvc *elasticsearch.ElasticService) gin.HandlerFunc {
	return func(c *gin.Context) {
		filter := &elasticsearch.MetadataFilter{
			PatientID:      c.Query("patient_id"),
			AssetID:        c.Query("asset_id"),
			Query:          c.Query("q"),
			ExcludeDeleted: true, // Always exclude soft-deleted records
		}

		// Parse date range
		if from := c.Query("from"); from != "" {
			if t, err := time.Parse(time.RFC3339, from); err == nil {
				filter.CreatedFrom = &t
			}
		}
		if to := c.Query("to"); to != "" {
			if t, err := time.Parse(time.RFC3339, to); err == nil {
				filter.CreatedTo = &t
			}
		}

		// Parse pagination
		if limit := c.Query("limit"); limit != "" {
			if v, err := strconv.Atoi(limit); err == nil {
				filter.Limit = v
			}
		}
		if offset := c.Query("offset"); offset != "" {
			if v, err := strconv.Atoi(offset); err == nil {
				filter.Offset = v
			}
		}

		results, total, err := elasticSvc.Search(c.Request.Context(), esIndexName, filter)
		if err != nil {
			utils.SendError(c, http.StatusInternalServerError, "Failed to fetch metadata from search engine: "+err.Error())
			return
		}

		if observer != nil {
			observer.OnList(len(results))
		}

		utils.SendSuccess(c, "Metadata list retrieved successfully", gin.H{
			"items": results,
			"total": total,
		})
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

// GetHealthHandler handles GET /api/health/elasticsearch
//
//	@Summary		Get Elasticsearch cluster health
//	@Tags			health
//	@Produce		json
//	@Success		200	{object}	utils.SuccessResponse{data=map[string]interface{}}
//	@Failure		500	{object}	utils.ErrorResponse
//	@Router			/health/elasticsearch [get]
func GetHealthHandler(elasticSvc *elasticsearch.ElasticService) gin.HandlerFunc {
	return func(c *gin.Context) {
		health, err := elasticSvc.GetHealth(c.Request.Context())
		if err != nil {
			utils.SendError(c, http.StatusInternalServerError, "Failed to get Elasticsearch health: "+err.Error())
			return
		}
		utils.SendSuccess(c, "Elasticsearch health retrieved", health)
	}
}
