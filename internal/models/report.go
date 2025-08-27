package models

import "time"

type ExpenseReport struct {
	ID        uint      `json:"id" gorm:"primarykey"`
	UserID    uint      `json:"user_id" gorm:"not null;index"`
	Title     string    `json:"title" gorm:"not null"`
	Status    string    `json:"status" gorm:"type:varchar(32);default:'draft';index"`
	Total     float64   `json:"total"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// Relations
	User     User      `json:"user" gorm:"foreignKey:UserID"`
	Expenses []Expense `json:"expenses" gorm:"many2many:report_expenses;"`
}
