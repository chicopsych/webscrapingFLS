// Package converter fornece a conversão de HTML para Markdown formatado.
//
// Separa a responsabilidade de transformação de formato do crawler (que extrai HTML)
// e do writer (que persiste Markdown), seguindo o princípio de Single Responsibility.
//
// # Por que uma biblioteca externa e não implementação própria?
//
// Converter HTML arbitrário para Markdown é complexo: envolve headings aninhados,
// listas ordenadas/não-ordenadas, tabelas, código inline/bloco, links absolutos/relativos,
// imagens, negrito/itálico sobreposto e entidades HTML. Implementar isso corretamente
// exigiria um parser HTML completo + engine de transformação. A biblioteca
// JohannesKaufmann/html-to-markdown (v1.6) resolve todos esses casos com testes extensivos.
//
// # Resolução de URLs Relativas
//
// Quando baseURL é fornecido, o conversor resolve links relativos (ex: "/sobre") e
// src de imagens em URLs absolutas (ex: "https://example.com/sobre").
// Isso torna o Markdown gerado autossuficiente, sem depender do domínio de origem.
package converter

import (
	"strings"

	md "github.com/JohannesKaufmann/html-to-markdown"
)

// HTMLToMarkdown converte uma string de HTML em Markdown formatado.
//
// Parâmetros:
//   - htmlContent: HTML interno do <body> da página (sem a tag <body> em si).
//   - baseURL: URL base da página para resolução de links e imagens relativos.
//     Pode ser string vazia; nesse caso, links relativos são mantidos como estão.
//
// Retorna o Markdown convertido. Em caso de erro interno do conversor, retorna
// o HTML original sem transformação (degradação graciosa / graceful degradation).
func HTMLToMarkdown(htmlContent string, baseURL string) string {
	conv := md.NewConverter(baseURL, true, nil)

	markdown, err := conv.ConvertString(htmlContent)
	if err != nil {
		// Degradação graciosa: se a conversão falhar, retorna o conteúdo sem formatação.
		return strings.TrimSpace(htmlContent)
	}

	return strings.TrimSpace(markdown)
}
