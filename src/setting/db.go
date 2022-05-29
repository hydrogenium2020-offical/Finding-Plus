//©2022 hydrogenium2020,AGPLv3协议.
package setting

import (
	"database/sql"
	"path"

	"gitee.com/hydrogenium2020/findingplus/src/util"
	_ "github.com/mattn/go-sqlite3"
)

//打开数据库
func OpenDB(p *string) *sql.DB {
	db_path := path.Dir(*p)
	db_path = path.Join("findingplus.db")

	db, e := sql.Open("sqlite3", db_path)
	util.ChkErr("打开数据库", e)

	//以WAL模式运行以提高并发性能
	_, e = db.Exec("PRAGMA journal_mode=WAL")
	util.ChkErr("切换数据库为WAL模式", e)

	return db
}

//初始化数据库
func InitDB(cpath *string) (s bool) {
	db, e := sql.Open("sqlite3", *cpath)
	util.ChkErr("打开数据库", e)

	var t []string = []string{}
	/*
		创建表指令
	*/
	//屏蔽网站
	t = append(t, `CREATE TABLE SITES(
		Id INTEGER PRIMARY KEY AUTOINCREMENT,
		Tag        CHAR(20)     NOT NULL,
		Content        CHAR(150)     NOT NULL
	 );`)
	//屏蔽词
	t = append(t, `CREATE TABLE WORDS(
		Id INTEGER PRIMARY KEY AUTOINCREMENT,
		Tag        CHAR(20)     NOT NULL,
		Content        CHAR(150)     NOT NULL
	 );`)
	//搜索结果(网页爬虫)
	t = append(t, `CREATE TABLE WEB_CATCH(
		Id INTEGER PRIMARY KEY AUTOINCREMENT,
		Date INT NOT NULL,
		Title        CHAR(150)     NOT NULL,
		Link        CHAR(150)     NOT NULL,
		Content        CHAR(1000)     NOT NULL
	 );`)

	for _, i := range t {
		_, e := db.Exec(i)
		util.ChkErr("创建数据库\n指令为"+i, e)
	}

	defer db.Close()
	return s
}
