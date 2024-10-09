package controller

import (
	"file_store/lib"
	"file_store/model"
	"file_store/util"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func Upload(c *gin.Context) {
	openId, _ := c.Get("openId")
	fId := c.DefaultQuery("fId", "0")
	user := model.GetUserInfo(openId)
	currentFolder := model.GetCurrentFolder(fId)
	fileFolders := model.GetFileFolder(fId, user.FileStoreId)
	parentFolder := model.GetParentFolder(fId)
	currentAllParent := model.GetCurrentAllParent(parentFolder, make([]model.FileFolder, 0))
	fileDetailUse := model.GetFileDetailUse(user.FileStoreId)

	c.HTML(http.StatusOK, "upload.html", gin.H{
		"user":             user,
		"currUpload":       "active",
		"fId":              currentFolder.Id,
		"fName":            currentFolder.FileFolderName,
		"fileFolders":      fileFolders,
		"parentFolder":     parentFolder,
		"currentAllParent": currentAllParent,
		"fileDetailUse":    fileDetailUse,
	})
}

func HandlerUpload(c *gin.Context) {
	openId, _ := c.Get("openId")
	user := model.GetUserInfo(openId)

	Fid := c.GetHeader("id")
	// conf := lib.LoadServerConfig()
	conf := lib.Config
	file, head, err := c.Request.FormFile("file")
	defer file.Close()

	if ok := model.CurrFileExists(Fid, head.Filename); !ok {
		c.JSON(http.StatusOK, gin.H{
			"code": 501,
		})
	}
	if ok := model.CapacityIsEnough(head.Size, user.FileStoreId); !ok {
		c.JSON(http.StatusOK, gin.H{
			"code": 503,
		})
	}
	if err != nil {
		fmt.Println("file upload err ,", err.Error())
		return
	}

	fileFolder := model.GetCurrentFolder(Fid)
	var filePath = "/"
	var location string
	if Fid == "0" {
		location = conf.UploadLocation + "/" + head.Filename
	} else {
		filePath = fileFolder.FolderPath + "/" + fileFolder.FileFolderName
		location = conf.UploadLocation + filePath + "/" + head.Filename
	}

	newFile, err := os.Create(location)
	if err != nil {
		fmt.Println("file create failed : ", err.Error())
		return
	}
	defer newFile.Close()

	fileSize, err := io.Copy(newFile, file)
	if err != nil {
		fmt.Println("file copy failed : ", err.Error())
		return
	}

	_, _ = newFile.Seek(0, 0)
	fileHash := util.GetSHA256HashCode(newFile)

	model.CreateFile(head.Filename, fileHash, fileSize, Fid, user.FileStoreId, filePath)
	model.SubtractSize(fileSize, user.FileStoreId)
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
	})

}
