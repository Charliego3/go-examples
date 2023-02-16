package websocket

//go:generate stringer -type=Status -trimprefix=Status -output=status_string.go
type Status uint8

const (
	StatusDisconnected Status = iota
	StatusWaiting
	StatusConnecting
	StatusDisconnecting
	StatusEstablish // unused
	StatusInactive  // unused
	StatusConnected
	StatusReConnecting
)
