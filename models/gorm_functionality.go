package models

import (
	"time"

	"github.com/jinzhu/gorm"
)

// DB connection handle
var db *gorm.DB

//InitializeModels will initialize the models and also the db connectivity
func InitializeModels(inDB *gorm.DB) {
	db = inDB
	//db.DropTableIfExists(&Account{}, &Character{})
	db.AutoMigrate(&Account{}, &Character{})
}

// BaseModel contains the basic stuff from gorm.
// I prefer usign this over the default gorm.Model because I prefer my ID to be at the top
// and the recording data at the very bottom
type BaseModel struct {
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `sql:"index"`
}

// AccResult is a return type for username/password lookups
type AccResult int

const (
	// EmailPasswordFit indicates that the combination was found.
	EmailPasswordFit AccResult = iota
	// EmailFit indicates that the email was found, but the password was wrong.
	EmailFit
	// EmailPasswordNotFit indicates that neither email nor password (which means that email is available for auto register)
	EmailPasswordNotFit
)
