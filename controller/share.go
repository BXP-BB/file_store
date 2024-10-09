package controller

import (
	"file_store/lib"
	"file_store/model"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/lifei6671/gocaptcha"
)

func ShareFile(c *gin.Context) {
	openId, _ := c.Get("openId")
	user := model.GetUserInfo(openId)

	fId := c.Query("id")
	url := c.Query("url")
	code := gocaptcha.RandText(4)
	fileId, _ := strconv.Atoi(fId)
	hash := model.CreateShare(code, user.UserName, fileId)

	c.JSON(http.StatusOK, gin.H{
		"url":  url + "?f=" + hash,
		"code": code,
	})
}

func SharePass(c *gin.Context) {
	f := c.Query("f")

	shareInfo := model.GetShareInfo(f)
	file := model.GetFileInfo(strconv.Itoa(shareInfo.FileId))

	c.HTML(http.StatusOK, "share.html", gin.H{
		"id":       shareInfo.FileId,
		"username": shareInfo.Username,
		"fileType": file.Type,
		"filename": file.FileName + file.Postfix,
		"hash":     shareInfo.Hash,
	})
}

func DownloadShareFile(c *gin.Context) {
	fileId := c.Query("id")
	code := c.Query("code")
	hash := c.Query("hash")
	username := c.Query("username")

	fileInfo := model.GetFileInfo(fileId)
	if ok := model.VerifyShareCode(fileId, strings.ToLower(code)); !ok {
		c.Redirect(http.StatusMovedPermanently, "/file/share?f="+hash)
		return
	}

	conf := lib.LoadServerConfig()
	// conf := lib.Config

	filePath := conf.UploadLocation + "/" + username + "/" + fileInfo.FilePath + "/" + fileInfo.FileName + fileInfo.Postfix
	dFile, err := os.Open(filePath)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "文件未找到", "file": fileInfo.FileName})
		return
	}
	defer dFile.Close()

	// 设置响应的Header
	c.Writer.Header().Set("Content-Type", "application/octet-stream")
	c.Writer.Header().Set("Content-Disposition", "attachment; filename=\""+fileInfo.FileName+fileInfo.Postfix+"\"")

	// 将文件内容发送给客户端
	_, err = io.Copy(c.Writer, dFile)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "内部服务器错误"})
		return
	}

	model.DownloadNumAdd(fileId)

	c.Status(http.StatusOK)
}
