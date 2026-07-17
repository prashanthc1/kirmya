package connections

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"workspace-app/internal/common"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

type requestConnectionInput struct {
	TargetUserID string            `json:"target_user_id" binding:"required"`
	Note         *string           `json:"note"`
	Source       *ConnectionSource `json:"source"`
}

func (h *Handler) RequestConnection(c *gin.Context) {
	var input requestConnectionInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    "VALIDATION_FAILED",
				"message": "Invalid request payload: " + err.Error(),
			},
		})
		return
	}

	currentUserID := common.UserIDFromContext(c.Request.Context())
	if currentUserID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": gin.H{
				"code":    "UNAUTHORIZED",
				"message": "Authentication required",
			},
		})
		return
	}

	err := h.service.SendRequest(c.Request.Context(), currentUserID, input.TargetUserID, input.Note, input.Source)
	if err != nil {
		h.writeError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"success": true,
			"message": "Connection request sent",
		},
	})
}

func (h *Handler) AcceptConnection(c *gin.Context) {
	connectionID := c.Param("id")
	if connectionID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    "VALIDATION_FAILED",
				"message": "Connection ID is required",
			},
		})
		return
	}

	currentUserID := common.UserIDFromContext(c.Request.Context())
	err := h.service.AcceptRequest(c.Request.Context(), connectionID, currentUserID)
	if err != nil {
		h.writeError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"success": true,
			"message": "Connection request accepted",
		},
	})
}

func (h *Handler) DeclineConnection(c *gin.Context) {
	connectionID := c.Param("id")
	if connectionID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    "VALIDATION_FAILED",
				"message": "Connection ID is required",
			},
		})
		return
	}

	currentUserID := common.UserIDFromContext(c.Request.Context())
	err := h.service.DeclineRequest(c.Request.Context(), connectionID, currentUserID)
	if err != nil {
		h.writeError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"success": true,
			"message": "Connection request declined",
		},
	})
}

func (h *Handler) RemoveConnection(c *gin.Context) {
	connectionID := c.Param("id")
	if connectionID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    "VALIDATION_FAILED",
				"message": "Connection ID is required",
			},
		})
		return
	}

	currentUserID := common.UserIDFromContext(c.Request.Context())
	err := h.service.RemoveConnection(c.Request.Context(), connectionID, currentUserID)
	if err != nil {
		h.writeError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"success": true,
			"message": "Connection removed",
		},
	})
}

type blockUserInput struct {
	Reason string `json:"reason"`
}

func (h *Handler) BlockUser(c *gin.Context) {
	targetUserID := c.Param("user_id")
	if targetUserID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    "VALIDATION_FAILED",
				"message": "User ID is required",
			},
		})
		return
	}

	var input blockUserInput
	_ = c.ShouldBindJSON(&input) // Reason is optional

	currentUserID := common.UserIDFromContext(c.Request.Context())
	err := h.service.BlockUser(c.Request.Context(), currentUserID, targetUserID, input.Reason)
	if err != nil {
		h.writeError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"success": true,
			"message": "User blocked",
		},
	})
}

func (h *Handler) UnblockUser(c *gin.Context) {
	targetUserID := c.Param("user_id")
	if targetUserID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    "VALIDATION_FAILED",
				"message": "User ID is required",
			},
		})
		return
	}

	currentUserID := common.UserIDFromContext(c.Request.Context())
	err := h.service.UnblockUser(c.Request.Context(), currentUserID, targetUserID)
	if err != nil {
		h.writeError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"success": true,
			"message": "User unblocked",
		},
	})
}

func (h *Handler) GetConnections(c *gin.Context) {
	currentUserID := common.UserIDFromContext(c.Request.Context())

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	conns, err := h.service.repo.GetConnections(c.Request.Context(), currentUserID, page, limit)
	if err != nil {
		h.writeError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": conns,
	})
}

func (h *Handler) GetPendingRequests(c *gin.Context) {
	currentUserID := common.UserIDFromContext(c.Request.Context())
	direction := c.Query("direction") // "incoming" or "outgoing"
	if direction != "incoming" && direction != "outgoing" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    "VALIDATION_FAILED",
				"message": "Direction must be 'incoming' or 'outgoing'",
			},
		})
		return
	}

	reqs, err := h.service.repo.GetPendingRequests(c.Request.Context(), currentUserID, direction)
	if err != nil {
		h.writeError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": reqs,
	})
}

func (h *Handler) GetMutualConnections(c *gin.Context) {
	currentUserID := common.UserIDFromContext(c.Request.Context())
	targetUserID := c.Param("user_id")
	if targetUserID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    "VALIDATION_FAILED",
				"message": "User ID is required",
			},
		})
		return
	}

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if limit < 1 {
		limit = 10
	}

	users, totalCount, err := h.service.repo.GetMutualConnections(c.Request.Context(), currentUserID, targetUserID, limit)
	if err != nil {
		h.writeError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"users": users,
			"total": totalCount,
		},
	})
}

func (h *Handler) GetSuggestions(c *gin.Context) {
	currentUserID := common.UserIDFromContext(c.Request.Context())

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if limit < 1 {
		limit = 10
	}

	suggestions, err := h.service.repo.GetSuggestions(c.Request.Context(), currentUserID, limit)
	if err != nil {
		h.writeError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": suggestions,
	})
}

func (h *Handler) writeError(c *gin.Context, err error) {
	mapped := MapError(err)
	if appErr, ok := mapped.(*common.AppError); ok {
		// Set Retry-After header for rate limits and cooldowns if applicable
		switch appErr.Code {
		case "RATE_LIMITED":
			c.Header("Retry-After", "86400") // 24h cooldown retry after
		case "COOLDOWN_ACTIVE":
			c.Header("Retry-After", "2592000") // 30 days cooldown retry after
		}

		c.JSON(appErr.Status, gin.H{
			"error": gin.H{
				"code":    appErr.Code,
				"message": appErr.Message,
			},
		})
		return
	}

	c.JSON(http.StatusInternalServerError, gin.H{
		"error": gin.H{
			"code":    "INTERNAL_ERROR",
			"message": err.Error(),
		},
	})
}
