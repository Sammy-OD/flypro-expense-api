package services

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"flypro/internal/dto"
	"flypro/internal/models"
	"flypro/internal/repository"

	"github.com/redis/go-redis/v9"
)

type ReportService interface {
	CreateReport(ctx context.Context, req dto.CreateReportRequest) (*models.ExpenseReport, error)
	ListReports(ctx context.Context, userID uint, q dto.ListQuery) ([]models.ExpenseReport, int64, error)
	AddExpenses(ctx context.Context, reportID uint, expenseIDs []uint, userID uint) (*models.ExpenseReport, error)
	Submit(ctx context.Context, reportID uint, userID uint) (*models.ExpenseReport, error)
}

type reportService struct {
	repo        repository.ReportRepository
	expenseRepo repository.ExpenseRepository
	currency    CurrencyService
	cache       *redis.Client
}

func NewReportService(repo repository.ReportRepository, expenseRepo repository.ExpenseRepository, currency CurrencyService, cache *redis.Client) ReportService {
	return &reportService{repo: repo, expenseRepo: expenseRepo, currency: currency, cache: cache}
}

func (s *reportService) CreateReport(ctx context.Context, req dto.CreateReportRequest) (*models.ExpenseReport, error) {
	rep := &models.ExpenseReport{UserID: req.UserID, Title: req.Title, Status: "draft"}
	if err := s.repo.Create(rep); err != nil {
		return nil, err
	}
	return rep, nil
}

func (s *reportService) ListReports(ctx context.Context, userID uint, q dto.ListQuery) ([]models.ExpenseReport, int64, error) {
	key := fmt.Sprintf("report_summaries:%d:%d:%d", userID, q.Page, q.PageSize)
	if b, err := s.cache.Get(ctx, key).Bytes(); err == nil {
		var cached struct {
			Items []models.ExpenseReport
			Total int64
		}
		if json.Unmarshal(b, &cached) == nil {
			return cached.Items, cached.Total, nil
		}
	}
	items, total, err := s.repo.List(userID, q.Page, q.PageSize)
	if err != nil {
		return nil, 0, err
	}
	if b, _ := json.Marshal(struct {
		Items []models.ExpenseReport
		Total int64
	}{items, total}); b != nil {
		_ = s.cache.Set(ctx, key, b, 30*time.Minute).Err()
	}
	return items, total, nil
}

func (s *reportService) AddExpenses(ctx context.Context, reportID uint, expenseIDs []uint, userID uint) (*models.ExpenseReport, error) {
	// validate ownership
	for _, id := range expenseIDs {
		e, err := s.expenseRepo.GetByID(id)
		if err != nil {
			return nil, err
		}
		if e.UserID != userID {
			return nil, fmt.Errorf("expense %d not owned by user", id)
		}
	}
	if err := s.repo.AddExpenses(reportID, expenseIDs); err != nil {
		return nil, err
	}
	return s.repo.GetByID(reportID)
}

func (s *reportService) Submit(ctx context.Context, reportID uint, userID uint) (*models.ExpenseReport, error) {
	rep, err := s.repo.GetByID(reportID)
	if err != nil {
		return nil, err
	}
	if rep.UserID != userID {
		return nil, fmt.Errorf("not owner")
	}
	if rep.Status != "draft" {
		return nil, fmt.Errorf("report not in draft state")
	}

	// Recompute total in USD
	totalUSD := 0.0
	for _, e := range rep.Expenses {
		amtUSD, err := s.currency.ConvertCurrency(ctx, e.Amount, e.Currency, "USD")
		if err != nil {
			return nil, err
		}
		totalUSD += amtUSD
	}
	rep.Total = totalUSD
	rep.Status = "submitted"
	if err := s.repo.Update(rep); err != nil {
		return nil, err
	}
	return rep, nil
}
