package notbearparser

var testStr = `<div class="test_class" id="test_id" data-pk="test-pk" disable test>`
var queryStr = `div[id=test_id, class=test_class, disable, hidden]`

// func TestRegexp(t *testing.T) {
// 	attrGroup := queryAttrRe.FindStringSubmatch(queryStr)
// 	_, tagAttrStr := attrGroup[1], attrGroup[2]
// 	allAttrs := queryAllAttrRe.FindAllString(tagAttrStr, -1)
// 	for _, attrStr := range allAttrs {
// 		if queryEquationAttrRe.MatchString(attrStr) {
// 			eqAttrGroup := queryEquationAttrRe.FindStringSubmatch(attrStr)
// 			fmt.Println(eqAttrGroup[1], eqAttrGroup[2])
// 		} else {
// 			sAttr := querySingleAttrRe.FindString(attrStr)
// 			fmt.Println(sAttr)
// 		}
// 	}
// }
