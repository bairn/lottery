/**
 * 并发编程,切片的线程安全问题
 */
package main

import (
	"fmt"
	"sync"
	"time"
)

var list []int = []int{}
var wgList sync.WaitGroup = sync.WaitGroup{}
var muList sync.Mutex = sync.Mutex{}

func main() {
	max := 100000
	fmt.Printf("list add num=%d\n", max)

	wgList.Add(max)
	time1 := time.Now().UnixNano()
	for i:=0 ; i<max ; i++ {
		go addNotSafe()
	}
	wgList.Wait()
	time2 := time.Now().UnixNano()
	fmt.Printf("list len=%d, time=%d ms\n", len(list), (time2-time1)/1000000)


	list = []int{}
	fmt.Printf("new list add num=%d\n", max)
	wgList.Add(max)
	time3 := time.Now().UnixNano()
	for i:=0; i< max; i++ {
		go addSafe()
	}
	wgList.Wait()

	time4 := time.Now().UnixNano()
	fmt.Printf("new list len=%d, time=%d ms\n", len(list), (time4-time3)/1000000)
}

func addNotSafe() {
	list = append(list, 1)
	wgList.Done()
}

func addSafe() {
	muList.Lock()
	list = append(list, 1)
	muList.Unlock()
	wgList.Done()
}

