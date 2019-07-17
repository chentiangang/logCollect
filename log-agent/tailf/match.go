package tailf

import (
	"fmt"
	"regexp"

	"github.com/chentiangang/xlog"
)

// const logRe = `^\^\[\[m\^\[\[\d\dm.*`
const colorRe = `(.*)\[\d\d:\d\d:\d\d.\d\d\d\]`

const fmtRe = `.*\[(\d\d:\d\d:\d\d.\d\d\d)\]\[([A-Z]+) \]\[([a-z-]+)\]\[([A-Z]+)\]\[([a-z/]+)\] (.*)`

func FmtLog(re *regexp.Regexp, s string) string {
	result := re.FindAllStringSubmatch(s, -1)
	if len(result) == 0 {
		xlog.LogWarn("match failed, str is nil. src string: %s", s)
		return fmt.Sprintf("{\"time\":\"%s\",\"level\":\"%s\",\"uid\":\"%s\",\"method\":\"%s\",\"apiName\":\"%s\",\"data\":%q,\"ip\":\"%s\"}\n",
			"", "", "", "", "", s, ip)
	}
	submatch := result[0]
	return fmt.Sprintf("{\"time\":\"%s\",\"level\":\"%s\",\"uid\":\"%s\",\"method\":\"%s\",\"apiName\":\"%s\",\"data\":%q,\"ip\":\"%s\"}\n",
		submatch[1], submatch[2], submatch[3], submatch[4], submatch[5], submatch[6], ip)
}
