package myhtmlparser

import (
	"errors"
	"fmt"
)

var NotValidHTMLErr error = errors.New("The string is not a valid HTML text")
var EOFErr error = errors.New("EOF")
var EmptyTagListErr = errors.New("There is no tag in cursor's tag list")
var NotValidQueryErr = errors.New("The query string is not valid")

type Cursor struct {
	// RawHTML         string
	Root            *Node
	CurrentNode     *Node
	CurrentDepth    int
	CurrentPosition int
	TagHandler      *TagHandler
}

func NewCursor(html string) *Cursor {
	return &Cursor{
		// RawHTML:    html,
		TagHandler: NewTagHandler(html),
	}
}

func (c *Cursor) GenerateRoot() error {
	switch len(c.TagHandler.TagList) {
	case 0:
		return EmptyTagListErr
	case 1:
		if c.TagHandler.TagList[0].Type != VOID_TAG {
			return NotValidHTMLErr
		}
		c.Root = GenerateNode(c.TagHandler.TagList[0], c.TagHandler.TagList[0], c.TagHandler.HTML)
		c.TagHandler.TagList = make([]*Tag, 0)
		return nil
	default:
		nRoot := 0
		for _, tag := range c.TagHandler.TagList {
			if tag.Depth == 1 {
				if tag.Type == VOID_TAG {
					nRoot += 2
				} else {
					nRoot += 1
				}
			}
		}
		if nRoot%2 != 0 {
			return NotValidHTMLErr
		}
		if nRoot > 2 {
			c.Root = &Node{Name: "VirtualNode"}
			return nil
		} else {
			if !IsPairs(c.TagHandler.TagList[0], c.TagHandler.TagList[len(c.TagHandler.TagList)-1]) {
				return NotValidHTMLErr
			}
			root := GenerateNode(c.TagHandler.TagList[0], c.TagHandler.TagList[len(c.TagHandler.TagList)-1], c.TagHandler.HTML)
			c.Root = root
			c.TagHandler.TagList = c.TagHandler.TagList[1 : len(c.TagHandler.TagList)-1]
			return nil
		}
	}
}

func (c *Cursor) Parse() error {
	err := c.TagHandler.Feed()
	if err != nil {
		return err
	}
	err = c.GenerateRoot()
	if err != nil {
		return err
	}
	GenerateNodeTree(c.Root, c.TagHandler.TagList, c.TagHandler.HTML)
	return nil
}

func GenerateNodeTree(currentNode *Node, tagList []*Tag, html string) {
	if len(tagList) == 0 {
		return
	}
	tag := tagList[0]
	for index, nextTag := range tagList {
		if IsPairs(tag, nextTag) {
			node := GenerateNode(tag, nextTag, html)
			node.Parent = currentNode
			currentNode.Children = append(currentNode.Children, node)
			GenerateNodeTree(node, tagList[1:index], html)
			if index != len(tagList)-1 {
				newTagList := tagList[index+1:]
				GenerateNodeTree(currentNode, newTagList, html)
			}
			break
		}
	}
}

func FilterTree(root *Node, filter Filter) []*Node {
	nodeList := make([]*Node, 0, 128)
	if ok := filter(root); ok {
		nodeList = append(nodeList, root)
	}
	for _, child := range root.Children {
		matchChildrenList := FilterTree(child, filter)
		nodeList = append(nodeList, matchChildrenList...)
	}
	return nodeList
}

func PrintTree(root *Node) {
	fmt.Println(root.Name, root.Attrs)
	for _, child := range root.Children {
		PrintTree(child)
	}
}

// func Search(root *Node, queryStr string) ([]*Node, error) {
// 	switch {
// 	case queryAttrRe.MatchString(queryStr):
// 		nodeGroup := queryAttrRe.FindStringSubmatch(queryStr)
// 		nodeName, nodeAttrStr := nodeGroup[1], nodeGroup[2]
// 		queryMap := FindQueryAttrs(nodeAttrStr)
// 		matchList := FilterTree(root, AttrFilter(nodeName, queryMap))
// 		return matchList, nil
// 	default:
// 		return []*Node{}, NotValidQueryErr
// 	}
// }

func SearchByQueryList(root *Node, queryList []*Query) []*Node {
	var query *Query
	var remianingQueryList []*Query
	nodeList := make([]*Node, 0, 16)
	switch len(queryList) {
	case 1:
		query, remianingQueryList = queryList[0], make([]*Query, 0, 0)
	default:
		query, remianingQueryList = queryList[0], queryList[1:]
	}
	switch query.Target {
	case ALL:
		nodeList = query.SearchAll(root)
	case ALL_CHILDREN:
		nodeList = query.SearchChildren(root)
	case DIRECT_CHILDREN:
		nodeList = query.SearchDirectChildren(root)
	}
	if len(remianingQueryList) > 0 {
		nextList := make([]*Node, 0, 16)
		for _, node := range nodeList {
			l := SearchByQueryList(node, remianingQueryList)
			nextList = append(nextList, l...)
		}
		return nextList
	}
	return nodeList
}

func Search(root *Node, queryStr string) ([]*Node, error) {
	queryList, err := NewQueryList(queryStr)
	if err != nil {
		return []*Node{}, err
	}
	nodeList := SearchByQueryList(root, queryList)
	return nodeList, nil
}
