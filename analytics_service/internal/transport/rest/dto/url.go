package dto

type TopURLData struct {
	LongURL     string `json:"long_url"`
	ShortURL    string `json:"short_url"`
	FollowCount int64  `json:"follow_count"`
	CreateCount int64  `json:"create_count"`
}
