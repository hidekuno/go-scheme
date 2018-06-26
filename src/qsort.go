/*
   Go lang 2nd study program.
   This is thing for go routine (concurrent program).

   hidekuno@gmail.com
*/
package main

import (
	"flag"
	"fmt"
	"math/rand"
	"runtime"
	"time"
)

type (
	Nums []int
)

const (
	TEST_DATA_COUNT = 80000000
	THRESHOLD       = 1024
	INT_MAX_VALUE   = 10000000
)

var (
	goroutine = flag.Bool("g", false, "exec go routine")
)

func qsort(data Nums, low, high int) {

	mid := data[(low+high)/2]
	l, r := low, high
	for l <= r {
		for ; data[l] < mid; l++ {
		}
		for ; mid < data[r]; r-- {
		}
		if l <= r {
			data[r], data[l] = data[l], data[r]
			l++
			r--
		}
	}
	if (high-low < THRESHOLD) || (*goroutine == false) {
		if low < r {
			qsort(data, low, r)
		}
		if l < high {
			qsort(data, l, high)
		}
	} else {
		finish_ch := make(chan bool, 2)
		go func() {
			if low < r {
				qsort(data, low, r)
			}
			finish_ch <- true
		}()
		go func() {
			if l < high {
				qsort(data, l, high)
			}
			finish_ch <- true
		}()
		<-finish_ch
		<-finish_ch
	}
}

func main() {

	runtime.GOMAXPROCS(runtime.NumCPU())
	flag.Parse()

	data := make(Nums, TEST_DATA_COUNT)
	for i, _ := range data {
		data[i] = rand.Intn(INT_MAX_VALUE)
	}
	fmt.Println(len(data), "counts data generate done.")
	t0 := time.Now()
	qsort(data, 0, len(data)-1)
	t1 := time.Now()

	fmt.Println(t1.Sub(t0))
}
