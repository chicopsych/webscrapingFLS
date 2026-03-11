package cli

import (
	"fmt"
	"log/slog"

	"github.com/chicopsych/webscrapingFLS/internal/service/scrape"
)

// Run executa o modo CLI usando o mesmo caso de uso compartilhado com a GUI.
// Quando apenas uma URL é fornecida, executa scraping individual.
// Com múltiplas URLs, executa em batch concorrente (até 3 simultâneas).
func Run(urls []string, outDir string, log *slog.Logger) error {
	if len(urls) == 0 {
		return fmt.Errorf("nenhuma URL fornecida")
	}

	service := scrape.New(log)

	if len(urls) == 1 {
		log.Info("Iniciando Web Scraper (CLI)", "url", urls[0])
		filePath, err := service.ExecuteScrape(urls[0], outDir)
		if err != nil {
			return err
		}
		log.Info("Extração concluída com sucesso e arquivo salvo", "path", filePath)
		return nil
	}

	log.Info("Iniciando Web Scraper (CLI batch)", "total_urls", len(urls))
	results := service.ExecuteBatch(urls, outDir, 3)

	var failed int
	for _, r := range results {
		if r.Err != nil {
			log.Error("Falha na extração", "url", r.URL, "error", r.Err)
			failed++
		} else {
			log.Info("Extração concluída e arquivo salvo", "url", r.URL, "path", r.FilePath)
		}
	}

	if failed > 0 {
		return fmt.Errorf("%d de %d URLs falharam durante o batch", failed, len(urls))
	}
	return nil
}
