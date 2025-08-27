-- +goose Up
CREATE INDEX IF NOT EXISTS idx_expenses_user_id ON expenses(user_id);
CREATE INDEX IF NOT EXISTS idx_expenses_category ON expenses(category);
CREATE INDEX IF NOT EXISTS idx_expenses_status ON expenses(status);

CREATE INDEX IF NOT EXISTS idx_reports_user_id ON expense_reports(user_id);
CREATE INDEX IF NOT EXISTS idx_reports_status ON expense_reports(status);

-- +goose Down
DROP INDEX IF EXISTS idx_expenses_user_id;
DROP INDEX IF EXISTS idx_expenses_category;
DROP INDEX IF EXISTS idx_expenses_status;
DROP INDEX IF EXISTS idx_reports_user_id;
DROP INDEX IF EXISTS idx_reports_status;
