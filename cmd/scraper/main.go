package main

import (
	"bufio"
	"flag"
	"log/slog"
	"os"
	"strings"

	"github.com/chicopsych/webscrapingFLS/internal/cli"
	"github.com/chicopsych/webscrapingFLS/internal/gui"
	"github.com/chicopsych/webscrapingFLS/internal/logger"
)

// main é o composition root da aplicação.
//
// Responsabilidades intencionais (SoC):
//   - processar flags de inicialização
//   - configurar logger
//   - decidir qual interface será executada (CLI ou GUI)
//
// O ponto de entrada NÃO contém regra de negócio de scraping, naming de arquivo,
// HTTP handlers ou serialização. Esses detalhes vivem em pacotes internos
// especializados para reduzir acoplamento e facilitar manutenção.
func main() {
	urlFlag := flag.String("url", "", "URL alvo para realizar o scraping")
	urlsFlag := flag.String("urls", "", "URLs separadas por vírgula para scraping em batch")
	urlsFileFlag := flag.String("urls-file", "", "Arquivo com uma URL por linha para scraping em batch")
	debugFlag := flag.Bool("debug", false, "Habilita log em nível debug")
	outDirFlag := flag.String("out", "data", "Diretório de saída para os arquivos .md")
	guiFlag := flag.Bool("gui", false, "Abre a interface gráfica")
	flag.Parse()

	log := logger.InitLogger(*debugFlag)
	noArgs := len(os.Args) == 1

	urls := collectURLs(*urlFlag, *urlsFlag, *urlsFileFlag, log)

	if len(urls) > 0 {
		if err := cli.Run(urls, *outDirFlag, log); err != nil {
			log.Error("Falha no modo CLI", "error", err)
			os.Exit(1)
		}
		return
	}

	if *guiFlag || noArgs {
		gui.Run(*outDirFlag, log)
		return
	}

	log.Error("Parâmetros inválidos. Use -url/-urls/-urls-file para CLI ou -gui/sem argumentos para GUI")
	os.Exit(1)
}

// collectURLs agrega URLs das três fontes de entrada: -url, -urls e -urls-file.
// Linhas em branco e linhas começando com '#' são ignoradas no arquivo.
func collectURLs(single, multi, filePath string, log *slog.Logger) []string {
	var urls []string

	if single != "" {
		urls = append(urls, strings.TrimSpace(single))
	}

	if multi != "" {
		for _, u := range strings.Split(multi, ",") {
			u = strings.TrimSpace(u)
			if u != "" {
				urls = append(urls, u)
			}
		}
	}

	if filePath != "" {
		f, err := os.Open(filePath)
		if err != nil {
			log.Error("Falha ao abrir arquivo de URLs", "path", filePath, "error", err)
			os.Exit(1)
		}
		defer f.Close()

		scanner := bufio.NewScanner(f)
		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())
			if line == "" || strings.HasPrefix(line, "#") {
				continue
			}
			urls = append(urls, line)
		}
		if err := scanner.Err(); err != nil {
			log.Error("Erro ao ler arquivo de URLs", "path", filePath, "error", err)
			os.Exit(1)
		}
	}

	return urls
}
