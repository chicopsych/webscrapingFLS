# webscrapingFLS

Web scraper em Go para extração de conteúdo de páginas web, com saída em Markdown e metadados em YAML Front Matter.

## Funcionalidades

- Extração de título e conteúdo de qualquer URL via linha de comando
- Interface GUI local (web) iniciada com `-gui` ou sem argumentos
- Saída em arquivo `.md` com YAML Front Matter (título, URL, timestamp, tags)
- Nome de arquivo baseado no título da página com sanitização cross-platform
- Proteção contra sobrescrita com sufixo incremental automático (`_1`, `_2`, ...)
- Logging estruturado em JSON com suporte a nível `debug` (via `log/slog`)
- Tratamento de erros de rede categorizados (rate limit, acesso negado, seletor não encontrado, etc.)
- Requisições assíncronas com limite de paralelismo e delay entre chamadas (via Colly)

## Pré-requisitos

- [Go](https://golang.org/) >= 1.21

## Instalação

```bash
git clone https://github.com/chicopsych/webscrapingFLS.git
cd webscrapingFLS
go mod download
```

## Uso

```bash
# Compilar
go build -o webscrapingfls ./cmd/scraper

# Modo CLI
./webscrapingfls -url https://example.com

# Com opções adicionais
./webscrapingfls -url https://example.com -out ./resultados -debug

# Modo GUI local (abre no navegador)
./webscrapingfls -gui
```

No Windows (PowerShell), o equivalente é:

```powershell
.\webscrapingfls.exe -url https://example.com -out .\data
.\webscrapingfls.exe -gui
```

### Flags disponíveis

| Flag      | Padrão  | Descrição                                    |
|-----------|---------|----------------------------------------------|
| `-url`    | (vazio) | URL alvo para o scraping (modo CLI)          |
| `-out`    | `data`  | Diretório de saída para os arquivos `.md`    |
| `-debug`  | `false` | Habilita log em nível debug                  |
| `-gui`    | `false` | Inicia interface GUI local (web)             |

## Estrutura do Projeto

```
webscrapingFLS/
├── cmd/
│   └── scraper/
│       └── main.go              # Composition root (flags + escolha CLI/GUI)
├── internal/
│   ├── cli/
│   │   └── run.go               # Adaptador de interface CLI
│   ├── converter/
│   │   └── converter.go         # Conversão HTML → Markdown (html-to-markdown v1.6)
│   ├── crawler/
│   │   ├── crawler.go           # Lógica de scraping com Colly
│   │   └── events.go            # Eventos assíncronos de scraping
│   ├── errors/
│   │   └── errors.go            # Sentinel errors do domínio
│   ├── filesystem/
│   │   ├── naming.go            # Sanitização segura de nomes de arquivos
│   │   ├── uniqueness.go        # Geração de nome único sem sobrescrita
│   │   └── windows_reserved.go  # Nomes reservados do Windows/NTFS
│   ├── gui/
│   │   ├── run.go               # Adaptador de interface GUI (HTTP local)
│   │   └── templates/
│   │       └── index.html       # Template HTML embutido via go:embed
│   ├── logger/
│   │   └── logger.go            # Logger estruturado (slog/JSON)
│   ├── models/
│   │   └── models.go            # DTO PageData (conteúdo + metadados)
│   ├── service/
│   │   └── scrape/
│   │       └── service.go       # Use case ExecuteScrape (orquestrador)
│   └── writer/
│       └── writer.go            # Persistência Markdown + Front Matter YAML
├── data/                    # Diretório padrão de saída dos arquivos .md gerados
├── go.mod
├── go.sum
└── README.md
```

## Arquitetura e Fluxo de Dados

O projeto segue Separation of Concerns e princípios de Clean Architecture: o `main` apenas compõe dependências e escolhe a interface (CLI/GUI), enquanto a regra de negócio fica concentrada no caso de uso `ExecuteScrape`.

Fluxo atual:

```
main (flags/logger)
	-> cli.Run() ou gui.Run()
	-> service/scrape.ExecuteScrape()
	-> crawler.Scrape()
		-> converter.HTMLToMarkdown()  ← HTML bruto → Markdown formatado
	-> filesystem.SanitizeFileNameFromTitle()
	-> filesystem.UniqueFileName()
	-> writer.SaveMarkdown()
	-> data/<titulo_sanitizado>.md
```

1. **`cmd/scraper/main.go`** — Composition root: processa flags e delega.
2. **`internal/cli` e `internal/gui`** — Camada de apresentação (input/output).
3. **`internal/service/scrape`** — Camada de aplicação (orquestra regra de negócio).
4. **`internal/crawler` e `internal/writer`** — Infraestrutura de scraping e persistência.
5. **`internal/filesystem`** — Segurança e integridade de nomes (NTFS/EXT4 + anti-path traversal).

## Exemplo de Saída

O scraper gera um arquivo `.md` no diretório de saída com o seguinte formato:

```markdown
---
title: Example Domain
url: https://example.com
timestamp: 2026-03-10T12:49:34-03:00
---

# Example Domain

This domain is for use in documentation examples without needing permission.

[Learn more](https://iana.org/domains/example)
```

Headings (`#`, `##`...), listas, links, negrito/itálico, código e tabelas presentes na página
original são preservados como Markdown válido.

### Política de nome de arquivo

- Nome base é derivado de `title` com sanitização de caracteres inválidos.
- O sistema impede path traversal removendo separadores de caminho e forçando basename.
- Nomes reservados do Windows (ex.: `CON`, `PRN`, `COM1`) são neutralizados com prefixo seguro.
- Em colisão de nome, gera sufixo incremental automático (`_1`, `_2`, ...).

## Dependências Principais

| Pacote | Versão | Uso |
|--------|--------|-----|
| [gocolly/colly](https://github.com/gocolly/colly) | v2.3.0 | Motor de scraping HTTP/HTML |
| [JohannesKaufmann/html-to-markdown](https://github.com/JohannesKaufmann/html-to-markdown) | v1.6.0 | Conversão HTML → Markdown formatado |
| [gopkg.in/yaml.v3](https://pkg.go.dev/gopkg.in/yaml.v3) | v3.0.1 | Serialização do YAML Front Matter |

## Tratamento de Erros

O projeto usa o padrão **Sentinel Error** do Go: variáveis de pacote do tipo `error` declaradas com `errors.New()`. Cada sentinel carrega uma identidade única (ponteiro) que permite comparação semântica via `errors.Is()` — mesmo após wrapping com `fmt.Errorf("%w", ...)`.

Esse padrão é preferido em Go ao invés de exceções (como `try/catch` em Java/Python) porque:
- O chamador é **obrigado** a tratar o erro no fluxo de controle (não pode ignorar silenciosamente)
- Erros são **valores comuns** que podem ser retornados, armazenados e comparados
- A cadeia de wrapping preserva contexto sem perder a identidade do erro original

| Erro | Situação | Camada |
|------|----------|--------|
| `ErrNetworkUnreachable` | Falha de DNS ou TCP handshake (host não respondeu ao SYN) | Rede/Transporte |
| `ErrRateLimited` | HTTP 429 — servidor impôs limite de requisições | Aplicação (HTTP) |
| `ErrAccessDenied` | HTTP 401 (autenticação) ou 403 (autorização/WAF) | Aplicação (HTTP) |
| `ErrSelectorNotFound` | Seletor CSS não encontrou elementos no DOM (possível SPA) | Parsing (DOM) |
| `ErrInvalidMetadata` | Dados insuficientes para gerar Front Matter YAML válido | Validação |
| `ErrFileWriteAccess` | Permissão negada ou falha de I/O (POSIX 0644 / NTFS ACL) | Sistema de Arquivos |

## Contrato de Dados (PageData)

A estrutura interna `PageData` é o DTO (Data Transfer Object) central do projeto, definida no pacote `models`. Ela conecta o crawler (produtor) ao writer (consumidor) sem acoplamento direto.

### Campos do Front Matter YAML

| Campo | Tipo | Struct Tag | Descrição |
|-------|------|------------|-----------|
| `title` | `string` | `yaml:"title"` | Conteúdo da tag `<title>` do HTML |
| `url` | `string` | `yaml:"url"` | URL completa da página extraída |
| `timestamp` | `time.Time` | `yaml:"timestamp"` | Momento da extração (RFC 3339) |
| `tags` | `[]string` | `yaml:"tags,omitempty"` | Categorização (omitido se vazio) |

### Corpo do Markdown

O campo canônico é `Content` (com struct tag `yaml:"-"` para exclusão do YAML). O conteúdo é separado do front matter pelo writer, que monta a estrutura `--- metadata --- body`.

### Struct Tags e Reflection

As anotações `` `yaml:"..."` `` são lidas em runtime via reflection pela biblioteca `gopkg.in/yaml.v3`. Destaques:
- `omitempty` → campo omitido na serialização se for zero value (nil, "", 0)
- `yaml:"-"` → campo completamente excluído da serialização

### Migração Gradual (Backward Compatibility)

Durante a transição, o projeto mantém compatibilidade com o campo legado `MarkdownBody`:

- O crawler preenche `Content` e também replica em `MarkdownBody`.
- O writer prioriza `Content`; se estiver vazio, usa fallback para `MarkdownBody`.

Isso permite migração gradual sem quebrar consumidores já existentes.

## Conceitos de Go Aplicados

Este projeto aplica diversos idioms e padrões recomendados pela comunidade Go:

| Conceito | Onde é Aplicado | Descrição |
|----------|----------------|-----------|
| **Sentinel Errors** | `internal/errors/` | Erros como valores comparáveis via `errors.Is()` |
| **Error Wrapping** | `internal/writer/` | `fmt.Errorf("%w", err)` preserva contexto + identidade |
| **Struct Tags** | `internal/models/` | Controle de serialização YAML via reflection |
| **slog (Structured Logging)** | `internal/logger/` | Logger JSON nativo do Go 1.21+ com níveis e atributos |
| **Separation of Concerns** | `cmd/` + `internal/*` | Fronteiras claras entre entrada, interface e negócio |
| **Use Case / Interactor** | `internal/service/scrape/` | Regra central reutilizada por CLI e GUI |
| **go:embed** | `internal/gui/run.go` | Template HTML externo embutido no binário |
| **Adaptador de Biblioteca** | `internal/converter/` | Encapsula dependência externa com API interna limpa |
| **Path Traversal Defense** | `internal/filesystem/` | Sanitização + basename para conter caminhos maliciosos |
| **NTFS/EXT4 Portability** | `internal/filesystem/` | Nome de arquivo compatível em múltiplos filesystems |
| **Goroutines + Channels** | `internal/crawler/events.go` | Comunicação assíncrona entre componentes |
| **select non-blocking** | `internal/crawler/events.go` | `select { case ch <- v: default: }` evita deadlock |
| **Injeção de Dependência** | todos os pacotes | Logger passado como parâmetro, não como global |
| **internal/ packages** | todo o projeto | Pacotes não exportáveis — encapsulamento a nível de módulo |
| **Zero Values úteis** | `internal/models/` | Structs inicializadas sem construtor são seguras de usar |
| **Async com Colly** | `internal/crawler/` | `colly.Async(true)` + `LimitRule` para paralelismo controlado |

## Estado Atual

- [x] Estrutura de pacotes com SoC/Clean Architecture (`cli`, `gui`, `service`, `filesystem`, `crawler`, `writer`, `logger`, `errors`, `models`)
- [x] Scraping assíncrono com Colly
- [x] GUI local via navegador com template embutido (`go:embed`)
- [x] Logger estruturado em JSON
- [x] Saída em Markdown com YAML Front Matter
- [x] Erros tipados e categorizados
- [x] Nome de arquivo seguro por título com sanitização NTFS/EXT4
- [x] Proteção contra sobrescrita com sufixo incremental
- [x] Conversão real de HTML para Markdown (headings, listas, links, código, tabelas)
- [ ] Suporte a múltiplas URLs (batch)
- [ ] Extração de tags/metadados automática
- [ ] Testes unitários

## Autor

**Chicopsych** — [github.com/chicopsych](https://github.com/chicopsych)
