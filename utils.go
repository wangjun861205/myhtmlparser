package notbearparser

import "fmt"

type AttrList []string

func NewAttrList() *AttrList {
	attrList := AttrList(make([]string, 0, 16))
	return &attrList
}

func (attrList *AttrList) String() string {
	return fmt.Sprintf("%s", *attrList)
}

func (attrList *AttrList) Pop() (string, bool) {
	switch len(*attrList) {
	case 0:
		return "", false
	case 1:
		attr, newAttrList := (*attrList)[0], AttrList(make([]string, 0))
		attrList = &newAttrList
		return attr, true
	default:
		attr, newAttrList := (*attrList)[len(*attrList)-1], (*attrList)[:len(*attrList)-1]
		attrList = &newAttrList
		return attr, true
	}
}

func (attrList *AttrList) LeftPop() (string, bool) {
	switch len(*attrList) {
	case 0:
		return "", false
	case 1:
		attr, newAttrList := (*attrList)[0], AttrList(make([]string, 0))
		attrList = &newAttrList
		return attr, true
	default:
		attr, newAttrList := (*attrList)[0], (*attrList)[1:]
		attrList = &newAttrList
		return attr, true
	}
}

func (attrList *AttrList) Splice(attr string) bool {
	attrIndex := -1
	for i, attrValue := range *attrList {
		if attrValue == attr {
			attrIndex = i
			break
		}
	}
	switch attrIndex {
	case -1:
		return false
	case 0:
		attrList.LeftPop()
		return true
	case len(*attrList) - 1:
		attrList.Pop()
		return true
	default:
		pList, lList := (*attrList)[:attrIndex], (*attrList)[attrIndex+1:]
		newList := append(pList, lList...)
		attrList = &newList
		return true
	}
}

func (attrList *AttrList) Exists(attrString string) bool {
	for _, attr := range *attrList {
		if attr == attrString {
			return true
		}
	}
	return false
}

func (attrList *AttrList) Append(s ...string) {
	*attrList = append(*attrList, s...)
}

func (attrList *AttrList) Contains(otherList *AttrList) bool {
	if len(*attrList) < len(*otherList) {
		return false
	}
	for _, attr := range *otherList {
		if !attrList.Exists(attr) {
			return false
		}
	}
	return true
}

type AttrMap map[string]*AttrList

func NewAttrMap() *AttrMap {
	attrMap := AttrMap(make(map[string]*AttrList))
	return &attrMap
}

func (attrMap *AttrMap) Exists(key, value string) bool {
	if attrList, ok := (*attrMap)[key]; ok {
		isExists := attrList.Exists(value)
		return isExists
	}
	return false
}

func (attrMap *AttrMap) Out(key string) (*AttrList, bool) {
	if attrList, ok := (*attrMap)[key]; ok {
		delete(*attrMap, key)
		return attrList, true
	}
	return &AttrList{}, false

}

func (attrMap *AttrMap) Add(key, value string) {
	if attrList, ok := (*attrMap)[key]; ok {
		attrList.Append(value)
	} else {
		attrList = NewAttrList()
		attrList.Append(value)
		(*attrMap)[key] = attrList
	}
}

func (attrMap *AttrMap) Contains(otherMap *AttrMap) bool {
	if len(*attrMap) < len(*otherMap) {
		return false
	}
	for queryKey, queryList := range *otherMap {
		if attrList, ok := (*attrMap)[queryKey]; !ok {
			return false
		} else {
			if !attrList.Contains(queryList) {
				return false
			}
		}
	}
	return true
}

func (attrMap *AttrMap) Get(key string) (AttrList, bool) {
	if attrList, ok := (*attrMap)[key]; ok {
		return *attrList, true
	} else {
		return AttrList{}, false
	}
}

type TagList []*Tag

func (tl *TagList) LeftPop() (*Tag, bool) {
	switch len(*tl) {
	case 0:
		return nil, false
	case 1:
		tag := (*tl)[0]
		*tl = TagList{}
		return tag, true
	default:
		tag := (*tl)[0]
		*tl = (*tl)[1:]
		return tag, true
	}
}

type NBString string

func (nbs NBString) EndsWith(endString string) bool {
	endLen := len(endString)
	return nbs[len(nbs)-endLen:] == NBString(endString)
}
