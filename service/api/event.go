package api

type EventCategory string

const (
	EventCategorySYSLOG  EventCategory = "CAT_SYSLOG_0"
	EventCategorySTATUS  EventCategory = "CAT_STATUS_1"
	EventCategoryNE107   EventCategory = "CAT_NE107_2"
	EventCategoryGENERIC EventCategory = "CAT_GENERIC_99"
)

type Event struct {
	Number      uint32        `json:"Number"`
	Description string        `json:"Description,omitempty"`
	Category    EventCategory `json:"Category"`
	Details     interface{}   `json:"Details"`
}

func (e *Event) Payload() any {
	return e
}
