package main

// EventType is used to uniquely identify each type of event
type EventType int

// Event is the base type for all events
type Event interface {
	GetType() EventType
}

const (
	EV_NONE     = iota // Nothing happened
	EV_QUIT            // Quit
	EV_KEYPRESS        // Key press
)

// EventNone is the nil/nothing happened event
type EventNone struct{}

func (e EventNone) GetType() EventType {
	return EV_NONE
}

// EventQuit is sent to indicate the application should exit
type EventQuit struct{}

func (e EventQuit) GetType() EventType {
	return EV_QUIT
}

// EventKeypress is sent when a key press is detected
type EventKeypress struct {
	Key Keypress
}

func (e EventKeypress) GetType() EventType {
	return EV_KEYPRESS
}
