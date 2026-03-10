package filesystem

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// UniqueFileName devolve um nome de arquivo não utilizado dentro de outputDir.
//
// Estratégia:
//   - Se o nome base já existir, aplica sufixo incremental: _1, _2, _3...
//   - Retorna apenas o nome do arquivo (sem diretório), para manter fronteira
//     clara entre camada de naming e camada de persistência.
//
// # Segurança
//
// Mesmo recebendo um nome já sanitizado, esta função aplica filepath.Base de novo
// para eliminar qualquer segmento de diretório acidental e reduzir risco de
// path traversal por uso incorreto da API.
func UniqueFileName(outputDir, fileName string) (string, error) {
	candidate := filepath.Base(strings.TrimSpace(fileName))
	candidate = strings.Trim(candidate, " .")
	if candidate == "" || candidate == "." || candidate == ".." {
		return "", fmt.Errorf("nome de arquivo inválido")
	}

	ext := filepath.Ext(candidate)
	if ext == "" {
		ext = ".md"
		candidate += ext
	}

	base := strings.TrimSuffix(candidate, ext)

	for i := 0; ; i++ {
		tryName := candidate
		if i > 0 {
			tryName = fmt.Sprintf("%s_%d%s", base, i, ext)
		}

		fullPath := filepath.Join(outputDir, tryName)
		_, err := os.Stat(fullPath)
		if os.IsNotExist(err) {
			return tryName, nil
		}
		if err != nil {
			return "", fmt.Errorf("falha ao verificar existência de %q: %w", fullPath, err)
		}
	}
}
