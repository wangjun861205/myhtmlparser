package myhtmlparser

type QUERY_TARGET int

const (
	ALL QUERY_TARGET = iota
	DIRECT_CHILDREN
	ALL_CHILDREN
)

var TargetMap = map[string]QUERY_TARGET{
	"":  ALL,
	",": ALL,
	">": DIRECT_CHILDREN,
	" ": ALL_CHILDREN,
}

type Query struct {
	Target   QUERY_TARGET
	NodeName string
	AttrMap  *AttrMap
}

func NewQueryList(queryStr string) ([]*Query, error) {
	queryList := make([]*Query, 0, 8)
	allQuery := queryRe.FindAllStringSubmatch(queryStr, -1)
	if len(allQuery) == 0 {
		return queryList, NotValidQueryErr
	}
	for _, subquery := range allQuery {
		query := &Query{AttrMap: NewAttrMap()}
		query.Target = TargetMap[subquery[1]]
		switch subquery[2] {
		case "":
			query.NodeName = subquery[3]
		case ".":
			query.AttrMap.Add("class", subquery[3])
		case "#":
			query.AttrMap.Add("id", subquery[3])
		}
		if subquery[4] != "" {
			allAttrs := queryAttrRe.FindAllStringSubmatch(subquery[4], -1)
			for _, attr := range allAttrs {
				query.AttrMap.Add(attr[1], attr[2])
			}
		}
		queryList = append(queryList, query)
	}
	return queryList, nil
}

func (query *Query) Match(node *Node) bool {
	if query.NodeName != "" && query.NodeName != node.Name {
		return false
	}
	return node.Attrs.Contains(query.AttrMap)
}

func (query *Query) SearchAll(node *Node) []*Node {
	matchList := make([]*Node, 0, 16)
	if isMatch := query.Match(node); isMatch {
		matchList = append(matchList, node)
	}
	for _, child := range node.Children {
		childrenMatchList := query.SearchAll(child)
		matchList = append(matchList, childrenMatchList...)
	}
	return matchList
}

func (query *Query) SearchChildren(node *Node) []*Node {
	matchList := make([]*Node, 0, 16)
	for _, child := range node.Children {
		childrenMatchList := query.SearchAll(child)
		matchList = append(matchList, childrenMatchList...)
	}
	return matchList
}

func (query *Query) SearchDirectChildren(node *Node) []*Node {
	matchList := make([]*Node, 0, 16)
	for _, child := range node.Children {
		if ok := query.Match(child); ok {
			matchList = append(matchList, child)
		}
	}
	return matchList
}
