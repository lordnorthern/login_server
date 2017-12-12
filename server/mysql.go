package server

import (
	"github.com/jinzhu/gorm"
	"github.com/lordnorthern/login_server/models"
)

// MySQLInstance blah tzif
type MySQLInstance struct {
	DB         *gorm.DB
	Definition string
}

var MySQL *MySQLInstance

// InitializeMySQL will initialize mysql connection
func InitializeMySQL(definition string) (*MySQLInstance, error) {
	newMySQLObject := new(MySQLInstance)
	dbCon, err := gorm.Open("mysql", models.Conf.MySQL.DbUsername+":"+models.Conf.MySQL.DbPassword+"@tcp("+models.Conf.MySQL.DbHost+":"+models.Conf.MySQL.DbPort+")/"+models.Conf.MySQL.DbName+"?charset=utf8&parseTime=true")
	if err != nil {
		return nil, err
	}
	newMySQLObject.DB = dbCon
	newMySQLObject.Definition = definition
	newMySQLObject.AddToList(&Terminatables)
	models.InitializeModels(newMySQLObject.DB)
	return newMySQLObject, nil
}

func (s *MySQLInstance) Terminate() {
	(*s.DB).Close()
}
func (s *MySQLInstance) GetName() string {
	return s.Definition
}

func (s *MySQLInstance) AddToList(list *[]models.Terminator) {
	*list = append(*list, s)
}
