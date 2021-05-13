package main

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/kataras/golog"
	"net/http"
	"strconv"
	"strings"
)

type FalconAndWinterSoldier struct {
	Updater
}

func (fws FalconAndWinterSoldier) Update() {
	resp, err := http.Get(fws.URL)
	logger := golog.Child(fws.Name)
	if err != nil {
		logger.Errorf("Can't fetch from URL[%s]", fws.URL)
		return
	}

	defer resp.Body.Close()
	document, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		logger.Errorf("Create document from URL[%s] error: %v", fws.URL, err)
		return
	}

	var max int
	document.Find(".container .play-content .play-item").Each(func(i int, div *goquery.Selection) {
		div.Find("li").Each(func(i int, a *goquery.Selection) {
			text := a.Find("a").Text()
			text = strings.ReplaceAll(text, "第", "")
			text = strings.ReplaceAll(text, "集", "")
			text = strings.ReplaceAll(text, "合", "")
			logger.Info(text)

			if text != "" {
				num, err := strconv.Atoi(text)
				if err != nil {
					golog.Errorf("Can't parse number from text[%div], error: %v", text, err)
					return
				}
				if num > max {
					max = num
				}
			}
		})
	})

	UpdateCache(fws.Name, max)
	logger.Info(max)
}
