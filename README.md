# Market Scraping

Web scraping API to search products across multiple e-commerce sites with price range filtering.

## 📋 Features

- Search products by name and price range
- Concurrent scraping across multiple sites
- Smart caching with Redis
- Rate limiting per site
- Simple and efficient REST API

## 🚀 Technologies

- **Go** - Main language (standard library)
- **Redis** - Caching and rate limiting
- **Goroutines** - Concurrency with worker pool
- **Scraping** - Colly
- **net/http** - HTTP server and client

## 📦 Installation

```bash
# Clone the repository
git clone https://github.com/your-username/market-scraping.git
cd market-scraping

# Initialize Go module
go mod init market-scraping


# Run the project
go run main.go
```

## 🔧 Usage

### Search products

```bash
POST /api/search
{
    "product": "smartphone",
    "min_price": 500.00,
    "max_price": 1500.00
}
```

### Response

```json
{
    "products": [
        {
            "name": "iPhone 14",
            "price": 1299.99,
            "rating": 4.5,
            "sold_count": 1500,
            "link": "https://...",
            "site": "mercadolivre"
        }
    ],
    "total_found": 25,
    "search_time": "2.3s"
}
```

## 🎯 Supported Sites

- Mercado Livre
- Americanas
- Casas Bahia
- Magazine Luiza
- Amazon 

## 📝 License

MIT License