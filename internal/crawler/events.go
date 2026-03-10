package crawler

import "time"

// EventType representa os tipos de eventos emitidos durante o scraping.
type EventType string

const (
	EventStarted        EventType = "started"
	EventRequestSent    EventType = "request_sent"
	EventTitleExtracted EventType = "title_extracted"
	EventCompleted      EventType = "completed"
	EventError          EventType = "error"
)

// Event encapsula mensagens de progresso e erros para CLI/GUI.
type Event struct {
	Type      EventType
	Message   string
	URL       string
	Timestamp time.Time
	Progress  float64
	Err       error
}

func emitEvent(events chan<- Event, event Event) {
	if events == nil {
		return
	}

	select {
	case events <- event:
	default:
		// Evita bloqueio quando o consumidor estiver atrasado.
	}
}
