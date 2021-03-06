package notbearparser

import (
	"regexp"
	"strings"
)

var scriptRe = regexp.MustCompile(`(?ms)<script.*?>.*?</script>`)
var commentRe = regexp.MustCompile(`(?ms)<!--.*?-->`)
var styleRe = regexp.MustCompile(`(?ms)<style.*?>.*?</style>`)

var optionAutoCompleteRe = regexp.MustCompile(`(?s)<option.*?>.*?(</option>)?\n`)

var valueQuoteRe = regexp.MustCompile(`\svalue=(\d+)`)
var tagRe = regexp.MustCompile(`<(/?)([\w][\w-]*)\s?([^>]*)>`)
var tagAllAttrsRe = regexp.MustCompile(`\s?([\w-_]+)=?(?:["'](.*?)["'])?\s?`)

var queryRe = regexp.MustCompile(`([>,\s])?([\.#])?([\w_-]*)(?:\[(.*?)\])?`)
var queryAttrRe = regexp.MustCompile(`([\w-_]+)=?(?:["'](.*?)["'])?,?\s?`)

var queryAllAttrRe = regexp.MustCompile(`[^\s,]+`)
var queryEquationAttrRe = regexp.MustCompile(`([^\s]+)=([^\s,]+),?`)
var querySingleAttrRe = regexp.MustCompile(`([^\s,]+),?`)
var queryClassRe = regexp.MustCompile(`^\..*?`)
var queryAllClassRe = regexp.MustCompile(`\.([\w-_]+),?`)

func FindTag(s string) (closeToken, name, attrStr string, index []int, valid bool) {
	tag := tagRe.FindStringSubmatch(s)
	if len(tag) != 0 {
		index = tagRe.FindStringIndex(s)
		closeToken = tag[1]
		name = tag[2]
		attrStr = tag[3]
		valid = true
		return
	}
	return
}

func FindAttrs(s string) (*AttrMap, error) {
	attrMap := NewAttrMap()
	allAttrs := tagAllAttrsRe.FindAllStringSubmatch(s, -1)
	for _, attr := range allAttrs {
		attrName, attrValues := attr[1], attr[2]
		if attrValues != "" {
			attrValueList := strings.Split(attrValues, " ")
			for _, value := range attrValueList {
				attrMap.Add(attrName, value)
			}
		}
	}
	return attrMap, nil
}
