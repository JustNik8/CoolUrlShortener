package shortener

import "strings"

const (
	alphabet = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
)

var (
	alphabetLen = len(alphabet)
)

//go:generate go run github.com/vektra/mockery/v2@v2.42.1 --name URLShortener
type URLShortener interface {
	ShortenURL(id uint32) string
}

type base62UrlShortener struct {
}

func NewBase62UrlShortener() URLShortener {
	return &base62UrlShortener{}
}

func (s *base62UrlShortener) ShortenURL(id uint32) string {
	nums := make([]int, 0)
	for id > 0 {
		rem := int(id) % alphabetLen
		nums = append(nums, rem)

		id /= uint32(alphabetLen)
	}

	var sb strings.Builder
	for i := range nums {
		idx := nums[i]
		sb.WriteByte(alphabet[idx])
	}

	return sb.String()
}
