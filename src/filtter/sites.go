//©2022 hydrogenium2020,AGPLv3协议.
package filtter

import (
	"database/sql"
	"fmt"
	"strings"

	"gitee.com/hydrogenium2020/findingplus/src/util"
	"github.com/IGLOU-EU/go-wildcard"
)

//屏蔽表数据
type TableData struct {
	Id      int    `json:"id"`
	Tag     string `json:"tag"`
	Content string `json:"content"`
}

func (t TableData) toStr() string {
	return fmt.Sprintf("Id=%d,Tag=%s,Content=%s", t.Id, t.Tag, t.Content)
}

func (t TableData) isNil() bool {
	return t.Tag == "" || t.Content == ""
}

//添加屏蔽网站
func AddSite(k, tag string, db *sql.DB) bool {
	k = strings.ReplaceAll(k, "https://", "")
	k = strings.ReplaceAll(k, "http://", "")

	c, e := db.Prepare(`INSERT INTO SITES(TAG,CONTENT) VALUES(?,?)`)
	util.ChkErr("添加屏蔽站--数据库准备", e)

	_, e = c.Exec(tag, k)
	util.Hlog(fmt.Sprintf("在%s中添加[%s]屏蔽词", tag, k))
	return !util.ChkErr("添加屏蔽站--数据库执行命令", e)
}

//删除屏蔽网站
func DelSite(k string, db *sql.DB) bool {
	c, e := db.Prepare(`DELETE FROM SITES WHERE CONTENT=?`)
	util.ChkErr("删除屏蔽站--数据库准备", e)

	_, e = c.Exec(k)
	util.Hlog(fmt.Sprintf("删除[%s]屏蔽词", k))
	return !util.ChkErr("删除屏蔽站--数据库准备", e)
}

//获取时间限制记录
func GetSite(tag string, pn, ps int, db *sql.DB) (r []TableData) {
	c := "select Id,Tag,Content from SITES"

	if tag != "" {
		c = c + fmt.Sprintf(` where tag='%s'`, tag)
	}
	if pn != -1 {
		c = c + fmt.Sprintf(` limit %d offset %d`, ps, (pn-1)*ps)
	}

	row, e := db.Query(c)
	util.ChkErr("获取屏蔽词--执行命令", e)

	if row != nil {
		for row.Next() {

			var t TableData
			e := row.Scan(&t.Id, &t.Tag, &t.Content)
			util.ChkErr("获取屏蔽词--读取表", e)
			r = append(r, t)
		}
	}

	return
}

//检查域名是否匹配
func CheckRule(url string, rules []TableData) (match bool, match_rule map[int]string) {
	match_rule = make(map[int]string)
	for _, rule := range rules {
		if rule.isNil() {
			util.Elog(fmt.Sprintf("%s规则参数缺少，少了Id参数或者Tag参数，将被忽略", rule.toStr()))
			continue
		}

		//如果匹配到屏蔽先等等,看有没有忽略规则
		if CheckStr("||", url, rule.Content) {

			match = true
			match_rule[rule.Id] = rule.Content
		}
		//一票否决（
		if CheckStr("@@||", url, rule.Content) && !rule.isNil() {

			match_rule[rule.Id] = rule.Content
			match = false
			return
		}
	}
	return
}

//通配符匹配
//该函数不遵循AGPL-3协议
//import from https://github.com/LearnToRunFast/Leetcode-solutions/blob/master/code/44.%E9%80%9A%E9%85%8D%E7%AC%A6%E5%8C%B9%E9%85%8D.go
func CheckStr(prefix string, str string, rule string) bool {
	if strings.HasPrefix(rule, prefix) {
		rule = strings.TrimPrefix(rule, prefix)
		if strings.HasSuffix(rule, "^") {
			rule = strings.TrimSuffix(rule, "^")
			return strings.Contains(str, rule)
		} else {
			return wildcard.MatchSimple(rule, str)
		}

	} else {
		return false
	}
}
func HasStr(l []string, s string) bool {
	for _, i := range l {
		if strings.Contains(s, i) {
			return true

		}
	}
	return false
}
