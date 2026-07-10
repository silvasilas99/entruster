package thirdpartyappevent

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hyperledger/fabric-gateway/pkg/client"
	"github.com/silvasilas99/entruster/app/core/audit"
	"github.com/silvasilas99/entruster/app/core/elasticsearch"
	"github.com/silvasilas99/entruster/app/core/middleware"
	"github.com/silvasilas99/entruster/utils"
)

// GetAll handles GET /api/thirdpartyappevent/
//
//	@Summary		List all thirdpartyappevent assets
//	@Tags			thirdpartyappevent
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
//	@Router			/thirdpartyappevent/ [get]
func GetAll(contract *client.Contract, observer *ThirdPartyAppEventObserver, elasticSvc *elasticsearch.ElasticService) gin.HandlerFunc {
	return func(c *gin.Context) {
		filter := &elasticsearch.MetadataFilter{
			PatientID:      c.Query("app_id"), // Mapping app_id to PatientID field temporarily for Elasticsearch
			AssetID:        c.Query("event_type"),
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
			utils.SendError(c, http.StatusInternalServerError, "Failed to fetch thirdpartyappevent from search engine: "+err.Error())
			return
		}

		if observer != nil {
			actor := "system"
			if user, exists := c.Get("currentUser"); exists {
				if userInfo, ok := user.(*middleware.UserInfo); ok {
					actor = userInfo.Name
				}
			}
			observer.OnList(len(results), actor)
		}

		utils.SendSuccess(c, "ThirdPartyAppEvent list retrieved successfully", gin.H{
			"items": results,
			"total": total,
		})
	}
}

// GetThirdPartyAppEventByIDHandler handles GET /api/thirdpartyappevent/:id
//
//	@Summary		Get thirdpartyappevent by ID
//	@Description	Evaluates GetThirdPartyAppEventById on the Fabric ledger and returns the matching thirdpartyappevent asset.
//	@Tags			thirdpartyappevent
//	@Produce		json
//	@Param			id	path		string	true	"Asset ID (numeric string, e.g. \"1\")"
//	@Success		200	{object}	utils.SuccessResponse{data=ThirdPartyAppEventModel}
//	@Failure		400	{object}	utils.ErrorResponse
//	@Failure		500	{object}	utils.ErrorResponse
//	@Router			/thirdpartyappevent/{id} [get]
func GetThirdPartyAppEventByIDHandler(contract *client.Contract) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		if id == "" {
			utils.SendError(c, http.StatusBadRequest, "ID parameter is required")
			return
		}
		service := NewThirdPartyAppEventService(contract, nil)
		m, err := service.GetThirdPartyAppEventByID(id)
		if err != nil {
			utils.SendError(c, http.StatusInternalServerError, err.Error())
			return
		}
		utils.SendSuccess(c, "ThirdPartyAppEvent retrieved successfully", m)
	}
}

// GetThirdPartyAppEventAuditoryByIDHandler handles GET /api/thirdpartyappevent/:id/auditory
//
//	@Summary		Get audit trail for a thirdpartyappevent asset
//	@Description	Returns the full immutable history of the asset from the audit service.
//	@Tags			thirdpartyappevent
//	@Produce		json
//	@Param			id	path		string	true	"Asset ID (numeric string)"
//	@Success		200	{object}	utils.SuccessResponse{data=[]audit.AuditModel}
//	@Failure		400	{object}	utils.ErrorResponse
//	@Router			/thirdpartyappevent/{id}/auditory [get]
func GetThirdPartyAppEventAuditoryByIDHandler(auditSvc *audit.AuditService) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		if id == "" {
			utils.SendError(c, http.StatusBadRequest, "ID parameter is required")
			return
		}
		history := auditSvc.GetByEntityID("ThirdPartyAppEvent", id)
		utils.SendSuccess(c, "ThirdPartyAppEvent audit trail retrieved successfully", history)
	}
}

// GetThirdPartyAppEventNativeHistoryByIDHandler handles GET /api/thirdpartyappevent/:id/history
//
//	@Summary		Get native transaction history for a thirdpartyappevent asset
//	@Description	Returns the immutable transaction history of the asset directly from the Hyperledger ledger.
//	@Tags			thirdpartyappevent
//	@Produce		json
//	@Param			id	path		string	true	"Asset ID (numeric string)"
//	@Success		200	{object}	utils.SuccessResponse{data=[]audit.HistoryRecord}
//	@Failure		400	{object}	utils.ErrorResponse
//	@Router			/thirdpartyappevent/{id}/history [get]
func GetThirdPartyAppEventNativeHistoryByIDHandler(auditSvc *audit.AuditService) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		if id == "" {
			utils.SendError(c, http.StatusBadRequest, "ID parameter is required")
			return
		}
		history, err := auditSvc.GetNativeHistory(id)
		if err != nil {
			utils.SendError(c, http.StatusInternalServerError, "Failed to retrieve native history: "+err.Error())
			return
		}

		var parsedHistory []map[string]interface{}
		json.Unmarshal(history, &parsedHistory)

		utils.SendSuccess(c, "ThirdPartyAppEvent native ledger history retrieved successfully", parsedHistory)
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
