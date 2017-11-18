package myhtmlparser

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"
)

func TestParser(t *testing.T) {
	htmlFile, err := os.Open("testhtml.html")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer htmlFile.Close()
	html, err := ioutil.ReadAll(htmlFile)
	if err != nil {
		fmt.Println(err)
		return
	}
	cursor := NewCursor(string(html[:]))
	err = cursor.Parse()
	// PrintTree(cursor.Root)
	l, err := Search(cursor.Root, "li")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(len(l))
	for _, n := range l {
		attrs, _ := n.Attrs.Out("class")
		fmt.Println(n.Name, attrs)
	}
}
