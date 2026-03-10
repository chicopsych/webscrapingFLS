package main

import (
	"flag"
	"os"

	"github.com/chicopsych/webscrapingFLS/internal/crawler"
	"github.com/chicopsych/webscrapingFLS/internal/logger"
	"github.com/chicopsych/webscrapingFLS/internal/writer"
)

func main() {
	urlFlag := flag.String("url", "", "URL alvo para realizar o scraping")
	debugFlag := flag.Bool("debug", false, "Habilita log em nível debug")
	outDirFlag := flag.String("out", "data", "Diretório de saída para os arquivos .md")
	flag.Parse()

	log := logger.InitLogger(*debugFlag)

	if *urlFlag == "" {
		log.Error("A flag -url é obrigatória")
		os.Exit(1)
	}

	log.Info("Iniciando Web Scraper", "url", *urlFlag)

	// Iniciar extração
	pageData, err := crawler.Scrape(*urlFlag, log)
	if err != nil {
		log.Error("Falha ao extrair dados da página", "error", err)
		os.Exit(1)
	}

	// Salvar resultado
	err = writer.SaveMarkdown(pageData, *outDirFlag, "scraped_page.md")
	if err != nil {
		log.Error("Falha ao salvar arquivo", "error", err)
		os.Exit(1)
	}

	log.Info("Extração concluída com sucesso e arquivo salvo.")
}
