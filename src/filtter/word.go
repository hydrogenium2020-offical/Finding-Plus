//©2022 hydrogenium2020,AGPLv3协议.
package filtter

import (
	"database/sql"
	"fmt"
	"strings"

	"gitee.com/hydrogenium2020/findingplus/src/util"
)

//检测屏蔽词
func ChkWord(str string, words *[]string, num int) (match bool, match_words []string, count int) {
	for _, word := range *words {

		if i := strings.Count(str, word); i > 0 {
			count += i
			match_words = append(match_words, word)
		}
		if count >= num {
			return true, match_words, count
		}
	}
	return false, nil, -1
}

//添加屏蔽词
func AddWord(k, tag string, db *sql.DB) bool {

	c, e := db.Prepare(`INSERT INTO WORDS(TAG,CONTENT) VALUES(?,?)`)
	util.ChkErr("添加屏蔽词--数据库准备", e)
	_, e = c.Exec(tag, k)
	util.Hlog(fmt.Sprintf("在%s中添加[%s]屏蔽词", tag, k))
	return !util.ChkErr("添加屏蔽词--数据库执行", e)
}

//删除屏蔽词
func DelWord(k string, db *sql.DB) bool {
	c, e := db.Prepare(`DELETE FROM WORDS WHERE CONTENT=?`)
	util.ChkErr("删除屏蔽词--数据库准备", e)

	_, e = c.Exec(k)
	util.Hlog(fmt.Sprintf("删除[%s]屏蔽词", k))
	return !util.ChkErr("删除屏蔽词--数据库执行", e)
}

//获取时间限制记录
func GetWord(tag string, pn, ps int, db *sql.DB) (r []TableData) {
	c := "select Id,Tag,Content from WORDS"

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
