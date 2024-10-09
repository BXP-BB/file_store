package controller

import (
	"file_store/lib"
	"file_store/model"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
)

func Files(c *gin.Context) {
	openId, _ := c.Get("openId")
	fId := c.DefaultQuery("fId", "0")
	user := model.GetUserInfo(openId)

	files := model.GetUserFile(fId, user.FileStoreId)
	fileFolder := model.GetFileFolder(fId, user.FileStoreId)
	parentFolder := model.GetParentFolder(fId)
	currentAllParent := model.GetCurrentAllParent(parentFolder, make([]model.FileFolder, 0))
	currentFolder := model.GetCurrentFolder(fId)
	fileDetailUse := model.GetFileDetailUse(user.FileStoreId)

	c.HTML(http.StatusOK, "files.html", gin.H{
		"currAll":          "active",
		"user":             user,
		"fId":              currentFolder.Id,
		"fName":            currentFolder.FileFolderName,
		"files":            files,
		"fileFolder":       fileFolder,
		"parentFolder":     parentFolder,
		"currentAllParent": currentAllParent,
		"fileDetailUse":    fileDetailUse,
	})

}

func AddFolder(c *gin.Context) {
	openId, _ := c.Get("openId")
	user := model.GetUserInfo(openId)

	folderName := c.PostForm("fileFolderName")
	parentId := c.DefaultPostForm("parentFolderId", "0")
	// model.CreateFolder(folderName, parentId, user.FileStoreId)
	parent := model.GetParentFolder(parentId)

	// conf := lib.LoadServerConfig()
	conf := lib.Config
	folderPath := "/"
	var localFolder string
	if parentId == "0" {
		localFolder = conf.UploadLocation + "/" + folderName
	} else {
		folderPath = parent.FolderPath + "/" + parent.FileFolderName
		localFolder = conf.UploadLocation + folderPath + "/" + folderName
	}

	err := os.Mkdir(localFolder, 0755)
	if err != nil {
		fmt.Println("Error creating directory:", err)
	} else {
		model.CreateFolder(folderName, parentId, user.FileStoreId)
	}

	c.Redirect(http.StatusMovedPermanently, "/cloud/files?fId="+parentId+"&fName="+parent.FileFolderName)
}

func DownloadFile(c *gin.Context) {
	fId := c.Query("fId")
	file := model.GetFileInfo(fId)
	if file.FileHash == "" {
		return
	}

	// conf := lib.LoadServerConfig()
	conf := lib.Config

	filePath := conf.UploadLocation + file.FilePath + "/" + file.FileName + file.Postfix
	dFile, err := os.Open(filePath)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "文件未找到", "file": file.FileName})
		return
	}
	defer dFile.Close()

	// 设置响应的Header
	c.Writer.Header().Set("Content-Type", "application/octet-stream")
	c.Writer.Header().Set("Content-Disposition", "attachment; filename=\""+file.FileName+file.Postfix+"\"")

	// 将文件内容发送给客户端
	_, err = io.Copy(c.Writer, dFile)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "内部服务器错误"})
		return
	}

	model.DownloadNumAdd(fId)

	c.Status(http.StatusOK)

	// 返回文本内容
	// 待完善 fileData
	// c.Header("Content-disposition","attachment;filenam\""+file.FileName+file.Postfix+"\"")
	// c.Data(http.StatusOK,"application/octect-stream", fileData)
}

func DeleteFile(c *gin.Context) {
	openId, _ := c.Get("openId")
	user := model.GetUserInfo(openId)

	fId := c.DefaultQuery("fId", "")
	folderId := c.Query("folder")
	if fId == "" {
		return
	}

	// conf := lib.LoadServerConfig()
	conf := lib.Config
	file := model.GetFileInfo(fId)
	filePath := file.FilePath
	location := conf.UploadLocation + filePath + "/" + file.FileName + file.Postfix
	err := os.Remove(location)
	if err != nil {
		fmt.Printf("Error deleting file: %s\n", err)
	} else {
		model.DeleteUseFile(fId, folderId, user.FileStoreId)
	}

	c.Redirect(http.StatusMovedPermanently, "/cloud/files?fid="+folderId)
}

func DeleteFileFolder(c *gin.Context) {
	fId := c.DefaultQuery("fId", "")
	if fId == "" {
		return
	}

	folderInfo := model.GetCurrentFolder(fId)
	// conf := lib.LoadServerConfig()
	conf := lib.Config

	dirName := conf.UploadLocation + folderInfo.FolderPath + "/" + folderInfo.FileFolderName
	err := os.RemoveAll(dirName)
	if err != nil {
		fmt.Printf("Error deleting folder: %s\n", err)
	} else {
		model.DeleteFileFolder(fId)
	}

	c.Redirect(http.StatusMovedPermanently, "/cloud/files?fid="+strconv.Itoa(folderInfo.ParentFolderId))
}

func UpdateFileFolder(c *gin.Context) {
	fileFolderName := c.PostForm("fileFolderName")
	fileFolderId := c.PostForm("fileFolderId")
	fileFolder := model.GetCurrentFolder(fileFolderId)
	model.UpdateFolderName(fileFolderId, fileFolderName)
	c.Redirect(http.StatusMovedPermanently, "/cloud/files?fid="+strconv.Itoa(fileFolder.ParentFolderId))
}
