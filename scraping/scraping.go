package scraping

type Scrapper interface {
}


type Product struct {
	Name      string  `json:"name"`
	Price     float32 `json:"price"`
	Rating    float32 `json:"rating"`
	SoldCount uint    `json:"sold_count"`
	Link      string  `json:"link"`
	Site      string  `json:"site"`
}
