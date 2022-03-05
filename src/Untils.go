package findingplus

import (
	"fmt"
	"strings"
)

func UrlProgress(url string, parm map[string]string) string {
	for _, m := range parm {
		fmt.Println(parm[m])
		url = strings.ReplaceAll(m, fmt.Sprintf(`{%s}`, parm), parm[m])
		fmt.Println(url)
	}
	return url
}
func CheckErr(e error) bool {
	if e != nil {
		fmt.Printf(e.Error())
		return true
	}
	return false
}
