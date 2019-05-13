package dao

import (
	"os"
	"strings"

	"github.com/jinzhu/gorm"
)

//OpenTestMysqlDatabaseConnection open mysql database connection
func OpenTestMysqlDatabaseConnection() (db *gorm.DB) {
	errMsgms := []string{}

	var username = os.Getenv("dbUsername")
	var password = os.Getenv("dbPassword")
	var dbHost = os.Getenv("dbHost")
	var dbSchema = os.Getenv("dbSchema")
	var dbPort = os.Getenv("dbPort")

	if len(username) == 0 {
		errMsgms = append(errMsgms, "key : dbUsername not set on env envirotnment")
	}
	if len(password) == 0 {
		errMsgms = append(errMsgms, "key : dbPassword not set on env envirotnment")
	}
	if len(dbHost) == 0 {
		errMsgms = append(errMsgms, "key : dbHost not set on env envirotnment")
	}
	if len(dbSchema) == 0 {
		errMsgms = append(errMsgms, "key : dbSchema not set on env envirotnment")
	}
	if len(dbPort) == 0 {
		errMsgms = append(errMsgms, "key : dbPort not set on env envirotnment")
	}
	if len(errMsgms) > 0 {
		panic("Database parameter not found : \n" + strings.Join(errMsgms, "\n"))
	}
	conQuery := username + ":" + password
	conQuery = conQuery + "@tcp(" + dbHost + ":" + dbPort + ")/" + dbSchema + "?charset=utf8&parseTime=True&loc=Local"
	db, err := gorm.Open("mysql", conQuery)
	if err != nil {
		panic(err.Error())
	}
	return db
}
