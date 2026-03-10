// Package sanitizer extrai o conteúdo principal de uma página HTML,
// removendo elementos de ruído como menus, banners, rodapés e imagens.
// Suporta configuração por domínio com fallback para um perfil genérico.
package sanitizer

import (
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// SanitizeConfig define os seletores CSS usados durante a sanitização.
type SanitizeConfig struct {
	// ContentSelectors são tentados em ordem; o primeiro que encontrar um elemento
	// com conteúdo de texto não-vazio é utilizado como região do artigo.
	ContentSelectors []string
	// RemoveSelectors são removidos do DOM antes da conversão para Markdown
	// (usado quando nenhum ContentSelector encontrar o artigo).
	RemoveSelectors []string
}

// DefaultConfig retorna a configuração genérica aplicada a qualquer domínio
// sem entrada específica no registry.
func DefaultConfig() SanitizeConfig {
	return SanitizeConfig{
		ContentSelectors: []string{
			"article",
			"main",
			`[role="main"]`,
			".article-content",
			".post-content",
			".entry-content",
			".content-body",
			"#content",
			"#main-content",
			".content",
		},
		RemoveSelectors: []string{
			"nav",
			"header",
			"footer",
			"aside",
			`[class*="sidebar"]`,
			`[class*="banner"]`,
			`[class*="advertisement"]`,
			`[id*="banner"]`,
			`[id*="ad-"]`,
			`[class*="social"]`,
			`[class*="share"]`,
			`[class*="related"]`,
			`[class*="recommended"]`,
			`[class*="newsletter"]`,
			`[class*="cookie"]`,
			`[class*="popup"]`,
		},
	}
}

// domainRegistry mapeia hostname → configuração específica.
// Adicione novos domínios aqui para refinamento de extração.
var domainRegistry = map[string]SanitizeConfig{
	// Joomla CMS — Instituto Newton C. Braga
	"www.newtoncbraga.com.br": {
		ContentSelectors: []string{
			".itemFullText",
			".item-page",
			"#content",
			".article-content",
			"article",
		},
		RemoveSelectors: []string{
			".moduletable",
			".mostReadBox",
			".tagBox",
			".bannerGroup",
			"#bottom",
			".readMore",
			`[class*="sidebar"]`,
			`[class*="social"]`,
			`[class*="share"]`,
			"nav",
			"header",
			"footer",
			"aside",
		},
	},
}

// ConfigForURL retorna a SanitizeConfig adequada para a URL fornecida.
// Retorna DefaultConfig() quando o domínio não está no registry.
func ConfigForURL(rawURL string) SanitizeConfig {
	u, err := url.Parse(rawURL)
	if err != nil {
		return DefaultConfig()
	}
	host := strings.ToLower(u.Hostname())
	if cfg, ok := domainRegistry[host]; ok {
		return cfg
	}
	return DefaultConfig()
}

// ExtractMainContent sanitiza o HTML do body e retorna HTML limpo pronto
// para conversão em Markdown. A lógica é híbrida:
//
//  1. Remove elementos de mídia (imagens, iframes, scripts, estilos).
//  2. Tenta localizar o artigo pelos ContentSelectors do perfil ativo.
//  3. Se nenhum seletor encontrar conteúdo, usa o DOM inteiro mas remove
//     os RemoveSelectors.
//
// Retorna o HTML original em caso de erro de parsing.
func ExtractMainContent(htmlContent, rawURL string) string {
	cfg := ConfigForURL(rawURL)

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlContent))
	if err != nil {
		// Não conseguimos parsear — retornar original para não perder dados.
		return htmlContent
	}

	// --- Passo 1: remover elementos indesejáveis independente de estratégia ---
	removeAlways := []string{
		"img", "picture", "figure", "figcaption", // imagens
		"script", "noscript", "style",             // código e estilos inline
		"iframe", "embed", "object", "video",      // mídia incorporada
		"form", "input", "button", "select",       // formulários
		"svg",                                      // vetores
	}
	for _, sel := range removeAlways {
		doc.Find(sel).Remove()
	}

	// --- Passo 2: tentativa semântica ---
	for _, sel := range cfg.ContentSelectors {
		node := doc.Find(sel)
		if node.Length() == 0 {
			continue
		}
		// Exige que o nó tenha texto real para evitar selecionar containers vazios.
		if strings.TrimSpace(node.Text()) == "" {
			continue
		}
		html, err := node.First().Html()
		if err == nil && strings.TrimSpace(html) != "" {
			return html
		}
	}

	// --- Passo 3: fallback — remover ruído do DOM completo ---
	for _, sel := range cfg.RemoveSelectors {
		doc.Find(sel).Remove()
	}

	html, err := doc.Find("body").Html()
	if err != nil || strings.TrimSpace(html) == "" {
		return htmlContent
	}
	return html
}
