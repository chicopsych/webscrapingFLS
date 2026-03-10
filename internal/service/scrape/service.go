package scrape

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/chicopsych/webscrapingFLS/internal/crawler"
	"github.com/chicopsych/webscrapingFLS/internal/filesystem"
	"github.com/chicopsych/webscrapingFLS/internal/writer"
)

// Service implementa o caso de uso ExecuteScrape.
//
// Na Clean Architecture, este componente representa a camada de aplicação
// (use case / interactor): ele coordena dependências de infraestrutura
// (crawler + writer + filesystem naming) sem expor detalhes à camada de interface
// (CLI/GUI). Assim, qualquer entrada (terminal, web, API) reutiliza a mesma regra.
type Service struct {
	logger *slog.Logger
}

// New cria um Service com logger injetado.
func New(logger *slog.Logger) *Service {
	return &Service{logger: logger}
}

// ExecuteScrape executa o fluxo completo de negócio:
//  1. Extrai dados da URL via crawler.
//  2. Constrói nome de arquivo seguro a partir do título.
//  3. Garante nome único para evitar sobrescrita.
//  4. Persiste Markdown + Front Matter via writer.
func (s *Service) ExecuteScrape(targetURL, outDir string) (string, error) {
	pageData, err := crawler.Scrape(targetURL, s.logger, nil)
	if err != nil {
		return "", fmt.Errorf("falha ao extrair dados da página: %w", err)
	}

	safeName := filesystem.SanitizeFileNameFromTitle(pageData.Title)

	// Mitigação prática de TOCTOU: se outra rotina criar o arquivo no meio do
	// processo, o writer falha com EEXIST (O_EXCL) e tentamos novo sufixo.
	const maxAttempts = 5
	for attempt := 1; attempt <= maxAttempts; attempt++ {
		uniqueName, uniqueErr := filesystem.UniqueFileName(outDir, safeName)
		if uniqueErr != nil {
			return "", fmt.Errorf("falha ao gerar nome de arquivo único: %w", uniqueErr)
		}

		err = writer.SaveMarkdown(pageData, outDir, uniqueName)
		if err == nil {
			return filepath.Join(outDir, uniqueName), nil
		}

		if os.IsExist(err) {
			s.logger.Warn("Colisão de nome detectada, tentando novo sufixo", "attempt", attempt, "file", uniqueName)
			continue
		}

		return "", fmt.Errorf("falha ao salvar arquivo: %w", err)
	}

	return "", fmt.Errorf("não foi possível reservar nome único após %d tentativas", maxAttempts)
}
