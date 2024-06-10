package dto

type LongURLData struct {
	LongURL string `json:"long_url" validate:"required"`
}

type URlData struct {
	LongURL  string `json:"long_url"`
	ShortURL string `json:"short_url"`
}
