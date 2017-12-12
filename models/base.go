package models

// Terminator is the interface for all objects that get terminated when the server shuts down
type Terminator interface {
	Terminate()
	GetName() string
	AddToList(*[]Terminator)
}

// Command is the main interface for commands that are being sent from and to the server
// This interface is implemented by commands exchanged by the public server and login client, and
// also by the internal server and an external game server.
type Command interface {
	GetSessionID() string
	SetSessionID(string)
	Serialize() ([]byte, error)
	GetCommandName() string
	SetCommandName(string)
}

// ParsedCommand contains the bytes slice of a command after it has been stitched together
// and also its length
type ParsedCommand struct {
	CmdBytes  []byte
	CmdLength int
}
