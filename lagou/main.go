package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

const (
	domain    = "https://www.lagou.com"
	userAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_5) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/83.0.4103.116 Safari/537.36"
)

var (
	path     string = "/jobs/list_"
	language string = "PHP"
	city     string = "北京"
	px       string = "new"
)

func main() {
	request()
}

func request() {
	client := &http.Client{}

	value := &url.Values{}
	//value.Set("companyIds", "62,2474,30608,28882,11293,5026,31064,294773,6040,97,82932533,422254,31261,68758,32123")
	req, _ := http.NewRequest("GET", "https://www.lagou.com/c/approve.json", nil)
	req.Header.Add("User-Agent", userAgent)
	resp, _ := client.Do(req)
	cookie := resp.Header.Get("Set-Cookie")
	fmt.Println(cookie)

	value = &url.Values{}
	value.Set("px", "new")
	value.Set("city", "北京")
	value.Set("needAddtionalResult", "false")
	uri, err := url.Parse("https://www.lagou.com/jobs/positionAjax.json")
	if err != nil {
		log.Fatalln(err)
	}
	uri.RawQuery = value.Encode()
	req, err = http.NewRequest("POST", uri.String(), nil)
	if err != nil {
		log.Fatalln(err)
	}
	req.Header.Add("User-Agent", userAgent)
	req.Header.Add("cookie", cookie)
	req.PostForm = url.Values{
		"first": {"true"},
		"kd":    {"PHP"},
		"pn":    {"1"},
	}
	fmt.Printf("%v", req)
	resp, err = client.Do(req)
	if err != nil {
		log.Fatalln(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		log.Fatalf("status code : %d , status : %s", resp.StatusCode, resp.Status)
	}
	body, err := ioutil.ReadAll(resp.Body)
	fmt.Printf("%s", body)
	// {"status":false,"msg":"您操作太频繁,请稍后再访问","clientIp":"","state":2402} 不知道接口校验了啥
}

// ioutil.WriteFile(fmt.Sprintf(`on_exit_%d.txt`, now.UnixNano()), []byte(now.String()), 0644)
