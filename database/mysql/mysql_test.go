package mysql

import (
	"fmt"
	"testing"
)

func TestSave(t *testing.T) {
	rs, err := Save("http://www.baidu.com", "baidu")
	t.Log(rs, err)
}

func TestGetQueryResult(t *testing.T) {
	cols, rs := GetQueryResult()
	fmt.Println(cols)
	fmt.Println(rs)
}
