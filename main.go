package main

import (
	"log"
	"net/http"

	"github.com/PuerkitoBio/goquery"
	"github.com/johnmccabe/go-bitbar"
)

const (
	url         = "https://s.weibo.com"
	realTimeHot = "/top/summary?cate=realtimehot"
	userAgent   = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_5) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/83.0.4103.116 Safari/537.36"
)

func main() {
	app := bitbar.New()
	app.StatusLine("热搜")
	submenu := app.NewSubMenu()
	list := getRealTimeHot()
	for k, v := range list {
		submenu.Line(k).Href(v)
	}
	app.Render()
}

func getRealTimeHot() map[string]string {
	list := make(map[string]string)
	client := &http.Client{}
	req, err := http.NewRequest("GET", url+realTimeHot, nil)
	if err != nil {
		log.Fatalln("not found", err)
	}
	req.Header.Add("User-Agent", userAgent)
	res, err := client.Do(req)
	//res, err = http.Get(url + realTimeHot)
	if err != nil {
		log.Fatalln("not found", err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Fatalf("there is a error %d, %s \n", res.StatusCode, res.Status)
	}
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatalln("doc error", err)
	}
	doc.Find("tr").Each(func(i int, s *goquery.Selection) {
		if i == 0 {
			return
		}
		title := s.Find("a").Text()
		path, _ := s.Find("a").Attr("href")
		uri := url + path
		list[title] = uri
	})
	return list
}
