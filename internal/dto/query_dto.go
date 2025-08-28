package dto

type ListQuery struct {
	Page     int    `form:"page,default=1" binding:"min=1"`
	PageSize int    `form:"page_size,default=20" binding:"min=1,max=100"`
	Category string `form:"category" binding:"omitempty,category"`
	Status   string `form:"status" binding:"omitempty,oneof=pending approved rejected"`
}
