package main

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha1"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/wonderivan/logger"
	"gopkg.in/gomail.v2"
	"hash"
	"html/template"
	"io"
	rand2 "math/rand"
	"os"
	"os/exec"
	"runtime"
	"testgin/api/config"
	"testgin/api/database"
	"testgin/api/response"
	"time"
)

var (
	//mysql connect
	conn = database.NewDataBaseOperation()

	//upload video limit for 50M
	videoSize = 1024 * 1024 * 50
	////upload video limit for 1M
	imageSize = 1024 * 1024 * 1

	limitNumber = 10 //upload limit number

	//1440minute for 24 hour .
	//default: 1 minute can upload 10 video
	resetLimit time.Duration = 1

	//global variable configure
	globalVar = make(map[string]interface{})

	//QQ email
	emailBody = bytes.Buffer{}
	//ready send email
	email = gomail.NewMessage()

	//send email body
	body = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>JayRain</title>
</head>
<body>
<div class="top"></div>
<div class="main">
    <div>
        <h2 >Dear sir. /Madam: </h2>
    <p>Thank You for Your use . Your Verification Code:&nbsp;&nbsp;&nbsp;
        <span style="color: lightskyblue">{{.Code}}</span></p>
    </div>

    <div>
     <p>Thanks, The JayRain Team</p>
    </div>
</div>
</body>
</html>
`

)

func init() {
	if config.Set.OpenElastic {
		if err = conn.InitElasticData(); err != nil {
			logger.Warn(err)
		}

	}
}

func encryption(pwd string) (string, error) {
	h := hmac.New(func() hash.Hash {
		return sha1.New()
	}, []byte("jaykey"))

	if _, err := io.WriteString(h, pwd); err != nil {
		return "", err
	}
	ok := h.Sum(nil)
	return fmt.Sprintf("%X", ok), nil
}

func responseError(c *gin.Context, resp response.Resp) {
	c.JSON(resp.ResponseCode, gin.H{
		"error": resp,
	})
}

func responseNormal(c *gin.Context, resp interface{}) {
	c.JSON(200, gin.H{
		"normal": resp,
	})
}

func removeResource(vid string) {
	_ = os.Remove(config.Set.VideoSavePath + vid)
	_ = os.Remove(config.Set.VideoSavePath + "output.mp4")
}

//transcoding to mp4
func Transcoding(vid string) error {
	var err bytes.Buffer
	command := "ffmpeg -i " + config.Set.VideoSavePath + vid +
		" -vcodec libx264 -threads 6 -preset ultrafast -b:v 2000k -acodec copy " +
		config.Set.VideoSavePath + "output.mp4"

	if runtime.GOOS == "linux" {
		cmd := exec.Command("/bin/bash", "-c", command)
		cmd.Stderr = &err
		_ = cmd.Run()
		if e := index(err.Bytes(), []byte("未找到命令")); e != nil {
			logger.Debug(err.String())
			return errors.New("Not Found ffmpeg . Please install ffmpeg Tools")
		}
	} else {
		cmd := exec.Command("cmd", "/C", command)
		cmd.Stderr = &err
		_ = cmd.Run()
		if e := index(err.Bytes(), []byte("不是内部或外部命令，也不是可运行的程序或批处理文件")); e != nil {
			logger.Debug(err.String())
			return errors.New("Not Found ffmpeg . Please install ffmpeg Tools")
		}
	}

	return nil
}

func sendEmail(emailAddr string) (map[string]interface{}, error) {
	result := database.Re.FindAllString(emailAddr, -1)
	var em string
	for _, i := range result {
		em += i
	}
	if em != "" {
		t, err := template.ParseFiles(body)
		code := randCode()
		_ = t.Execute(&emailBody, code)

		email.SetHeader(`From`, config.Set.EmailUser)
		email.SetHeader(`To`, emailAddr)
		email.SetHeader(`Subject`, config.Set.EmailTitle)
		email.SetBody(`text/html`, emailBody.String())
		err = gomail.NewDialer(config.Set.EmailHost, int(config.Set.EmailPort),
			config.Set.EmailUser, config.Set.EmailAccesskey).
			DialAndSend(email)

		if err != nil {
			return nil, err
		}
		return code, nil
	} else {
		return nil, errors.New("invalid email:" + em)
	}
}

func randCode() map[string]interface{} {
	code := make(map[string]interface{})
	rand2.Seed(time.Now().UnixNano())
	code["Code"] = rand2.Intn(999999)
	return code
}

func index(p, s []byte) error {

	var slice []byte
	for _, b := range p {
		for i := 0; i < len(s); i++ {
			if b == s[i] {
				slice = append(slice, s[i])
			}
		}
	}
	if len(slice) > 0 {
		return errors.New("happened err")
	} else {
		return nil
	}
}
