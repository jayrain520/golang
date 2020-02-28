Go Streaming media server 
=====

API interface: 项目的简单展示地址：http://jayrain.cn/static/video.html
-----
    func main() {
	//gin.SetMode(gin.ReleaseMode)
	server := gin.Default() //startup debug
	api := server.Group("/api")
	{

		api.POST("/set/global/variable", SetGlobalVariable)
		api.POST("/add/comment", AddComment)
		api.POST("/login", login)
		api.POST("/register", register)
		api.POST("/add/tag", addTag)
		api.POST("/check/email",checkEmail)


		//post file to server local storage
		api.POST("/add/video", addVideoMiddleware, addVideo)

		//post file to aliyun oss storage serve
		api.POST("/add/video/to/oss", addVideoMiddleware, AddVideoToOss)

		//delete aliyun oss storage file
		api.DELETE("/delete/oss/file", DeleteFileOss)
    
		//delete server local video and  image file
		api.DELETE("/delete/local/video/", DeleteLocalVideo)
    
		//delete comment
		api.DELETE("/delete/comment", DeleteComment)
    
		api.DELETE("/delete/user", DeleteUser)


		api.GET("/all/src/:src", Src)
		api.GET("/load/tag", loadTag)

		api.GET("/load/user/all/video/:user-name/:from/:to", loadUserAllVideos)
		api.GET("/load/all/video/:number", loadAllVideo)
		api.GET("/check/session/:user-name", checkSession)

		//get video comment
		api.GET("/load/video/comment", LoadVideoComment)
		api.GET("/load/self/comment", LoadSelfComment)

		api.GET("/video/info", ElasticSearch)

	}
	server.Run(":" + config.Set.BindPort)
    }
    
    
Streaming media server json config 
----

日志的等级
![log_level](https://github.com/jayrain520/golang/blob/master/image/logger.jpg)


     {
  
    "logger": {
       "level": "DEBG"                      //日志的等级
       },
     
     "mysql": {
      "db_connect_address":"127.0.0.1",     //链接地址
     "set_max_open_conns": 10,              //最大连接
     "set_max_idle_conns": 20,              //最大空闲数
      "db_name":"video_server",
      "user_name":"root",
      "port":"3306",
      "password":"",
      "charset":"utf-8"
      },
      
     "local_storage_path": {
     "local_video_path": "./videos/"     //本地视频文件的储存位置
     },

      "server_address": {
        "src_bind": "127.0.0.1:8081",   //分开就写公网地址和端口，域名也行
    "server_run_port": "8081"           //服务运行的端口
     },
      
     "aliyun_oss_storage":{             //阿里云的oss,存储服务
        "open_oss_serve": false,        //是否启用oss
        "oss_videos_path": "videos/",   //bucket里的路径
        "bucket_name": "",
        "resource_address": "https://jay-video.oss-cn-hongkong.aliyuncs.com/videos/",
        "end_point": "",                //oss储存的的地址
        "access_key_id":"",             //oss储存的授权码
        "access_key_secret": ""         //oss储存的授权码
     },
  
    "elastic_search": {
    "open_elastic": false,              //是否启用elasticsearch
    "init_data_number": 10,             //把数据库里的数据缓存到elasticsearch
    "set_sniff": false,
    "address": "http://127.0.0.1:9200",
    "index": "video",
    "type": "info"
     },

	//使用邮箱发送验证码，需要配置自己的QQ邮箱
     "QQ_Email": {
    "email_title": "Your Verification Code:",
    "email_user": "",                  //邮箱的账户
    "email_accesskey": "",             //邮箱的授权码     
    "email_host": "smtp.qq.com",
    "email_port": 25
      }
    }






Project tree iamge
-----
![流程图](https://github.com/jayrain520/golang/blob/master/video_serve/project.png)
