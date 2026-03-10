package cli

import (
	"log/slog"

	"github.com/chicopsych/webscrapingFLS/internal/service/scrape"
)

// Run executa o modo CLI usando o mesmo caso de uso compartilhado com a GUI.
//
// Go Idiom: a interface (CLI) fica fina e delega regra de negócio ao serviço.
// Isso reduz duplicação e facilita testes de integração por camada.
func Run(targetURL, outDir string, log *slog.Logger) error {
	log.Info("Iniciando Web Scraper (CLI)", "url", targetURL)

	service := scrape.New(log)
	filePath, err := service.ExecuteScrape(targetURL, outDir)
	if err != nil {
		return err
	}

	log.Info("Extração concluída com sucesso e arquivo salvo", "path", filePath)
	return nil
}
