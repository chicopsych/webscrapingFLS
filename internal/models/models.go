package models

import "time"

// PageData representa o conteúdo extraído e estruturado
type PageData struct {
	Title        string    `json:"title" yaml:"title"`
	URL          string    `json:"url" yaml:"url"`
	Timestamp    time.Time `json:"timestamp" yaml:"timestamp"`
	Author       string    `json:"author,omitempty" yaml:"author,omitempty"`
	RawContent   string    `json:"-" yaml:"-"`        // HTML original ou texto limpo
	MarkdownBody string    `json:"markdown" yaml:"-"` // Texto convertido para MD
	Tags         []string  `json:"tags,omitempty" yaml:"tags,omitempty"`
}
