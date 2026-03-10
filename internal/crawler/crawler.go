package crawler

import (
	"log/slog"
	"strings"
	"time"

	"github.com/chicopsych/webscrapingFLS/internal/errors"
	"github.com/chicopsych/webscrapingFLS/internal/models"

	"github.com/gocolly/colly/v2"
)

// Scrape executa a extração em uma dada URL e retorna a estrutura populada.
func Scrape(url string, logger *slog.Logger) (models.PageData, error) {
	logger.Info("Iniciando requisição", "url", url)
	start := time.Now()

	c := colly.NewCollector(
		colly.Async(true),
	)

	c.Limit(&colly.LimitRule{
		DomainGlob:  "*",
		Parallelism: 2,
		Delay:       1 * time.Second,
	})

	var pageData models.PageData
	pageData.URL = url
	pageData.Timestamp = time.Now()

	var scrapeErr error

	// Tratamento de Erros de Rede/HTTP
	c.OnError(func(r *colly.Response, err error) {
		logger.Error("Erro na requisição", "url", r.Request.URL.String(), "status", r.StatusCode, "error", err)
		if r.StatusCode == 429 {
			scrapeErr = errors.ErrRateLimited
		} else if r.StatusCode == 403 || r.StatusCode == 401 {
			scrapeErr = errors.ErrAccessDenied
		} else {
			scrapeErr = errors.ErrNetworkUnreachable
		}
	})

	// Parsing do HTML
	c.OnHTML("html", func(e *colly.HTMLElement) {
		title := e.ChildText("title")
		if title == "" {
			scrapeErr = errors.ErrSelectorNotFound
			logger.Warn("Seletor de título não encontrado", "url", e.Request.URL.String())
		} else {
			pageData.Title = strings.TrimSpace(title)
			logger.Info("Título extraído com sucesso", "title", pageData.Title)
		}

		body := e.ChildText("body")
		if body != "" {
			pageData.RawContent = e.Text
			pageData.MarkdownBody = strings.TrimSpace(body) // Simulação de html para md
		}
	})

	c.OnRequest(func(r *colly.Request) {
		logger.Debug("Enviando requisição", "url", r.URL.String())
	})

	c.OnScraped(func(r *colly.Response) {
		duration := time.Since(start)
		logger.Info("Requisição finalizada", "url", r.Request.URL.String(), "duration_ms", duration.Milliseconds())
	})

	c.Visit(url)
	c.Wait()

	if scrapeErr != nil {
		return pageData, scrapeErr
	}

	return pageData, nil
}
