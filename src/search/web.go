//©2022 hydrogenium2020,AGPLv3协议.
package search

import (
	"bytes"
	"database/sql"
	"fmt"
	"os"
	"strings"
	"time"

	"gitee.com/hydrogenium2020/findingplus/src/filtter"
	"gitee.com/hydrogenium2020/findingplus/src/util"
	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
	"github.com/gocolly/colly/debug"
)

type WebCatchResult struct {
	Id      int64  `json:"id"`
	Date    int64  `json:"data"`
	Title   string `json:"title"`
	Link    string `json:"link"`
	Content string `json:"content"`
}

//获取时间限制记录
func GetWebCatchResult(pn, ps int, db *sql.DB) (r []WebCatchResult) {
	c := "select Id,Date,Title,Link,Content from WEB_CATCH"

	if pn != -1 {
		c = c + fmt.Sprintf(` limit %d offset %d`, ps, (pn-1)*ps)
	}

	row, e := db.Query(c)
	util.ChkErr("查询网页抓取技术", e)

	if row != nil {
		for row.Next() {

			var t WebCatchResult
			e := row.Scan(&t.Id, &t.Date, &t.Title, &t.Link, &t.Content)
			util.ChkErr("查询网页抓取技术", e)
			r = append(r, t)
		}
	}

	return
}

func (r *WebCatchResult) DbAdd(db *sql.DB) (ok bool) {
	r.Date = time.Now().Unix()
	c, e := db.Prepare(`INSERT INTO WEB_CATCH(Date,Title,Link,Content) VALUES(?,?,?,?)`)
	if s := util.ChkErr("添加网页抓取记录--数据库准备", e); s {
		return s
	}
	_, e = c.Exec(r.Date, r.Title, r.Link, r.Content)

	return !util.ChkErr("添加网页抓取记录--数据库准备", e)
}

func WebCatch(urls []string, db *sql.DB) (result string) {
	if len(urls) > 0 {

		rc := make(chan WebCatchResult)
		close := make(chan int)
		for i, url := range urls {
			go WebSpider(url, i, rc, close)

		}

		i := 0
		for {
			select {
			case r := <-rc:
				//r.DbAdd(db)
				strings.Count(r.Content, "")
			case <-close:
				i++
				if i >= len(urls) {
					break
				}
			}

		}
	}
	return result
}
func WebSpider(url string, i int, r chan WebCatchResult, close chan int) {
	util.Hlog(fmt.Sprintf("正在爬取%d URL=%s", i, url))

	if !strings.HasPrefix(url, "http") {
		url = "https://" + url
	}
	writer, err := os.OpenFile(fmt.Sprintf("spider-%d.log", i), os.O_RDWR|os.O_CREATE, 0666)
	util.ChkErr(fmt.Sprintf("写入/读取文件spider-%d.log", i), err)

	c := colly.NewCollector(colly.Debugger(&debug.LogDebugger{Output: writer}), colly.MaxDepth(3))

	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		e.Request.Visit(e.Attr("href"))
	})
	c.OnResponse(func(rsp *colly.Response) {
		w, _ := filtter.Sync()
		// 读取数据
		doc, err := goquery.NewDocumentFromReader(bytes.NewReader(rsp.Body))
		util.ChkErr("读取相应数据", err)

		if rsp != nil {
			ctx := strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(doc.Find("p").Text(), " ", ""), "\t", ""), "\n", "")
			link := rsp.Request.URL.String()
			title := doc.Find("title").Text()

			m, w1, _ := filtter.ChkWord(ctx, &w, 5)
			if m {
				println("屏蔽内容", ctx)
				//屏蔽词
				for _, i := range w1 {
					print(i)
				}
				print("\n")
				//to do addsites
			}

			res := WebCatchResult{Content: ctx, Link: link, Title: title}
			r <- res
			close <- 200
		} else {
			close <- -1
		}

	})

	if err := c.Visit(url); err != nil {
		util.ChkErr("爬取数据"+url, err)
	}
}
