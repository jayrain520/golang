package main

import (
	"fmt"
	_ "fmt"
	"github.com/gin-gonic/gin"
	"github.com/wonderivan/logger"
	"html/template"
	"net/http"
	"testgin/api/config"
	"testgin/api/response"
	"time"
)

func main() {
	//gin.SetMode(gin.ReleaseMode)
	server := gin.Default()

	
	api := server.Group("/api")
	{

		//api.POST("/set/global/variable", SetGlobalVariable)
		api.POST("/add/comment", AddComment)
		api.POST("/login", login)
		api.POST("/register", register)
		//api.POST("/add/tag", addTag)
		//api.POST("/check/email", checkEmail)

		//post file to server local storage
		api.POST("/add/video", addVideoMiddleware, addVideo)

		//post file to aliyun oss storage serve
		//api.POST("/add/video/to/oss", addVideoMiddleware, AddVideoToOss)

		//delete aliyun oss storage file
		//api.DELETE("/delete/oss/file", DeleteFileOss)
		//delete server local video and  image file
		//api.DELETE("/delete/local/video/", DeleteLocalVideo)
		//delete comment
		//api.DELETE("/delete/comment", DeleteComment)

		//api.DELETE("/delete/user", DeleteUser)

		api.GET("/all/src/:src", Src)
		//api.GET("/load/tag", loadTag)

		//api.GET("/load/user/all/video/:user-name/:from/:to", loadUserAllVideos)
		api.GET("/load/all/video/:number", loadAllVideo)
		//api.GET("/check/session/:user-name", checkSession)

		//get video comment
		api.GET("/load/video/comment", LoadVideoComment)
		//api.GET("/load/self/comment", LoadSelfComment)

		//api.GET("/video/info", ElasticSearch)
	}



	server.Run(":" + config.Set.BindPort)
}

var once bool
var total = 0

func addVideoMiddleware(c *gin.Context) {
	if total == limitNumber {
		logger.Warn("Coming upload video limit: %v", total)
		responseError(c, response.Resp{
			ResponseCode: 400,
			ResponseMsg:  "Coming upload video limit: " + fmt.Sprintf("%d", total),
		})
		//reset total ,the is simple function
		//so I'm  don't use lock
		if !once {
			once = true
			go func(limit *int) {
				time.Sleep(time.Minute * resetLimit)
				*limit = 0
				once = false
			}(&total)
		}
		//stop next function
		c.Abort()
		return
	}
	total++
	logger.Debug(" %v Video Upload to Server", total)
}
