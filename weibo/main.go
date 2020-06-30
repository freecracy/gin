package main

import (
	"log"
	"net/http"
	"sync"
	"sync/atomic"

	"github.com/PuerkitoBio/goquery"
	"github.com/getlantern/systray"
	"github.com/skratchdot/open-golang/open"
)

const (
	url         = "https://s.weibo.com"
	realTimeHot = "/top/summary?cate=realtimehot"
	userAgent   = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_5) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/83.0.4103.116 Safari/537.36"
)

var wg sync.WaitGroup
var l sync.Mutex

func main() {
	systray.Run(menu, nil)
}

func menu() {
	systray.SetTitle("热搜")
	list := getRealTimeHot()
	n := make(map[int32]*systray.MenuItem)
	for k, v := range list {
		wg.Add(1)
		go func(k, v string) {
			m := systray.AddMenuItem(k, v)
			l.Lock()
			n[m.GetId()] = m
			l.Unlock()
			wg.Done()
			<-m.ClickedCh
			open.Run(v)
		}(k, v)
	}
	wg.Wait()
	systray.AddSeparator()
	refresh := systray.AddMenuItem("刷新", "刷新")
	quit := systray.AddMenuItem("退出", "退出")
	for {
		select {
		case <-refresh.ClickedCh:
			update(n)
		case <-quit.ClickedCh:
			systray.Quit()
		}
	}

	// go func() {
	// 	refresh := systray.AddMenuItem("刷新", "刷新")
	// 	<-refresh.ClickedCh
	// 	update(n)
	// }()
	// go func() {
	// 	quit := systray.AddMenuItem("退出", "退出")
	// 	<-quit.ClickedCh
	// 	systray.Quit()
	// }()
}

func update(n map[int32]*systray.MenuItem) {
	new := getRealTimeHot()
	i := int32(0)
	for k, v := range new {
		l.Lock()
		n[i].SetTitle(k)
		n[i].SetTooltip(v)
		<-n[i].ClickedCh
		atomic.AddInt32(&i, 1)
		l.Unlock()
	}
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

/** 修改systray.go
func (item *MenuItem) GetId() int32 {
	return item.id
}
**/
