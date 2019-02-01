/*
   Go lang 3rd study program.
   This is prototype program mini scheme subset.

   hidekuno@gmail.com
*/
package web

import (
	"flag"
	"log"
	"net/http"
)

// Start Wasm
func StartWebAseembly() {
	listen := flag.String("listen", ":9000", "listen address")
	dir := flag.String("dir", "./wasm", "directory to serve")
	flag.Parse()

	log.Printf("listening on %q...", *listen)
	log.Fatal(http.ListenAndServe(*listen, http.FileServer(http.Dir(*dir))))
}
