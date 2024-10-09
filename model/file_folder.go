package model

import (
	"file_store/model/mysql"
	"fmt"
	"strconv"
	"time"
)

type FileFolder struct {
	Id             int
	FileFolderName string
	ParentFolderId int
	FileStoreId    int
	FolderPath     string
	Time           string
}

func CreateFolder(folderName, parentId string, fileStoreId int) {
	parentIdInt, err := strconv.Atoi(parentId)
	if err != nil {
		fmt.Println("parent id is error")
		return
	}

	var fileFolder FileFolder

	if parentIdInt == 0 {
		fileFolder = FileFolder{
			FileFolderName: folderName,
			ParentFolderId: parentIdInt,
			FileStoreId:    fileStoreId,
			Time:           time.Now().Format("2006-01-02 15:04:05"),
		}
	} else {
		parentFolder := GetParentFolder(parentId)
		fileFolder = FileFolder{
			FileFolderName: folderName,
			ParentFolderId: parentIdInt,
			FileStoreId:    fileStoreId,
			FolderPath:     parentFolder.FolderPath + "/" + parentFolder.FileFolderName,
			Time:           time.Now().Format("2006-01-02 15:04:05"),
		}
	}

	mysql.DB.Create(&fileFolder)
}

func GetParentFolder(fId string) (fileFolder FileFolder) {
	mysql.DB.Find(&fileFolder, "id = ?", fId)
	return
}

func GetFileFolder(parentId string, fileStoreId int) (fileFolders []FileFolder) {
	mysql.DB.Order("time desc").Find(&fileFolders, "parent_folder_id = ? and file_store_id = ?", parentId, fileStoreId)
	return
}

func GetCurrentFolder(fId string) (fileFolder FileFolder) {
	mysql.DB.Find(&fileFolder, "id = ?", fId)
	return
}

func GetCurrentAllParent(folder FileFolder, folders []FileFolder) []FileFolder {
	var parentFolder FileFolder
	if folder.ParentFolderId != 0 {
		mysql.DB.Find(&parentFolder, "id = ?", folder.ParentFolderId)
		folders = append(folders, parentFolder)
		return GetCurrentAllParent(parentFolder, folders)
	}
	for i, j := 0, len(folders)-1; i < j; i, j = i+1, j-1 {
		folders[i], folders[j] = folders[j], folders[i]
	}
	return folders
}

func GetUserFileFolderCount(fileStoreId int) (fileFolderCount int) {
	var fileFolder []FileFolder
	// mysql.DB.Find(&fileFolder, "file_store_id = ?", fileStoreId).Count(&fileFolderCount)
	mysql.DB.Where("file_store_id = ?", fileStoreId).Find(&fileFolder).Count(&fileFolderCount)
	return
}

func DeleteFileFolder(fId string) bool {
	var fileFolder FileFolder
	var fileFolder2 FileFolder
	mysql.DB.Where("id = ?", fId).Delete(FileFolder{})
	mysql.DB.Where("parent_folder_id = ?", fId).Delete(MyFile{})
	mysql.DB.Find(&fileFolder, "parent_folder_id = ?", fId)
	mysql.DB.Where("parent_folder_id = ?", fId).Delete(FileFolder{})

	mysql.DB.Find(&fileFolder2, "parent_folder_id = ?", fileFolder.Id)
	if fileFolder2.Id != 0 {
		return DeleteFileFolder(strconv.Itoa(fileFolder.Id))
	}

	return true
}

func UpdateFolderName(fId, fName string) {
	var fileFolder FileFolder
	mysql.DB.Model(&fileFolder).Where("id = ?", fId).Update("file_folder_name", fName)
}
