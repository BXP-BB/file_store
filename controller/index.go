package controller

import (
	"file_store/model"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Index(c *gin.Context) {
	openid, _ := c.Get("openId")
	user := model.GetUserInfo(openid)
	userFileStore := model.GetUserFileStore(user.Id)

	fileCount := model.GetUserFileCount(user.FileStoreId)
	fileFolderCount := model.GetUserFileFolderCount(user.FileStoreId)
	fileDetailUse := model.GetFileDetailUse(user.FileStoreId)

	c.HTML(http.StatusOK, "index.html", gin.H{
		"user":            user,
		"currIndex":       "active",
		"userFileStore":   userFileStore,
		"fileCount":       fileCount,
		"fileFolderCount": fileFolderCount,
		"fileDetailUse":   fileDetailUse,
	})
}
