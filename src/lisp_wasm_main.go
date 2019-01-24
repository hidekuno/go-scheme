/*
   Go lang 3rd study program.
   This is prototype program mini scheme subset.

   hidekuno@gmail.com
*/
package main

import (
	"runtime"
	"scheme/web"
)

// Main
func main() {

	runtime.GOMAXPROCS(runtime.NumCPU())

	web.StartWebAseembly()
}
