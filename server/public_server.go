package server

import (
	"bytes"
	"encoding/binary"

	"github.com/lordnorthern/login_server/helpers"
	"github.com/lordnorthern/login_server/models"
)

// PublicConnectionHandler will handle a connection
func PublicConnectionHandler(user *models.User) {
	connectionAlive := true
	for connectionAlive {
		cmdChan := make(chan models.ParsedCommand)
		go func() {
			cmdsBuffer := make(map[int32]models.ParsedCommand)
			for {
				var step int32
				var Serial int32

				buffer := make([]byte, 1400)
				pSize, err := (*user.Conn).Read(buffer)
				if err != nil || pSize == 0 {
					helpers.LogError(err)
					connectionAlive = false
					break
				}

				bufLength := bytes.NewReader(buffer[0:4])
				err = binary.Read(bufLength, binary.LittleEndian, &step)
				if err != nil {
					helpers.LogError(err)
					connectionAlive = false
					break
				}
				bufSerial := bytes.NewReader(buffer[4:8])
				err = binary.Read(bufSerial, binary.LittleEndian, &Serial)
				if err != nil {
					helpers.LogError(err)
					connectionAlive = false
					break
				}
				parsedCmd := cmdsBuffer[Serial]
				parsedCmd.CmdLength += (pSize - 8)
				parsedCmd.CmdBytes = buffer[8:pSize]
				if step == 0 {
					delete(cmdsBuffer, Serial)
					cmdChan <- parsedCmd
					continue
				} else {
					cmdsBuffer[Serial] = parsedCmd
				}
			}
		}()

		for parsedCmd := range cmdChan {
			completeCommand := parsedCmd.CmdBytes
			packetSize := parsedCmd.CmdLength
			var Result []byte
			if len(user.EncryptionKey) > 0 {
				trimmedPacket, err := helpers.Decrypt(completeCommand[:packetSize], user.EncryptionKey)
				if err != nil {
					helpers.LogError(err)
					continue
				}
				Result = trimmedPacket
			} else {
				Result = completeCommand[:packetSize]
			}
			cmd := models.NewPublicCommand()
			err := helpers.Unserialize(&Result, &cmd)
			if err != nil {
				helpers.LogError(err)
				connectionAlive = false
			}
			action := newAction(cmd, user)
			action.rerouteCommands()
		}
	}
}
