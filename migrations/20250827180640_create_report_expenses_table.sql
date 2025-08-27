-- +goose Up
CREATE TABLE IF NOT EXISTS report_expenses (
  expense_report_id INT NOT NULL REFERENCES expense_reports(id) ON DELETE CASCADE,
  expense_id INT NOT NULL REFERENCES expenses(id) ON DELETE CASCADE,
  PRIMARY KEY (expense_report_id, expense_id)
);

-- +goose Down
DROP TABLE IF EXISTS report_expenses;
