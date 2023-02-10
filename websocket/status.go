package websocket

//go:generate stringer -type=WebsocketStatus -trimprefix=WebsocketStatus -output=status_string.go
type WebsocketStatus uint8

const (
	WebsocketStatusDisconnected WebsocketStatus = iota
	WebsocketStatusWaiting
	WebsocketStatusConnecting
	WebsocketStatusDisconnecting
	WebsocketStatusEstablish // unused
	WebsocketStatusInactive  // unused
	WebsocketStatusConnected
	WebsocketStatusReConnecting
)
