package main

import (
	"fmt"
	"lru"
)

type String string

func (s String) Len() int {
	return len(s)
}

func main() {
	cache := lru.New(2, nil)
	cache.Set("hello", String("World"))
	fmt.Println(cache.Get("hello"))
}
