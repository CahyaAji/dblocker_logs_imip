package handlers

import (
	"dblocker_logs_server/internal/models"
	"dblocker_logs_server/internal/repository"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ActionLogHandler struct {
	Repo *repository.ActionLogRepository
}

func NewActionLogHandler(repo *repository.ActionLogRepository) *ActionLogHandler {
	return &ActionLogHandler{Repo: repo}
}

func (h *ActionLogHandler) CreateActionLog(c *gin.Context) {
	var input models.ActionLog
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.Repo.Create(&input); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": input})
}

func (h *ActionLogHandler) GetActionLogs(c *gin.Context) {
	logs, err := h.Repo.FindAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": logs})
}

func (h *ActionLogHandler) GetActionLogByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID format"})
		return
	}

	actionLog, err := h.Repo.FindByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Action log not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": actionLog})
}

func (h *ActionLogHandler) UpdateActionLog(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var input models.ActionLog
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	input.ID = uint(id) // Ensure ID matches URL

	if err := h.Repo.Update(&input); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": input})
}

func (h *ActionLogHandler) DeleteActionLog(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	if err := h.Repo.Delete(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Action log deleted successfully"})
}
