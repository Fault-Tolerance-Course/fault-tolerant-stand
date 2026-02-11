package order_events

import "encoding/json"

const correlationIDHeader = "x-correlation-id"

type BaseEvent struct {
	EventType string          `json:"event_type"`
	EntityID  string          `json:"entity_id"`
	Payload   json.RawMessage `json:"payload"`
}
