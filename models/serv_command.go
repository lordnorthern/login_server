package models

import (
	"bytes"
	"encoding/binary"
	"errors"
)

// ServCommand is used to hold the data exchanged by the login server, and a game server
type ServCommand struct {
	Command string
	Params  []string
}

// NewInternalCommand initializes a new ServCommand
func NewInternalCommand() *ServCommand {
	cmd := new(ServCommand)
	return cmd
}

func (c *ServCommand) SetSessionID(newSessionID string) {}
func (c *ServCommand) GetSessionID() string {
	return ""
}

func (c *ServCommand) GetCommandName() string {
	return c.Command
}

func (c *ServCommand) SetCommandName(newCommandName string) {
	c.Command = newCommandName
}

// Serialize will serialize the command.
// Use this instead of the global serialze in helpers
// This is necessary because different command types may (and do in fact) have different
// serialization algorithms
func (c ServCommand) Serialize() (OutBytes []byte, err error) {
	command := []byte(c.Command)
	command = append(command, 0)
	var commandLength []byte
	commandLength, err = encodeInt32(int32(len(command)))
	if err != nil {
		return
	}
	OutBytes = commandLength
	OutBytes = append(OutBytes, command...)
	var paramsLength []byte
	paramsLength, err = encodeInt32(int32(len(c.Params)))
	if err != nil {
		return
	}

	OutBytes = append(OutBytes, paramsLength...)
	for _, val := range c.Params {
		var param []byte
		param = []byte(val)
		param = append(param, 0)
		var paramLength []byte
		paramLength, err = encodeInt32(int32(len(param)))
		if err != nil {
			return
		}
		OutBytes = append(OutBytes, paramLength...)
		OutBytes = append(OutBytes, param...)
	}
	return
}

// Decode uses a custom decoding algorithm to speak with Unreal Engine.
// This is a work in progress and is probably a shitty way of doing things.
func (c *ServCommand) Decode(inBytes []byte) {
	var err error
	var point int32
	var cmdLength int32
	cmdLength, err = decodeInt32(inBytes[point:4])
	if err != nil {
		return
	}
	point = 4
	c.Command = string(inBytes[point : point+cmdLength-1])
	point = point + cmdLength
	var sliceLength int32
	sliceLength, err = decodeInt32(inBytes[point : point+4])
	if err != nil {
		return
	}
	point += 4
	c.Params = make([]string, sliceLength)
	for ind, _ := range c.Params {
		var paramLength int32
		paramLength, err = decodeInt32(inBytes[point : point+4])
		point += 4
		if err != nil {
			continue
		}
		c.Params[ind] = string(inBytes[point : point+paramLength-1])
		point += paramLength
	}
}

func decodeInt32(inBytes []byte) (out int32, err error) {
	if len(inBytes) != 4 {
		err = errors.New("InvalidInBytes")
	} else {
		err = binary.Read(bytes.NewReader(inBytes), binary.LittleEndian, &out)
	}
	return
}

func encodeInt32(inNum int32) (out []byte, err error) {
	buf := new(bytes.Buffer)
	err = binary.Write(buf, binary.LittleEndian, inNum)
	out = buf.Bytes()
	return
}
