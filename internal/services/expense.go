package services

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"flypro/internal/dto"
	"flypro/internal/models"
	"flypro/internal/repository"
	"flypro/internal/utils"

	"github.com/redis/go-redis/v9"
)

type ExpenseService interface {
	CreateExpense(ctx context.Context, req dto.CreateExpenseRequest) (*models.Expense, error)
	ListExpenses(ctx context.Context, userID uint, q dto.ListQuery) ([]models.Expense, int64, error)
	GetExpense(ctx context.Context, id uint) (*models.Expense, error)
	UpdateExpense(ctx context.Context, id uint, req dto.UpdateExpenseRequest) (*models.Expense, error)
	DeleteExpense(ctx context.Context, id uint) error
}

type expenseService struct {
	repo     repository.ExpenseRepository
	currency CurrencyService
	cache    *redis.Client
}

func NewExpenseService(repo repository.ExpenseRepository, currency CurrencyService, cache *redis.Client) ExpenseService {
	return &expenseService{repo: repo, currency: currency, cache: cache}
}

func (s *expenseService) CreateExpense(ctx context.Context, req dto.CreateExpenseRequest) (*models.Expense, error) {
	e := &models.Expense{
		UserID:      req.UserID,
		Amount:      req.Amount,
		Currency:    req.Currency,
		Category:    req.Category,
		Description: req.Description,
		Receipt:     req.Receipt,
		Status:      "pending",
	}
	if err := s.repo.Create(e); err != nil {
		return nil, err
	}
	_ = s.cache.Del(ctx, s.listCacheKey(req.UserID, dto.ListQuery{})).Err()
	return e, nil
}

func (s *expenseService) ListExpenses(ctx context.Context, userID uint, q dto.ListQuery) ([]models.Expense, int64, error) {
	key := s.listCacheKey(userID, q)
	if b, err := s.cache.Get(ctx, key).Bytes(); err == nil {
		var cached struct {
			Items []models.Expense
			Total int64
		}
		if json.Unmarshal(b, &cached) == nil {
			return cached.Items, cached.Total, nil
		}
	}
	items, total, err := s.repo.List(userID, q.Page, q.PageSize, q.Category, q.Status)
	if err != nil {
		return nil, 0, err
	}
	if b, _ := json.Marshal(struct {
		Items []models.Expense
		Total int64
	}{items, total}); b != nil {
		_ = s.cache.Set(ctx, key, b, 5*time.Minute).Err()
	}
	return items, total, nil
}

func (s *expenseService) GetExpense(ctx context.Context, id uint) (*models.Expense, error) {
	return s.repo.GetByID(id)
}

func (s *expenseService) UpdateExpense(ctx context.Context, id uint, req dto.UpdateExpenseRequest) (*models.Expense, error) {
	e, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if req.Amount != nil {
		e.Amount = *req.Amount
	}
	if req.Currency != nil {
		e.Currency = *req.Currency
	}
	if req.Category != nil {
		e.Category = *req.Category
	}
	if req.Description != nil {
		e.Description = *req.Description
	}
	if req.Receipt != nil {
		e.Receipt = *req.Receipt
	}
	if req.Status != nil {
		e.Status = *req.Status
	}
	if err := s.repo.Update(e); err != nil {
		return nil, err
	}
	_ = s.cache.FlushDB(ctx).Err() // simple invalidation demo
	return e, nil
}

func (s *expenseService) DeleteExpense(ctx context.Context, id uint) error {
	return s.repo.Delete(id)
}

func (s *expenseService) listCacheKey(userID uint, q dto.ListQuery) string {
	return fmt.Sprintf("expenses:%d:%d:%d:%s:%s", userID, q.Page, q.PageSize, q.Category, q.Status)
}

var _ = utils.ErrNotFound // keep linter happy for utils import
