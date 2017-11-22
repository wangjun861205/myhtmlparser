package notbearparser

import "fmt"

type Node struct {
	Name          string
	Attrs         *AttrMap
	Content       string
	Parent        *Node
	Children      []*Node
	PrevSibling   *Node
	NextSibling   *Node
	Siblings      []*Node
	StartPosition int
	EndPosition   int
}

func (node *Node) String() string {
	return fmt.Sprintf("name: %s, attrs: %v", node.Name, node.Attrs)
}

func GenerateNode(startTag, endTag *Tag, html string) *Node {
	if startTag.Type == VOID_START_TAG {
		node := &Node{
			Name:     startTag.Name,
			Attrs:    startTag.Attrs,
			Content:  "",
			Children: make([]*Node, 0, 8),
			Siblings: make([]*Node, 0, 8),
		}
		return node
	} else {
		node := &Node{
			Name:     startTag.Name,
			Attrs:    startTag.Attrs,
			Content:  html[startTag.Position[1]:endTag.Position[0]],
			Children: make([]*Node, 0, 8),
			Siblings: make([]*Node, 0, 8),
		}
		return node
	}
}

func FromStartTag(tag *Tag) *Node {
	return &Node{
		Name:          tag.Name,
		Attrs:         tag.Attrs,
		StartPosition: tag.Position[1],
		Children:      make([]*Node, 0, 8),
		Siblings:      make([]*Node, 0, 8),
	}
}
