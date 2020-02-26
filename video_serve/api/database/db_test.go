package database

import (
	"database/sql"
	"github.com/wonderivan/logger"
	"strconv"
	"testgin/api/models"
	"testing"
	"time"
)

//Duplicate entry 'jay-1806245714@qq.com' for key 'unique' 唯一key

var DB models.DataBase

func TestMain(m *testing.M) {
	//clearTable()
	DB = NewDataBaseOperation()
	m.Run()
	//clearTable()

}

func clearTable() {
	_, _ = dbConn.Exec("truncate users")
	_, _ = dbConn.Exec("truncate comments")
	_, _ = dbConn.Exec("truncate video_info")
	_, _ = dbConn.Exec("truncate sessions")
}

//test user table
func TestUser(t *testing.T) {

	t.Run("add", testAddUser)
	t.Run("getUser", testGetUser)
	t.Run("delete", testDelete)
	t.Run("repeatGet", testRepeatGetUser)

}
func testAddUser(t *testing.T) {
	err := DB.AddNewUser("jay", "shijie", "1806245714@qq.com")
	if err != nil {
		t.Errorf("%s", err)
	}
}
func testGetUser(t *testing.T) {
	pwd, err := DB.CheckPassword("jay")
	if err != nil {
		if err == sql.ErrNoRows {
			logger.Debug(err)
			return
		}
		t.Errorf("%s", err)

	}
	t.Logf("passworld: %s", pwd)
}
func testDelete(t *testing.T) {
	err := DB.DeleteUser("jay", "shijie")
	if err != nil {
		t.Errorf("%s", err)

	}
}
func testRepeatGetUser(t *testing.T) {
	testGetUser(t)
}

//test video_info table
func TestVideoInfo(t *testing.T) {
	t.Run("AddUser", testAddUser)
	t.Run("AddVideoInfo", testAddVideoInfo)
	t.Run("GetUserAllVideoInfo", testGetUserAllVideoInfo)
	//主页加载
	t.Run("GetAllVideoInfo", testGetAllVideoInfo)
	t.Run("DeleteVideoInfo", testDeleteVideoInfo)

}
func testAddVideoInfo(t *testing.T) {
	//for i := 1; i < 6; i++ {
	//	s := strconv.Itoa(i)
	//	err = DB.AddNewVideoInfo("1234"+s, "jay", "测试"+s, "http://image"+s, i, 22+i)
	//	if err != nil {
	//		t.Fatal(err)
	//	}
	//}
}
func testGetUserAllVideoInfo(t *testing.T) {
	vInfo, err := DB.LoadUserAllVideos("jay", 1, 3)
	if err != nil {
		t.Fatal(err)
	}
	for _, v := range *vInfo {
		t.Logf("%+v", v)
	}
}

//加载最新的视频一般用在主页
func testGetAllVideoInfo(t *testing.T) {
	lVideo, err := DB.LoadAllVideos(4, 1, time.Now().Unix())
	if err != nil {
		t.Fatal(err)
	}
	for _, v := range *lVideo {
		t.Logf("%+v\n", v)
	}
}

//添加到将要被删除的表单里
func testDeleteVideoInfo(t *testing.T) {

}

//test comments table
func TestComments(t *testing.T) {
	t.Run("AddComment", testAddComment)
	t.Run("GetAllComment", testGetAllComment)
	t.Run("GetSelfCommentRcord", testGetSelfCommentRcord)
	t.Run("DeleteCommets", testDeleteCommets)
}
func testAddComment(t *testing.T) {
	testAddVideoInfo(t) //vid:1234, name:jay, title: 测试, 视频封面:http://image, class:电影
	for i := 1; i < 5; i++ {
		str := strconv.Itoa(i)
		DB.AddNewComments("1234"+str, "jay", "I like this video")
		time.Sleep(time.Second)
	}
}
func testGetAllComment(t *testing.T) {
	//lComment, err := DB.LoadAllComments("12344", 1, time.Now().Unix())
	//if err != nil {
	//	t.Fatal(err)
	//}
	//for _, c := range *lComment {
	//	t.Logf("video id: %+v", c)
	//}

}
func testGetSelfCommentRcord(t *testing.T) {
	//测试数据
	//vid:1234, name:jay, title: 测试, 视频封面:http://image, class:电影
	lComment, err := DB.LoadSelfCommentsRecord("jay", 3, 500)
	if err != nil {
		t.Fatal(err)
	}
	for _, c := range *lComment {
		t.Logf("%+v", c)
	}
}
func testDeleteCommets(t *testing.T) {
	////测试数据
	////vid:1234, name:jay, title: 测试, 视频封面:http://image, class:电影
	//err := DB.DeleteComment("jay")
	//if err != nil {
	//	t.Fatal(err)
	//}
}

//test session

func TestSession(t *testing.T) {
	t.Run("AddSession", testAddSession)
	t.Run("LoadAllSession", testLoadAllSession)
	t.Run("DeleteSession", testDeleteSession)
}

func testAddSession(t *testing.T) {

	DB.AddNewSessionID(models.NewUUID(), time.Now().Unix()+1000000, "jay")
	DB.AddNewSessionID(models.NewUUID(), time.Now().Unix()+1000000, "jays")
	DB.AddNewSessionID(models.NewUUID(), time.Now().Unix()+1000000, "jayss")
	DB.AddNewSessionID(models.NewUUID(), time.Now().Unix()+1000000, "jaysss")

}
func testLoadAllSession(t *testing.T) {
	Lid, err := DB.LoadALLSessionID()
	if err != nil {
		t.Fatal(err)
	}
	for _, id := range *Lid {
		t.Log(id)
	}
}
func testDeleteSession(t *testing.T) {
	//response:=DB.DeleteExpireSessionID(id)
	//if response!=nil{
	//	t.Fatal(response)
	//}
}

func TestElastic(t *testing.T) {
	err := DB.InitElasticData()
	if err != nil {
		t.Fatal(err)
	}

}
