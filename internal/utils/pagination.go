package utils

type PageResult[T any] struct {
	Items []T `json:"items"`
	Page int `json:"page"`
	PageSize int `json:"page_size"`
	Total int64 `json:"total"`
}
