package server

import (
	"errors"
	"fmt"

	"github.com/asaskevich/govalidator"
	"github.com/lordnorthern/login_server/helpers"
	"github.com/lordnorthern/login_server/models"
)

func addNewCharacter(character *models.Character, account *models.Account) error {
	var rowCount int
	var err error
	MySQL.DB.Where(&models.Character{
		Name: character.Name,
	}).First(&models.Character{}).Count(&rowCount)
	if rowCount > 0 {
		err = errors.New("NameTaken")
	} else {
		character.AccountID = account.ID
		if err := MySQL.DB.Save(character).Error; err != nil {
			fmt.Println("The error", err)
			err = errors.New("OtherError")
		}
	}
	MySQL.DB.Find(account).Related(&account.Characters)
	return err
}

func createAccount(account *models.Account) error {
	var rowCount int
	var err error
	if !govalidator.IsEmail(account.Email) {
		return errors.New("InvalidEmail")
	}
	MySQL.DB.Where(&models.Account{
		Email: account.Email,
	}).First(&models.Account{}).Count(&rowCount)
	if rowCount > 0 {
		err = errors.New("EmailTaken")
	} else {
		if err := MySQL.DB.Save(account).Error; err != nil {
			fmt.Println("The error", err)
			err = errors.New("OtherError")
		}
	}
	MySQL.DB.Find(account).Related(&account.Characters)
	return err
}

func checkLoginCredentials(creds *models.EmailPassword) (*models.Account, models.AccResult) {
	creds.HashPassword()
	userAccunt := models.Account{}
	chkAcc := models.Account{
		Email: creds.Email,
	}

	MySQL.DB.Where(&chkAcc).Find(&userAccunt)
	if userAccunt.ID != 0 && userAccunt.Password == creds.Password {
		return &userAccunt, models.EmailPasswordFit
	} else if userAccunt.ID != 0 && userAccunt.Password != creds.Password {
		return nil, models.EmailFit
	}
	return &userAccunt, models.EmailPasswordNotFit
}

func authenticateUser(userAccunt *models.Account) {
	userAccunt.SessionID = helpers.GenerateRandomString(64, false)
	MySQL.DB.Save(&userAccunt)
}

func listOfCharacters(userAccunt *models.Account) bool {
	if userAccunt.ID == 0 {
		return false
	}
	MySQL.DB.Find(&userAccunt).Related(&userAccunt.Characters)
	return true
}

func findCharByUserSession(characterID uint, sessionID string) *models.Account {
	var findAccount *models.Account
	MySQL.DB.First(findAccount)
	return findAccount
}
