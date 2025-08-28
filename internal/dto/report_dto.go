package dto

type CreateReportRequest struct {
	UserID uint   `json:"user_id" binding:"required"`
	Title  string `json:"title" binding:"required,min=2"`
}

type AddExpensesToReportRequest struct {
	ExpenseIDs []uint `json:"expense_ids" binding:"required,dive,gt=0"`
}
