package database

import (
	"fmt"
	"log"
	"os"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/lib/pq"
)

type ()

type DBStruct struct {
	db *gorm.DB
}

type InTransaction func(tx *gorm.DB) error

var DBInstance *gorm.DB

func NewDB() (error, bool) {
	host := os.Getenv("PG_HOST")
	port := os.Getenv("PG_PORT")
	user := os.Getenv("PG_USER")
	password := os.Getenv("PG_PASSWD")
	dbname := os.Getenv("PG_DB")
	ssl := "disable"

	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s "+
		"dbname=%s sslmode=%s password=%s",
		host, port, user, dbname, ssl, password)

	log.Println(psqlInfo)
	db, err := gorm.Open("postgres", psqlInfo)
	if err != nil {
		log.Println("ERR_GRM_DB_INSTANCE: ", err)
		return err, false
	}

	DBInstance = db

	db.LogMode(true)

	//
	db.BlockGlobalUpdate(true)

	//do migration
	db.AutoMigrate(

		&NotifyMethod{},
		&AuditLog{},
		&User{},
		&Role{},
		&APIToken{},
	)

	//assoc

	//staff
	db.Model(&User{}).AddForeignKey("role_id", "roles(role_id)", "CASCADE", "RESTRICT")
	// db.Model(&Employee{}).AddForeignKey("status_id", "status_types(status_id)", "CASCADE", "RESTRICT")
	// db.Model(&Employee{}).AddForeignKey("role_id", "emp_roles(role_id)", "CASCADE", "RESTRICT")
	// db.Model(&Employee{}).AddForeignKey("pos_id", "position_roles(pos_id)", "CASCADE", "RESTRICT")
	//next of kin
	// db.Model(&NextKin{}).AddForeignKey("emp_id", "employees(emp_id)", "CASCADE", "RESTRICT")
	//notify method
	// db.Model(&NotifyMethod{}).AddForeignKey("emp_id", "employees(emp_id)", "CASCADE", "RESTRICT")

	return nil, true
}

func GetInstance() *gorm.DB {
	return DBInstance
}

func DoTransaction(fn InTransaction) error {

	tx := DBInstance.Begin()
	if tx.Error != nil {
		log.Println(tx.Error)
	}

	err := fn(tx)

	if err != nil {
		xerr := tx.Rollback().Error
		if xerr != nil {
			return xerr
		}
		return err
	}

	if err = tx.Commit().Error; err != nil {
		return err
	}

	return nil

}
