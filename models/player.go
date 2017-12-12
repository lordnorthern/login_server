package models

import (
	"errors"
)

// Account will operate against the player table
type Account struct {
	ID         uint        `gorm:"primary_key"`
	Email      string      `gorm:"size:64;unique_index"`
	Password   string      `gorm:"size:128"`
	Level      int         `gorm:"default:1"`
	Characters []Character `gorm:"ForeignKey:AccountID"`
	SessionID  string      `gorm:"size:64;unique_index"`
	BaseModel
}

// Character table will hold the characters data
type Character struct {
	ID        uint   `gorm:"primary_key"`
	AccountID uint   `gorm:"index"`
	Nickname  string `gorm:"size:64"`
	Name      string `gorm:"size:64;unique_index"`
	Gender    uint
	BaseModel
}

// Validate will validate that all the information is proper
func (c *Character) Validate() map[string]error {
	errs := make(map[string]error)
	if len(c.Nickname) < 3 || len(c.Nickname) > 16 {
		errs["Nickname"] = errors.New("Character nickname must be between 3 and 16 characters long")
	}
	if len(c.Name) < 5 || len(c.Name) > 32 {
		errs["Name"] = errors.New("Character name must be between 5 and 32 characters long")
	}
	if c.Gender != 1 && c.Gender != 2 {
		errs["Gender"] = errors.New("Please, select gender")
	}

	return errs
}
