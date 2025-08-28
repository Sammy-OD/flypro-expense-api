package repository

import (
	"flypro/internal/models"

	"gorm.io/gorm"
)

type ExpenseRepository interface {
	Create(expense *models.Expense) error
	List(userID uint, page, pageSize int, category, status string) ([]models.Expense, int64, error)
	GetByID(id uint) (*models.Expense, error)
	Update(expense *models.Expense) error
	Delete(id uint) error
}

type expenseRepository struct{ db *gorm.DB }

func NewExpenseRepository(db *gorm.DB) ExpenseRepository { return &expenseRepository{db: db} }

func (r *expenseRepository) Create(e *models.Expense) error {
	return r.db.Create(e).Error
}

func (r *expenseRepository) List(userID uint, page, pageSize int, category, status string) ([]models.Expense, int64, error) {
	var items []models.Expense
	q := r.db.Model(&models.Expense{}).Where("user_id = ?", userID)
	if category != "" {
		q = q.Where("category = ?", category)
	}
	if status != "" {
		q = q.Where("status = ?", status)
	}
	q = q.Preload("User")
	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	offset := (page - 1) * pageSize
	if err := q.Order("created_at DESC").Limit(pageSize).Offset(offset).Find(&items).Error; err != nil {
		return nil, 0, err
	}
	return items, total, nil
}

func (r *expenseRepository) GetByID(id uint) (*models.Expense, error) {
	var e models.Expense
	if err := r.db.Preload("User").First(&e, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, err
	}
	return &e, nil
}

func (r *expenseRepository) Update(e *models.Expense) error {
	return r.db.Save(e).Error
}

func (r *expenseRepository) Delete(id uint) error {
	return r.db.Delete(&models.Expense{}, id).Error
}
