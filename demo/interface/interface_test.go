package _interface

import (
	"fmt"
	"testing"
)

// 接口类型断言
func TestInterface(t *testing.T) {
	var empty interface{} = "key"
	i, ok := empty.(string)
	fmt.Println(i, ok)
}
