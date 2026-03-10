package crawler

import (
	"log/slog"
	"strings"
	"time"

	"github.com/chicopsych/webscrapingFLS/internal/converter"
	"github.com/chicopsych/webscrapingFLS/internal/errors"
	"github.com/chicopsych/webscrapingFLS/internal/models"

	"github.com/gocolly/colly/v2"
)

// Scrape executa a extração em uma dada URL e retorna a estrutura populada.
func Scrape(url string, logger *slog.Logger, events chan<- Event) (models.PageData, error) {
	logger.Info("Iniciando requisição", "url", url)
	start := time.Now()
	emitEvent(events, Event{
		Type:      EventStarted,
		Message:   "Iniciando extração",
		URL:       url,
		Timestamp: time.Now(),
		Progress:  0.05,
	})

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

		var classifiedErr error
		if r.StatusCode == 429 {
			classifiedErr = errors.ErrRateLimited
		} else if r.StatusCode == 403 || r.StatusCode == 401 {
			classifiedErr = errors.ErrAccessDenied
		} else {
			classifiedErr = errors.ErrNetworkUnreachable
		}

		scrapeErr = classifiedErr
		emitEvent(events, Event{
			Type:      EventError,
			Message:   "Erro durante a requisição",
			URL:       r.Request.URL.String(),
			Timestamp: time.Now(),
			Progress:  1,
			Err:       classifiedErr,
		})
	})

	// Parsing do HTML
	c.OnHTML("html", func(e *colly.HTMLElement) {
		title := e.ChildText("title")
		if title == "" {
			scrapeErr = errors.ErrSelectorNotFound
			logger.Warn("Seletor de título não encontrado", "url", e.Request.URL.String())
			emitEvent(events, Event{
				Type:      EventError,
				Message:   "Seletor de título não encontrado",
				URL:       e.Request.URL.String(),
				Timestamp: time.Now(),
				Progress:  1,
				Err:       errors.ErrSelectorNotFound,
			})
		} else {
			pageData.Title = strings.TrimSpace(title)
			logger.Info("Título extraído com sucesso", "title", pageData.Title)
			emitEvent(events, Event{
				Type:      EventTitleExtracted,
				Message:   "Título extraído com sucesso",
				URL:       e.Request.URL.String(),
				Timestamp: time.Now(),
				Progress:  0.6,
			})
		}

		body := e.ChildText("body")
		if body != "" {
			// Extrai o HTML interno do <body> para conversão estruturada.
			// e.ChildText() retorna apenas texto plano (sem tags); e.DOM.Find().Html()
			// retorna o HTML bruto preservando headings, listas, links e código.
			htmlContent, htmlErr := e.DOM.Find("body").Html()
			if htmlErr == nil && strings.TrimSpace(htmlContent) != "" {
				pageData.RawContent = strings.TrimSpace(body) // texto plano para debug
				markdown := converter.HTMLToMarkdown(htmlContent, url)
				pageData.Content = markdown
				pageData.MarkdownBody = markdown // compatibilidade
			} else {
				// Fallback para texto plano se o DOM não estiver disponível.
				pageData.RawContent = strings.TrimSpace(body)
				pageData.Content = strings.TrimSpace(body)
				pageData.MarkdownBody = pageData.Content
			}
		}
	})

	c.OnRequest(func(r *colly.Request) {
		logger.Debug("Enviando requisição", "url", r.URL.String())
		emitEvent(events, Event{
			Type:      EventRequestSent,
			Message:   "Requisição enviada",
			URL:       r.URL.String(),
			Timestamp: time.Now(),
			Progress:  0.2,
		})
	})

	c.OnScraped(func(r *colly.Response) {
		duration := time.Since(start)
		logger.Info("Requisição finalizada", "url", r.Request.URL.String(), "duration_ms", duration.Milliseconds())
		emitEvent(events, Event{
			Type:      EventCompleted,
			Message:   "Extração finalizada",
			URL:       r.Request.URL.String(),
			Timestamp: time.Now(),
			Progress:  1,
		})
	})

	c.Visit(url)
	c.Wait()

	if scrapeErr != nil {
		return pageData, scrapeErr
	}

	return pageData, nil
}
