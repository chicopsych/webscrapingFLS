package gui

import (
	"embed"
	"html/template"
	"log/slog"
	"net"
	"net/http"
	neturl "net/url"
	"os/exec"
	"runtime"
	"strings"

	"github.com/chicopsych/webscrapingFLS/internal/service/scrape"
)

//go:embed templates/index.html
var templateFS embed.FS

// scrapeResultData carrega o resultado de uma URL individual no batch da GUI.
type scrapeResultData struct {
	URL      string
	FilePath string
	Error    string
	OK       bool
}

type pageData struct {
	Message string
	Error   string
	OutDir  string
	Results []scrapeResultData
}

// Run inicia a interface web local (GUI) e delega o caso de uso ao serviço.
//
// Este pacote representa camada de apresentação. Ele lida com HTTP, templates
// e experiência do usuário; a regra de negócio permanece no pacote service/scrape.
func Run(outDir string, log *slog.Logger) {
	log.Info("Modo GUI selecionado", "out_dir", outDir)

	tmpl, err := template.ParseFS(templateFS, "templates/index.html")
	if err != nil {
		log.Error("Falha ao carregar template da GUI", "error", err)
		return
	}

	service := scrape.New(log)
	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		data := pageData{
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

		// Coletar e validar URLs do textarea (uma por linha).
		raw := strings.TrimSpace(r.FormValue("urls"))
		if raw == "" {
			http.Redirect(w, r, "/?err="+neturl.QueryEscape("Informe ao menos uma URL válida."), http.StatusSeeOther)
			return
		}

		var validURLs []string
		for _, line := range strings.Split(raw, "\n") {
			u := strings.TrimSpace(line)
			if u == "" || strings.HasPrefix(u, "#") {
				continue
			}
			if _, err := neturl.ParseRequestURI(u); err != nil {
				http.Redirect(w, r, "/?err="+neturl.QueryEscape("URL inválida: "+u), http.StatusSeeOther)
				return
			}
			validURLs = append(validURLs, u)
		}

		if len(validURLs) == 0 {
			http.Redirect(w, r, "/?err="+neturl.QueryEscape("Nenhuma URL válida encontrada."), http.StatusSeeOther)
			return
		}

		batchResults := service.ExecuteBatch(validURLs, outDir, 3)

		results := make([]scrapeResultData, len(batchResults))
		for i, br := range batchResults {
			rd := scrapeResultData{URL: br.URL, OK: br.Err == nil}
			if br.Err != nil {
				rd.Error = br.Err.Error()
			} else {
				rd.FilePath = br.FilePath
				log.Info("Arquivo salvo via GUI", "path", br.FilePath)
			}
			results[i] = rd
		}

		// Renderizar diretamente (sem redirect) para exibir resultados estruturados.
		data := pageData{OutDir: outDir, Results: results}
		if err := tmpl.Execute(w, data); err != nil {
			log.Error("Falha ao renderizar template da GUI", "error", err)
			http.Error(w, "erro interno ao renderizar página", http.StatusInternalServerError)
		}
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
