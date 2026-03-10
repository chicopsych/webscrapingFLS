package main

import (
	"flag"
	"fmt"
	"html/template"
	"log/slog"
	"net"
	"net/http"
	neturl "net/url"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/chicopsych/webscrapingFLS/internal/crawler"
	"github.com/chicopsych/webscrapingFLS/internal/logger"
	"github.com/chicopsych/webscrapingFLS/internal/writer"
)

const guiTemplate = `<!doctype html>
<html lang="pt-BR">
<head>
	<meta charset="utf-8" />
	<meta name="viewport" content="width=device-width, initial-scale=1" />
	<title>webscrapingFLS</title>
	<style>
		body { font-family: Segoe UI, Arial, sans-serif; max-width: 860px; margin: 40px auto; padding: 0 16px; }
		h1 { margin-bottom: 6px; }
		.muted { color: #555; margin-top: 0; }
		.box { border: 1px solid #ddd; border-radius: 8px; padding: 16px; background: #fafafa; }
		label { display: block; margin-bottom: 8px; font-weight: 600; }
		input[type="url"] { width: 100%; padding: 10px; font-size: 14px; border-radius: 6px; border: 1px solid #bbb; box-sizing: border-box; }
		button { margin-top: 12px; background: #1463ff; color: white; border: 0; border-radius: 6px; padding: 10px 14px; cursor: pointer; }
		.ok { margin-top: 14px; color: #0c6b2f; font-weight: 600; }
		.err { margin-top: 14px; color: #9a1c1c; font-weight: 600; }
		code { background: #f2f2f2; padding: 2px 4px; border-radius: 4px; }
	</style>
</head>
<body>
	<h1>webscrapingFLS</h1>
	<p class="muted">Interface GUI local para extração de páginas.</p>
	<div class="box">
		<form method="post" action="/scrape">
			<label for="url">URL alvo</label>
			<input id="url" name="url" type="url" placeholder="https://example.com" required />
			<button type="submit">Extrair e salvar</button>
		</form>
		{{if .Message}}<div class="ok">{{.Message}}</div>{{end}}
		{{if .Error}}<div class="err">{{.Error}}</div>{{end}}
		<p class="muted">Saída padrão: <code>{{.OutDir}}/scraped_page.md</code></p>
	</div>
</body>
</html>`

type guiPageData struct {
	Message string
	Error   string
	OutDir  string
}

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
	_, err := scrapeAndSave(url, outDir, log)
	if err != nil {
		return err
	}

	log.Info("Extração concluída com sucesso e arquivo salvo")
	return nil
}

func runGUI(log *slog.Logger, outDir string) {
	log.Info("Modo GUI selecionado", "out_dir", outDir)
	tmpl := template.Must(template.New("gui").Parse(guiTemplate))
	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		data := guiPageData{
			Message: r.URL.Query().Get("msg"),
			Error:   r.URL.Query().Get("err"),
			OutDir:  outDir,
		}
		if err := tmpl.Execute(w, data); err != nil {
			log.Error("Falha ao renderizar template da GUI", "error", err)
			http.Error(w, "erro interno ao renderizar página", http.StatusInternalServerError)
		}
	})

	mux.HandleFunc("/scrape", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		targetURL := strings.TrimSpace(r.FormValue("url"))
		if targetURL == "" {
			http.Redirect(w, r, "/?err="+neturl.QueryEscape("Informe uma URL válida."), http.StatusSeeOther)
			return
		}

		if _, err := neturl.ParseRequestURI(targetURL); err != nil {
			http.Redirect(w, r, "/?err="+neturl.QueryEscape("URL inválida: "+err.Error()), http.StatusSeeOther)
			return
		}

		filePath, err := scrapeAndSave(targetURL, outDir, log)
		if err != nil {
			http.Redirect(w, r, "/?err="+neturl.QueryEscape(err.Error()), http.StatusSeeOther)
			return
		}

		log.Info("Arquivo salvo via GUI", "path", filePath)
		http.Redirect(w, r, "/?msg="+neturl.QueryEscape("Extração concluída. Arquivo salvo em: "+filePath), http.StatusSeeOther)
	})

	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		log.Error("Falha ao iniciar listener da GUI", "error", err)
		return
	}

	guiURL := "http://" + ln.Addr().String()
	if err = openBrowser(guiURL); err != nil {
		log.Warn("Não foi possível abrir navegador automaticamente", "url", guiURL, "error", err)
	}

	log.Info("GUI web iniciada", "url", guiURL)
	if err = http.Serve(ln, mux); err != nil {
		log.Error("Servidor GUI finalizado com erro", "error", err)
	}
}

func scrapeAndSave(targetURL, outDir string, log *slog.Logger) (string, error) {
	pageData, err := crawler.Scrape(targetURL, log, nil)
	if err != nil {
		return "", fmt.Errorf("falha ao extrair dados da página: %w", err)
	}

	fileName := "scraped_page.md"
	if err = writer.SaveMarkdown(pageData, outDir, fileName); err != nil {
		return "", fmt.Errorf("falha ao salvar arquivo: %w", err)
	}

	return filepath.Join(outDir, fileName), nil
}

func openBrowser(targetURL string) error {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("rundll32", "url.dll,FileProtocolHandler", targetURL)
	case "darwin":
		cmd = exec.Command("open", targetURL)
	default:
		cmd = exec.Command("xdg-open", targetURL)
	}

	return cmd.Start()
}
