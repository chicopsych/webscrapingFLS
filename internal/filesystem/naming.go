package filesystem

import (
	"path/filepath"
	"strings"
)

const (
	defaultBaseName   = "scraped_page"
	maxBaseNameRunes  = 120
)

// SanitizeFileNameFromTitle converte o título da página em um nome de arquivo
// seguro e portável entre NTFS (Windows) e EXT4 (Linux).
//
// # Segurança e Integridade
//
// A entrada (title) pode conter caracteres maliciosos ou incompatíveis com
// o sistema de arquivos. Esta função aplica defesa em profundidade:
//  1. Remove separadores de diretório para bloquear path traversal.
//  2. Remove caracteres inválidos no NTFS (< > : " / \ | ? *).
//  3. Normaliza espaços para underscore e elimina sufixos inválidos (ponto/espaço).
//  4. Protege nomes reservados do Windows (CON, PRN, COM1...).
//  5. Força basename com filepath.Base, impedindo que o resultado contenha caminho.
//
// O retorno sempre termina com extensão .md.
func SanitizeFileNameFromTitle(title string) string {
	base := strings.TrimSpace(title)
	if base == "" {
		return defaultBaseName + ".md"
	}

	replacer := strings.NewReplacer(
		"<", "_",
		">", "_",
		":", "_",
		"\"", "_",
		"/", "_",
		"\\", "_",
		"|", "_",
		"?", "_",
		"*", "_",
	)

	base = replacer.Replace(base)
	base = strings.ReplaceAll(base, "..", "_")
	base = strings.Join(strings.Fields(base), "_")
	base = strings.Trim(base, " .")
	base = filepath.Base(base)

	if base == "" || base == "." || base == ".." {
		base = defaultBaseName
	}

	base = trimRunes(base, maxBaseNameRunes)
	base = strings.Trim(base, " .")
	if base == "" {
		base = defaultBaseName
	}

	if isWindowsReservedName(base) {
		base = "page_" + base
	}

	return base + ".md"
}

func isWindowsReservedName(base string) bool {
	upper := strings.ToUpper(strings.TrimSpace(base))
	if upper == "" {
		return false
	}

	if _, reserved := windowsReservedNames[upper]; reserved {
		return true
	}

	nameOnly := upper
	if idx := strings.IndexRune(upper, '.'); idx > 0 {
		nameOnly = upper[:idx]
	}

	_, reserved := windowsReservedNames[nameOnly]
	return reserved
}

func trimRunes(s string, max int) string {
	r := []rune(s)
	if len(r) <= max {
		return s
	}
	return string(r[:max])
}
