package handlers

import (
	"net/http"
	"strconv"

	"flypro/internal/dto"

	"github.com/gin-gonic/gin"
)

// ---- Expenses ----
func (h *Handler) CreateExpense(c *gin.Context) {
	var req dto.CreateExpenseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	e, err := h.expenses.CreateExpense(c, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, e)
}

func (h *Handler) ListExpenses(c *gin.Context) {
	var q dto.ListQuery
	if err := c.ShouldBindQuery(&q); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// For demo, get user_id from query; in real app from auth context
	userIDStr := c.Query("user_id")
	id, _ := strconv.Atoi(userIDStr)

	items, total, err := h.expenses.ListExpenses(c, uint(id), q)
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

func (h *Handler) GetExpense(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	e, err := h.expenses.GetExpense(c, uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	c.JSON(http.StatusOK, e)
}

func (h *Handler) UpdateExpense(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var req dto.UpdateExpenseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	e, err := h.expenses.UpdateExpense(c, uint(id), req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, e)
}

func (h *Handler) DeleteExpense(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	if err := h.expenses.DeleteExpense(c, uint(id)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}
