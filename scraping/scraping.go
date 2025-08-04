package scraping

import (
	"net/http"
	"time"
)

const AGENT = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36"

type Scrapper interface {
	Search(productName string, minPrice, maxPrice float64) ([]Product, error)
	GetSiteName() string
}

// ScrapperBase contém propriedades comuns a todos os scrapers
type ScrapperBase struct {
	Client    *http.Client
	BaseURL   string
	RateLimit time.Duration
	SiteName  string
}

// NewScrapperBase cria uma nova instância base com configurações padrão
func NewScrapperBase(baseURL, siteName string) ScrapperBase {
	return ScrapperBase{
		Client:    &http.Client{Timeout: 30 * time.Second},
		BaseURL:   baseURL,
		RateLimit: 2 * time.Second,
		SiteName:  siteName,
	}
}

type Product struct {
	Name      string  `json:"name"`
	Price     float64 `json:"price"`
	Rating    float32 `json:"rating"`
	SoldCount uint    `json:"sold_count"`
	Link      string  `json:"link"`
	Site      string  `json:"site"`
}
