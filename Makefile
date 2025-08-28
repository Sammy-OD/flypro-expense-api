SHELL := /bin/bash

APP_NAME=flypro
DB_URL?=postgres://postgres:postgres@localhost:5432/flypro?sslmode=disable

run:
	go run ./cmd/server

seed:
	go run ./cmd/seed

migrate-up:
	goose -dir ./migrations postgres "$(DB_URL)" up

migrate-down:
	goose -dir ./migrations postgres "$(DB_URL)" down

migrate-status:
	goose -dir ./migrations postgres "$(DB_URL)" status

new-migration:
	@test -n "$(name)" || (echo "Usage: make new-migration name=descriptive_name" && exit 1)
	goose -dir ./migrations create $(name) sql
