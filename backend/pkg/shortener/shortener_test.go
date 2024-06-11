package shortener

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestBase62UrlShortener(t *testing.T) {
	base62UrlShortener := NewBase62UrlShortener()
	t.Run("is idempotent", func(t *testing.T) {

		id := uuid.New().ID()
		expectedLongURL := base62UrlShortener.ShortenURL(id)

		for i := 0; i < 1000; i++ {
			longURL := base62UrlShortener.ShortenURL(id)
			assert.Equal(t, expectedLongURL, longURL)
		}
	})
}

func BenchmarkNewBase62UrlShortener(b *testing.B) {
	base62UrlShortener := NewBase62UrlShortener()
	id := uuid.New().ID()

	for n := 0; n < b.N; n++ {
		_ = base62UrlShortener.ShortenURL(id)
	}
}
