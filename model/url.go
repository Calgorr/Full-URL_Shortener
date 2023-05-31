package model

import (
	"crypto/md5"
	"encoding/base64"
	"regexp"
	"strings"
	"time"
)

type URL struct {
	ID        int64     `json:"id"`
	UserID    int64     `json:"userid"`
	LongURL   string    `json:"long_url"`
	ShortURL  string    `json:"short_url"`
	UsedTimes int       `json:"usedtimes"`
	CreatedAt time.Time `json:"created_at"`
	LastUsed  time.Time `json:"last_used"`
}

func NewLink(url, customPath string) *URL {
	link := new(URL)
	link.LongURL = url
	link.LastUsed = time.Now()
	md5 := md5.Sum([]byte(url))
	if customPath != "" {
		link.ShortURL = customPath
	} else {
		link.ShortURL = strings.ReplaceAll(strings.ReplaceAll(base64.StdEncoding.EncodeToString(md5[:])[:6], "/", "_"), "+", "-")
	}
	return link
}
func (u *URL) IsValidURL(url string) bool {
	// Regular expression pattern for validating URL format
	pattern := `^(http(s)?:\/\/)?([^\s]+\.){1,2}[^\s]+(\/[^\s]*)?$`

	match, _ := regexp.MatchString(pattern, url)
	return match
}
