package types

import (
	"fmt"
)

// TransportType represents different types of transport
type TransportType struct {
	API string
}

// Transport represents a transport line with its details
type Transport struct {
	Type        TransportType
	Number      string
	Stop        string
	Destination string
}

func (t Transport) String() string {
	return fmt.Sprintf("%s %s @ %s -> %s", t.Type.API, t.Number, t.Stop, t.Destination)
}

// Line represents a transport line
type Line struct {
	Type TransportType
	ID   string
}

func (l Line) String() string {
	return fmt.Sprintf("%s %s", l.Type.API, l.ID)
}

// Request represents a request for transport information
type Request struct {
	RouteID   string
	StopID    string
	StopName  string
	Direction string
}

// Result represents a transport timing result
type Result struct {
	Dest string `json:"dest"`
	Time string `json:"time"`
}

func (r Result) String() string {
	return fmt.Sprintf("%s: %s", r.Dest, r.Time)
}
