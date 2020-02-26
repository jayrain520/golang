package main

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/wonderivan/logger"
	"net/http"
	"os"
	"regexp"
	"runtime"
	"strconv"
	"testgin/api/config"
	"testgin/api/database"
	"testgin/api/models"
	"testgin/api/response"
	"time"
)

//api

func SetGlobalVariable(c *gin.Context) {
	globalVar["GoroutineNumber"] = runtime.NumGoroutine()
	if vs, _ := strconv.Atoi(c.PostForm("video_size")); vs != 0 {
		vs *= 1024 * 1024
		videoSize = vs
		globalVar["LimitVideoSize"] = strconv.Itoa(videoSize/1024/1024) + "M"
	} else {
		globalVar["LimitVideoSize"] = strconv.Itoa(videoSize/1024/1024) + "M"
	}

	if is, _ := strconv.Atoi(c.PostForm("image_size")); is != 0 {
		is *= 1024 * 1024
		imageSize = is
		globalVar["limitImageSize"] = strconv.Itoa(imageSize/1024/1024) + "M"
	} else {
		globalVar["limitImageSize"] = strconv.Itoa(imageSize/1024/1024) + "M"
	}

	if ln, _ := strconv.Atoi(c.PostForm("limit_upload")); ln != 0 {
		limitNumber = ln
		globalVar["LimitUpload"] = strconv.Itoa(int(resetLimit)) + " Minute Upload " + strconv.Itoa(limitNumber) + " once"
	} else {
		globalVar["LimitUpload"] = strconv.Itoa(int(resetLimit)) + " Minute Upload " + strconv.Itoa(limitNumber) + " once"
	}

	if rl, _ := strconv.Atoi(c.PostForm("reset_limit")); rl != 0 {

		resetLimit = time.Duration(int64(rl))
		globalVar["ResetTime"] = strconv.Itoa(int(resetLimit)) + " Minute"
	} else {
		globalVar["ResetTime"] = strconv.Itoa(int(resetLimit)) + " Minute"
	}

	if mpid, _ := strconv.Atoi(c.PostForm("max_parent_id")); mpid != 0 {
		if mpid > 20 {
			globalVar["MaxParentID"] = "The Parent ID Not More 20"
		} else {
			database.MaxParentID = mpid
			globalVar["MaxParentID"] = "1-" + strconv.Itoa(database.MaxParentID)
		}

	} else {
		globalVar["MaxParentID"] = "1-" + strconv.Itoa(database.MaxParentID)
	}

	if sbid, _ := strconv.Atoi(c.PostForm("max_sub_id")); sbid != 0 {
		if sbid <= 20 {
			globalVar["MaxSubID"] = "The Sub ID Not be smaller 21"
		} else {
			database.MaxSubID = sbid
			globalVar["MaxSubID"] = "21-" + strconv.Itoa(database.MaxSubID)
		}
	} else {
		globalVar["MaxSubID"] = "21-" + strconv.Itoa(database.MaxSubID)
	}
	c.JSON(200, gin.H{
		"GlobalVariable": globalVar,
	})

}

func checkSession(c *gin.Context) {
	userName := c.Param("user-name")
	if session, ok := isExists(userName); ok {
		c.JSON(200, gin.H{
			"normal": session,
		})
		return
	} else {
		responseError(c, response.ExpireSession)
		logger.Debug("Session Not Exists")
		return
	}
}

func login(c *gin.Context) {
	var newUser models.User
	err := c.ShouldBind(&newUser)
	if err != nil {
		logger.Error(err)
		responseError(c, response.IntervalErr)
		return
	}
	if newUser.Pwd, err = encryption(newUser.Pwd); err != nil {
		logger.Error(err)
		responseError(c, response.IntervalErr)
		return
	}

	//retrieve data in mysql,than check
	oldUser, err := conn.CheckPassword(newUser.Name)
	if err != nil {
		if err == sql.ErrNoRows {
			logger.Debug(err)
			responseError(c, response.UserNotFound)
			return
		}
		logger.Error(err)
		responseError(c, response.IntervalErr)
		return
	}

	//start check user name and password
	legal := make([]interface{}, 0)
	if newUser.Pwd == oldUser.Pwd && newUser.Name == oldUser.Name {

		if session, ok := isExists(newUser.Name); ok {
			legal = append(legal, oldUser)
			legal = append(legal, &session)
			responseNormal(c, legal)
			return
		}
		//if memory don't session id .then generate new session  to mysql
		sid := models.NewUUID()
		expire := Day14_STime()
		session := models.SessionID{
			ID:       sid,
			Expire:   expire,
			UserName: newUser.Name,
		}
		err = conn.AddNewSessionID(sid, expire, newUser.Name)
		if err != nil {
			logger.Alert(err)
			responseError(c, response.IntervalErr)
			return
		}

		legal = append(legal, oldUser)
		legal = append(legal, &session)

		allSession.Store(session.UserName, session)
		responseNormal(c, legal)
		return
	}

}

func register(c *gin.Context) {
	var newUser models.User
	err := c.ShouldBind(&newUser)
	if err != nil {
		logger.Error(err)
		responseError(c, response.IntervalErr)
		return
	}
	if newUser.Name == "" {
		logger.Error("User Name is empty")
		responseError(c, response.RequestInvalid)
		return
	}


	limitN := []rune(newUser.Name)
	if len(limitN) < 3 || len(limitN) > 12 {
		logger.Debug("User Name Not be Smaller 3 or More 12: %v", len(limitN))
		responseError(c, response.UserNameLimits)
		return
	}

	if newUser.Pwd, err = encryption(newUser.Pwd); err != nil {
		logger.Error(err)
		responseError(c, response.IntervalErr)
		return
	}

	err = conn.AddNewUser(newUser.Name, newUser.Pwd, newUser.Email)
	if err != nil {
		logger.Error(err)
		responseError(c, response.Resp{ResponseCode: 500, ResponseMsg: "User Name or Email Registered "})
		return
	} else {
		logger.Debug("one User Register Successfully")
		responseNormal(c, response.Success)
	}
}

func loadUserAllVideos(c *gin.Context) {
	userName := c.Param("user-name")
	from, err := strconv.Atoi(c.Param("from"))
	to, err := strconv.Atoi(c.Param("to"))
	if err != nil {
		logger.Debug("Parse to int Failed: %v", err)
		responseError(c, response.IntervalErr)
		return
	}

	if _, ok := isExists(userName); ok {
		listVideo, err := conn.LoadUserAllVideos(userName, from, to)
		if err != nil {
			if err == sql.ErrNoRows {
				logger.Debug(err)
				responseError(c, response.IntervalErr)
				return
			}
			logger.Alert(err)
			responseError(c, response.IntervalErr)
			return
		}
		responseNormal(c, listVideo)
		return
	}
	responseError(c, response.ExpireSession)

}

//load video from two weeks before to now
func loadAllVideo(c *gin.Context) {
	number, err := strconv.Atoi(c.Param("number"))
	if number == -1 {
		logger.Debug("The Number is empty")
		responseError(c, response.RequestInvalid)
		return
	}
	listVideo, err := conn.LoadAllVideos(number, time.Now().Unix()-1000000, time.Now().Unix())
	if err != nil {
		logger.Error(err)
		responseError(c, response.IntervalErr)
		return
	}
	responseNormal(c, listVideo)
}

//local video and image store
func addVideo(c *gin.Context) {
	authorNmae := c.PostForm("author_name")
	title := c.PostForm("title")
	parentID, err := strconv.Atoi(c.PostForm("parentid"))
	subID, err := strconv.Atoi(c.PostForm("subid"))
	if err != nil || parentID < 0 || parentID > database.MaxParentID || subID > database.MaxSubID || subID < 21 {
		responseError(c, response.Resp{
			ResponseCode: 400,
			ResponseMsg:  fmt.Sprintf("%s", database.InvalidID),
		})
		return
	}

	//session check
	if _, ok := isExists(authorNmae); !ok {
		responseError(c, response.ExpireSession)
		return
	}


	file, err := c.FormFile("file")
	if err != nil {
		logger.Error("Get Post File Failed: %v%v", err)
		responseError(c, response.IntervalErr)
		return
	}
	if file.Size > int64(videoSize) {
		logger.Debug("Video Size more: %d M", videoSize/1024/1024)
		responseError(c, response.Resp{
			ResponseCode: 400,
			ResponseMsg:  "Video size be smaller: " + strconv.Itoa(videoSize/1024/1024) + "M",
		})
		return
	}


	vid := models.NewUUID()
	err = c.SaveUploadedFile(file, config.Set.VideoSavePath+vid)
	if err!= nil {
		removeResource(vid)
		logger.Error(err)
		responseError(c, response.IntervalErr)
		return
	}

	//start transcoding to mp4
	err = Transcoding(vid)
	if err != nil {
		removeResource(vid)
		logger.Error(err)
		responseError(c, response.Resp{ResponseCode: 500, ResponseMsg: fmt.Sprintf("%s", err)})
		return
	}

	_ = os.Rename(config.Set.VideoSavePath+"output.mp4", config.Set.VideoSavePath+vid)
	video_link := "http://" + config.Set.Bind + "/api/all/src/" + vid
	err = conn.AddNewVideoInfo(video_link, authorNmae, title, parentID, subID)
	if err != nil {
		removeResource(vid)
		logger.Error(err)
		responseError(c, response.IntervalErr)
		return
	}

	//save to elastic search , if you want use elastic search please open and configure it
	if config.Set.OpenElastic {
		info, err := conn.CheckVideoInfo(video_link)
		if err != nil {
			logger.Error("One Video Save to Elastic Failed: %v", err)
		}
		_ = saveToElastic(info)
	}
	responseNormal(c, response.Success)

}

//aliyun oss storage serve
func AddVideoToOss(c *gin.Context) {
	if !config.Set.OpenOssServe || config.Set.AccessKeySecret == "" ||
		config.Set.AccessKeyID == "" || config.Set.EndPoint == "" {
		logger.Warn(response.NotOpenOss.ResponseMsg)
		responseError(c, response.NotOpenOss)
		return
	}

	authorNmae := c.PostForm("author_name")
	title := c.PostForm("title")
	parentID, err := strconv.Atoi(c.PostForm("parentid"))
	subID, err := strconv.Atoi(c.PostForm("subid"))
	if err != nil || parentID < 0 || parentID > database.MaxParentID || subID > database.MaxSubID || subID < 21 {
		responseError(c, response.Resp{
			ResponseCode: 400,
			ResponseMsg:  fmt.Sprintf("%s", database.InvalidID),
		})
		return
	}

	//session check
	if _, ok := isExists(authorNmae); !ok {
		responseError(c, response.ExpireSession)
		return
	}

	//start download video and image
	file, err := c.FormFile("file")
	if err !=  nil {
		logger.Error(err)
		responseError(c, response.IntervalErr)
		return
	}
	if file.Size > int64(videoSize) {
		logger.Debug("Video Size more: %d M", videoSize/1024/1024)
		responseError(c, response.Resp{
			ResponseCode: 400,
			ResponseMsg:  "Video size be smaller: " + strconv.Itoa(videoSize/1024/1024) + "M",
		})
		return
	}


	vid := models.NewUUID()
	err = c.SaveUploadedFile(file, config.Set.VideoSavePath+vid)
	if err!=nil {
		removeResource(vid)
		logger.Error(err)
		responseError(c, response.IntervalErr)
		return
	}

	//start transcoding to mp4 , if have error will delete local resource
	err = Transcoding(vid)
	if err != nil {
		removeResource(vid)
		logger.Debug(err)
		responseError(c, response.Resp{ResponseCode: 500, ResponseMsg: "Video Transcoding Failed"})
		return
	}
	//_ = os.Rename(config.Set.VideoSavePath+"output.mp4", config.Set.VideoSavePath+vid)

	//start upload video and image to oss serve
	ok := UploadToOss(config.Set.OssVideoPath+vid+".mp4", config.Set.VideoSavePath+"output.mp4", config.Set.BucketName)
	if !ok {
		removeResource(vid)
		responseError(c, response.IntervalErr)
		return
	}

	//delete local resource
	removeResource(vid)
	video_link := "https://jay-video.oss-cn-hongkong.aliyuncs.com/videos/" + vid + ".mp4"
	err = conn.AddNewVideoInfo(video_link, authorNmae, title ,parentID, subID)
	if err != nil {
		logger.Error(err)
		responseError(c, response.IntervalErr)
		return
	}

	//save to elastic search . if you want use elastic search please open and configure it
	if config.Set.OpenElastic {
		info, err := conn.CheckVideoInfo(video_link)
		if err != nil {
			logger.Error("one video save to elastic failed: %v", err)
		}
		err = saveToElastic(info)
		if err != nil {
			logger.Error(err)
		}
	}
	responseNormal(c, response.Success)

}

func loadTag(c *gin.Context) {
	tag, err := conn.LoadVideoTag()
	if err != nil {
		logger.Error(err)
		responseError(c, response.IntervalErr)
		return
	}
	responseNormal(c, tag)
}

func addTag(c *gin.Context) {
	var tag models.VideoTag
	err := c.ShouldBind(&tag)
	if err != nil {
		logger.Error(err)
		responseError(c, response.IntervalErr)
		return
	}

	err = conn.AddNewVideoTag(tag.ID, tag.ClassName)
	if err != nil {
		logger.Error(err)
		responseError(c, response.Resp{
			ResponseCode: 500,
			ResponseMsg:  "The Tag is Exists",
		})
		return
	}
	if id, _ := strconv.Atoi(tag.ID); id > 20 {
		database.MaxSubID++
	} else {
		database.MaxParentID++
	}
	responseNormal(c, response.Success)
}

//supply local video and image resource path
func Src(c *gin.Context) {
	src := c.Param("src")
	file, err := os.Open(config.Set.VideoSavePath + src)
	if err != nil {
		logger.Error(err)
		responseError(c, response.IntervalErr)
		return
	}
	http.ServeContent(c.Writer, c.Request, "", time.Now(), file)
}

func DeleteFileOss(c *gin.Context) {
	if !config.Set.OpenOssServe || config.Set.AccessKeySecret == "" ||
		config.Set.AccessKeyID == "" || config.Set.EndPoint == "" {
		logger.Warn(response.NotOpenOss.ResponseMsg)
		responseError(c, response.NotOpenOss)
		return
	}
	vid := c.PostForm("vid")
	if vid == "" {
		responseError(c, response.RequestInvalid)
		return
	}
	info, err := conn.CheckVideoInfo(vid)
	if err != nil {
		responseError(c, response.RequestInvalid)
		return
	}

	re := regexp.MustCompile(config.Set.ResourceAddress + "([a-zA-z0-9-]+)")
	deleteVid := re.FindAllSubmatch([]byte(info.Vid), -1)

	ok := DeleteOssFile(config.Set.OssVideoPath+string(deleteVid[0][1])+".mp4", config.Set.BucketName)
	if !ok {
		responseError(c, response.IntervalErr)
		return
	}

	if err=conn.DeleteVideo(info.Vid);err!=nil{
		logger.Error(err)
		responseError(c,response.IntervalErr)
		return
	}


	if config.Set.OpenElastic {
		if err = deleteElasticData(info.Vid); err != nil {
			logger.Error(err)
		}
	}

	responseNormal(c, response.Success)
}

func DeleteLocalVideo(c *gin.Context) {
	vid := c.PostForm("vid")
	if vid == "" {
		logger.Debug("The Delete vid is empty")
		responseError(c, response.RequestInvalid)
		return
	}
	info, err := conn.CheckVideoInfo(vid)
	if err != nil {
		logger.Error(err)
		responseError(c, response.RequestInvalid)
		return
	}
	//match id ,not link
	re := regexp.MustCompile("http://" + config.Set.Bind + "/api/all/src/([a-zA-z0-9-]+.*)")
	deleteVid := re.FindAllSubmatch([]byte(info.Vid), -1)
	removeResource(string(deleteVid[0][1]))

	err = conn.DeleteVideo(info.Vid)
	if err != nil {
		logger.Error("Database Delete Video Info Failed : %v", err)
		responseError(c, response.IntervalErr)
		return
	}
	//delete local resource by id

	if config.Set.OpenElastic {
		if err = deleteElasticData(info.Vid); err != nil {
			logger.Error(err)
		}
	}

	responseNormal(c, response.Success)
}

func DeleteComment(c *gin.Context) {
	vid := c.PostForm("vid")
	content:=c.PostForm("content")
	if vid == "" {
		logger.Debug("Delete vid is empty")
		responseError(c, response.RequestInvalid)
		return
	}
	info,err:=conn.CheckVideoInfo(vid)
	if err!=nil{
		logger.Error(err)
		responseError(c,response.RequestInvalid)
		return
	}
	err = conn.DeleteComment(info.ID,content)
	if err != nil {
		logger.Error(err)
		responseError(c, response.RequestInvalid)
		return
	}
	responseNormal(c, response.Success)

}

func DeleteUser(c *gin.Context) {
	userName := c.PostForm("user_name")
	pwd := c.PostForm("pwd")
	if userName == "" || pwd == "" {
		logger.Debug("User Name or Password is Empty")
		responseError(c, response.RequestInvalid)
		return
	}
	pwd, _ = encryption(pwd)
	if session, ok := isExists(userName); ok {
		err := conn.DeleteUser(session.UserName, pwd)
		if err != nil {
			logger.Debug(err)
			responseError(c, response.IntervalErr)
			return
		}
		//delete local and database session
		deleteSession(&session)
		responseNormal(c, response.Success)
		return

	}
	responseError(c, response.RequestInvalid)
}

func LoadVideoComment(c *gin.Context) {
	vid := c.Query("vid")
	from:=c.Query("from")
	to := c.Query("to")
	if err != nil {
		logger.Error("Parse to int Failed: %v", err)
		responseError(c, response.IntervalErr)
		return
	}

	info,err:=conn.CheckVideoInfo(vid)

	if err!=nil{
		logger.Error("Parse to int Failed: %v", err)
		responseError(c, response.RequestInvalid)
		return
	}
	//Always load comments to now from two weeks ago
	listComment, err := conn.LoadAllComments(info.ID, from,to)
	if err != nil {
		logger.Error(err)
		responseError(c, response.IntervalErr)
		return
	}

	responseNormal(c, listComment)

}

func AddComment(c *gin.Context) {

	content := c.PostForm("content")
	userName := c.PostForm("user_name")
	vid := c.PostForm("vid")

	if content == "" || userName == "" || vid == "" {
		logger.Debug(content,userName,vid)
		logger.Debug("Submit Content or user name or vid is empty")
		responseError(c, response.RequestInvalid)
		return
	}
	info,err:=conn.CheckVideoInfo(vid)
	if err!=nil{
		responseError(c,response.RequestInvalid)
		return
	}
	if session,ok:=isExists(userName);!ok{
		responseError(c,response.RequestInvalid)
		return
	}else {
		err = conn.AddNewComments(info.ID, session.UserName, content)
	if err != nil {
		logger.Error(err)
		responseError(c, response.IntervalErr)
		return
	}

	responseNormal(c, response.Success)
	}

}

func LoadSelfComment(c *gin.Context) {
	userName := c.Query("user_name")
	from, err := strconv.Atoi(c.Query("from"))
	to, err := strconv.Atoi(c.Query("to"))
	if err != nil || userName ==""{
		logger.Debug("Parse to int Failed: %v", err)
		responseError(c, response.RequestInvalid)
		return
	}

	selfComment, err := conn.LoadSelfCommentsRecord(userName, from, to)
	if err != nil {
		logger.Error(err)
		responseError(c, response.IntervalErr)
		return
	}

	responseNormal(c, selfComment)

}

//use multi field match
func ElasticSearch(c *gin.Context) {
	data := c.Query("search")
	if data == "" {
		responseError(c, response.RequestInvalid)
		return
	}
	result, err := searchElastic(data)
	if err != nil {
		logger.Error(err)
		responseError(c, response.IntervalErr)
		return
	}
	responseNormal(c, result)

}

func checkEmail(c *gin.Context)  {
	emailAddr:=c.PostForm("email_addr")
	result,err:=sendEmail(emailAddr)
	if err!=nil{
		responseError(c,response.Resp{
			ResponseCode: 400,
			ResponseMsg:  fmt.Sprintf("%s",err),
		})
		return
	}
	responseNormal(c,result)
}
