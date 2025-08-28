package handlers

import (
	"flypro/internal/services"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	users    services.UserService
	expenses services.ExpenseService
	reports  services.ReportService
}

func NewHandler(u services.UserService, e services.ExpenseService, r services.ReportService) *Handler {
	return &Handler{users: u, expenses: e, reports: r}
}

func (h *Handler) RegisterRoutes(r *gin.Engine) {
	api := r.Group("/api")

	api.POST("/users", h.CreateUser)
	api.GET("/users/:id", h.GetUser)

	api.POST("/expenses", h.CreateExpense)
	api.GET("/expenses", h.ListExpenses)
	api.GET("/expenses/:id", h.GetExpense)
	api.PUT("/expenses/:id", h.UpdateExpense)
	api.DELETE("/expenses/:id", h.DeleteExpense)

	api.POST("/reports", h.CreateReport)
	api.GET("/reports", h.ListReports)
	api.POST("/reports/:id/expenses", h.AddExpensesToReport)
	api.PUT("/reports/:id/submit", h.SubmitReport)
}
