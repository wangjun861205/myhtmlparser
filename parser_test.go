package notbearparser

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"
)

func TestParser(t *testing.T) {
	htmlFile, err := os.Open("html.html")
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
	if err != nil {
		fmt.Println(err)
		return
	}
	// PrintTree(cursor.Root)
	l, err := Search(cursor.Root, "option[value='2015']")
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, n := range l {
		fmt.Println(n.Name, n.Attrs, n.Content)
	}
}
