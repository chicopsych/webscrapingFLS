package models

import "time"

// PageData centraliza o conteúdo e os metadados para o YAML Front Matter.
type PageData struct {
	Title     string    `yaml:"title"`
	URL       string    `yaml:"url"`
	Timestamp time.Time `yaml:"timestamp"`
	Tags      []string  `yaml:"tags,omitempty"`
	Content   string    `yaml:"-"` // O conteúdo não vai no YAML, mas no corpo do MD.

	// Campos legados (compatibilidade temporária): manter até migrar todos os consumidores para Content.
	RawContent   string `yaml:"-"`
	MarkdownBody string `yaml:"-"`
}
