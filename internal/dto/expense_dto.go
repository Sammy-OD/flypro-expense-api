package dto

type CreateExpenseRequest struct {
	Amount      float64 `json:"amount" binding:"required,gt=0"`
	Currency    string  `json:"currency" binding:"required,len=3,currency"`
	Category    string  `json:"category" binding:"required,category"`
	Description string  `json:"description" binding:"max=500"`
	Receipt     string  `json:"receipt"`
	UserID      uint    `json:"user_id" binding:"required"`
}

type UpdateExpenseRequest struct {
	Amount      *float64 `json:"amount" binding:"omitempty,gt=0"`
	Currency    *string  `json:"currency" binding:"omitempty,len=3,currency"`
	Category    *string  `json:"category" binding:"omitempty,category"`
	Description *string  `json:"description" binding:"omitempty,max=500"`
	Receipt     *string  `json:"receipt"`
	Status      *string  `json:"status" binding:"omitempty,oneof=pending approved rejected"`
}
