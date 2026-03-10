package main

import (
	"flag"
	"os"

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
	debugFlag := flag.Bool("debug", false, "Habilita log em nível debug")
	outDirFlag := flag.String("out", "data", "Diretório de saída para os arquivos .md")
	guiFlag := flag.Bool("gui", false, "Abre a interface gráfica")
	flag.Parse()

	log := logger.InitLogger(*debugFlag)
	noArgs := len(os.Args) == 1

	if *urlFlag != "" {
		if err := cli.Run(*urlFlag, *outDirFlag, log); err != nil {
			log.Error("Falha no modo CLI", "error", err)
			os.Exit(1)
		}
		return
	}

	if *guiFlag || noArgs {
		gui.Run(*outDirFlag, log)
		return
	}

	log.Error("Parâmetros inválidos. Use -url para CLI ou -gui/sem argumentos para GUI")
	os.Exit(1)
}
