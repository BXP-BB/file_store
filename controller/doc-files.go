package controller

import (
	"file_store/model"
	"net/http"

	"github.com/gin-gonic/gin"
)

func DocFiles(c *gin.Context) {
	openId, _ := c.Get("openId")
	user := model.GetUserInfo(openId)

	fileDetailUse := model.GetFileDetailUse(user.FileStoreId)
	docFiles := model.GetTypeFile(1, user.FileStoreId)

	c.HTML(http.StatusOK, "doc-files.html", gin.H{
		"user":          user,
		"fileDetailUse": fileDetailUse,
		"docFiles":      docFiles,
		"docCount":      len(docFiles),
		"currDoc":       "active",
		"currClass":     "active",
	})
}
