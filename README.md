# webscrapingFLS

Web scraper em Go para extração de conteúdo de páginas web, com saída em Markdown e metadados em YAML Front Matter.

## Funcionalidades

- Extração de título e conteúdo de qualquer URL via linha de comando
- Saída em arquivo `.md` com YAML Front Matter (título, URL, timestamp, tags)
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
go build -o scraper ./cmd/scraper

# Executar
./scraper -url https://example.com

# Com opções adicionais
./scraper -url https://example.com -out ./resultados -debug
```

### Flags disponíveis

| Flag      | Padrão  | Descrição                                    |
|-----------|---------|----------------------------------------------|
| `-url`    | (vazio) | **Obrigatória.** URL alvo para o scraping    |
| `-out`    | `data`  | Diretório de saída para os arquivos `.md`    |
| `-debug`  | `false` | Habilita log em nível debug                  |

## Estrutura do Projeto

```
webscrapingFLS/
├── cmd/
│   └── scraper/
│       └── main.go          # Entrypoint da aplicação (flags, orquestração)
├── internal/
│   ├── crawler/
│   │   └── crawler.go       # Lógica de scraping com Colly (async, limites, parsing)
│   ├── errors/
│   │   └── errors.go        # Erros tipados e categorizados do domínio
│   ├── logger/
│   │   └── logger.go        # Inicialização do logger estruturado (slog/JSON)
│   ├── models/
│   │   └── models.go        # Struct PageData (dados extraídos e metadados)
│   └── writer/
│       └── writer.go        # Serialização para Markdown + YAML Front Matter
├── data/                    # Diretório padrão de saída dos arquivos .md gerados
├── main.go                  # Arquivo raiz (comentário de versão/projeto)
├── go.mod
├── go.sum
└── README.md
```

## Arquitetura e Fluxo de Dados

O projeto segue uma arquitetura modular onde cada pacote interno tem responsabilidade única. O fluxo de dados é linear e unidirecional:

```
┌──────────────────────────────────────────────────────────────────────┐
│                        FLUXO DE DADOS                              │
│                                                                    │
│  ┌─────────┐     ┌───────────┐     ┌───────────┐     ┌─────────┐  │
│  │  main   │────▶│  crawler   │────▶│ PageData  │────▶│ writer  │  │
│  │ (flags) │     │ (Colly)    │     │ (models)  │     │ (.md)   │  │
│  └─────────┘     └───────────┘     └───────────┘     └─────────┘  │
│       │                │                                   │       │
│       │           ┌────┴────┐                              │       │
│       ▼           │ events  │                              ▼       │
│  ┌─────────┐      │ (chan)  │                        ┌──────────┐  │
│  │ logger  │◀─────┴────────▶│                        │ arquivo  │  │
│  │ (slog)  │                                         │   .md    │  │
│  └─────────┘                                         └──────────┘  │
└──────────────────────────────────────────────────────────────────────┘
```

1. **`main.go`** — Faz parsing dos flags CLI (`-url`, `-out`, `-debug`) e orquestra o pipeline.
2. **`logger`** — Inicializa o `slog` com handler JSON. Todos os módulos recebem o logger por injeção de dependência.
3. **`crawler`** — Usa o Colly para fazer a requisição HTTP, parsear o DOM e popular a struct `PageData`. Emite eventos via channel para feedback de progresso.
4. **`models`** — Define `PageData`, o contrato de dados central. Não contém lógica, apenas a estrutura.
5. **`writer`** — Serializa `PageData` em YAML Front Matter + corpo Markdown e grava no disco.
6. **`errors`** — Define erros sentinela do domínio, usados por todos os módulos para categorização.

## Exemplo de Saída

O scraper gera um arquivo `.md` no diretório de saída com o seguinte formato:

```markdown
---
title: Example Domain
url: https://example.com
timestamp: 2026-03-10T10:56:54-03:00
---

Example Domain
This domain is for use in documentation examples...
```

## Dependências Principais

| Pacote | Versão | Uso |
|--------|--------|-----|
| [gocolly/colly](https://github.com/gocolly/colly) | v2.3.0 | Motor de scraping HTTP/HTML |
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
| **Goroutines + Channels** | `internal/crawler/events.go` | Comunicação assíncrona entre crawler e UI/CLI |
| **select non-blocking** | `internal/crawler/events.go` | `select { case ch <- v: default: }` evita deadlock |
| **Injeção de Dependência** | todos os pacotes | Logger passado como parâmetro, não como global |
| **internal/ packages** | todo o projeto | Pacotes não exportáveis — encapsulamento a nível de módulo |
| **Zero Values úteis** | `internal/models/` | Structs inicializadas sem construtor são seguras de usar |
| **Async com Colly** | `internal/crawler/` | `colly.Async(true)` + `LimitRule` para paralelismo controlado |

## Estado Atual

- [x] Estrutura de pacotes internos (`crawler`, `writer`, `logger`, `errors`, `models`)
- [x] Scraping assíncrono com Colly
- [x] Logger estruturado em JSON
- [x] Saída em Markdown com YAML Front Matter
- [x] Erros tipados e categorizados
- [ ] Conversão real de HTML para Markdown (atualmente usa texto plano)
- [ ] Suporte a múltiplas URLs (batch)
- [ ] Extração de tags/metadados automática
- [ ] Testes unitários

## Autor

**Chicopsych** — [github.com/chicopsych](https://github.com/chicopsych)
