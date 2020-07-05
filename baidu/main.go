package main

import (
	"log"
	"net/http"

	"github.com/PuerkitoBio/goquery"
	"github.com/axgle/mahonia"
	"github.com/freecracy/gin/systray"
)

const (
	domain    = "http://top.baidu.com/buzz?b=1&fr=topindex"
	uri       = "/buzz?b=1&fr=topindex"
	userAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_5) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/83.0.4103.116 Safari/537.36"
)

func main() {
	systray.Run(run, nil)
}

func run() {
	systray.SetTitle("baidu")
	list := getData()
	menu := make(chan *systray.MenuItem, len(list))
	refresh := systray.AddMenuItem("刷新", "刷新")
	refresh.SetOnClick(func(s string) {
		list = getData()
		for k, v := range list {
			m := <-menu
			m.SetTitle(k)
			m.SetTooltip(v)
			menu <- m
		}
	})
	// go func(r *systray.MenuItem) {
	// 	for {
	// 		<-r.ClickedCh
	// 		list = getData()
	// 		for k, v := range list {
	// 			m := <-menu
	// 			m.SetTitle(k)
	// 			m.SetTooltip(v)
	// 			menu <- m
	// 		}
	// 	}
	// }(refresh)
	systray.AddSeparator()
	for k, v := range list {
		m := systray.AddMenuItem(k, v)
		menu <- m
		// go func(m *systray.MenuItem) {
		// 	for {
		// 		<-m.ClickedCh
		// 		open.Run(m.GetTooltip())
		// 	}
		// }(m)
	}
	systray.AddSeparator()
	quit := systray.AddMenuItem("退出", "退出")
	quit.SetOnClick(func(s string) {
		systray.Quit()
	})
	// go func(q *systray.MenuItem) {
	// 	<-quit.ClickedCh
	// 	systray.Quit()
	// }(quit)
}

func getData() map[string]string {
	ans := make(map[string]string)
	req, err := http.NewRequest("GET", domain+uri, nil)
	if err != nil {
		log.Fatalln(err)
	}
	req.Header.Add("User-Agent", userAgent)
	req.Header.Add("Accept-Charset", "utf-8")
	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalln(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		log.Fatalf("%d : %s \n", resp.StatusCode, resp.Status)
	}
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	doc.Find("tbody").Each(func(_ int, s *goquery.Selection) {
		s.Find("tr").Each(func(i int, m *goquery.Selection) {
			if i == 0 || i == 2 || i == 4 || i == 6 {
				return
			}
			a := m.Find("a").Eq(0)
			title := a.Text()
			href, _ := a.Attr("href")
			ans[convertToString(title, "GBK", "utf-8")] = href
		})
	})
	return ans
}

func convertToString(src string, srcCode string, tagCode string) string {
	srcCoder := mahonia.NewDecoder(srcCode)
	srcResult := srcCoder.ConvertString(src)
	tagCoder := mahonia.NewDecoder(tagCode)
	_, cdata, _ := tagCoder.Translate([]byte(srcResult), true)
	result := string(cdata)
	return result
}
