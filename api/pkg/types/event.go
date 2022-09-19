package types

type EventCategory string

const (
	EventCategory_SYSLOG  EventCategory = "CAT_SYSLOG_0"
	EventCategory_STATUS  EventCategory = "CAT_STATUS_1"
	EventCategory_NE107   EventCategory = "CAT_NE107_2"
	EventCategory_GENERIC EventCategory = "CAT_GENERIC_99"
)

type Event struct {
	Number      uint32        `json:"Number"`
	Description string        `json:"Description,omitempty"`
	Category    EventCategory `json:"Category"`
	Details     interface{}   `json:"Details"`
}
