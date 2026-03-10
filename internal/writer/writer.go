package writer

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/chicopsych/webscrapingFLS/internal/errors"
	"github.com/chicopsych/webscrapingFLS/internal/models"

	"gopkg.in/yaml.v3"
)

// SaveMarkdown salva o objeto PageData em um arquivo Markdown (.md) com Front Matter YAML
func SaveMarkdown(data models.PageData, outputDir string, filename string) error {
	// Garante que o diretório exista
	if err := os.MkdirAll(outputDir, os.ModePerm); err != nil {
		return fmt.Errorf("%w: %v", errors.ErrFileWriteAccess, err)
	}

	// Serializa o metadado (Front Matter)
	frontMatter, err := yaml.Marshal(&data)
	if err != nil {
		return fmt.Errorf("%w: falha ao gerar front-matter YAML: %v", errors.ErrInvalidMetadata, err)
	}

	body := data.Content
	if body == "" {
		body = data.MarkdownBody
	}

	// Monta o conteúdo final do arquivo
	content := fmt.Sprintf("---\n%s---\n\n%s\n", string(frontMatter), body)

	// Caminho completo
	filePath := filepath.Join(outputDir, filename)

	// Salva no disco
	if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
		return fmt.Errorf("%w: %v", errors.ErrFileWriteAccess, err)
	}

	return nil
}
