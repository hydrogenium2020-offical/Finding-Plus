package findingplus

import (
	"fmt"
	"os"
	"strings"

	"github.com/BurntSushi/toml"
)

func LoadEngine() (w []Info) {

	d, e := os.ReadDir("./engine/text")
	if !CheckErr(e) {
		for _, de := range d {

			if strings.HasSuffix(de.Name(), ".toml") {
				var www Info
				toml.DecodeFile("./engine/text/"+de.Name(), &www)
				fmt.Println(www)
				w = append(w, www)
			}
		}
	}
	return w
}

type Config struct {
	Port string
}

type Info struct {
	Name string
	Desc string
	UA   string
	Type string
	Url  string
	//Header       map[string]string
	TextProgress TextProgress
}
type TextProgress struct {
	Title string
	Link  string
	P     string
}
