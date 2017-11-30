package notbearparser

import (
	"errors"
	"fmt"
)

var NotValidHTMLErr error = errors.New("The string is not a valid HTML text")
var EOFErr error = errors.New("EOF")
var EmptyTagListErr = errors.New("There is no tag in cursor's tag list")
var NotValidQueryErr = errors.New("The query string is not valid")

type Cursor struct {
	Root            *Node
	CurrentNode     *Node
	CurrentDepth    int
	CurrentPosition int
	TagHandler      *TagHandler
}

func NewCursor(html string) *Cursor {
	return &Cursor{
		TagHandler: NewTagHandler(html),
	}
}

func (c *Cursor) Parse() error {
	err := c.TagHandler.Feed()
	if err != nil {
		return err
	}
	c.Root = &Node{Name: "root", Attrs: NewAttrMap()}
	c.TagHandler.parseTagList(c.Root)
	return nil
}

func (c *Cursor) Search(queryString string) (NodeList, error) {
	return Search(c.Root, queryString)
}

func PrintTree(root *Node) {
	fmt.Println(root.Name, root.Attrs)
	for _, child := range root.Children {
		PrintTree(child)
	}
}

// func SearchByQueryList(root *Node, queryList []*Query) []*Node {
// 	var query *Query
// 	var remianingQueryList []*Query
// 	nodeList := make([]*Node, 0, 16)
// 	switch len(queryList) {
// 	case 1:
// 		query, remianingQueryList = queryList[0], make([]*Query, 0, 0)
// 	default:
// 		query, remianingQueryList = queryList[0], queryList[1:]
// 	}
// 	switch query.Target {
// 	case ALL:
// 		nodeList = query.SearchAll(root)
// 	case ALL_CHILDREN:
// 		nodeList = query.SearchChildren(root)
// 	case DIRECT_CHILDREN:
// 		nodeList = query.SearchDirectChildren(root)
// 	}
// 	if len(remianingQueryList) > 0 {
// 		nextList := make([]*Node, 0, 16)
// 		for _, node := range nodeList {
// 			l := SearchByQueryList(node, remianingQueryList)
// 			nextList = append(nextList, l...)
// 		}
// 		return nextList
// 	}
// 	return nodeList
// }

func SearchByQueryList(root *Node, queryList []*Query) NodeList {
	var query *Query
	var remianingQueryList []*Query
	nodeList := make(NodeList, 0, 16)
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

// func Search(root *Node, queryStr string) ([]*Node, error) {
// 	queryList, err := NewQueryList(queryStr)
// 	if err != nil {
// 		return []*Node{}, err
// 	}
// 	nodeList := SearchByQueryList(root, queryList)
// 	return nodeList, nil
// }

func Search(root *Node, queryStr string) (NodeList, error) {
	queryList, err := NewQueryList(queryStr)
	if err != nil {
		return []*Node{}, err
	}
	nodeList := SearchByQueryList(root, queryList)
	return nodeList, nil
}
