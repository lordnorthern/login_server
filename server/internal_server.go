package server

import (
	"github.com/lordnorthern/login_server/helpers"
	"github.com/lordnorthern/login_server/models"
)

// InternalConnectionHandler will handle a connection
func InternalConnectionHandler(user *models.User) {
	for {
		buffer := make([]byte, 1400)
		packetSize, err := (*user.Conn).Read(buffer)
		if err != nil {
			helpers.LogError(err)
			break
		}
		_ = packetSize
		cmd := models.NewInternalCommand()
		cmd.Decode(buffer)
		action := newAction(cmd, user)
		action.rerouteCommands()
	}
}
