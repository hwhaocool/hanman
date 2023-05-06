package main

import "github.com/PuerkitoBio/goquery"

type MyWeb interface {
	// 漫画名称
	ComicName(doc *goquery.Document) string

	// 子页面url和标题
	PageUrl(doc *goquery.Document) []OnePage

	// 图片url
	ImgList(doc *goquery.Document) []string
}

type OnePage struct {
	Url  string
	Name string
}
