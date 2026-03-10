package errors

import "errors"

var (
	// Erros Críticos (Interrompem o fluxo / demandam retry)
	ErrNetworkUnreachable = errors.New("network unreachable: falha de DNS ou conexão TCP")
	ErrFileWriteAccess    = errors.New("file write access: falha ao salvar arquivo .md no disco")

	// Erros de Acesso/Bloqueio
	ErrRateLimited  = errors.New("rate limited: servidor retornou HTTP 429 (Muitas requisições)")
	ErrAccessDenied = errors.New("access denied: servidor retornou HTTP 403/401 (Possível Captcha ou sem permissão)")

	// Erros de Integridade de Dados
	ErrSelectorNotFound = errors.New("selector not found: elemento HTML (CSS/XPath) não encontrado")
	ErrInvalidMetadata  = errors.New("invalid metadata: metadados extraídos estão incompletos ou corrompidos")
)
