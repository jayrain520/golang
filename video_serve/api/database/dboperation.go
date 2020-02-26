package database

import (
	"context"
	"database/sql"
	"errors"
	"github.com/wonderivan/logger"
	"regexp"
	"strconv"
	"testgin/api/config"
	"testgin/api/models"
	"time"
)

var (

	MaxParentID = 20
	MaxSubID    = 31
	InvalidID   = errors.New("MaxParentID:1-" + strconv.Itoa(MaxParentID) + " , MaxSubID:21-" + strconv.Itoa(MaxSubID))
	Re 			= regexp.MustCompile("[0-9]+@[0-9a-zA-Z]+.com")
)

type DataBaseOperation struct {
	dbConnect *sql.DB
}

func init() {
	initializationID()

}

func initializationID() {
	var count int = 0
	stmtQuery, _ := dbConn.Prepare("select type_id from video_type where type_id > 20")
	rows, _ := stmtQuery.Query()
	for rows.Next() {
		count++
	}
	MaxSubID = 20 + count
	count = 0
	stmtQuery, _ = dbConn.Prepare("select type_id from video_type where type_id < 21")
	rows, _ = stmtQuery.Query()
	for rows.Next() {
		count++
	}
	MaxParentID = count
}

func NewDataBaseOperation() models.DataBase {
	return &DataBaseOperation{dbConnect: dbConn}
}

func (d *DataBaseOperation) CheckPassword(username string) (*models.User, error) {
	stmtQuery, err := d.dbConnect.Prepare("select user_name,pwd,email from users where user_name = ?")
	if err != nil {
		return nil, err
	}
	defer stmtQuery.Close()
	var u models.User
	err = stmtQuery.QueryRow(username).Scan(&u.Name, &u.Pwd, &u.Email)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (d *DataBaseOperation) CheckSessionID(sid string) bool {
	panic("implement me")
}

func (d *DataBaseOperation) CheckVideoInfo(vid string) (*models.VInfo, error) {
	stmtQuery, err := d.dbConnect.Prepare(`
	select id,vid,author_name,title,t.type_name,s.type_name from video_info 
	inner join video_type as t on parent_tag = t.type_id
	inner join video_type as s on sub_tag = s.type_id  
	where vid = ?
`)
	if err != nil {
		return nil, err
	}
	defer stmtQuery.Close()
	var v models.VInfo
	err = stmtQuery.QueryRow(vid).Scan(&v.ID,&v.Vid, &v.AuthorName, &v.Title, &v.Class, &v.SubClass)
	if err != nil {
		return &v, err
	}
	return &v, err
}

func (d *DataBaseOperation) LoadAllVideos(limiter int, from, to int64) (*[]models.VideoInfo, error) {
	stmtQuery, err := d.dbConnect.Prepare(`
	select vid,author_name,title,t.type_name,s.type_name from video_info 
	inner join video_type as t on parent_tag = t.type_id
	inner join video_type as s on sub_tag = s.type_id 
	where create_time > from_unixtime(?) and create_time  <= from_unixtime(?) 
	order by create_time desc  limit ?`)
	if err != nil {
		return nil, err
	}
	defer stmtQuery.Close()
	var v []models.VideoInfo

	rows, err := stmtQuery.Query(from, to, limiter)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		lv := models.VideoInfo{}
		err = rows.Scan(&lv.Vid, &lv.AuthorName, &lv.Title, &lv.Class, &lv.SubClass)
		if err != nil {
			return &v, err
		}
		v = append(v, lv)
	}
	return &v, nil

}

func (d *DataBaseOperation) LoadUserAllVideos(username string, from, to int) (*[]models.VideoInfo, error) {

	stmtQuery, err := d.dbConnect.Prepare(`
	select vid,author_name,title,p.type_name,s.type_name from video_info as v 
	inner join video_type as p on v.parent_tag = p.type_id 
	inner join video_type as s on v.sub_tag = s.type_id 
	where author_name = ? order by create_time desc limit ?,?
`)
	if err != nil {
		return nil, err
	}
	defer stmtQuery.Close()
	var v []models.VideoInfo
	rows, err := stmtQuery.Query(username, from, to)
	if err != nil {
		return &v, err
	}

	for rows.Next() {
		lv := models.VideoInfo{}
		err = rows.Scan(&lv.Vid, &lv.AuthorName, &lv.Title, &lv.Class, &lv.SubClass)
		if err != nil {
			return &v, err
		}
		v = append(v, lv)
	}

	return &v, nil
}

func (d *DataBaseOperation) LoadAllComments(vid string, from, to string) (*[]models.Comments, error) {
	stmtQuery, err := d.dbConnect.Prepare(`
	select user_name,content,ctime from comments 
	where video_id = ? order by ctime desc limit ?,?
`)
	if err != nil {
		return nil, err
	}
	defer stmtQuery.Close()
	rows, err := stmtQuery.Query(vid,from,to)
	if err != nil {
		return nil, err
	}
	var comments []models.Comments
	for rows.Next() {
		c := models.Comments{}
		err = rows.Scan(&c.UserName, &c.Content, &c.CTime)
		if err != nil {
			return &comments, err
		}
		comments = append(comments, c)
	}
	return &comments, nil
}

func (d *DataBaseOperation) LoadALLSessionID() (*[]models.SessionID, error) {
	stmtQuery, err := d.dbConnect.Prepare("select session_id,expire,user_name from sessions ")
	if err != nil {
		return nil, err
	}
	defer stmtQuery.Close()
	rows, err := stmtQuery.Query()
	if err != nil {
		return nil, err
	}
	var sessions []models.SessionID
	for rows.Next() {
		s := models.SessionID{}
		err := rows.Scan(&s.ID, &s.Expire, &s.UserName)
		if err != nil {
			return &sessions, err
		}
		sessions = append(sessions, s)
	}

	return &sessions, nil
}

func (d *DataBaseOperation) LoadSelfCommentsRecord(userName string, from, to int) (*[]models.Comments, error) {

	stmtQuery, err := d.dbConnect.Prepare(`
select info.vid,user_name,content,ctime from comments as c
inner join video_info as info on c.video_id = info.id 
where user_name = ? order by ctime desc limit ?,?
`)

	if err != nil {
		return nil, err
	}
	defer stmtQuery.Close()

	rows, err := stmtQuery.Query(userName, from, to)
	if err != nil {
		return nil, err
	}
	var lc []models.Comments
	for rows.Next() {
		c := models.Comments{}
		err := rows.Scan(&c.VideoID, &c.UserName, &c.Content, &c.CTime)
		if err != nil {
			return &lc, err
		}
		lc = append(lc, c)
	}

	return &lc, nil

}

func (d *DataBaseOperation) LoadVideoTag() (*[]models.VideoTag, error) {
	stmtQuery, err := d.dbConnect.Prepare("select type_id,type_name from video_type")
	if err != nil {
		return nil, err
	}
	defer stmtQuery.Close()
	rows, err := stmtQuery.Query()
	if err != nil {
		return nil, err
	}
	var tag []models.VideoTag
	for rows.Next() {
		t := models.VideoTag{}
		err := rows.Scan(&t.ID, &t.ClassName)
		if err != nil {
			return &tag, err
		}
		tag = append(tag, t)
	}

	return &tag, nil

}

func (d *DataBaseOperation) AddNewComments(vid string, userName, content string) error {
	stmtIns, err := d.dbConnect.Prepare("insert into comments(video_id,user_name,content) values (?,?,?)")
	if err != nil {
		return err
	}
	defer stmtIns.Close()
	_, err = stmtIns.Exec(vid, userName, content)
	if err != nil {
		return err
	}
	return nil
}

func (d *DataBaseOperation) AddNewVideoInfo(vid string, authorName string, title string,  ParentID, SubID int) error {
	if ParentID < 0 || ParentID > MaxParentID || SubID > MaxSubID || SubID < 21 {
		return InvalidID
	}
	stmtIns, err := d.dbConnect.Prepare(`
	insert into video_info(vid, author_name, title,parent_tag,sub_tag) values(?,?,?,?,?)
`)
	if err != nil {
		return err
	}
	defer stmtIns.Close()

	_, err = stmtIns.Exec(vid, authorName, title, ParentID, SubID)
	if err != nil {
		return err
	}
	return nil

}

func (d *DataBaseOperation) AddNewSessionID(sessionID string, expire int64, username string) error {
	stmtIns, err := d.dbConnect.Prepare(`
	insert into sessions(session_id, expire, user_name) values (?,?,?)
`)
	if err != nil {
		return err
	}
	defer stmtIns.Close()
	_, err = stmtIns.Exec(sessionID, expire, username)
	if err != nil {
		return err
	}

	return nil

}

func (d *DataBaseOperation) AddNewUser(userName, passWord string, email string) error {
	stmtIns, err := d.dbConnect.Prepare("insert into users(user_name, pwd, email) values (?,?,?)")
	if err != nil {
		return nil
	}
	defer stmtIns.Close()

	result := Re.FindAllString(email, -1)
	var em string
	for _, i := range result {
		em += i
	}
	_, err = stmtIns.Exec(userName, passWord, em)
	if err != nil {
		return err
	}

	return nil
}

func (d *DataBaseOperation) AddNewVideoTag(vid string, className string) error {
	stmtIns, err := d.dbConnect.Prepare("insert into video_type(type_id,type_name) values (?,?)")
	if err != nil {
		return err
	}

	defer stmtIns.Close()
	_, err = stmtIns.Exec(vid, className)
	if err != nil {
		return err
	}
	return nil
}

func (d *DataBaseOperation) DeleteUser(username string, pwd string) error {
	stmtDel, err := d.dbConnect.Prepare("delete from users where user_name = ? and pwd = ?")
	if err != nil {
		return err
	}
	defer stmtDel.Close()
	_, err = stmtDel.Exec(username, pwd)
	if err != nil {
		return err
	}
	return nil
}

func (d *DataBaseOperation) DeleteVideo(vid string) error {
	stmtDel, err := d.dbConnect.Prepare("delete from video_info where vid = ?")
	if err != nil {
		return err
	}
	defer stmtDel.Close()
	_, err = stmtDel.Exec(vid)
	if err != nil {
		return err
	}
	return nil
}

func (d *DataBaseOperation) DeleteExpireSessionID(sid string) error {
	stmtDel, err := d.dbConnect.Prepare("delete from sessions where session_id = ?")
	if err != nil {
		return err
	}
	defer stmtDel.Close()
	_, err = stmtDel.Exec(sid)
	if err != nil {
		return err
	}

	return nil

}

func (d *DataBaseOperation) DeleteComment(vid ,content string) error {
	stmtDel, err := d.dbConnect.Prepare("delete from comments where video_id = ? and content = ?")
	if err != nil {
		return err
	}
	defer stmtDel.Close()
	_, err = stmtDel.Exec(vid,content)
	if err != nil {
		return err
	}
	return nil
}

func (d *DataBaseOperation) InitElasticData() error {
	listVideoInfo, err := d.LoadAllVideos(int(config.Set.InitDataNumber), int64(1), time.Now().Unix())
	if err != nil {
		logger.Emer(err)
		return err
	}
	//logger.Error(len(*listVideoInfo))
	for _, info := range *listVideoInfo {
		indexServe := config.Set.ElsClient.Index().
			Index(config.Set.Index).Type(config.Set.Type).Id(info.Vid).
			BodyJson(info)
		_, err = indexServe.Do(context.Background())
		if err != nil {
			logger.Warn(err)
			continue
		}

	}
	return nil
}
