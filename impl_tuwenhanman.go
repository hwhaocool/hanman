package main

import (
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// tuwenhuanman 网站，实现了 MyWeb 接口

type TuWenHanMan struct {
}

func (x TuWenHanMan) ComicName(doc *goquery.Document) string {

	return doc.Find("h1.fed-part-eone.fed-font-xvi a").First().Text()
}

func (x TuWenHanMan) PageUrl(doc *goquery.Document) []OnePage {

	var ret []OnePage

	doc.Find("a.fed-btns-info.fed-rims-info.fed-part-eone").Each(func(i int, s *goquery.Selection) {
		href, _ := s.Attr("href")

		itemUrl := "https://www.tuwenhanman.com" + href

		ret = append(ret, OnePage{
			Url:  itemUrl,
			Name: s.Text(),
		})
	})

	return ret
}

func (x TuWenHanMan) ImgList(doc *goquery.Document) []string {

	var ret []string

	doc.Find("img").Each(func(i int, s *goquery.Selection) {

		src, _ := s.Attr("src")

		if strings.HasPrefix(src, "http") {

			ret = append(ret, src)
		}

	})

	return ret
}
