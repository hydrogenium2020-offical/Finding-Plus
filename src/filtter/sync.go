//©2022 hydrogenium2020,AGPLv3协议.
package filtter

import (
	"net/http"
	"strings"

	"gitee.com/hydrogenium2020/findingplus/src/setting"
	"gitee.com/hydrogenium2020/findingplus/src/util"
	"github.com/PuerkitoBio/goquery"
)

func Sync() (word, sites []string) {
	rsp, err := http.Get(setting.Conf.WordUrl)
	util.ChkErr("同步屏蔽列表", err)
	doc, err := goquery.NewDocumentFromReader(rsp.Body)
	util.ChkErr("读取同步响应内容", err)

	for _, r := range strings.Split(doc.Text(), "\n") {
		if r != "" {
			word = append(word, r)
		}
	}

	rsp, err = http.Get(setting.Conf.SitesUrl)
	util.ChkErr("同步屏蔽列表", err)
	doc, err = goquery.NewDocumentFromReader(rsp.Body)
	util.ChkErr("读取同步响应内容", err)

	for _, r := range strings.Split(doc.Text(), "\n") {
		if r != "" {
			sites = append(sites, r)

		}
	}

	return
}
