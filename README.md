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

| Erro | Situação |
|------|----------|
| `ErrNetworkUnreachable` | Rede inacessível ou falha de DNS |
| `ErrRateLimited` | Limite de requisições atingido (HTTP 429) |
| `ErrAccessDenied` | Acesso negado ou bloqueio por firewall (HTTP 401/403) |
| `ErrSelectorNotFound` | Seletor HTML não encontrado na página |
| `ErrInvalidMetadata` | Metadados extraídos inválidos ou incompletos |
| `ErrFileWriteAccess` | Falha de permissão ou escrita no sistema de arquivos |

## Contrato de Dados (PageData)

A estrutura interna `PageData` usa os campos abaixo no front matter YAML:

- `title`
- `url`
- `timestamp`
- `tags` (opcional)

Para o corpo do markdown, o campo canônico é `Content`.

Durante a transição, o projeto mantém compatibilidade com o campo legado `MarkdownBody`:

- O crawler preenche `Content` e também replica em `MarkdownBody`.
- O writer prioriza `Content`; se estiver vazio, usa fallback para `MarkdownBody`.

Isso permite migração gradual sem quebrar consumidores já existentes.

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
