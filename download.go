package main

import (
	"fmt"
	"math/rand"
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
	successCount int
	stopFlag     bool
	myweb        MyWeb
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

	ret := x.task()
	if !ret {
		ch <- false
		return
	}

	x.log("done")

	if x.failCount > 0 {

		x.log("重试1次")

		if !x.stopFlag {
			x.failCount = 0
			x.task()

		}
	}

	x.log(fmt.Sprintf("成功下载图片数量:%d, 失败数量:%d", x.successCount, x.failCount))

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
		OnBeforeRequest(func(client *req.Client, req *req.Request) error {

			// 随机 andorid UA，避免 cf 限流
			ua := fmt.Sprintf("Mozilla/5.0 (Linux; Android 4.4.2; Nexus 4 Build/KOT49H) AppleWebKit/%d.%d (KHTML, like Gecko) Chrome/112.0.0.0 Safari/%d.%d",
				rand.Intn(50)+5, rand.Intn(50)+5, rand.Intn(50)+5, rand.Intn(50)+5)

			// fmt.Println("ua,", ua)

			req.SetHeader("user-agent", ua)
			return nil
		}).
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

func (x *Download) task() bool {

	if strings.Contains(x.Url, "tuwenhanman.com") {
		fmt.Println("use TuWenHanMan")
		x.myweb = TuWenHanMan{}

	} else if strings.Contains(x.Url, "bingmh.com") {
		fmt.Println("use Bingmh")
		x.myweb = Bingmh{}
	} else {
		x.log("此网站暂未支持解析，请联系作者")
		return false
	}

	// summary
	err := x.crawl(x.Url, func(doc *goquery.Document) error {

		fmt.Println("crawl summary ok", x.Url)

		name := x.myweb.ComicName(doc)

		x.log("漫画名称:" + name)

		ret := x.myweb.PageUrl(doc)
		for i, onepage := range ret {
			itemUrl := onepage.Url
			title := onepage.Name
			x.onePage(itemUrl, name, title, i)
		}

		return nil
	})

	x.wg.Wait()

	if err != nil {
		x.log(fmt.Sprint(err))
	}

	return true
}

func (x *Download) onePage(url string, name string, title string, chapter int) {

	if x.stopFlag {
		return
	}

	err := x.crawl(url, func(doc *goquery.Document) error {

		x.log("章节name: " + title)

		ret := x.myweb.ImgList(doc)

		for i, src := range ret {

			if x.stopFlag {
				break
			}

			x.wg.Add(1)
			go x.download(src, name, title, chapter, i)
		}

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
	// fmt.Println(newPath)

	dir := filepath.Dir(newPath)

	err := os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		fmt.Println(err)
	}

	_, err = os.Stat(newPath)
	if os.IsNotExist(err) {
		// File does not exist

		_, err = x.globalClient.R().SetOutputFile(newPath).
			Get(url)

		if err != nil {
			x.log("err, " + newPath)
			os.Remove(newPath)

			x.failCount++
		} else {
			x.successCount++
		}
	}

}
