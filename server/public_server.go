package server

import (
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
				buffer := make([]byte, models.BufferSize)
				trueBuffer := make([]byte, 0)
				var totalPSize int
				for {
					pSize, err := (*user.Conn).Read(buffer)
					if err != nil || pSize == 0 {
						helpers.LogError(err)
						connectionAlive = false
						break
					}
					trueBuffer = append(trueBuffer, buffer...)
					totalPSize += pSize
					if pSize < int(models.BufferSize) {
						break
					}
				}
				if !connectionAlive {
					break
				}
				rChunks, err := helpers.HandleChunk(trueBuffer, totalPSize)
				if err != nil {
					continue
				}
				for _, sChunk := range rChunks {
					if cmdPacket, found := cmdsBuffer[sChunk.Serial]; !found {
						cmdPacket.CmdBytes = make([][]byte, sChunk.TotalChunks)
						cmdPacket.CmdBytes[sChunk.Step] = sChunk.ChunkBytes
						cmdPacket.CmdLength = +int(sChunk.Length)
						cmdsBuffer[sChunk.Serial] = cmdPacket
					} else {
						cmdPacket.CmdBytes[sChunk.Step] = sChunk.ChunkBytes
						cmdsBuffer[sChunk.Serial] = cmdPacket
					}

					full := true
					for _, chk := range cmdsBuffer[sChunk.Serial].CmdBytes {
						if len(chk) == 0 {
							full = false
						}
					}
					if full {
						cmdChan <- cmdsBuffer[sChunk.Serial]
					}
				}

			}
		}()

		for parsedCmd := range cmdChan {
			var completeCommand []byte
			for _, sChunk := range parsedCmd.CmdBytes {
				completeCommand = append(completeCommand, sChunk...)
			}

			var Result []byte
			if len(user.EncryptionKey) > 0 {
				trimmedPacket, err := helpers.Decrypt(completeCommand, user.EncryptionKey)
				if err != nil {
					helpers.LogError(err)
					continue
				}
				Result = trimmedPacket
			} else {
				Result = completeCommand
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
