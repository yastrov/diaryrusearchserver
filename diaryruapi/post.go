package diaryruapi

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"strconv"
)

/*
diarytype
"diary"
"favorites"
"quotes"
*/
func Post_get(sid, shortname, diarytype, from, src string) ([]*PostStruct, error) {
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
