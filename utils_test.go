package notbearparser

import (
	"fmt"
	"testing"
)

func TestNBString(t *testing.T) {
	nbs := NBString("<option>test</option>")
	result := nbs.EndsWith("</option>")
	fmt.Println(result)
}
