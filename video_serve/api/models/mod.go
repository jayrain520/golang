package models

type User struct {
	Name  string `json:"user_name" form:"user_name"`
	Pwd   string `json:"pwd" form:"pwd"`
	Email string `json:"email" form:"email"`
}

type VideoInfo struct {
	Vid        string `json:"vid" form:"vid"`
	AuthorName string `json:"author_name" form:"author_name"`
	Title      string `json:"title",form:"title"`
	Class      string `json:"class" form:"class"`
	SubClass   string `json:"sub_class" form:"sub_class"`
}

type Comments struct {
	VideoID  string `json:"video_id",form:"video_id"`
	UserName string `json:"user_name" form:"user_name"`
	Content  string `json:"content" form:"content"`
	CTime    string `json:"comment_time" form:"comment_time"`
}

type SessionID struct {
	ID       string `json:"id" form:"id"`
	Expire   int64  `json:"expire" form:"expire"`
	UserName string `json:"user_name" form:"user_name"`
}

type VideoTag struct {
	ID        string `json:"id" form:"id"`
	ClassName string `json:"class_name" form:"class_name"`
}

//check video info
type VInfo struct {
	ID string
	VideoInfo
}
