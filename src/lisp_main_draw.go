/*
   Go lang 3rd study program.
   This is prototype program mini scheme subset.

   hidekuno@gmail.com
*/
package main

import (
	"runtime"
	"scheme"
)

// Main
func main() {

	runtime.GOMAXPROCS(runtime.NumCPU())

	scheme.BuildFunc()
	scheme.BuildGtkFunc()

	cui := make(chan bool)
	go func() {
		scheme.DoInteractive()
		cui <- true
	}()
	<-cui
}
