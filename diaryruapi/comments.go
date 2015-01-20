package diaryruapi

import (
	"encoding/json"
	"errors"
	_ "log"
	"net/http"
	"net/url"
)

func Comment_get(sid, postid string) ([]*CommentStruct, error) {
	values := url.Values{}
	values.Add("sid", sid)
	values.Add("method", "comment.get")
	values.Add("postid", postid)
	resp, err := diaryGet(values)
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(resp.Status)
	}
	var message DiaryAPICommentGet
	decoder := json.NewDecoder(resp.Body)
	if err = decoder.Decode(&message); err != nil {
		return nil, err
	}
	if message.Result != 0 {
		return nil, errors.New(message.Error)
	}
	result := make([]*CommentStruct, 0, len(message.Comments))
	for _, comment_unit := range message.Comments {
		comment_unit.Postid = postid
		result = append(result, comment_unit)
	}
	return result, nil
}

func Comment_get_for_post(sid string, post *PostStruct) ([]*CommentStruct, error) {
	post.CountComments()
	if post.CountComments() == 0 {
		return make([]*CommentStruct, 0), nil
	}
	values := url.Values{}
	values.Add("sid", sid)
	values.Add("method", "comment.get")
	values.Add("postid", post.Postid)
	resp, err := diaryGet(values)
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(resp.Status)
	}
	var message DiaryAPICommentGet
	decoder := json.NewDecoder(resp.Body)
	if err = decoder.Decode(&message); err != nil {
		return nil, err
	}
	if message.Result != 0 {
		return nil, errors.New(message.Error)
	}
	result := make([]*CommentStruct, 0, len(message.Comments))
	for _, comment_unit := range message.Comments {
		comment_unit.Postid = post.Postid
		result = append(result, comment_unit)
	}
	return result, nil
}

func Comment_get_for_posts(sid string, posts []*PostStruct) (result []*CommentStruct, err error) {
	result = make([]*CommentStruct, 0, 20)
	for _, post := range posts {
		if post.Comments_count_data != "" {
			comments, err := Comment_get_for_post(sid, post)
			if err != nil {
				return nil, err
			}
			for _, comment := range comments {
				result = append(result, comment)
			}
		}
	}
	return result, nil
}
