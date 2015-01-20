package diaryruapi

import (
	"encoding/json"
	"errors"
	_ "log"
	"net/http"
	"net/url"
)

func Auth(user, password string) (string, error) {
	values := url.Values{}
	values.Add("username", makeDiaryRuLogin(user))
	values.Add("password", makeDiaryRuPassword(password))
	values.Add("method", "user.auth")
	values.Add("appkey", appkey)

	resp, err := diaryPostForm(values)
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", errors.New(resp.Status)
	}
	var message DiaryAPIAuthResponse
	decoder := json.NewDecoder(resp.Body)
	if err = decoder.Decode(&message); err != nil {
		return "", err
	}
	sid := message.SID
	if message.Result != 0 {
		return "", errors.New(message.Error)
	}
	return sid, nil
}

func JournalGet(sid, userid, shortname string) (*JournalStruct, error) {
	var message DiaryAPIJournalGet
	values := url.Values{}
	values.Add("sid", sid)
	if userid != "" {
		values.Add("userid", userid)
	}
	if shortname != "" {
		values.Add("shortname", shortname)
	}
	values.Add("method", "journal.get")

	resp, err := diaryGet(values)
	defer resp.Body.Close()
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(resp.Status)
	}

	decoder := json.NewDecoder(resp.Body)
	if err = decoder.Decode(&message); err != nil {
		return nil, err
	}
	if message.Result != 0 {
		return nil, errors.New(message.Error)
	}
	return message.Journal, nil
}
