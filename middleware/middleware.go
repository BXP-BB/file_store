package middleware

import (
	"file_store/lib"
	"file_store/model"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CheckLogin(c *gin.Context) {
	fmt.Println("*******************middleware**********************")
	token, err := c.Cookie("Token")
	if err != nil {
		fmt.Println("cookie ->", err.Error())
		c.Redirect(http.StatusFound, "/")
		c.Abort()
	}

	openId, err := lib.GetKey(token)
	if err != nil {
		fmt.Println(" get redis err:", err.Error())
		c.Redirect(http.StatusFound, "/")
		c.Abort()
	}

	user := model.GetUserInfo(openId)

	if user.Id == 0 {
		c.Redirect(http.StatusFound, "/")
	} else {
		c.Set("openId", openId)
		c.Next()
	}
}
