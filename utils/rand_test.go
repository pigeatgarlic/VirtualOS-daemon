package utils

import (
	"fmt"
	"testing"
)

func TestRand(t *testing.T) {
	for i := 0; i < 10; i++ {
		str := CreateRandomString(20)
		fmt.Printf("%dth time :  %s \n",i,str)
		
	}
}