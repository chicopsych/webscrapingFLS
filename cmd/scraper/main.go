package main

import (
	"flag"
	"fmt"
	"log/slog"
	"os"

	"github.com/chicopsych/webscrapingFLS/internal/crawler"
	"github.com/chicopsych/webscrapingFLS/internal/logger"
	"github.com/chicopsych/webscrapingFLS/internal/writer"
)

func main() {
	urlFlag := flag.String("url", "", "URL alvo para realizar o scraping")
	debugFlag := flag.Bool("debug", false, "Habilita log em nível debug")
	outDirFlag := flag.String("out", "data", "Diretório de saída para os arquivos .md")
	guiFlag := flag.Bool("gui", false, "Abre a interface gráfica")
	flag.Parse()

	log := logger.InitLogger(*debugFlag)
	noArgs := len(os.Args) == 1

	if *urlFlag != "" {
		if err := runCLI(*urlFlag, *outDirFlag, log); err != nil {
			log.Error("Falha no modo CLI", "error", err)
			os.Exit(1)
		}
		return
	}

	if *guiFlag || noArgs {
		runGUI(log, *outDirFlag)
		return
	}

	log.Error("Parâmetros inválidos. Use -url para CLI ou -gui/sem argumentos para GUI")
	os.Exit(1)
}

func runCLI(url, outDir string, log *slog.Logger) error {
	log.Info("Iniciando Web Scraper (CLI)", "url", url)
	pageData, err := crawler.Scrape(url, log, nil)
	if err != nil {
		return fmt.Errorf("falha ao extrair dados da página: %w", err)
	}

	if err = writer.SaveMarkdown(pageData, outDir, "scraped_page.md"); err != nil {
		return fmt.Errorf("falha ao salvar arquivo: %w", err)
	}

	log.Info("Extração concluída com sucesso e arquivo salvo")
	return nil
}

func runGUI(log *slog.Logger, outDir string) {
	log.Info("Modo GUI selecionado", "out_dir", outDir)

	// Etapa 1: apenas esqueleto de fluxo e roteamento de eventos.
	// O layout Fyne (campo URL, botão, spinner, área de log) será implementado após feedback.
	events := make(chan crawler.Event, 32)

	go consumeEvents(log, events)

	// Exemplo de execução assíncrona para evitar bloqueio da UI:
	// go startExtractionAsync(urlFromInput, outDir, log, events)

	log.Info("Esqueleto GUI pronto para receber layout e bindings")
}

func startExtractionAsync(url, outDir string, log *slog.Logger, events chan<- crawler.Event) {
	go func() {
		pageData, err := crawler.Scrape(url, log, events)
		if err != nil {
			crawlerEvent := crawler.Event{
				Type:    crawler.EventError,
				Message: "Falha no scraping",
				URL:     url,
				Err:     err,
			}
			select {
			case events <- crawlerEvent:
			default:
			}
			return
		}

		if err = writer.SaveMarkdown(pageData, outDir, "scraped_page.md"); err != nil {
			crawlerEvent := crawler.Event{
				Type:    crawler.EventError,
				Message: "Falha ao salvar arquivo",
				URL:     url,
				Err:     err,
			}
			select {
			case events <- crawlerEvent:
			default:
			}
		}
	}()
}

func consumeEvents(log *slog.Logger, events <-chan crawler.Event) {
	for event := range events {
		if event.Err != nil {
			log.Error(event.Message, "url", event.URL, "error", event.Err, "progress", event.Progress)
			continue
		}

		log.Info(event.Message, "url", event.URL, "type", event.Type, "progress", event.Progress)
	}
}
