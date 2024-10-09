package model

import "file_store/model/mysql"

type FileStore struct {
	Id          int
	UserId      int
	CurrentSize int64
	MaxSize     int64
}

func GetUserFileStore(userId int) (fileStore FileStore) {
	mysql.DB.Find(&fileStore, "user_id = ?", userId)
	return
}

func CapacityIsEnough(fileSize int64, fileStoreId int) bool {
	var fileStore FileStore
	mysql.DB.First(&fileStore, fileStoreId)

	return fileStore.MaxSize-(fileSize/1024) >= 0
}
