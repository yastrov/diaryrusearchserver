package diaryruapi

import (
	"fmt"
	"log"
	"testing"
)

var (
	user := ""
	password := ""
	shortname := ""
	diarytype := "diary"
)

func TestGetPostsAndComments(t *testing.T) {
	sid, err := Auth(user, password)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(sid)
	posts, err := Post_get(sid, shortname, diarytype, "", "")
	fmt.Println(len(posts))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(posts)
	for _, post := range posts {
		fmt.Println(post.Avatar_path)
		if post == nil {
			//log.Fatal("post is nnil")
			fmt.Println("nil", post)
		}
		comments, err := Comment_get_for_post(sid, post)
		if err != nil {
			log.Fatal(err)
		}
		for _, c := range comments {
			fmt.Println(c.Postid)
		}
	}
}
