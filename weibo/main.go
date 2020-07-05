package main

import (
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/freecracy/gin/systray"
)

const (
	url         = "https://s.weibo.com"
	realTimeHot = "/top/summary?cate=realtimehot"
	userAgent   = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_5) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/83.0.4103.116 Safari/537.36"
)

var wg sync.WaitGroup

type menu struct {
	title string
	tips  string
}

func main() {
	systray.Run(run, nil)
}

func run() {
	systray.SetTitle("▤")
	list := getRealTimeHot()
	bar := make(chan *systray.MenuItem, len(list))
	for k, v := range list {
		wg.Add(1)
		go func(k, v string) {
			m := systray.AddMenuItem(k, v)
			bar <- m
			wg.Done()
			// for {
			// 	<-m.ClickedCh
			// 	open.Run(m.GetTooltip()) // 需要修改systray
			// }
		}(k, v)
	}
	wg.Wait()
	systray.AddSeparator()
	quit := systray.AddMenuItem("退出", "退出")
	quit.SetOnClick(func(s string) {
		systray.Quit()
	})
	// go func(quit *systray.MenuItem) {
	// 	select {
	// 	case <-quit.ClickedCh:
	// 		systray.Quit()
	// 	}
	// }(quit)

	t1 := time.NewTicker(30 * time.Minute)
	go func(t *time.Ticker) {
		for {
			<-t.C
			now2 := getRealTimeHot()
			now1 := make(chan menu, len(now2))
			for k, v := range now2 {
				now1 <- menu{
					title: k,
					tips:  v,
				}
			}
			for i := 0; i < len(bar); i++ {
				l := <-now1
				m := <-bar
				m.SetTitle(l.title)
				m.SetTooltip(l.tips)
				bar <- m
			}
			close(now1)
		}
	}(t1)

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
	if res.StatusCode != http.StatusOK {
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
