package models

import (
	"crypto/rand"
	"fmt"
	"io"
)

type DataBase interface {
	//校验
	CheckPassword(username string) (*User, error) //登陆校验
	CheckSessionID(sid string) bool
	CheckVideoInfo(vid string) (*VInfo, error)
	//加载

	LoadAllVideos(limit int, from, to int64) (*[]VideoInfo, error)
	LoadAllComments(vid string, from, to string) (*[]Comments, error)
	LoadUserAllVideos(username string, from, to int) (*[]VideoInfo, error)

	//把数据库里的sessionID加载到内存里
	LoadALLSessionID() (*[]SessionID, error)
	//加载自己的评论过的记录，方便删除
	LoadSelfCommentsRecord(userName string, from, to int) (*[]Comments, error)
	LoadVideoTag() (*[]VideoTag, error)

	//	LoadVideoByClass()

	//添加
	AddNewComments(vid string, userName, content string) error
	AddNewVideoInfo(vid string, authName string, title string, ParentID, SubID int) error
	AddNewSessionID(sessionID string, expire int64, username string) error
	AddNewUser(userName, passWord string, email string) error
	AddNewVideoTag(id string, className string) error

	//删除
	DeleteUser(username string, pwd string) error
	DeleteVideo(vid string) error
	DeleteExpireSessionID(sid string) error
	DeleteComment(vid ,content string) error

	InitElasticData() error
}

func NewUUID() string {
	uuid := make([]byte, 16)
	n, _ := io.ReadFull(rand.Reader, uuid) //把 rand.Reader 读到 uuid  这个buffer里
	if n != len(uuid) {
		return ""
	}
	// variant bits; see section 4.1.1
	uuid[8] = uuid[8]&^0xc0 | 0x80
	// version 4 (pseudo-random); see section 4.1.3
	uuid[6] = uuid[6]&^0xf0 | 0x40
	return fmt.Sprintf("%x-%x-%x-%x-%x", uuid[0:4], uuid[4:6], uuid[6:8], uuid[8:10], uuid[10:])
}
