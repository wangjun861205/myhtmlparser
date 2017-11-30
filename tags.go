package notbearparser

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

type TAG_TYPE int

const (
	START_TAG TAG_TYPE = iota
	END_TAG
	VOID_START_TAG
	VOID_END_TAG
	INVALID_TAG
)

var VoidEleMap = map[string]bool{
	"area":    true,
	"base":    true,
	"br":      true,
	"col":     true,
	"command": true,
	"embed":   true,
	"hr":      true,
	"img":     true,
	"input":   true,
	"keygen":  true,
	"link":    true,
	"meta":    true,
	"param":   true,
	"source":  true,
	"track":   true,
	"wbr":     true,
}

var AutoCompleteReMap = map[string]*regexp.Regexp{
	"option": optionAutoCompleteRe,
}

var EscapeCharMap = map[string]string{
	"&nbsp;": " ",
	"&amp;":  "&",
	"&lt;":   "<",
	"&gt;":   ">",
	"&quot;": `""`,
	"&#39;":  "'",
}

var TypeMap = map[string]func(string) TAG_TYPE{
	"": func(tagName string) TAG_TYPE {
		if _, ok := VoidEleMap[tagName]; ok {
			return VOID_START_TAG
		} else {
			return START_TAG
		}
	},
	"/": func(tagName string) TAG_TYPE {
		if _, ok := VoidEleMap[tagName]; ok {
			return INVALID_TAG
		} else {
			return END_TAG
		}
	},
}

func CheckTagType(closeToken, tagName string) TAG_TYPE {
	return TypeMap[closeToken](tagName)
}

var NoValidElementErr = errors.New("No valid HTML element")
var NegativeDepthErr = errors.New("Tag depth must be positive number")

type Tag struct {
	Name     string
	Attrs    *AttrMap
	Position []int
	Type     TAG_TYPE
}

func (tag *Tag) String() string {
	return fmt.Sprintf("tag name: %s, tag attribute map: %v, tag pisition: %d-%d, tag type: %d",
		tag.Name, tag.Attrs, tag.Position[0], tag.Position[1], tag.Type)
}

type TagHandler struct {
	TagList TagList
	HTML    string
}

func NewTagHandler(html string) *TagHandler {
	return &TagHandler{
		HTML: html,
	}
}

func (tg *TagHandler) ClearScript() {
	cleanHTML := scriptRe.ReplaceAllString(tg.HTML, "")
	tg.HTML = cleanHTML
}

func (tg *TagHandler) ClearComment() {
	cleanHTML := commentRe.ReplaceAllString(tg.HTML, "")
	tg.HTML = cleanHTML
}

func (tg *TagHandler) ClearStyle() {
	cleanHTML := styleRe.ReplaceAllString(tg.HTML, "")
	tg.HTML = cleanHTML
}

func (tg *TagHandler) AutoComplete() {
	for tagName, re := range AutoCompleteReMap {
		newHTML := re.ReplaceAllStringFunc(tg.HTML, func(s string) string {
			s = strings.Trim(s, "\n\r ")
			if NBString(s).EndsWith(fmt.Sprintf("</%s>", tagName)) {
				return s + "\n"
			} else {
				return fmt.Sprintf("%s%s\n", s, fmt.Sprintf("</%s>", tagName))
			}
		})
		tg.HTML = newHTML
	}
}

func (tg *TagHandler) AddValueQuote() {
	newHTML := valueQuoteRe.ReplaceAllStringFunc(tg.HTML, func(arg1 string) string {
		value := valueQuoteRe.FindStringSubmatch(arg1)[1]
		return fmt.Sprintf(` value="%s"`, value)
	})
	tg.HTML = newHTML
}

func (tg *TagHandler) Escape() {
	newHTML := tg.HTML
	for k, v := range EscapeCharMap {
		newHTML = strings.Replace(newHTML, k, v, -1)
	}
	tg.HTML = newHTML
}

func (tg *TagHandler) Feed() error {
	tg.ClearScript()
	tg.ClearStyle()
	tg.ClearComment()
	tg.AutoComplete()
	tg.AddValueQuote()
	tg.Escape()
	// f, _ := os.OpenFile("processed_html.html", os.O_CREATE|os.O_WRONLY, 0664)
	// f.WriteString(tg.HTML)
	allTagsStr := tagRe.FindAllStringSubmatch(tg.HTML, -1)
	allTagsIndex := tagRe.FindAllStringIndex(tg.HTML, -1)
	for i, tagStr := range allTagsStr {
		closeToken := tagStr[1]
		tagName := strings.ToLower(tagStr[2])
		tagAttrStr := tagStr[3]
		switch CheckTagType(closeToken, tagName) {
		case VOID_START_TAG:
			attrMap, err := FindAttrs(tagAttrStr)
			if err != nil {
				return err
			}
			voidStartTag := &Tag{Name: tagName, Position: allTagsIndex[i], Type: VOID_START_TAG, Attrs: attrMap}
			voidEndTag := &Tag{Name: tagName, Position: allTagsIndex[i], Type: VOID_END_TAG, Attrs: &AttrMap{}}
			tg.TagList = append(tg.TagList, voidStartTag, voidEndTag)
		case START_TAG:
			attrMap, err := FindAttrs(tagAttrStr)
			if err != nil {
				return err
			}
			tag := &Tag{Name: tagName, Position: allTagsIndex[i], Type: START_TAG, Attrs: attrMap}
			tg.TagList = append(tg.TagList, tag)
		case END_TAG:
			tag := &Tag{Name: tagName, Position: allTagsIndex[i], Type: END_TAG, Attrs: &AttrMap{}}
			tg.TagList = append(tg.TagList, tag)
		case INVALID_TAG:
			continue
		}
	}
	return nil
}

func (th *TagHandler) parseTagList(parentNode *Node) {
	if len(th.TagList) == 0 {
		return
	}
	currentTag, _ := th.TagList.LeftPop()
	switch currentTag.Type {
	case START_TAG, VOID_START_TAG:
		node := FromStartTag(currentTag)
		parentNode.Children = append(parentNode.Children, node)
		node.Parent = parentNode
		th.parseTagList(node)
	case END_TAG, VOID_END_TAG:
		if currentTag.Name == parentNode.Name {
			if currentTag.Type == END_TAG {
				parentNode.EndPosition = currentTag.Position[0]
			} else {
				parentNode.EndPosition = parentNode.StartPosition
			}
			parentNode.Content = th.HTML[parentNode.StartPosition:parentNode.EndPosition]
			th.parseTagList(parentNode.Parent)
		} else {
			_, _ = th.TagList.LeftPop()
			th.parseTagList(parentNode)
		}
	}
}
