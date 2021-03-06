/*
   Go lang 3rd study program.
   This is prototype program mini scheme subset.

   hidekuno@gmail.com
*/
package main

/*
  curl -v -c /tmp/cookie.txt -X POST --header "Content-Type: application/json" --data '{"user": "admin", "password": "hogehoge"}' 'http://localhost:9000/login'
  curl -b /tmp/cookie.txt  -c /tmp/cookie.txt --data "(+ 1 2)" 'http://localhost:9000/lisp'
  curl -v -b /tmp/cookie.txt  -c /tmp/cookie.txt 'http://localhost:9000/logout'
*/
import (
	"runtime"

	"local.packages/web"
)

// Main
func main() {

	runtime.GOMAXPROCS(runtime.NumCPU())

	web.StartApiService()
}
