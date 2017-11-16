package myhtmlparser

import (
	"regexp"
)

var scriptRe = regexp.MustCompile(`(?ms)<script.*?>.*?</script>`)
var tagRe = regexp.MustCompile(`<(/?)([\w-]+)\s?(.*?)>`)
var tagAllAttrsRe = regexp.MustCompile(`\s?([\w-_]+)(=["'].*?["'])?\s?`)
var equationAttrRe = regexp.MustCompile(`=\s?["'](.*?)["']`)
var singleAttrRe = regexp.MustCompile(`[^\s\\]+`)

// var queryRe = regexp.MustCompile(`([>,\s])?([\.#])?([\w_-]*)(?:\[([\w_-]+)=["'](.*?)["']\])?`)
var queryRe = regexp.MustCompile(`([>,\s])?([\.#])?([\w_-]*)(?:\[(.*?)\])?`)
var queryAttrRe = regexp.MustCompile(`([\w-_]+)=?(?:["'](.*?)["'])?,?\s?`)

// var queryAttrRe = regexp.MustCompile(`([\w-]*?)\[(.+?)\]`)
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

func FindAttrs(s string) (map[string]string, error) {
	attrMap := make(map[string]string)
	// allAttrs := tagAttrsRe.FindAllString(s, -1)
	allAttrs := tagAllAttrsRe.FindAllStringSubmatch(s, -1)
	for _, attr := range allAttrs {
		value := ""
		if attr[2] != "" {
			value = equationAttrRe.FindStringSubmatch(attr[2])[1]
		}
		attrMap[attr[1]] = value
	}
	return attrMap, nil
}

func FindQueryAttrs(s string) map[string]string {
	attrMap := make(map[string]string)
	allAttrs := queryAllAttrRe.FindAllString(s, -1)
	for _, attr := range allAttrs {
		if queryEquationAttrRe.MatchString(attr) {
			group := queryEquationAttrRe.FindStringSubmatch(attr)
			attrMap[group[1]] = group[2]
		} else {
			group := querySingleAttrRe.FindStringSubmatch(attr)
			attrMap[group[1]] = ""
		}
	}
	return attrMap
}
