package myhtmlparser

import (
	"errors"
	"fmt"
)

type TAG_TYPE int

const (
	VOID_TAG TAG_TYPE = iota
	START_TAG
	END_TAG
	INVALID_TAG
	VOID_START_TAG
	VOID_END_TAG
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

var TypeMap = map[string]func(string) TAG_TYPE{
	"": func(tagName string) TAG_TYPE {
		if _, ok := VoidEleMap[tagName]; ok {
			return VOID_TAG
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
	Name string
	// Classes  []string
	// Attrs    map[string]string
	Attrs    *AttrMap
	Depth    int
	Position []int
	Type     TAG_TYPE
}

func (tag *Tag) String() string {
	return fmt.Sprintf("tag name: %s, tag attribute map: %v, tag depth: %d, tag pisition: %d-%d, tag type: %d",
		tag.Name, tag.Attrs, tag.Depth, tag.Position[0], tag.Position[1], tag.Type)
}

type TagHandler struct {
	TagList []*Tag
	Depth   int
	HTML    string
}

func NewTagHandler(html string) *TagHandler {
	return &TagHandler{
		Depth: 0,
		HTML:  html,
	}
}

func (tg *TagHandler) DecDepth() error {
	if tg.Depth > 0 {
		tg.Depth -= 1
		return nil
	}
	return NegativeDepthErr
}

func (tg *TagHandler) ClearScript() {
	cleanHTML := scriptRe.ReplaceAllString(tg.HTML, "")
	tg.HTML = cleanHTML
}

func (tg *TagHandler) Feed() error {
	tg.ClearScript()
	allTagsStr := tagRe.FindAllStringSubmatch(tg.HTML, -1)
	allTagsIndex := tagRe.FindAllStringIndex(tg.HTML, -1)
	for i, tagStr := range allTagsStr {
		closeToken := tagStr[1]
		tagName := tagStr[2]
		tagAttrStr := tagStr[3]
		switch CheckTagType(closeToken, tagName) {
		case VOID_TAG:
			tg.Depth += 1
			attrMap, err := FindAttrs(tagAttrStr)
			if err != nil {
				return err
			}
			voidStartTag := &Tag{Name: tagName, Depth: tg.Depth, Position: allTagsIndex[i], Type: VOID_START_TAG, Attrs: attrMap}
			voidEndTag := &Tag{Name: tagName, Depth: tg.Depth, Position: allTagsIndex[i], Type: VOID_END_TAG}
			tg.TagList = append(tg.TagList, voidStartTag, voidEndTag)
			err = tg.DecDepth()
			if err != nil {
				return err
			}
		case START_TAG:
			tg.Depth += 1
			attrMap, err := FindAttrs(tagAttrStr)
			if err != nil {
				return err
			}
			tag := &Tag{Name: tagName, Depth: tg.Depth, Position: allTagsIndex[i], Type: START_TAG, Attrs: attrMap}
			tg.TagList = append(tg.TagList, tag)
		case END_TAG:
			tag := &Tag{Name: tagName, Depth: tg.Depth, Position: allTagsIndex[i], Type: END_TAG}
			tg.TagList = append(tg.TagList, tag)
			err := tg.DecDepth()
			if err != nil {
				fmt.Println(closeToken, tagName, allTagsIndex[i])
				return err
			}
		case INVALID_TAG:
			continue
		}
	}
	return nil
}

func IsPairs(pTag, lTag *Tag) bool {
	return pTag.Name == lTag.Name && pTag.Depth == lTag.Depth && ((pTag.Type == START_TAG && lTag.Type == END_TAG) || (pTag.Type == VOID_START_TAG && lTag.Type == VOID_END_TAG))

}
