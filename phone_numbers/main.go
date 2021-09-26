package main

import (
	"bytes"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/go-resty/resty/v2"
	"github.com/jmoiron/sqlx"
	"github.com/kataras/golog"
	_ "github.com/lib/pq"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
	"io"
	"regexp"
	"time"
)

var (
	host       = "http://www.haianw.com"
	client     *resty.Client
	phoneRegex *regexp.Regexp
	nameRegex  *regexp.Regexp
	db         *sqlx.DB

	maxPage = 220
	//maxPage = 1
)

func init() {
	client = resty.New()
	client.SetHostURL(host)

	phoneRegex = regexp.MustCompile(`(?m)((13[0-9])|(14[5-9])|(15([0-3]|[5-9]))|(16[6-7])|(17[1-8])|(18[0-9])|(19[1|3])|(19[5|6])|(19[8|9]))\d{8}`)
	nameRegex = regexp.MustCompile("(.先生|.女士)")

	var err error
	db, err = sqlx.Open("postgres", "postgres://charlie:root@127.0.0.1:5432/charlie?sslmode=disable")
	if err != nil {
		golog.Fatal(err)
	}
}

func main() {
	page := 1

	for page <= maxPage {
		requestURL := "/forum-28-%d.html"
		resp, err := client.R().Get(fmt.Sprintf(requestURL, page))
		if err != nil {
			panic(err)
		}

		reader := getReader(resp)
		err = resp.RawBody().Close()
		if err != nil {
			golog.Errorf("关闭body异常", err)
			return
		}

		document, err := goquery.NewDocumentFromReader(reader)
		if err != nil {
			golog.Errorf("Create document from URL[%s] error: %v", requestURL, err)
			return
		}

		document.Find("#threadlisttableid").Each(func(i int, table *goquery.Selection) {
			table.Find("tbody tr th .xst").Each(func(j int, a *goquery.Selection) {
				attr, _ := a.Attr("href")
				//println(attr, exists)
				if attr == "" {
					return
				}

				pageText := getNumber(attr)
				time.Sleep(time.Second)

				names := nameRegex.FindStringSubmatch(pageText)
				var name string
				if len(names) > 0 {
					name = names[0]
				}
				phones := phoneRegex.FindStringSubmatch(pageText)
				if len(phones) > 0 {
					for _, phone := range phones {
						if len(phone) < 11 {
							continue
						}
						var logName string
						if name != "" {
							logName = name + " --> "
						}
						fmt.Printf("[%d-%d]: %s%q\n", page, j, logName, phone)

						insert := true
						row := db.QueryRowx("SELECT 1 FROM phone_numbers WHERE phone = $1", phone)
						if row.Err() == nil {
							var existsPhone string
							_ = row.Scan(&existsPhone)
							if existsPhone == "1" {
								insert = false
							}
						}

						if insert {
							_, err := db.Exec("INSERT INTO phone_numbers(name, phone) VALUES ($1, $2)", name, phone)
							if err != nil {
								golog.Error(err)
							}
						}
					}
				}
			})
		})
		page++
	}
}

func getReader(resp *resty.Response) io.Reader {
	return transform.NewReader(bytes.NewReader(resp.Body()), simplifiedchinese.GBK.NewDecoder())
}

func getNumber(url string) string {
	response, err := client.R().Get(url)
	if err != nil {
		golog.Errorf("获取内容失败, URL:%s", url)
		return ""
	}

	reader := getReader(response)
	document, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		golog.Errorf("格式化内容失败, URL:%s", url)
		return ""
	}

	return document.Find("#postlist").Text()
}
