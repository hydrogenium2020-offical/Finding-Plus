//©2022 hydrogenium2020,AGPLv3协议.

package net

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path"
	"runtime/debug"
	"strings"

	"gitee.com/hydrogenium2020/findingplus/src/filtter"
	"gitee.com/hydrogenium2020/findingplus/src/search"
	"gitee.com/hydrogenium2020/findingplus/src/setting"
	"gitee.com/hydrogenium2020/findingplus/src/util"
	"github.com/gin-gonic/gin"
)

var port int
var conf_path string

// 开启跨域函数
func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Expose-Headers", "Access-Control-Allow-Origin")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		defer func() {
			if err := recover(); err != nil {
				util.Elog(fmt.Sprintf("Panic info is: %v", err))
				util.Elog(fmt.Sprintf("Panic info is: %s", debug.Stack()))
			}
		}()

		c.Next()
	}
}
func Http() {

	flag.StringVar(&conf_path, "conf", "./findingplus.toml", "端口")
	flag.IntVar(&port, "port", 0, "端口")
	flag.Parse()

	setting.Conf = setting.LoadConf(&conf_path)
	os.Chdir(setting.Conf.WorkDir)
	db := setting.OpenDB(&setting.Conf.WorkDir)

	/*
		gin相关
	*/
	//设置日记
	f, err := os.Create("findingplus.log")
	util.CheckPanic("创建日记文件", err)
	defer f.Close()
	gin.DefaultWriter = io.MultiWriter(f)

	//初始化
	r := gin.Default()

	r.Use(Cors())

	// 加载Html
	r.LoadHTMLGlob(path.Join(setting.Conf.WorkDir, "html/*"))
	//加载静态文件
	r.Static("/static", "static")
	//首页
	r.GET("/", func(c *gin.Context) {
		c.HTML(200, "index.html", "w")
	})
	r.GET("/index", func(c *gin.Context) {
		c.HTML(200, "index.html", "w")
	})

	//搜索
	r.GET("/search", func(c *gin.Context) {
		//请求参数
		q := c.Query("query")
		// 搜索偏移项
		o := util.AtoiDef(c.Query("offset"), 1)

		w, s := filtter.Sync()
		site := []filtter.TableData{}
		for i, r := range s {
			site = append(site, filtter.TableData{Id: i + 1, Content: r, Tag: "Finding+"})
		}

		r := Search(q, o)

		t, n, count := filtter.ChkWord(strings.ReplaceAll(q, " ", ""), &w, 1)
		if t {
			c.String(200, "请不要访问屏蔽内容!!!", n, fmt.Sprintf("count=%d", count))
			return
		}

		l := []string{}
		for i, data := range r {
			r3 := []WebData{}

			for _, r2 := range data.Data {
				m, match := filtter.CheckRule(r2.Link, site)
				if m {
					println("屏蔽", r2.Link)
					println("屏蔽域名规则")
					for i, j := range match {
						println(i, j)
					}
					continue
				}
				m, w, _ = filtter.ChkWord(r2.Title+r2.P, &w, 3)

				if m {
					println("屏蔽内容", r2.P)
					for _, i := range w {
						print(i)
					}
					print("\n")
					//to do addsites

					continue
				}
				l = append(l, r2.Link)
				r3 = append(r3, r2)
			}
			r[i].Data = r3
			//go search.WebCatch(l, db)
		}

		c.HTML(200, "search.tmpl", gin.H{
			"query":  q,
			"result": r,
		})
	})
	/*
		//搜索
		r.GET("/api/search", func(c *gin.Context) {
			//请求参数
			q := c.Query("query")

			w, s := filtter.Sync()
			site := []filtter.TableData{}
			for i, r := range s {
				site = append(site, filtter.TableData{Id: i + 1, Content: r, Tag: "Finding+"})
			}

			if q != "" {
				// 搜索偏移项
				o := util.AtoiDef(q, 1)
				r := Search(q, o)

				for _, data := range r {
					l := []string{}
					for _, r := range data.Data {

						m, _ := filtter.CheckRule(r.Link, site)
						if m {
							continue
						}

						m, _, _ = filtter.ChkWord(r.P, &w, 3)
						if m {
							continue
						}

						l = append(l, r.Link)
					}

					//go search.WebCatch(l, db)
				}

				c.JSON(CODE_SEC, gin.H{
					"code":    CODE_SEC,
					"message": SEC,
					"data":    r,
				})
			} else {

				c.JSON(CODE_SEC, gin.H{
					"code":    CODE_SEC,
					"message": SEC,
					"data":    nil,
				})
			}

		})*/

	//搜索
	r.GET("/api/list", func(c *gin.Context) {

		pn := util.AtoiDef(c.Query("pn"), -1)
		ps := util.AtoiDef(c.Query("ps"), -1)
		r := search.GetWebCatchResult(pn, ps, db)

		c.JSON(CODE_SEC, gin.H{
			"code":    CODE_SEC,
			"message": SEC,
			"data":    r,
		})

	})

	//搜索
	r.GET("/api/catch", func(c *gin.Context) {
		q := c.Query("query")
		o := util.AtoiDef(c.Query("offset"), 1)
		url := c.Query("url")
		t := c.Query("t")
		l := c.Query("l")
		p := c.Query("p")

		r := SearchWithInfo(q, o, EngineInfo{Type: "Text", UA: "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.114 Safari/537.36", Url: url, TextProgress: TextProgress{Title: t, P: p, Link: l}})

		c.HTML(200, "catch.tmpl", gin.H{
			"query":  q,
			"url":    url,
			"offset": o,
			"t":      t,
			"l":      l,
			"p":      p,
			"result": r,
		})
	})

	if port == 0 {
		port = setting.Conf.Port
	}

	url := fmt.Sprintf("%s:%d", setting.Conf.Addr, port)
	util.Hlog(fmt.Sprintf("监听%d端口 %s", port, "http://"+url))
	r.Run(url)
}
