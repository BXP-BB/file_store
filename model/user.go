package model

import (
	"file_store/model/mysql"
	"time"
)

type User struct {
	Id           int
	OpenId       string
	FileStoreId  int
	UserName     string
	RegisterTime time.Time
	ImagePath    string
}

func CreateUser(openId, username, image string) {
	user := User{
		OpenId:       openId,
		FileStoreId:  0,
		UserName:     username,
		RegisterTime: time.Now(),
		ImagePath:    image,
	}
	mysql.DB.Create(&user)

	fileStore := FileStore{
		UserId:      user.Id,
		CurrentSize: 0,
		MaxSize:     1048576,
	}
	mysql.DB.Create(&fileStore)

	user.FileStoreId = fileStore.Id
	mysql.DB.Save(&user)
}

func QueryUserExists(OpenId string) bool {
	var user User
	mysql.DB.Find(&user, "open_id = ?", OpenId)

	return user.Id != 0
}

func GetUserInfo(openId interface{}) (user User) {
	mysql.DB.Find(&user, "open_id = ?", openId)
	return
}
