package server

import (
	"fmt"
	"strconv"

	"github.com/lordnorthern/login_server/helpers"
	"github.com/lordnorthern/login_server/models"
)

type action struct {
	user *models.User
	cmd  *models.Command
}

func newAction(cmd models.Command, user *models.User) *action {
	newAction := new(action)
	newAction.cmd = &cmd
	newAction.user = user
	return newAction
}

func (a *action) runAuthenticated(callback func()) bool {
	if (*a.cmd).GetSessionID() != a.user.Account.SessionID {
		(*a.user.Conn).Close()
		return false
	}
	callback()
	return true
}

// rerouteCommands reroutes commands to appropriate functions
func (a *action) rerouteCommands() {
	switch (*a.cmd).GetCommandName() {
	case "Hello":
		a.helloCmd()
	case "Login":
		a.loginCmd()
	case "NewCharacter":
		a.runAuthenticated(a.newCharacterCmd)
	case "ServUserLogin":
		a.ServUserLogin()
	}
}

func (a *action) ServUserLogin() {
	cmd, castErr := (*a.cmd).(*models.ServCommand)
	if !castErr {
		return
	}
	tmpCharID, convErr := strconv.Atoi(cmd.Params[0])
	characterID := uint(tmpCharID)
	if convErr != nil {
		return
	}
	account := findCharByUserSession(characterID, cmd.Params[1])
	fmt.Println(account)
	resp := models.NewInternalCommand()
	resp.Command = "Response"
	resp.Params = append(resp.Params, "Test")
	a.user.SendCommand(resp, false)
}

func (a *action) newCharacterCmd() {
	cmd, castErr := (*a.cmd).(*models.PublicCommand)
	if !castErr {
		return
	}
	newCharacter := models.Character{}
	if err := helpers.Unserialize(&cmd.Details, &newCharacter); err != nil {
		return
	}
	response := models.NewPublicCommand()
	response.Command = "NewCharacterResponse"

	if errs := newCharacter.Validate(); len(errs) > 0 {
		response.Parameters["Result"] = "Failed"
	} else if err := addNewCharacter(&newCharacter, (*a.user).Account); err != nil {
		response.Parameters["Result"] = "Failed"
	} else {
		response.Parameters["Result"] = "Success"
	}
	a.user.SendCommand(response, true)

	charsCmd := models.NewPublicCommand()
	charsCmd.Command = "CharactersList"
	charsCmd.Details, _ = helpers.Serialize(a.user.Account.Characters)
	a.user.SendCommand(charsCmd, true)
}

func (a *action) helloCmd() {
	a.user.CreateEncryptionCode()
	response := models.NewPublicCommand()
	response.Command = "HelloResponse"
	response.Parameters["EncryptionKey"] = string(a.user.EncryptionKey)
	a.user.SendCommand(response, false)
}

func (a *action) loginCmd() {
	cmd, castErr := (*a.cmd).(*models.PublicCommand)
	if !castErr {
		return
	}
	var emailPassword models.EmailPassword
	err := helpers.Unserialize(&cmd.Details, &emailPassword)
	if err != nil {
		return
	}

	ResponseCmd := models.NewPublicCommand()
	ResponseCmd.Command = "LoginResponse"
	userAccunt, res := checkLoginCredentials(&emailPassword)
	switch res {
	case models.EmailPasswordFit:
		authenticateUser(userAccunt)
		ResponseCmd.Parameters["Result"] = "Success"
		listOfCharacters(userAccunt)
		a.user.Account = userAccunt
	case models.EmailFit:
		ResponseCmd.Parameters["Result"] = "Failed"
	case models.EmailPasswordNotFit:
		userAccunt.Email = emailPassword.Email
		userAccunt.Password = emailPassword.Password
		userAccunt.SessionID = helpers.GenerateRandomString(64, false)
		if err := createAccount(userAccunt); err != nil {
			ResponseCmd.Parameters["Result"] = "Failed"
		} else {
			ResponseCmd.Parameters["Result"] = "Success"
			listOfCharacters(userAccunt)
			a.user.Account = userAccunt
		}
	}
	if userAccunt != nil {
		ResponseCmd.Details, _ = helpers.Serialize(userAccunt)
	}
	a.user.SendCommand(ResponseCmd, true)
}
