// Package models define as estruturas de dados (DTOs) compartilhadas entre os módulos do scraper.
//
// # Por que um pacote separado para modelos?
//
// Em Engenharia de Software, o princípio de Separação de Responsabilidades (Separation of
// Concerns) recomenda que estruturas de dados de domínio vivam em seu próprio pacote.
// Isso evita dependências circulares: tanto o crawler quanto o writer importam models,
// mas não precisam importar um ao outro. Esse padrão é essencial em projetos Go com
// múltiplos pacotes internos.
//
// # Structs em Go vs Classes em OOP
//
// Go não possui classes, herança ou construtores. A unidade fundamental de agregação
// de dados é a struct — um tipo de valor (value type) que agrupa campos tipados.
//
// Em linguagens OOP (Java, Python), você usaria uma classe com getters/setters e herança.
// Em Go, a filosofia é diferente:
//   - Composição sobre herança: embed structs dentro de outras structs (não usado aqui,
//     mas é o mecanismo Go para "reutilizar" comportamento).
//   - Sem construtores: structs são inicializadas com zero values seguros (string → "",
//     int → 0, slice → nil, time.Time → zero time). Isso elimina NullPointerException.
//   - Métodos são funções com receiver: func (p PageData) Validate() error { ... }
//
// # Fluxo de Dados no Projeto
//
// PageData é o contrato central que conecta todos os módulos:
//
//	              ┌─────────────┐
//	URL input ──▶ │  crawler     │ ── Scrape() retorna PageData
//	              └──────┬──────┘
//	                     │
//	                     ▼
//	              ┌─────────────┐
//	              │  PageData   │ ← Struct definida neste pacote
//	              └──────┬──────┘
//	                     │
//	                     ▼
//	              ┌─────────────┐
//	              │  writer     │ ── SaveMarkdown() consome PageData
//	              └─────────────┘
//
// Essa separação garante que o crawler produz dados e o writer os consome,
// sem acoplamento direto entre os dois módulos.
package models

import "time"

// PageData centraliza o conteúdo extraído e os metadados de uma página web
// para serialização em arquivo Markdown com YAML Front Matter.
//
// O YAML Front Matter é um bloco de metadados no topo de arquivos Markdown,
// delimitado por "---". Ferramentas como Hugo, Jekyll e Obsidian o utilizam
// para extrair título, tags, data e outros metadados sem parsear o corpo do texto.
//
// # Struct Tags e Reflection
//
// As anotações `yaml:"..."` são struct tags — metadados que a biblioteca gopkg.in/yaml.v3
// lê em tempo de execução via reflection (pacote reflect). Cada tag controla como o campo
// é serializado:
//   - `yaml:"title"`     → o campo Title vira a chave "title" no YAML
//   - `yaml:"tags,omitempty"` → "tags" aparece no YAML apenas se Tags não for nil/vazio
//   - `yaml:"-"`         → o campo é EXCLUÍDO da serialização YAML
//
// Go Idiom: struct tags são strings brutas (raw strings) dentro de backticks.
// O compilador não as valida — erros de digitação só aparecem em runtime.
// Use `go vet` para detectar tags malformadas.
//
// # Go Idiom: Zero Values
//
// Uma PageData recém-declarada (var p models.PageData) é imediatamente utilizável:
//   - Title, URL, Content → "" (string vazia)
//   - Timestamp → time.Time{} (1 de janeiro do ano 1, 00:00:00 UTC)
//   - Tags → nil (slice nil, que se comporta como slice vazio em range/len/append)
//
// Isso significa que NÃO precisamos de um construtor NewPageData().
// Em Go, se o zero value é útil, não crie um construtor — esse é um idiom fundamental.
type PageData struct {
	// Title armazena o conteúdo da tag <title> do HTML.
	// Extraído pelo crawler via e.ChildText("title") (goquery/Colly).
	// É o primeiro campo do Front Matter YAML, essencial para identificação do documento.
	Title string `yaml:"title"`

	// URL armazena o endereço completo da página extraída.
	// Usamos string (não *url.URL) por simplicidade: a serialização YAML de uma string
	// é direta, enquanto url.URL exigiria um marshaler customizado. Como o URL já chega
	// validado pelo Colly (que faz parsing internamente), não há necessidade de re-validar.
	URL string `yaml:"url"`

	// Timestamp registra o momento exato da extração, usando time.Time do Go.
	// A serialização YAML produz formato RFC 3339 (ex: "2026-03-10T10:56:54-03:00"),
	// que é o padrão ISO 8601 usado em logs, APIs REST e bancos de dados.
	// Esse campo é crucial para rastreabilidade: permite saber QUANDO o conteúdo foi capturado.
	Timestamp time.Time `yaml:"timestamp"`

	// Tags é um slice de strings para categorização futura do conteúdo extraído.
	//
	// Go Idiom: slice nil vs slice vazio
	//   - var tags []string       → nil (len=0, cap=0) — serializa como NADA (omitempty)
	//   - tags := []string{}      → não-nil (len=0, cap=0) — serializa como "tags: []"
	//   - tags := make([]string, 0) → equivalente ao anterior
	//
	// A tag `omitempty` garante que, enquanto Tags for nil (não populado pelo crawler),
	// o campo simplesmente não aparece no Front Matter YAML, mantendo a saída limpa.
	Tags []string `yaml:"tags,omitempty"`

	// Content armazena o corpo principal da página, extraído do <body> do HTML.
	//
	// A tag `yaml:"-"` EXCLUI este campo da serialização YAML. Isso é intencional:
	// no formato Markdown com Front Matter, o conteúdo vai ABAIXO do bloco "---",
	// separado dos metadados. O writer monta manualmente essa estrutura:
	//
	//   ---
	//   title: Example Domain      ← Front Matter (campos com tags yaml)
	//   url: https://example.com
	//   ---
	//
	//   Example Domain...           ← Body (campo Content, excluído do YAML)
	//
	// O conteúdo é convertido de HTML para Markdown formatado pelo pacote internal/converter,
	// preservando headings, listas, links, código e tabelas da página original.
	Content string `yaml:"-"`

	// --- Campos Legados (Compatibilidade Temporária) ---
	//
	// Best Practice: Migração Gradual (Backward Compatibility)
	//
	// Quando um campo precisa ser renomeado ou substituído em um sistema em produção,
	// a abordagem segura é manter AMBOS os campos durante a transição:
	//   1. O produtor (crawler) preenche o campo novo (Content) E o legado (MarkdownBody)
	//   2. O consumidor (writer) prioriza Content, com fallback para MarkdownBody
	//   3. Após migrar todos os consumidores, os campos legados são removidos
	//
	// Esse padrão evita breaking changes e é amplamente usado em APIs versionadas.

	// RawContent armazena o texto bruto (não processado) do HTML.
	// Mantido para debug e futura pipeline de processamento HTML→Markdown.
	RawContent string `yaml:"-"`

	// MarkdownBody é o campo legado para o corpo do documento.
	// Depreciado em favor de Content. O writer faz fallback para este campo
	// caso Content esteja vazio (ver writer.SaveMarkdown).
	MarkdownBody string `yaml:"-"`
}
