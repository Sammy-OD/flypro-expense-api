package models

import "time"

type Expense struct {
	ID          uint      `json:"id" gorm:"primarykey"`
	UserID      uint      `json:"user_id" gorm:"not null;index"`
	Amount      float64   `json:"amount" gorm:"not null"`
	Currency    string    `json:"currency" gorm:"type:char(3);not null;index"`
	Category    string    `json:"category" gorm:"type:varchar(32);not null;index"`
	Description string    `json:"description"`
	Receipt     string    `json:"receipt"`
	Status      string    `json:"status" gorm:"type:varchar(32);default:'pending';index"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`

	// Relations
	User User `json:"user" gorm:"foreignKey:UserID"`
}
