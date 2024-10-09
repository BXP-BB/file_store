package model

import (
	"file_store/model/mysql"
	"file_store/util"
	"path"
	"strconv"
	"strings"
	"time"
)

type MyFile struct {
	Id             int
	FileName       string
	FileHash       string
	FileStoreId    int
	FilePath       string
	DownloadNum    int
	UploadTime     string
	ParentFolderId int
	Size           int64
	SizeStr        string
	Type           int
	Postfix        string
}

func CreateFile(fileName, fileHash string, fileSize int64, fId string, fileStoreId int, filePath string) {
	var sizeStr string
	fileSuffix := path.Ext(fileName)
	filePreFix := fileName[0 : len(fileName)-len(fileSuffix)]
	fid, _ := strconv.Atoi(fId)

	if fileSize < 1048576 {
		sizeStr = strconv.FormatInt(fileSize/1024, 10) + "KB"
	} else {
		sizeStr = strconv.FormatInt(fileSize/102400, 10) + "MB"
	}

	myFile := MyFile{
		FileName:       filePreFix,
		FileHash:       fileHash,
		FileStoreId:    fileStoreId,
		FilePath:       filePath,
		DownloadNum:    0,
		UploadTime:     time.Now().Format("2006-01-02 15:04:05"),
		ParentFolderId: fid,
		Size:           fileSize / 1024,
		SizeStr:        sizeStr,
		Type:           util.GetFileTypeInt(fileSuffix),
		Postfix:        strings.ToLower(fileSuffix),
	}
	mysql.DB.Create(&myFile)
}

func GetUserFile(parentId string, storeId int) (files []MyFile) {
	mysql.DB.Find(&files, "file_store_id = ? and parent_folder_id = ?", storeId, parentId)
	return
}

func SubtractSize(size int64, fileStoreId int) {
	var fileStore FileStore
	mysql.DB.First(&fileStore, fileStoreId)
	fileStore.CurrentSize = fileStore.CurrentSize + size/1024
	fileStore.MaxSize = fileStore.MaxSize - size/1024
	mysql.DB.Save(&fileStore)
}

func GetUserFileCount(fileStoreId int) (fileCount int) {
	var file []MyFile
	// mysql.DB.Find(&file, "file_store_id = ?", fileStoreId).Count(&fileCount)
	mysql.DB.Where("file_store_id = ?", fileStoreId).Find(&file).Count(&fileCount)
	return
}

func GetFileDetailUse(fileStoreId int) map[string]int64 {
	var files []MyFile
	var (
		docCount   int64
		imgCount   int64
		videoCount int64
		musicCount int64
		otherCount int64
	)

	fileDetailUseMap := make(map[string]int64, 0)

	docCount = mysql.DB.Find(&files, "file_store_id = ? AND type = ?", fileStoreId, 1).RowsAffected
	fileDetailUseMap["docCount"] = docCount

	imgCount = mysql.DB.Find(&files, "file_store_id = ? AND type = ?", fileStoreId, 2).RowsAffected
	fileDetailUseMap["imgCount"] = imgCount

	videoCount = mysql.DB.Find(&files, "file_store_id = ? AND type = ?", fileStoreId, 3).RowsAffected
	fileDetailUseMap["videoCount"] = videoCount

	musicCount = mysql.DB.Find(&files, "file_store_id = ? AND type = ?", fileStoreId, 4).RowsAffected
	fileDetailUseMap["musicCount"] = musicCount

	otherCount = mysql.DB.Find(&files, "file_store_id = ? AND type = ?", fileStoreId, 5).RowsAffected
	fileDetailUseMap["otherCount"] = otherCount

	return fileDetailUseMap

}

func GetTypeFile(fileType, fileStoreId int) (files []MyFile) {
	mysql.DB.Find(&files, "file_store_id = ? and type = ?", fileStoreId, fileType)
	return
}

func CurrFileExists(fId, fileName string) bool {
	var file MyFile
	fileSuffix := strings.ToLower(path.Ext(fileName))
	filePrefix := fileName[0 : len(fileName)-len(fileSuffix)]
	mysql.DB.Find(&file, "parent_folder_id = ? and file_name = ? and postfix = ?", fId, filePrefix, fileSuffix)

	return file.Size <= 0
}

func FileHashExists(fileHash string) bool {
	var file MyFile
	mysql.DB.Find(&file, "file_hash = ?", fileHash)

	return file.FileHash == ""
}

func GetFileInfo(fId string) (file MyFile) {
	mysql.DB.First(&file, fId)
	return
}

func DownloadNumAdd(fId string) {
	var file MyFile
	mysql.DB.First(&file, fId)
	file.DownloadNum = file.DownloadNum + 1
	mysql.DB.Save(&file)
}

func DeleteUseFile(fId, folderId string, storeId int) {
	mysql.DB.Where("id = ? and file_store_id = ? and parent_folder_id = ?", fId, storeId, folderId).Delete(MyFile{})
}
