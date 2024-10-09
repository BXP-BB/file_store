package router

import (
	"file_store/controller"
	"file_store/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRoute() *gin.Engine {
	router := gin.Default()
	router.GET("/", controller.Login)
	router.POST("/login", controller.HandleLogin)
	router.POST("/register", controller.Register)

	file := router.Group("file")
	{
		file.GET("/share", controller.SharePass)
		file.GET("/shareDownload", controller.DownloadShareFile)
	}

	cloud := router.Group("cloud")
	cloud.Use(middleware.CheckLogin)
	{
		cloud.GET("/index", controller.Index)
		cloud.GET("/logout", controller.Logout)

		cloud.GET("/files", controller.Files)
		cloud.GET("/upload", controller.Upload)
		cloud.GET("/downloadFile", controller.DownloadFile)
		cloud.GET("/deleteFile", controller.DeleteFile)
		cloud.GET("/deleteFolder", controller.DeleteFileFolder)

		cloud.GET("/doc-files", controller.DocFiles)
		cloud.GET("/image-files", controller.ImageFiles)
		cloud.GET("/video-files", controller.VideoFiles)
		cloud.GET("/music-files", controller.MusicFiles)
		cloud.GET("/other-files", controller.OtherFiles)

		cloud.GET("/help", controller.Help)
	}

	{
		cloud.POST("/uploadFile", controller.HandlerUpload)
		cloud.POST("/addFolder", controller.AddFolder)
		cloud.POST("/updateFolder", controller.UpdateFileFolder)
		cloud.POST("/getQrCode", controller.ShareFile)
	}
	return router
}
