package utils

import (
	"fmt"
	"testing"
)

func TestPath(t *testing.T) {
	path,err := FindProcessPath(nil,"where.exe")
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Printf("path found: %s\n",path)
}