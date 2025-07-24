package scraping

import (
	"time"
)

type Scrapper interface {
	Search(productName string, minPrice, maxPrice float64) ([]Product, error)
	GetSiteName() string
	GetRateLimit() time.Duration
	CanSearchProduct(productName string) bool
	IsAvailable() bool
	GetMaxRetries() int
}

type Product struct {
	Name      string  `json:"name"`
	Price     float32 `json:"price"`
	Rating    float32 `json:"rating"`
	SoldCount uint    `json:"sold_count"`
	Link      string  `json:"link"`
	Site      string  `json:"site"`
}
