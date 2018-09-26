/*
   Go lang 3rd study program.
   This is prototype program mini scheme subset.

   hidekuno@gmail.com
*/
package main

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"scheme"
)

func doLisp(w http.ResponseWriter, r *http.Request) {
	log.Print(r.URL)

	bufbody := new(bytes.Buffer)
	bufbody.ReadFrom(r.Body)

	e, err := scheme.DoCoreLogic(bufbody.String(), rootEnv)
	if err != nil {
		fmt.Fprintln(w, err.Error())
	} else {
		e.Fprint(w)
		fmt.Fprintln(w)
	}
}

var (
	rootEnv *scheme.SimpleEnv
)

// Main
func main() {
	scheme.BuildFunc()
	rootEnv = scheme.NewSimpleEnv(nil, nil)

	http.HandleFunc("/lisp", doLisp)

	if err := http.ListenAndServe(":9000", nil); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
