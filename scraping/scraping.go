package scraping

const AGENT = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36"

type Scrapper interface {
	Search(productName string, minPrice, maxPrice float64) ([]Product, error)

	GetSiteName() string
}

type Product struct {
	Name      string  `json:"name"`
	Price     float64 `json:"price"`
	Rating    float32 `json:"rating"`
	SoldCount uint    `json:"sold_count"`
	Link      string  `json:"link"`
	Site      string  `json:"site"`
}
