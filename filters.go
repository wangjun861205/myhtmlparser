package myhtmlparser

type Filter func(*Node) bool

func AttrFilter(nodeName string, attrMap map[string]string) Filter {
	return Filter(func(node *Node) bool {
		return node.Name == nodeName && MatchAttrs(node, attrMap)
	})
}
