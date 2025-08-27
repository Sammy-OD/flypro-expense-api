package handlers

import (
	"flypro/internal/services"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	users services.UserService
}

func NewHandler(u services.UserService) *Handler {
	return &Handler{users: u}
}

func (h *Handler) RegisterRoutes(r *gin.Engine) {
	api := r.Group("/api")

	api.POST("/users", h.CreateUser)
	api.GET("/users/:id", h.GetUser)
}
