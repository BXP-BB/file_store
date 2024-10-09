package mysql

import (
	"file_store/lib"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

var DB *gorm.DB

func InitDB(conf lib.ServerConfig) {
	var err error
	dbParams := fmt.Sprintf("%v:%v@tcp(%v)/%v?charset=utf8&parseTime=True&loc=Local",
		conf.User,
		conf.Password,
		conf.Host,
		conf.DBName,
	)

	fmt.Println("databse dbParams : ", dbParams)

	DB, err = gorm.Open("mysql", dbParams)
	if err != nil {
		log.Fatal(2, err)
	}

	DB.SingularTable(true)

	gorm.DefaultTableNameHandler = func(db *gorm.DB, defaultTableName string) string {
		return conf.TablePrefix + defaultTableName
	}

	DB.DB().SetMaxIdleConns(10)
	DB.DB().SetMaxOpenConns(100)
	fmt.Println("databse init on port ", conf.Host)

	// DB.Exec("INSERT INTO file_folder (id, file_folder_name, parent_folder_id, file_store_id, time) VALUES (?, ?, ?, ?, ?)", 1, conf.UploadLocation, 0, 0, time.Now().Format("2024-08-24 15:47:00"))
}
