//©2022 hydrogenium2020,AGPLv3协议.
package util

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

func UrlProgress(url string, parm map[string]string) string {
	for _, m := range parm {
		Hlog(parm[m])
		url = strings.ReplaceAll(m, fmt.Sprintf(`{%s}`, parm), parm[m])
		Hlog(url)
	}
	return url
}

// 判断错误
func ChkErr(prefix string, e error) bool {
	if e != nil {
		Elog(prefix)
		Elog(e.Error())
		return true
	}
	return false
}

// 判断错误
func CheckPanic(prefix string, e error) {
	if e != nil {
		Elog(prefix)
		Elog(e.Error())
		panic(-1)
	}
}

//正常日记
func Hlog(i ...any) {
	log.SetPrefix("[Finding+]")
	for _, i := range i {
		log.Println(i)
	}

}

//错误日记
func Elog(i ...any) {
	log.SetPrefix("!!!--->[Error]")
	for _, i := range i {
		log.Println(i)
	}
}

//str转换为int,带默认值
func AtoiDef(str string, def int) (i int) {
	if it, err := strconv.Atoi(str); err != nil {
		i = def
	} else {
		i = it
	}
	return
}

//str转换为int,带默认值
func QueryDef(c *gin.Context, key, def string) (i string) {
	if it := c.Query(key); it == "" {
		i = def
	} else {
		i = it
	}
	return
}
