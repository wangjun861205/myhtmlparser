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
	l, err := Search(cursor.Root, "#menuHorizontal>li>a")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(len(l))
	for _, n := range l {
		fmt.Println(n.Name, n.Attrs["href"], n.Content)
	}
}
