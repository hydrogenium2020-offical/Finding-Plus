//©2022 hydrogenium2020,AGPLv3协议.
package setting

import (
	"math/rand"
	"os"
	"path"

	"gitee.com/hydrogenium2020/findingplus/src/util"
	"github.com/BurntSushi/toml"
)

var Conf Config

type Config struct {
	//监听端口
	Port int `toml:"port"`
	//监听地址
	Addr string `toml:"addr"`
	//工作目录
	WorkDir string `toml:"work_dir"`
	//凭证
	Token string `toml:"token"`
	//数据库最大记录
	DB_MAX int `toml:"db_max"`
	//默认UA
	DefUA string `toml:"default_ua"`
	//屏蔽词列表地址
	WordUrl string `toml:"words_url"`
	//屏蔽站点地址
	SitesUrl string `toml:"sites_url"`
}

func (c *Config) isNil() bool {
	return c.Port == 0 || c.Addr == "" || c.WorkDir == "" || c.Token == "" || c.DB_MAX == 0 || c.DefUA == "" || c.WordUrl == "" || c.SitesUrl == ""
}

func LoadConf(path *string) (conf Config) {

	toml.DecodeFile(*path, &conf)
	if conf.isNil() {
		util.Elog("配置不完整")
		ReloadConf(path)
	}
	return conf
}

func ReloadConf(cpath *string) {
	util.Elog("正在重新生成配置文件")
	util.Elog("配置文件位于", *cpath)
	SaveConf(cpath, &Config{Port: 7000, Addr: "0.0.0.0", WorkDir: path.Dir(*cpath), Token: RandStr(100), DB_MAX: 1000000, DefUA: "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/101.0.4951.64 Safari/537.36", WordUrl: "ToRelaced", SitesUrl: "ToRelaced"})

	db_path := path.Dir(*cpath)
	db_path = path.Join("findingplus.db")
	InitDB(&db_path)
}

func RandStr(l int) string {
	var k = []rune("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	r := make([]rune, l)
	for i := range r {
		r[i] = k[rand.Intn(len(k))]
	}
	return string(r)
}

func SaveConf(path *string, conf *Config) {
	f, e := os.Create(*path)
	util.ChkErr("写入配置文件", e)

	t := toml.NewEncoder(f)
	e = t.Encode(*conf)

	defer f.Close()
}
