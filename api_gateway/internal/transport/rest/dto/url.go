package dto

type TopURLData struct {
	LongURL     string `json:"long_url"`
	ShortURL    string `json:"short_url"`
	FollowCount int64  `json:"follow_count"`
	CreateCount int64  `json:"create_count"`
}

type TopURLDataResponse struct {
	TopURLData []TopURLData `json:"top_url_data"`
	Pagination Pagination   `json:"pagination"`
}

type LongURLData struct {
	LongURL string `json:"long_url"`
}

type URlData struct {
	LongURL  string `json:"long_url"`
	ShortURL string `json:"short_url"`
}
