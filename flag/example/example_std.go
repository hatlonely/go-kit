package main

import (
	"fmt"

	"github.com/hatlonely/go-kit/flag"
)

func main() {
	key1 := flag.Int("key1", 123, "key1 usage")
	key2 := flag.String("key2", "val2", "key2 usage")
	if err := flag.Parse(); err != nil {
		fmt.Println(flag.Usage())
		return
	}
	fmt.Println("key1 =>", *key1)
	fmt.Println("key2 =>", *key2)
}
