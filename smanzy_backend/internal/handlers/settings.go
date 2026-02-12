package handlers

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ristep/smanzy_backend/internal/db"
)

type SettingsHandler struct {
	conn    *sql.DB
	queries *db.Queries
}

func NewSettingsHandler(conn *sql.DB, queries *db.Queries) *SettingsHandler {
	return &SettingsHandler{
		conn:    conn,
		queries: queries,
	}
}

type UpdateSettingRequest struct {
	Value string `json:"value" binding:"required"`
}

func (sh *SettingsHandler) GetSettingHandler(c *gin.Context) {
	key := c.Param("key")
	value, err := sh.queries.GetSetting(c.Request.Context(), key)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, ErrorResponse{Error: "Setting not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Database error"})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{Data: gin.H{"key": key, "value": value}})
}

func (sh *SettingsHandler) ListSettingsHandler(c *gin.Context) {
	settings, err := sh.queries.ListSettings(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Database error"})
		return
	}

	settingsMap := make(map[string]string)
	for _, s := range settings {
		settingsMap[s.Key] = s.Value
	}

	c.JSON(http.StatusOK, SuccessResponse{Data: settingsMap})
}

func (sh *SettingsHandler) UpdateSettingHandler(c *gin.Context) {
	key := c.Param("key")
	var req UpdateSettingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid input"})
		return
	}

	setting, err := sh.queries.UpsertSetting(c.Request.Context(), db.UpsertSettingParams{
		Key:   key,
		Value: req.Value,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to update setting"})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{Data: setting})
}
