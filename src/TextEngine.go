package findingplus

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func Search(query, offset string) (r []TextResult) {

	if query != "" {
		ic := make(chan Info)
		rc := make(chan []TextResult)

		a := LoadEngine()
		for i := range a {

			a[i].Url = strings.ReplaceAll(a[i].Url, `{query}`, url.QueryEscape(query))
			a[i].Url = strings.ReplaceAll(a[i].Url, `{offset}`, url.QueryEscape(offset))
			go spider(ic, rc)
			ic <- a[i]
		}

		for range a {
			r = append(r, <-rc...)
		}
	}
	return r
}
func spider(c chan Info, r chan []TextResult) {
	//接受数据信息
	info := <-c

	rsp, err := http.Get(info.Url)
	rsp.Header.Add("User-agent", info.UA)

	if err != nil {
		panic(err)
	}
	fmt.Println("开始请求", info.Url)
	CheckErr(err)

	if rsp.StatusCode != 200 {
		fmt.Println("Error", rsp.StatusCode, rsp.Body)
	}

	doc, err := goquery.NewDocumentFromReader(rsp.Body)
	CheckErr(err)

	fmt.Println(doc.Text())

	var rs []TextResult

	link := doc.Find(info.TextProgress.Link)
	p := doc.Find(info.TextProgress.P)

	for i := 0; i < link.Length(); i++ {
		l, _ := link.Eq(i).Attr("href")
		rs = append(rs, TextResult{Title: link.Eq(i).Text(), P: p.Eq(i).Text(), Link: l})
	}

	r <- rs
}

type TextResult struct {
	Title string `json:"title"`
	Link  string `json:"link"`
	P     string `json:"p"`
}
