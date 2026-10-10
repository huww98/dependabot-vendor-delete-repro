package main

import (
	"fmt"

	"github.com/cespare/xxhash/v2"
)

func main() {
	fmt.Println(xxhash.Sum64String("hello"))
}
