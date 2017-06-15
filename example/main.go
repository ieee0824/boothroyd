package main

import (
	"github.com/ieee0824/boothroyd"
	"fmt"
	"time"
)

func main() {
	q := boothroyd.New()
	q.Enqueue("test", "hoge")
	q.Enqueue("test", "hoge")
	q.Enqueue("test", "hoge")
	q.Enqueue("test", "hoge")
	q.Enqueue("test", "hoge")
	q.Enqueue("test", "hoge")
	fmt.Println(q)

	for {
		if q.IsEmpty() {
			fmt.Println(q)
			time.Sleep(25*time.Second)
			fmt.Println(q)
			break
		}
		select {
		case d := <- q.C:
			fmt.Println(d)
		}
	}
}