package scrape

import (
	"fmt"
	"log/slog"
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
	uniqueName, err := filesystem.UniqueFileName(outDir, safeName)
	if err != nil {
		return "", fmt.Errorf("falha ao gerar nome de arquivo único: %w", err)
	}

	if err = writer.SaveMarkdown(pageData, outDir, uniqueName); err != nil {
		return "", fmt.Errorf("falha ao salvar arquivo: %w", err)
	}

	return filepath.Join(outDir, uniqueName), nil
}
