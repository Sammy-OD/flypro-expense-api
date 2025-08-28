package handlers

import (
	"net/http"
	"strconv"

	"flypro/internal/dto"

	"github.com/gin-gonic/gin"
)

// ---- Reports ----
func (h *Handler) CreateReport(c *gin.Context) {
	var req dto.CreateReportRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	rep, err := h.reports.CreateReport(c, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, rep)
}

func (h *Handler) ListReports(c *gin.Context) {
	var q dto.ListQuery
	if err := c.ShouldBindQuery(&q); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	userIDStr := c.Query("user_id")
	id, _ := strconv.Atoi(userIDStr)

	items, total, err := h.reports.ListReports(c, uint(id), q)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"items":     items,
		"total":     total,
		"page":      q.Page,
		"page_size": q.PageSize,
	})
}

func (h *Handler) AddExpensesToReport(c *gin.Context) {
	reportID, _ := strconv.Atoi(c.Param("id"))
	var req dto.AddExpensesToReportRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	userIDStr := c.Query("user_id")
	uid, _ := strconv.Atoi(userIDStr)

	rep, err := h.reports.AddExpenses(c, uint(reportID), req.ExpenseIDs, uint(uid))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, rep)
}

func (h *Handler) SubmitReport(c *gin.Context) {
	reportID, _ := strconv.Atoi(c.Param("id"))
	userIDStr := c.Query("user_id")
	uid, _ := strconv.Atoi(userIDStr)

	rep, err := h.reports.Submit(c, uint(reportID), uint(uid))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, rep)
}
