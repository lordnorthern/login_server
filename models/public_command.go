package models

import (
	"bytes"
	"encoding/base64"
	"encoding/gob"
)

// PublicCommand is used to send data between the server and the login client
type PublicCommand struct {
	Command    string
	SessionID  string
	Parameters map[string]string
	Details    []byte
}

// NewPublicCommand initiates a new empty command object.
func NewPublicCommand() *PublicCommand {
	cmd := new(PublicCommand)
	cmd.Parameters = make(map[string]string)
	return cmd
}

func (c *PublicCommand) GetCommandName() string {
	return c.Command
}

func (c *PublicCommand) SetCommandName(newCommandName string) {
	c.Command = newCommandName
}

func (c *PublicCommand) SetSessionID(newSessionID string) {
	c.SessionID = newSessionID
}
func (c *PublicCommand) GetSessionID() string {
	return c.SessionID
}

// Serialize will serialize the command.
// Use this instead of the global serialze in helpers
// This is necessary because different command types may (and do in fact) have different
// serialization algorithms
func (c *PublicCommand) Serialize() ([]byte, error) {
	buffer := bytes.Buffer{}
	encoder := gob.NewEncoder(&buffer)
	err := encoder.Encode(c)
	if err != nil {
		return []byte{}, err
	}
	encoded := make([]byte, base64.StdEncoding.EncodedLen(buffer.Len()))
	base64.StdEncoding.Encode(encoded, buffer.Bytes())
	return encoded, nil
}
