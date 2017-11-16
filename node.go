package myhtmlparser

type Node struct {
	Name        string
	Attrs       map[string]string
	Content     string
	Parent      *Node
	Children    []*Node
	PrevSibling *Node
	NextSibling *Node
	Siblings    []*Node
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

func MatchAttrs(node *Node, attrMap map[string]string) bool {
	for k, v := range attrMap {
		if attrVal, ok := node.Attrs[k]; !ok || attrVal != v {
			return false
		}
	}
	return true
}
