package dto

type Pagination struct {
	Next          int `json:"next,omitempty"`
	Previous      int `json:"previous,omitempty"`
	RecordPerPage int `json:"record_per_page,omitempty"`
	CurrentPage   int `json:"current_page,omitempty"`
	TotalPage     int `json:"total_page,omitempty"`
}
