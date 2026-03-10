package errors

import "errors"

var (
	// Definição de erros tipados para facilitar o troubleshooting.
	ErrNetworkUnreachable = errors.New("rede inacessível ou falha de DNS")
	ErrRateLimited        = errors.New("limite de requisições atingido (HTTP 429)")
	ErrAccessDenied       = errors.New("acesso negado ou bloqueio por firewall (HTTP 401/403)")
	ErrSelectorNotFound   = errors.New("seletor HTML não encontrado na página")
	ErrInvalidMetadata    = errors.New("metadados extraídos são inválidos ou incompletos")
	ErrFileWriteAccess    = errors.New("falha de permissão ou escrita no sistema de arquivos")
)
