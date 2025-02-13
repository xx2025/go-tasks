package service

type Pagination struct {
	Page  int           `json:"page"`
	Total int64         `json:"total"`
	Limit int           `json:"limit"`
	Items []interface{} `json:"items"`
}
