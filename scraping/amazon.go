package scraping

import (
	"fmt"
	"strconv"
	"strings"
	"sync"

	"github.com/gocolly/colly/v2"
)

type amazonScrapper struct {
	ScrapperBase
}

func NewAmazonScrapper() *amazonScrapper {
	base := NewScrapperBase("https://www.amazon.com.br/s?k=", "Amazon")
	return &amazonScrapper{
		ScrapperBase: base,
	}
}

// extractProduct extrai dados de um produto de um elemento HTML
func (as *amazonScrapper) extractProduct(e *colly.HTMLElement, minPrice, maxPrice float64) *Product {
	// Extrair dados básicos
	name := as.extractText(e, []string{"h2 a span", "h2 span", "a[href*='/dp/'] span"})
	link := as.extractAttr(e, []string{"h2 a", "a[href*='/dp/']"}, "href")
	priceText := as.extractText(e, []string{"span.a-price-whole", "span.a-price span.a-offscreen", ".a-price .a-offscreen"})

	if name == "" || link == "" || priceText == "" {
		return nil
	}

	link = as.processLink(link)

	price, err := as.parsePrice(priceText)
	if err != nil || price < minPrice || price > maxPrice {
		return nil
	}

	rating := as.extractRating(e)
	soldCount := as.extractSoldCount(e)

	return &Product{
		Name:      strings.TrimSpace(name),
		Price:     price,
		Rating:    rating,
		SoldCount: soldCount,
		Link:      link,
		Site:      as.SiteName,
	}
}

func (as *amazonScrapper) extractText(e *colly.HTMLElement, selectors []string) string {
	for _, selector := range selectors {
		if text := e.ChildText(selector); text != "" {
			return text
		}
	}
	return ""
}

// extractAttr tenta extrair atributo usando múltiplos seletores
func (as *amazonScrapper) extractAttr(e *colly.HTMLElement, selectors []string, attr string) string {
	for _, selector := range selectors {
		if value := e.ChildAttr(selector, attr); value != "" {
			return value
		}
	}
	return ""
}

// processLink processa e limpa o link do produto
func (as *amazonScrapper) processLink(link string) string {
	if strings.HasPrefix(link, "/") {
		link = "https://www.amazon.com.br" + link
	}
	if idx := strings.Index(link, "?"); idx != -1 {
		link = link[:idx]
	}
	return link
}

// parsePrice converte texto de preço para float64
func (as *amazonScrapper) parsePrice(priceText string) (float64, error) {
	priceText = strings.ReplaceAll(priceText, "R$", "")
	priceText = strings.ReplaceAll(priceText, " ", "")
	priceText = strings.ReplaceAll(priceText, ".", "")
	priceText = strings.ReplaceAll(priceText, ",", ".")
	return strconv.ParseFloat(strings.TrimSpace(priceText), 64)
}

// extractRating extrai rating do produto
func (as *amazonScrapper) extractRating(e *colly.HTMLElement) float32 {
	ratingText := as.extractText(e, []string{
		"i.a-icon-star-small span",
		"span[aria-label*='estrelas']",
		"i.a-icon-star span",
	})

	if ratingText == "" {
		return 0
	}

	if r, err := strconv.ParseFloat(strings.Split(ratingText, " ")[0], 32); err == nil {
		return float32(r)
	}
	return 0
}

// extractSoldCount extrai quantidade vendida
func (as *amazonScrapper) extractSoldCount(e *colly.HTMLElement) uint {
	soldText := as.extractText(e, []string{
		"span[aria-label*='vendidos']",
		"span[aria-label*='comprados']",
	})

	if soldText == "" {
		return 0
	}

	if strings.Contains(soldText, "mil") {
		soldText = strings.ReplaceAll(soldText, "mil", "000")
	}

	if s, err := strconv.ParseUint(strings.Fields(soldText)[0], 10, 32); err == nil {
		return uint(s)
	}
	return 0
}

// scrapePage faz scraping de uma página específica
func (as *amazonScrapper) scrapePage(url string, minPrice, maxPrice float64) ([]Product, error) {
	var products []Product
	// ===== THREAD SAFETY =====
	// Mutex protege o slice 'products' de acesso concorrente
	// Necessário porque múltiplas goroutines podem acessar simultaneamente
	var mu sync.Mutex

	c := colly.NewCollector()
	c.UserAgent = AGENT

	// ===== RATE LIMITING =====
	// Configuração para evitar sobrecarregar o servidor da Amazon
	// DomainGlob: aplica limite apenas para domínios da Amazon
	// RandomDelay: adiciona delay aleatório entre requisições
	// Parallelism: 1 = apenas uma requisição por vez por collector
	c.Limit(&colly.LimitRule{
		DomainGlob:  "*amazon.com.br*",
		RandomDelay: as.RateLimit,
		Parallelism: 1,
	})

	// ===== CONFIGURAÇÃO DOS SELETORES =====
	// Múltiplos seletores para capturar diferentes formatos de produtos
	// Aumenta a chance de encontrar produtos mesmo se o HTML mudar
	selectors := []string{"div[data-component-type='s-search-result']", "div.s-result-item"}

	// ===== PROCESSAMENTO DOS PRODUTOS =====
	for _, selector := range selectors {
		c.OnHTML(selector, func(e *colly.HTMLElement) {
			if product := as.extractProduct(e, minPrice, maxPrice); product != nil {
				// ===== PROTEÇÃO CONCORRENTE =====
				// Lock garante que apenas uma goroutine modifique o slice por vez
				// Evita race conditions e corrupção de dados
				mu.Lock()
				products = append(products, *product)
				mu.Unlock()
			}
		})
	}

	// ===== EXECUÇÃO DO SCRAPING =====
	err := c.Visit(url)
	if err != nil {
		return nil, fmt.Errorf("erro ao visitar a página: %v", err)
	}

	return products, nil
}

func (as *amazonScrapper) Search(productName string, minPrice, maxPrice float64) ([]Product, error) {

	baseURL := as.BaseURL + strings.ReplaceAll(productName, " ", "+")
	urls := []string{
		baseURL,             // Página 1
		baseURL + "&page=2", // Página 2
		baseURL + "&page=3", // Página 3
	}

	results := make(chan []Product, len(urls))
	errors := make(chan error, len(urls))

	semaphore := make(chan struct{}, 2)

	var wg sync.WaitGroup

	for _, url := range urls {
		wg.Add(1)
		go func(pageURL string) {
			defer wg.Done()

			semaphore <- struct{}{}
			defer func() { <-semaphore }()

			products, err := as.scrapePage(pageURL, minPrice, maxPrice)
			if err != nil {

				errors <- err
				return
			}
			// Se sucesso, envia produtos para canal de resultados
			results <- products
		}(url) // Passa URL como parâmetro para evitar closure issues
	}

	go func() {
		wg.Wait()      // Bloqueia até todas as goroutines terminarem
		close(results) // Fecha canal de resultados
		close(errors)  // Fecha canal de erros
	}()

	var allProducts []Product
	for products := range results { // Range em canal bloqueia até fechar
		allProducts = append(allProducts, products...)
	}

	select {
	case err := <-errors:
		// Se houve erro, retorna produtos coletados + erro
		return allProducts, err
	default:
		// Se não houve erro, retorna apenas produtos
		return allProducts, nil
	}
}

func (as *amazonScrapper) GetSiteName() string {
	return "Amazon"
}
