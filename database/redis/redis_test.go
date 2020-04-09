package redis

import (
	"fmt"
	"testing"
)

func TestHSetNX(t *testing.T) {
	for i := 0; i < 5; i++ {
		re, err := HSetNX("k1", "f10")
		fmt.Println(re, ":", err)
	}
}

func TestHKExists(t *testing.T) {
	res, err := HKExists("k1", "f100")
	t.Log(res, err)
}
