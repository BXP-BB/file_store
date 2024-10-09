package controller

import (
	"file_store/lib"
	"file_store/model"
	"file_store/util"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

type PrivateInfo struct {
	AccessToken  string `json:"access_token"`
	ExpriesIn    string `json:"expries_in"`
	RefreshToken string `json:"refresh_token"`
	OpenId       string `json:"open_id"`
}

type LoginResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
	Url     string `json:"url"`
}

type Userinfo struct {
	User     string `form:"user" json:"user"`
	Password string `form:"password" json:"password"`
}

func Login(c *gin.Context) {
	fmt.Println("******************Login******************")
	c.HTML(http.StatusOK, "login.html", nil)
}

func HandleLogin(c *gin.Context) {
	fmt.Println("******************HandleLogin******************")
	var user Userinfo
	if err := c.ShouldBind(&user); err != nil {
		// 如果绑定失败（例如，字段名不匹配或格式错误），则返回错误
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	} else {
		openId := util.EncodeMd5(user.Password + "BXP")

		if ok := model.QueryUserExists(openId); !ok {
			c.JSON(http.StatusOK, LoginResponse{Success: false, Message: "登录失败", Url: "/"})
		} else {
			hashToken := util.EncodeMd5("token" + string(time.Now().Unix()) + openId)
			if err := lib.SetKey(hashToken, openId, 24*3600); err != nil {
				fmt.Println("Redis set err:", err.Error())
				return
			}

			c.SetCookie("Token", hashToken, 3600, "/", "10.150.253.121", false, true)

			conf := lib.LoadServerConfig()
			rootFolder := conf.UploadLocation + "/" + user.User
			lib.Config.UploadLocation = rootFolder

			c.JSON(http.StatusOK, LoginResponse{Success: true, Message: "登录成功", Url: "/cloud/index"})
		}

		// c.Redirect(http.StatusMovedPermanently, "/cloud/index")
		// c.String(http.StatusOK, "OK")
		// if loginForm.User == "admin" && loginForm.Password == "123456" {
		// 登录成功
		// c.JSON(http.StatusOK, LoginResponse{Success: true, Message: "登录成功"})
		// 如果需要，可以在这里添加更多的逻辑，比如设置 Cookie、生成 JWT 令牌等
		// } else {
		// 登录失败
		// c.JSON(http.StatusUnauthorized, LoginResponse{Success: false, Message: "用户名或密码错误"})
		// }
	}
}

func Register(c *gin.Context) {
	fmt.Println("******************Register******************")
	var user Userinfo
	if err := c.ShouldBind(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	} else {
		openId := util.EncodeMd5(user.Password + "BXP")
		if ok := model.QueryUserExists(openId); !ok {
			model.CreateUser(openId, user.User, "/static/img/login.png")

			conf := lib.LoadServerConfig()
			rootFolder := conf.UploadLocation + "/" + user.User
			err := os.Mkdir(rootFolder, 0755)
			if err != nil {
				fmt.Println("Error creating directory:", err)
			}

			c.JSON(http.StatusOK, LoginResponse{Success: true, Message: "注册成功", Url: "/"})
		} else {
			c.JSON(http.StatusOK, LoginResponse{Success: false, Message: "用户已存在", Url: "/"})
		}
	}
}

func Logout(c *gin.Context) {
	token, err := c.Cookie("Token")
	if err != nil {
		fmt.Println("cookie ", err.Error())
	}
	if err := lib.DelKey(token); err != nil {
		fmt.Println("Del redis err: ", err.Error())
	}

	c.SetCookie("Token", "", 0, "/", "10.150.253.121", false, false)
	c.Redirect(http.StatusFound, "/")
}
