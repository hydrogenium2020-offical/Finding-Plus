package findingplus

import (
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/gin-gonic/gin"
)

var port int

func Http() {
	flag.IntVar(&port, "port", 7000, "端口")
	flag.Parse()

	//设置日记
	f, err := os.Create("findingplus.log")
	if err != nil {
		fmt.Println(err)
	}
	gin.DefaultWriter = io.MultiWriter(f)

	//初始化
	r := gin.Default()

	// 加载Html
	r.LoadHTMLGlob("./html/*")
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
		q := c.Query("query")
		r := Search(q, "1")
		c.HTML(200, "search.tmpl", gin.H{
			"query":  q,
			"result": r,
		})
	})

	fmt.Println("监听", port, "端口")
	r.Run(fmt.Sprintf(":%d", port))
}
