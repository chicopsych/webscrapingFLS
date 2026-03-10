// Package errors define os erros sentinela (sentinel errors) do domínio do web scraper.
//
// Em Go, erros são tratados como valores comuns — não como exceções (Java/Python).
// Essa é uma decisão central da linguagem, resumida na máxima de Rob Pike:
// "Errors are values" (https://go.dev/blog/errors-are-values).
//
// # O Padrão Sentinel Error
//
// Sentinel errors são variáveis de pacote do tipo error, declaradas com errors.New().
// Cada variável carrega uma identidade única (ponteiro) que permite comparação semântica
// via errors.Is() em qualquer ponto da call stack — mesmo após wrapping com fmt.Errorf("%w", ...).
//
// Diferente de exceções (throw/catch), esse padrão obriga o chamador a tratar o erro
// explicitamente no fluxo de controle, resultando em código mais previsível e depurável.
//
// # Rastreabilidade e Debugging (Kali Linux)
//
// Todos os erros deste pacote são capturados pelo logger (slog) no formato JSON estruturado.
// Em um ambiente como o Kali Linux, onde o scraper pode ser usado para reconhecimento em
// pentests, esses erros categorizados permitem filtrar rapidamente eventos no terminal:
//
//	cat output.json | jq 'select(.error | contains("429"))'
//
// Essa rastreabilidade é essencial para correlacionar falhas de rede com regras de
// firewall (iptables/nftables) ou com políticas de rate limiting do alvo.
//
// # Go Idiom: por que "var" e não "const"?
//
// Sentinel errors são declarados com var (não const) porque errors.New() retorna
// um ponteiro (*errors.errorString) — ponteiros não são constantes em tempo de compilação
// em Go. Essa é uma sutileza importante: o valor do erro é a identidade do ponteiro,
// não a string em si. Dois errors.New("mesma msg") produzem erros DIFERENTES.
package errors

import "errors"

var (
	// ErrNetworkUnreachable indica falha na camada de rede antes de receber
	// qualquer resposta HTTP do servidor alvo.
	//
	// Cenários comuns:
	//   - Falha no TCP 3-Way Handshake (SYN → SYN-ACK → ACK): o host não respondeu
	//     ao pacote SYN inicial dentro do timeout do sistema operacional.
	//   - Falha na resolução DNS: o sistema não conseguiu traduzir o domínio para IP.
	//     No Kali Linux, a configuração DNS fica em /etc/resolv.conf; em ambientes
	//     com VPN ou pivoting, esse arquivo pode apontar para um DNS inacessível.
	//   - O Colly abstrai toda essa camada via net/http do Go. Quando o transport
	//     subjacente retorna um erro (ex: dial tcp: no route to host), o Colly
	//     repassa via callback OnError com StatusCode == 0.
	//
	// No pipeline do scraper, esse erro é logado pelo slog como JSON estruturado,
	// permitindo debugging com ferramentas como jq, grep ou ELK Stack.
	ErrNetworkUnreachable = errors.New("rede inacessível ou falha de DNS")

	// ErrRateLimited indica que o servidor respondeu com HTTP 429 (Too Many Requests).
	//
	// O código HTTP 429 é definido na RFC 6585 e sinaliza que o cliente ultrapassou
	// o limite de requisições permitido em um intervalo de tempo. Servidores podem
	// incluir o header Retry-After indicando quando o cliente pode tentar novamente.
	//
	// O Colly oferece mitigação preventiva via colly.LimitRule, configurando:
	//   - Parallelism: número máximo de goroutines simultâneas por domínio
	//   - Delay: intervalo mínimo entre requisições ao mesmo domínio
	//
	// Mesmo com essas proteções, servidores com rate limiting agressivo (ex: Cloudflare)
	// podem retornar 429. Esse erro permite ao chamador implementar backoff exponencial.
	ErrRateLimited = errors.New("limite de requisições atingido (HTTP 429)")

	// ErrAccessDenied indica que o servidor recusou a requisição por questões de
	// autenticação (HTTP 401) ou autorização (HTTP 403).
	//
	// A distinção entre os dois status codes é fundamental em segurança:
	//   - 401 Unauthorized: o servidor não sabe QUEM é o cliente (falta credencial).
	//     Apesar do nome "Unauthorized", trata-se de autenticação, não autorização.
	//   - 403 Forbidden: o servidor sabe quem é o cliente, mas ele NÃO TEM PERMISSÃO
	//     para acessar o recurso. Nenhuma credencial resolverá esse bloqueio.
	//
	// Em ambientes de pentest (Kali Linux), esse erro frequentemente indica:
	//   - WAF (Web Application Firewall) bloqueando User-Agent suspeito
	//   - IDS/IPS (Intrusion Detection/Prevention System) reagindo a padrões de scraping
	//   - Geoblocking via serviços como Cloudflare ou Akamai
	//
	// O scraper atualmente não envia headers customizados (ex: User-Agent real),
	// o que pode disparar 403 em sites com proteção anti-bot.
	ErrAccessDenied = errors.New("acesso negado ou bloqueio por firewall (HTTP 401/403)")

	// ErrSelectorNotFound indica que o seletor CSS usado para extrair dados do DOM
	// não encontrou nenhum elemento correspondente na página HTML.
	//
	// O Colly usa a biblioteca goquery (wrapper Go do jQuery) para navegar o DOM.
	// O parsing é feito sobre o HTML estático retornado pelo servidor — ou seja,
	// conteúdo renderizado via JavaScript (SPAs com React, Vue, Angular) NÃO estará
	// presente no DOM que o Colly analisa. Para esses casos seria necessário um
	// headless browser (ex: chromedp, Rod, Playwright).
	//
	// No projeto, o seletor principal é "title" (tag <title> do HTML), extraído
	// via e.ChildText("title") no callback OnHTML. Se a página não possuir essa tag,
	// ou se o HTML estiver malformado, este erro é disparado.
	ErrSelectorNotFound = errors.New("seletor HTML não encontrado na página")

	// ErrInvalidMetadata indica que os dados extraídos pelo crawler são inválidos
	// ou incompletos para gerar um arquivo Markdown válido com Front Matter YAML.
	//
	// A validação de metadados ocorre antes da serialização: se campos obrigatórios
	// como Title ou URL estão vazios, o pipeline não deve prosseguir para o writer,
	// evitando arquivos .md corrompidos ou sem informação útil.
	//
	// Esse erro protege a integridade do fluxo de dados:
	//   crawler.Scrape() → PageData (validação) → writer.SaveMarkdown()
	ErrInvalidMetadata = errors.New("metadados extraídos são inválidos ou incompletos")

	// ErrFileWriteAccess indica falha ao gravar o arquivo .md no sistema de arquivos.
	//
	// Esse erro encapsula problemas de I/O no nível do OS, incluindo:
	//   - Permissões insuficientes: em sistemas POSIX (Linux/Kali), o arquivo é criado
	//     com modo 0644 (rw-r--r--) via os.WriteFile. Se o diretório de saída pertencer
	//     a root e o scraper rodar como usuário comum, a escrita falhará.
	//   - Em NTFS (Windows): as permissões são gerenciadas por ACLs (Access Control Lists),
	//     não por bits POSIX. O Go abstrai isso — os.WriteFile funciona em ambos os
	//     filesystems — mas o modo 0644 é ignorado pelo NTFS. Caminhos com caracteres
	//     especiais (ex: ":", "*") são válidos em ext4 mas inválidos em NTFS.
	//   - Disco cheio ou path inexistente: o writer usa os.MkdirAll() antes de gravar,
	//     mas espaço em disco deve ser verificado externamente.
	//
	// O erro é wrappado com fmt.Errorf("%w: %v", ErrFileWriteAccess, err) no writer,
	// preservando a identidade do sentinel para errors.Is() e adicionando o erro
	// original do OS para contexto completo no log.
	ErrFileWriteAccess = errors.New("falha de permissão ou escrita no sistema de arquivos")
)
