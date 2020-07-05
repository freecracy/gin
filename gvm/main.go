package main

import (
	"archive/zip"
	"crypto/sha256"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/user"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

const (
	dl  = "https://github.com/golang/dl/archive/master.zip"
	gvm = ".gvm2"
)

func main() {
	downloaddl()
}

func verifySHA256() string {
	homedir, _ := homeDir()
	f, err := os.Open(filepath.Join(homedir, gvm, "dl", "dl.zip"))
	if err != nil {
		log.Fatalln(err)
	}
	defer f.Close()
	hash := sha256.New()
	if _, err := io.Copy(hash, f); err != nil {
		log.Fatalln(err)
	}
	return fmt.Sprintf("%x", hash.Sum(nil))
}

type userAgentTransport struct {
	rt http.RoundTripper
}

func (uat *userAgentTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	version := runtime.Version()
	if strings.Contains(version, "devel") {
		version = "devel"
	}
	r.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_5) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/83.0.4103.116 Safari/537.36")
	return uat.rt.RoundTrip(r)
}

type progressWriter struct {
	w     io.Writer
	n     int64
	total int64
	last  time.Time
}

func (p *progressWriter) update() {
	end := "..."
	if p.n == p.total {
		end = ""
	}
	f := func(i int64) int {
		var n int
		for ; i != 0; i /= 10 {
			n++
		}
		return n
	}
	fmt.Fprintf(os.Stderr, "\rdownload %5.1f%% (%*d / %d bytes)%s",
		(100.0*float64(p.n))/float64(p.total),
		f(p.total), p.n, p.total, end) // \r 和 \n 区别 ???
}

func f(i int64) int {
	var n int
	for ; i != 0; i /= 10 {
		n++
	}
	return n
}

func (p *progressWriter) Write(buf []byte) (n int, err error) {
	n, err = p.w.Write(buf)
	p.n += int64(n)
	if now := time.Now(); now.Unix() != p.last.Unix() {
		p.update()
		p.last = now
	}
	return
}

func downloaddl() {
	homedir, err := homeDir()
	if err != nil {
		log.Fatalln(err)
	}
	dldir := filepath.Join(homedir, gvm, "dl")
	if _, err := os.Stat(dldir); err != nil {
		os.MkdirAll(dldir, 0755)
	}
	if _, err := os.Stat(filepath.Join(dldir, "dl.zip")); err == nil {
		return
	}
	c := &http.Client{
		Transport: &userAgentTransport{&http.Transport{
			DisableCompression: true,
			DisableKeepAlives:  true,
			Proxy:              http.ProxyFromEnvironment,
		}},
	}
	resp, err := c.Get(dl)
	if err != nil {
		log.Fatalln(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusNotFound {
		log.Fatalf("%d : %s \n", resp.StatusCode, resp.Status)
	}
	if resp.StatusCode != http.StatusOK {
		log.Fatalf("%d : %s \n", resp.StatusCode, resp.Status)
	}
	f, err := os.Create(filepath.Join(dldir, "dl.zip"))
	if err != nil {
		log.Fatalln(err)
	}
	pw := &progressWriter{w: f, total: resp.ContentLength}
	n, err := io.Copy(pw, resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	if resp.ContentLength != -1 && resp.ContentLength != n {
		log.Fatalf("copied %d, expected %d bytes", n, resp.ContentLength)
	}
	pw.update()
	f.Close()
}

func homeDir() (string, error) {
	if dir := os.Getenv("HOME"); dir != "" {
		return dir, nil
	}
	if u, err := user.Current(); err != nil && u.HomeDir != "" {
		return u.HomeDir, nil
	}
	return "", errors.New("用户目录不存在")
}

func unpackArchive(targetDir, archiveFile string) error {
	switch {
	case strings.HasSuffix(archiveFile, ".zip"):
		return nil
	case strings.HasSuffix(archiveFile, ".tar.gz"):
		return nil
	default:
		return errors.New(".zip or .tar.gz")
	}
}

func unpackZip(targetDir, archiveFile string) error {
	f, err := zip.OpenReader(archiveFile)
	if err != nil {
		return err
	}
	defer f.Close()
	for _, f := range f.File {
		name := strings.TrimPrefix(f.Name, "dl-master/")
		out := filepath.Join(targetDir, name)
		if f.FileInfo().IsDir() {
			if err := os.MkdirAll(out, 0755); err != nil {
				return nil
			}
			continue
		}
		old, err := f.Open()
		if err != nil {
			return err
		}
		if err := os.MkdirAll(
			filepath.Join(filepath.Dir(out)),
			0755); err != nil {
			return err
		}
		new, err := os.OpenFile(out, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return err
		}
		_, err = io.Copy(new, old)
		old.Close()
		if err != nil {
			new.Close()
			return err
		}
		if err := new.Close(); err != nil {
			return err
		}
	}
	return nil
}
