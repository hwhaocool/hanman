package main

import (
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// https://bingmh.com/ 网站，实现了 MyWeb 接口

type Bingmh struct {
}

func (x Bingmh) ComicName(doc *goquery.Document) string {

	return doc.Find("h1.comic-name").First().Text()
}

func (x Bingmh) PageUrl(doc *goquery.Document) []OnePage {

	var ret []OnePage

	doc.Find("li.chapter-item a").Each(func(i int, s *goquery.Selection) {
		href, _ := s.Attr("href")

		itemUrl := "https://bingmh.com" + href

		ret = append(ret, OnePage{
			Url:  itemUrl,
			Name: s.Text(),
		})
	})

	// 逆序， 代码来自 chatgpt
	// for i := len(ret)/2 - 1; i >= 0; i-- {
	// 	opp := len(ret) - 1 - i
	// 	ret[i], ret[opp] = ret[opp], ret[i]
	// }

	return ret
}

func (x Bingmh) ImgList(doc *goquery.Document) []string {

	var ret []string

	doc.Find("li.comic-page img").Each(func(i int, s *goquery.Selection) {

		src, _ := s.Attr("src")

		if strings.HasPrefix(src, "http") {

			ret = append(ret, src)
		}

	})

	return ret
}
