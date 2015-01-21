package diaryruapi

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"sync"
)

/*
diarytype
"diary"
"favorites"
"quotes"
*/
func Post_get(sid, shortname, diarytype, from, src string) ([]*PostStruct, error) {
	log.Println("Post_get")
	values := url.Values{}
	values.Add("sid", sid)
	values.Add("type", diarytype)
	values.Add("method", "post.get")
	if shortname != "" {
		values.Add("shortname", shortname)
	}
	if from != "" {
		values.Add("from", from)
	}
	if src != "" {
		values.Add("src", src)
	}
	resp, err := diaryGet(values)
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(resp.Status)
	}
	var message DiaryAPIPostGet
	decoder := json.NewDecoder(resp.Body)
	if err = decoder.Decode(&message); err != nil {
		return nil, err
	}
	result := make([]*PostStruct, 0, len(message.Posts))
	for _, post_unit := range message.Posts {
		result = append(result, post_unit)
	}
	if message.Result != 0 {
		return nil, errors.New(message.Error)
	}
	return result, nil
}

func Post_create(sid, title, message string) (string, error) {
	log.Println("Post_create")
	values := url.Values{}
	values.Add("sid", sid)
	values.Add("message", message)
	values.Add("message_src", message)
	values.Add("method", "post.create")
	values.Add("title", title)
	values.Add("close_access_mode", "0")

	resp, err := diaryPostForm(values)
	defer resp.Body.Close()
	if err != nil {
		return "", err
	}
	if resp.StatusCode != http.StatusOK {
		return "", errors.New(resp.Status)
	}

	var mess DiaryAPIPostCreate
	decoder := json.NewDecoder(resp.Body)
	if err = decoder.Decode(&message); err != nil {
		return "", err
	}
	if mess.Result != 0 {
		return "", errors.New(mess.Error)
	}
	return mess.PostID, nil
}

/*
diarytype
"diary"
"favorites"
"quotes"
*/
func PostsAllGet(sid, diarytype string, journal *JournalStruct) ([]*PostStruct, error) {
	result := make([]*PostStruct, 0, journal.Posts)
	var i uint64 = 0
	for i < journal.Posts {
		r, err := Post_get(sid, journal.Shortname, diarytype, strconv.FormatUint(i, 10), "")
		if err != nil {
			return result, err
		}
		i += uint64(len(r))
		for _, post := range r {
			result = append(result, post)
		}
	}
	return result, nil
}

// Channels version
func Post_get_Channels(sid, shortname, diarytype, from, src string, post_chan chan *PostStruct, err_chan chan error) {
	values := url.Values{}
	values.Add("sid", sid)
	values.Add("type", diarytype)
	values.Add("method", "post.get")
	if shortname != "" {
		values.Add("shortname", shortname)
	}
	if from != "" {
		values.Add("from", from)
	}
	if src != "" {
		values.Add("src", from)
	}
	resp, err := diaryGet(values)
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		err_chan <- errors.New(resp.Status)
		return
	}
	var message DiaryAPIPostGet
	decoder := json.NewDecoder(resp.Body)
	if err = decoder.Decode(&message); err != nil {
		err_chan <- err
		return
	}
	for _, post_unit := range message.Posts {
		post_chan <- post_unit
	}
	if message.Result != 0 {
		err_chan <- errors.New(message.Error)
		return
	}
}

func PostsAllGetChannels(sid, diarytype string, journal *JournalStruct, post_chan chan *PostStruct, err chan error, wg *sync.WaitGroup) {
	log.Println("PostsAllGetChannels")
	var i uint64 = 0
	defer wg.Done()
	for i < /*journal.Posts*/ 1 {
		log.Println("Make sync Post_get: ", strconv.FormatUint(i, 10))
		r, er := Post_get(sid, journal.Shortname, diarytype, strconv.FormatUint(i, 10), "")
		if er != nil {
			err <- er
			return
		}
		for _, post := range r {
			post_chan <- post
		}
		i += uint64(len(r))
	}
}
