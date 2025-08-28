package main

import (
	"fmt"

	"flypro/internal/config"
	"flypro/internal/models"
	"flypro/internal/repository"

	"gorm.io/gorm/clause"
)

func main() {
	cfg := config.Load()
	db := repository.MustOpenGorm(cfg)

	fmt.Println("Seeding demo data...")

	// Seed user (only insert if doesn't exist)
	u := models.User{Email: "johndoe@flypro.io", Name: "John Doe"}
	db.Clauses(clause.OnConflict{DoNothing: true}).Create(&u)

	// If the user already exists, fetch it so UserID is correct
	db.Where("email = ?", "johndoe@flypro.io").First(&u)

	// Seed expenses (skip duplicates)
	exp := []models.Expense{
		{UserID: u.ID, Amount: 125.50, Currency: "USD", Category: "travel", Description: "Flight ticket", Status: "approved"},
		{UserID: u.ID, Amount: 40.00, Currency: "USD", Category: "meals", Description: "Lunch", Status: "pending"},
	}
	for _, e := range exp {
		db.Clauses(clause.OnConflict{DoNothing: true}).Create(&e)
	}

	// Seed report
	report := models.ExpenseReport{UserID: u.ID, Title: "NYC Trip", Status: "draft", Total: 0}
	db.Clauses(clause.OnConflict{DoNothing: true}).Create(&report)

	fmt.Println("Done.")
}
