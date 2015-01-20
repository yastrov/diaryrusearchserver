package diaryruapi

import (
	"crypto/md5"
	"encoding/hex"
	"github.com/yastrov/charmap"
	"net/http"
	"net/url"
	"strconv"
)

func makeDiaryRuPassword(password string) string {
	hash := md5.Sum([]byte(key + password))
	return hex.EncodeToString(hash[:])
}

func makeDiaryRuLogin(username string) string {
	strcp1251, _ := charmap.Encode(username, "cp-1251")
	return strcp1251
}

func diaryPostForm(values url.Values) (*http.Response, error) {
	dURL, _ := url.Parse(DiaryMainUrl)
	resp, err := http.PostForm(dURL.String(), values)
	return resp, err
}

func diaryGet(values url.Values) (*http.Response, error) {
	dURL, _ := url.Parse(DiaryMainUrl)
	dURL.RawQuery = values.Encode()
	resp, err := http.Get(dURL.String())
	return resp, err
}

func (post *PostStruct) CountComments() uint64 {
	if post.Comments_count_data == "" {
		return 0
	}
	i, err := strconv.ParseUint(post.Comments_count_data, 10, 64)
	if err != nil || i == 0 {
		return i
	}
	return i
}
