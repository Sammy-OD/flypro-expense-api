package repository

import (
	"errors"

	"flypro/internal/models"

	"gorm.io/gorm"
)

type ReportRepository interface {
	Create(report *models.ExpenseReport) error
	List(userID uint, page, pageSize int) ([]models.ExpenseReport, int64, error)
	GetByID(id uint) (*models.ExpenseReport, error)
	AddExpenses(reportID uint, expenseIDs []uint) error
	Update(report *models.ExpenseReport) error
}

type reportRepository struct{ db *gorm.DB }

func NewReportRepository(db *gorm.DB) ReportRepository { return &reportRepository{db: db} }

func (r *reportRepository) Create(rep *models.ExpenseReport) error {
	return r.db.Create(rep).Error
}

func (r *reportRepository) List(userID uint, page, pageSize int) ([]models.ExpenseReport, int64, error) {
	var items []models.ExpenseReport
	q := r.db.Model(&models.ExpenseReport{}).Where("user_id = ?", userID)
	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	offset := (page - 1) * pageSize
	if err := q.Preload("User").Preload("Expenses").Order("created_at DESC").Limit(pageSize).Offset(offset).Find(&items).Error; err != nil {
		return nil, 0, err
	}
	return items, total, nil
}

func (r *reportRepository) GetByID(id uint) (*models.ExpenseReport, error) {
	var rep models.ExpenseReport
	if err := r.db.Preload("Expenses").Preload("Expenses.User").First(&rep, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, err
	}
	return &rep, nil
}

func (r *reportRepository) AddExpenses(reportID uint, expenseIDs []uint) error {
	for _, eid := range expenseIDs {
		if err := r.db.Exec("INSERT INTO report_expenses (expense_report_id, expense_id) VALUES (?, ?) ON CONFLICT DO NOTHING", reportID, eid).Error; err != nil {
			return err
		}
	}
	return nil
}

func (r *reportRepository) Update(rep *models.ExpenseReport) error {
	return r.db.Save(rep).Error
}
