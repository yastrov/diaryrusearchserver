package diaryruapi

import (
	"bytes"
)

const (
	appkey       = "" // ok also
	key          = "" // pk also
	DiaryMainUrl = "http://www.diary.ru/api/"
)

type DiaryAPIAuthResponse struct {
	Result int    `json:"result,string"`
	SID    string `json:"sid"`
	Error  string `json:"error"`
}

type PostStruct struct {
	Avatar_path         string            `json:"avatar_path"`
	Postid              string            `json:"postid"`
	Juserid             string            `json:"juserid"`
	Shortname           string            `json:"shortname"`
	Journal_name        string            `json:"journal_name"`
	Message_src         string            `json:"message_src"`
	Message_html        string            `json:"message_html"`
	Author_userid       string            `json:"author_userid"`
	Author_shortname    string            `json:"author_shortname"`
	Author_username     string            `json:"author_username"`
	Author_title        string            `json:"author_title"`
	Title               string            `json:"title"`
	No_comments         string            `json:"no_comments"`         // Flag for no comments
	Comments_count_data string            `json:"comments_count_data"` //Count of comments
	Tags_data           map[string]string `json:"tags_data"`
	Subscribed          string            `json:"subscribed"`
	Can_edit            string            `json:"can_edit"`
	Avatarid            string            `json:"avatarid "`
	No_smile            string            `json:"no_smile"`
	Jaccess             string            `json:"jaccess"`
	Dateline_cdate      string            `json:"dateline_cdate"`
	Close_access_mode2  string            `json:"close_access_mode2"`
	Close_access_mode   string            `json:"close_access_mode"`
	Dateline_date       string            `json:"dateline_date"`
	Access              string            `json:"access"`
}

func (post *PostStruct) MakeUrl() string {
	buffer := bytes.NewBufferString("http://")
	buffer.WriteString(post.Author_shortname)
	buffer.WriteString(".diary.ru/p")
	buffer.WriteString(post.Postid)
	buffer.WriteString(".htm")
	return buffer.String()
}

type DiaryAPIPostGet struct {
	Result int                    `json:"result,string"`
	Posts  map[string]*PostStruct `json:"posts"`
	Error  string                 `json:"error"`
}

type DiaryAPICommentGet struct {
	Result   int                       `json:"result,string"`
	Comments map[string]*CommentStruct `json:"comments"`
	Error    string                    `json:"error"`
}

type CommentStruct struct {
	Avatar_path      string `json:"avatar_path"`
	Postid           string `json:"postid"` // For manually control
	Commentid        string `json:"commentid"`
	Shortname        string `json:"shortname"`
	Message_html     string `json:"message_html"`
	Author_userid    string `json:"author_userid"`
	Author_shortname string `json:"author_shortname"`
	Author_avatar    string `json:"author_avatar"`
	Author_username  string `json:"author_username"`
	Author_title     string `json:"author_title"`
	Can_edit         string `json:"can_edit"`
	Can_delete       string `json:"can_delete"`
	Dateline         string `json:"dateline"`
}

func (post *CommentStruct) MakeUrl() string {
	buffer := bytes.NewBufferString("http://")
	buffer.WriteString(post.Author_shortname)
	buffer.WriteString(".diary.ru/p")
	buffer.WriteString(post.Postid)
	buffer.WriteString(".htm#")
	buffer.WriteString(post.Commentid)
	return buffer.String()
}

type DiaryAPIPostCreate struct {
	Result  int    `json:"result,string"`
	Error   string `json:"error,string"`
	Message string `json:"message"`
	PostID  string `json:"postid"`
}

type DiaryAPIJournalGet struct {
	Result  int            `json:"result,string"`
	Journal *JournalStruct `json:"journal"`
	Error   string         `json:"error"`
}

type JournalStruct struct {
	Userid    string `json:"userid"`
	Ctime     string `json:"ctime"`
	Shortname string `json:"shortname"`
	Title     string `json:"title"`
	Count_pch string `json:"count_pch"`
	Access    string `json:"access"`
	//Tags          map[string]string `json:"tags"`
	Posts         uint64 `json:"posts,string"`
	Can_write     string `json:"can_write"`
	Count_members string `json:"count_members"`
	Last_post_id  string `json:"last_post_id"`
	Last_post     string `last_post"`
	Jtype         string `json:"jtype"`
}
