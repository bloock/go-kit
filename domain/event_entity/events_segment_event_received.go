package event_entity

import "time"

type SegmentEvenReceived struct {
	EventType  string     `json:"type"`
	UserId     string     `json:"userId"`
	Event      string     `json:"event"`
	Context    Context    `json:"context"`
	Properties Properties `json:"properties"`
	Timestamp  time.Time  `json:"timestamp"`
}

func NewSegmentEvenReceived(eventType string, userId string, event string, context Context,
	properties Properties, timestamp time.Time) *SegmentEvenReceived {
	return &SegmentEvenReceived{
		EventType:  eventType,
		UserId:     userId,
		Event:      event,
		Context:    context,
		Properties: properties,
		Timestamp:  timestamp,
	}
}

type Location struct {
	Country   string  `json:"country"`
	City      string  `json:"city"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type Context struct {
	Location Location `json:"location"`
	Ip       string   `json:"ip"`
}
type Properties struct {
	Records []string `json:"records"`
	Success bool     `json:"success"`
}
