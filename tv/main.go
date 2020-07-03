package main

import (
	"log"
	"net/http"

	"github.com/PuerkitoBio/goquery"
	"github.com/getlantern/systray"
	"github.com/skratchdot/open-golang/open"
)

const (
	qqURL     = "https://v.qq.com/channel/tv"
	iqiyiURL  = "https://www.iqiyi.com/dianshiju/"
	userAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_5) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/83.0.4103.116 Safari/537.36"
)

func main() {
	systray.Run(Run, nil)
}

func Run() {
	systray.SetTitle("热播")
	systray.AddMenuItem("腾讯|爱奇艺", "腾讯|爱奇艺")
	systray.AddSeparator()
	list := getQQData()
	addMenuItem(list)
	systray.AddSeparator()
	list = getIqiyi()
	addMenuItem(list)
	systray.AddSeparator()
	quit := systray.AddMenuItem("退出", "退出")
	<-quit.ClickedCh
	systray.Quit()
}

func addMenuItem(list map[string]string) {
	for k, v := range list {
		m := systray.AddMenuItem(k, v)
		go func(m *systray.MenuItem) {
			for {
				<-m.ClickedCh
				open.Run(m.GetTooltip())
			}
		}(m)
	}
}

func getIqiyi() map[string]string {
	ans := make(map[string]string)
	doc := curlPost(iqiyiURL)
	doc.Find(".focus_title_inner").Each(func(i int, s *goquery.Selection) {
		title := s.Find(".caption").Text()
		if title == "" {
			return
		}
		href, _ := s.Attr("href")
		ans[title] = "https:" + href
	})
	return ans
}

func getQQData() map[string]string {
	ans := make(map[string]string)
	doc := curlPost(qqURL)
	doc.Find(".slider_figure_inner").Each(func(_ int, s *goquery.Selection) {
		s.Find("a").Each(func(i int, m *goquery.Selection) {
			href, _ := m.Attr("href")
			title := m.Find(".slider_figure_title ").Text()
			ans[title] = href
		})
	})
	return ans
}

func curlPost(url string) *goquery.Document {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatalln(err)
	}
	req.Header.Add("User-Agent", userAgent)
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
	return doc
}
