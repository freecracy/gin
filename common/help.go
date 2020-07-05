package common

import (
	"io"
	"log"
	"net/http"
)

// CurlPost ...
func CurlPost(s string) io.ReadCloser {
	req, err := http.NewRequest("GET", s, nil)
	if err != nil {
		log.Fatalln(err)
	}
	req.Header.Add("User-Agent", UserAgent)
	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalln(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		log.Fatalf("%d : %s \n", resp.StatusCode, resp.Status)
	}
	return resp.Body
}
