//©2022 hydrogenium2020,AGPLv3协议.
package net

import (
	"bytes"
	"fmt"
	"net/url"
	"os"
	"strings"

	"gitee.com/hydrogenium2020/findingplus/src/util"
	"github.com/BurntSushi/toml"
	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
)

//搜索网页结果
type WebResult struct {
	Name       string
	Offset     int       `json:"offset"`
	LastOffset int       `json:"lastpage"`
	PageNum    int       `json:"pn"`
	Query      string    `json:"query"`
	Data       []WebData `json:"data"`
}

//搜索网页数据
type WebData struct {
	Title string `json:"title"`
	Link  string `json:"link"`
	P     string `json:"p"`
}

//搜索网页数据
type TextProgress struct {
	Title string `json:"title"`
	Link  string `json:"link"`
	P     string `json:"p"`
}

//引擎信息
type EngineInfo struct {
	Name string
	Desc string
	UA   string
	Type string
	Url  string
	//Header       map[string]string
	TextProgress TextProgress
}

func LoadEngine() (w []EngineInfo) {

	d, e := os.ReadDir("./engine/text")
	if !util.ChkErr("加载引擎配置文件", e) {
		for _, de := range d {

			if strings.HasSuffix(de.Name(), ".toml") {
				var www EngineInfo

				toml.DecodeFile("./engine/text/"+de.Name(), &www)
				w = append(w, www)
			}
		}
	}
	return w
}

//搜索
func Search(query string, offset int) (r []WebResult) {

	if query != "" {
		ic := make(chan EngineInfo)
		rc := make(chan WebResult)

		a := LoadEngine()
		for i := range a {

			a[i].Url = strings.ReplaceAll(a[i].Url, `{query}`, url.QueryEscape(query))
			a[i].Url = strings.ReplaceAll(a[i].Url, `{offset}`, fmt.Sprintf("%d", offset))
			go searchSpider(query, offset, ic, rc)
			ic <- a[i]
		}

		for range a {
			r = append(r, <-rc)
		}
	}
	return r
}

//自定义搜索
func SearchWithInfo(query string, offset int, i EngineInfo) (r []WebResult) {

	if query != "" {
		ic := make(chan EngineInfo)
		rc := make(chan WebResult)
		i.Url = strings.ReplaceAll(i.Url, `{query}`, url.QueryEscape(query))
		i.Url = strings.ReplaceAll(i.Url, `{offset}`, fmt.Sprintf("%d", offset))
		go searchSpider(query, offset, ic, rc)
		ic <- i

		r = append(r, <-rc)

	}
	return r
}

//搜索爬虫
func searchSpider(query string, offset int, c chan EngineInfo, r chan WebResult) {
	//接受数据信息
	EngineInfo := <-c

	if !strings.HasPrefix(EngineInfo.Url, "http") {
		EngineInfo.Url = "https://" + EngineInfo.Url

	}
	util.Hlog("开始请求" + EngineInfo.Url)

	co := colly.NewCollector(colly.MaxDepth(5))

	//响应后开始处理数据
	co.OnResponse(func(rsp *colly.Response) {

		// 读取数据
		doc, err := goquery.NewDocumentFromReader(bytes.NewReader(rsp.Body))
		util.ChkErr("读取相应数据", err)

		// 创造返回值
		var rs WebResult
		rs.Offset = offset
		rs.Name = EngineInfo.Name

		//选取HTML元素
		link := doc.Find(EngineInfo.TextProgress.Link) //链接
		p := doc.Find(EngineInfo.TextProgress.P)       //介绍

		for i := 0; i < link.Length(); i++ {
			rs.Offset++
			l, _ := link.Eq(i).Attr("href")

			rs.Data = append(rs.Data, WebData{Title: link.Eq(i).Text(), P: p.Eq(i).Text(), Link: l})
		}

		if a := offset - 10; a >= 0 {
			rs.LastOffset = a

		} else {
			rs.LastOffset = 0
		}

		rs.Query = query

		r <- rs

	})
	if err := co.Visit(EngineInfo.Url); err != nil {
		util.Elog(err)
	}

}
