package main

import (
	"fmt"
	"time"
)

func main() {
	foo := &Foo{val: 1}
	go func(foo *Foo) {
		foo.Change()
		time.Sleep(time.Millisecond * 100)
	}(foo)
	go func(foo *Foo) {
		for i := 0; len(foo.vals) > 0; {
			fmt.Println(i, foo.vals[i%len(foo.vals)])
			i++
		}
	}(foo)
	foo.Change()
	time.Sleep(time.Second)
	fmt.Println(foo.val)
	fmt.Println(foo.vals, len(foo.vals))
}

type Foo struct {
	vals []int
	val  int
}

func (f *Foo) Change() {
	f.val += 1
	f.vals = append(f.vals, f.val)
}
