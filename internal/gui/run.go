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

type pageData struct {
	Message string
	Error   string
	OutDir  string
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

		targetURL := strings.TrimSpace(r.FormValue("url"))
		if targetURL == "" {
			http.Redirect(w, r, "/?err="+neturl.QueryEscape("Informe uma URL válida."), http.StatusSeeOther)
			return
		}

		if _, err := neturl.ParseRequestURI(targetURL); err != nil {
			http.Redirect(w, r, "/?err="+neturl.QueryEscape("URL inválida: "+err.Error()), http.StatusSeeOther)
			return
		}

		filePath, err := service.ExecuteScrape(targetURL, outDir)
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
