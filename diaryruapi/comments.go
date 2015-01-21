package diaryruapi

import (
	"encoding/json"
	"errors"
	_ "log"
	"net/http"
	"net/url"
	"sync"
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

// Channels
func Comment_get_for_post_Channels(sid string, post *PostStruct, comment_chan chan *CommentStruct, err_chan chan error) {
	if post.CountComments() == 0 {
		return
	}
	values := url.Values{}
	values.Add("sid", sid)
	values.Add("method", "comment.get")
	values.Add("postid", post.Postid)
	resp, err := diaryGet(values)
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		err_chan <- errors.New(resp.Status)
		return
	}
	var message DiaryAPICommentGet
	decoder := json.NewDecoder(resp.Body)
	if err = decoder.Decode(&message); err != nil {
		err_chan <- err
		return
	}
	if message.Result != 0 {
		err_chan <- errors.New(message.Error)
		return
	}
	for _, comment_unit := range message.Comments {
		comment_unit.Postid = post.Postid
		comment_chan <- comment_unit
	}
}

func Comment_get_for_posts_Channels(sid string, posts []*PostStruct, comment_chan chan *CommentStruct, err chan error, wg *sync.WaitGroup) {
	defer wg.Done()
	for _, post := range posts {
		if post.CountComments() != 0 {
			Comment_get_for_post_Channels(sid, post, comment_chan, err)
		}
	}
}
