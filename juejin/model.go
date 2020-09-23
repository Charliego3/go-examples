package main

import (
	"encoding/json"
	"github.com/whimthen/kits/logger"
	"github.com/whimthen/kits/request"
	"net/url"
	"time"
)

type Section struct {
	Id        string `json:"_id"`
	Title     string `json:"title"`
	User      string `json:"user"`
	MetaId    string `json:"metaId"`
	PV        int    `json:"pv"`
	SectionId string `json:"sectionId"`
	UV        int    `json:"uv"`
}

type ListSection struct {
	S int       `json:"s"`
	M string    `json:"m"`
	D []Section `json:"d"`
}

type QueryBase struct {
	UID      string `json:"uid"`
	ClientId string `json:"client_id"`
	Token    string `json:"token"`
	Src      string `json:"src"`
}

type Content struct {
	S int    `json:"s"`
	M string `json:"m"`
	D struct {
		ID           string    `json:"_id"`
		Title        string    `json:"title"`
		IsFree       bool      `json:"isFree"`
		IsFinished   bool      `json:"isFinished"`
		User         string    `json:"user"`
		ViewCount    int       `json:"viewCount"`
		MetaID       string    `json:"metaId"`
		Content      string    `json:"content"`
		ContentSize  int       `json:"contentSize"`
		HTML         string    `json:"html"`
		CreatedAt    time.Time `json:"createdAt"`
		UpdatedAt    time.Time `json:"updatedAt"`
		IsDeleted    bool      `json:"isDeleted"`
		Pv           int       `json:"pv"`
		CommentCount int       `json:"commentCount"`
		SectionID    string    `json:"sectionId"`
	} `json:"d"`
}

func (q *QueryBase) ToUrlValues() url.Values {
	return url.Values{
		"uid":       {q.UID},
		"client_id": {q.ClientId},
		"token":     {q.Token},
		"src":       {q.Src},
	}
}

func (q *QueryBase) GetListSection(id string) *ListSection {
	values := q.ToUrlValues()
	values.Add("id", id)
	resp, err := request.GetWithHeader("https://xiaoce-cache-api-ms.juejin.im/v1/getListSection", values, request.ReqHeader{
		"User-Agent":   "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_5) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/13.1.1 Safari/605.1.15",
		"Referer":      "https://juejin.im/book/5bffcbc9f265da614b11b731/section/5c0374a06fb9a049d37ed783",
		"Content-Type": "application/json;charset=utf8",
	}).ResponseBytes()
	//req := request.HttpRequest{
	//	Url: "https://xiaoce-cache-api-ms.juejin.im/v1/getListSection",
	//}
	//req.JsonContentType().AddHeader("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_5) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/13.1.1 Safari/605.1.15")
	//resp, err := req.AddHeader("Referer", "https://juejin.im/book/5bffcbc9f265da614b11b731/section/5c0374a06fb9a049d37ed783").DoGet(values).ResponseBytes()
	if err != nil {
		logger.Error("Request ListSection Error: %+v", err)
		return nil
	}
	section := ListSection{}
	err = json.Unmarshal(resp, &section)
	if err != nil {
		logger.Error("Format ListSection Error: %+v", err)
		return nil
	}
	return &section
}

func (q *QueryBase) GetSection(sectionId string) *Content {
	values := q.ToUrlValues()
	values.Add("sectionId", sectionId)
	resp, err := request.GetWithHeader("https://xiaoce-cache-api-ms.juejin.im/v1/getListSection", values, request.ReqHeader{
		"User-Agent":   "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_5) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/13.1.1 Safari/605.1.15",
		"Referer":      "https://juejin.im/book/5bffcbc9f265da614b11b731/section/5c0374a06fb9a049d37ed783",
		"Content-Type": "application/json;charset=utf8",
	}).ResponseBytes()
	//req.JsonContentType().AddHeader("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_5) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/13.1.1 Safari/605.1.15")
	//resp, err := req.AddHeader("Referer", fmt.Sprintf("https://juejin.im/book/5bffcbc9f265da614b11b731/section/%s", sectionId)).DoGet(values).ResponseBytes()
	if err != nil {
		logger.Error("Request Section Error: %+v", err)
		return nil
	}
	content := Content{}
	err = json.Unmarshal(resp, &content)
	if err != nil {
		logger.Error("Format Section Error: %+v", err)
		return nil
	}
	return &content
}
