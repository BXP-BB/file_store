package model

import (
	"file_store/model/mysql"
	"file_store/util"
	"strings"
	"time"
)

type Share struct {
	Id       int
	Code     string
	FileId   int
	Username string
	Hash     string
}

func CreateShare(code, username string, fId int) string {
	share := Share{
		Code:     strings.ToLower(code),
		FileId:   fId,
		Username: username,
		Hash:     util.EncodeMd5(code + string(time.Now().Unix())),
	}
	mysql.DB.Create(&share)

	return share.Hash
}

func GetShareInfo(f string) (share Share) {
	mysql.DB.Find(&share, "hash = ?", f)
	return
}

func VerifyShareCode(fId, code string) bool {
	var share Share
	mysql.DB.Find(&share, "file_id = ? and code = ?", fId, code)

	return share.Id != 0
}
