package main

import (
	"github.com/PuerkitoBio/goquery"
	"log"
	"net/http"
)

func goQuery() {
	const nginxDownloadPageUrl = "https://nginx.org/en/download.html"
	resp, err := http.Get(nginxDownloadPageUrl)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	doc.Find("#content").Each(func(i int, s *goquery.Selection) {
		a1 := s.Find("table").Eq(1).Find("a").Eq(1)
		href, exists := a1.Attr("href")
		if exists {
			println(href)
		}
	})
}
