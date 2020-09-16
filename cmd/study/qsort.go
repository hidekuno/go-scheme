/*
   Go lang 2nd study program.
   This is thing for go routine (concurrent program).

   hidekuno@gmail.com
*/
package main

import (
	crand "crypto/rand"
	"flag"
	"fmt"
	"math"
	"math/big"
	"math/rand"
	"runtime"
	"time"
)

type (
	Nums []int
)

const (
	TestDataCount = 80000000
	THRESHOLD     = 1024
	IntMaxValue   = 10000000
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
		finish := make(chan bool, 2)
		go func() {
			if low < r {
				qsort(data, low, r)
			}
			finish <- true
		}()
		go func() {
			if l < high {
				qsort(data, l, high)
			}
			finish <- true
		}()
		<-finish
		<-finish
	}
}

func main() {

	runtime.GOMAXPROCS(runtime.NumCPU())
	flag.Parse()

	seed, _ := crand.Int(crand.Reader, big.NewInt(math.MaxInt64))
	rand.Seed(seed.Int64())

	data := make(Nums, TestDataCount)
	for i, _ := range data {
		data[i] = rand.Intn(IntMaxValue)
	}
	fmt.Println(len(data), "counts data generate done.")
	t0 := time.Now()
	qsort(data, 0, len(data)-1)
	t1 := time.Now()

	fmt.Println(t1.Sub(t0))
}
