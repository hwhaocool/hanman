package main

import (
	"fmt"
	"math/rand"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/lxn/walk"

	"github.com/PuerkitoBio/goquery"
	"github.com/imroc/req/v3"
)

type Download struct {
	Url   string
	Path  string
	Proxy string
	OutTE *walk.TextEdit

	globalClient *req.Client
	wg           sync.WaitGroup
	failCount    int
	stopFlag     bool
}

func NewDownload(url string, path string, proxy string, out *walk.TextEdit) *Download {

	x := &Download{
		Url:   url,
		Path:  path,
		Proxy: proxy,
		OutTE: out,
	}

	x.init()

	return x
}

func (x *Download) Start(ch chan bool) {

	x.task()

	x.log("done")

	if x.failCount > 0 {

		x.log(fmt.Sprintf("失败数量:%d", x.failCount))
		x.log("重试1次")

		if !x.stopFlag {
			x.task()
		}
	}

	ch <- true

}

func (x *Download) Stop() {
	x.stopFlag = true
}

func (x *Download) init() {

	x.globalClient = req.C().
		// Enable dump at the request-level for each request, and only
		// temporarily stores the dump content in memory, so we can call
		// resp.Dump() to get the dump content when needed in response
		// middleware.
		EnableDumpEachRequest().
		OnAfterResponse(func(client *req.Client, resp *req.Response) error {
			if resp.Err != nil { // Ignore when there is an underlying error, e.g. network error.
				return nil
			}
			// Treat non-successful responses as errors, record raw dump content in error message.
			if !resp.IsSuccessState() { // Status code is not between 200 and 299.
				resp.Err = fmt.Errorf("bad response, raw content:\n%s", resp.Dump())
			}
			return nil
		})

	if x.Proxy != "" {
		x.globalClient.SetProxyURL(x.Proxy)
	}
}

func (x *Download) crawl(url string, callback func(doc *goquery.Document) error) error {
	// Send request.
	resp, err := x.globalClient.R().Get(url)
	if err != nil {
		return err
	}

	// Pass resp.Body to goquery.
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil { // Append raw dump content to error message if goquery parse failed to help troubleshoot.
		return fmt.Errorf("failed to parse html: %s, raw content:\n%s", err.Error(), resp.Dump())
	}
	err = callback(doc)
	if err != nil {
		err = fmt.Errorf("%s, raw content:\n%s", err.Error(), resp.Dump())
	}
	return err
}

func (x *Download) log(str string) {
	x.OutTE.AppendText(str + "\r\n")
}

func (x *Download) task() {

	y, _ := url.Parse(x.Url)

	// summary

	err := x.crawl(x.Url, func(doc *goquery.Document) error {

		name := doc.Find("h1.fed-part-eone.fed-font-xvi a").First().Text()

		x.log("漫画名称:" + name)

		doc.Find("a.fed-btns-info.fed-rims-info.fed-part-eone").Each(func(i int, s *goquery.Selection) {
			href, _ := s.Attr("href")

			itemUrl := y.Scheme + "://" + y.Host + href

			// fmt.Println(itemUrl)

			x.onePage(itemUrl, name, i)

		})
		return nil
	})

	x.wg.Wait()

	if err != nil {
		x.log(fmt.Sprint(err))
	}
}

func (x *Download) onePage(url string, name string, chapter int) {

	if x.stopFlag {
		return
	}

	err := x.crawl(url, func(doc *goquery.Document) error {

		title := doc.Find("h2").First().Text()

		x.log("章节name: " + title)

		doc.Find("img").Each(func(i int, s *goquery.Selection) {

			if x.stopFlag {
				return
			}

			src, _ := s.Attr("src")

			if strings.HasPrefix(src, "http") {

				x.wg.Add(1)
				go x.download(src, name, title, chapter, i)
			}

		})

		return nil
	})

	if err != nil {
		x.log(fmt.Sprint(err))
	}
}

func (x *Download) download(url string, name string, title string, chapter int, index int) {

	defer x.wg.Done()

	ext := filepath.Ext(url)

	filename := fmt.Sprintf("%03d-%s-%03d%s", chapter, strings.ReplaceAll(title, " ", "_"), index, ext)

	newPath := filepath.Join(x.Path, name, filename)
	fmt.Println(newPath)

	dir := filepath.Dir(newPath)

	err := os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		fmt.Println(err)
	}

	_, err = os.Stat(newPath)
	if os.IsNotExist(err) {
		// File does not exist

		// 随机UA，避免 cf 限流
		ua := fmt.Sprintf("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/%d.%d (KHTML, like Gecko) Chrome/112.0.0.0 Safari/%d.%d",
			rand.Intn(50)+5, rand.Intn(50)+5, rand.Intn(50)+5, rand.Intn(50)+5)

		_, err = x.globalClient.R().SetOutputFile(newPath).
			SetHeader("user-agent", ua).
			Get(url)

		if err != nil {
			x.log("err, " + newPath)
			os.Remove(newPath)

			x.failCount++
		}
	}

}
