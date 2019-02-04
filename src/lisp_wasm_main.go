/*
   Go lang 3rd study program.
   This is prototype program mini scheme subset.

   hidekuno@gmail.com
*/
package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"runtime"
	"time"
	"web"
)

type Event struct {
	Type string `json:"type"`
	Text string `json:"text"`
	User string `json:"user"`
}

func broadcast_code(code string) {

	ev := &Event{Type: "MESSAGE", Text: code, User: "test"}
	data, _ := json.Marshal(ev)

	client := &http.Client{Timeout: time.Duration(10) * time.Second}

	req, _ := http.NewRequest("POST", "http://localhost:9000/message", bytes.NewBuffer(data))
	req.Header.Set("Content-Type", "application/json")
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	b, err := ioutil.ReadAll(res.Body)
	fmt.Println(string(b))
	if err != nil {
		fmt.Println(err)
	}
}
func do_client() {
	reader := bufio.NewReaderSize(os.Stdin, 1024)
	for {
		fmt.Print("[Websocket] ")
		b, _, err := reader.ReadLine()
		if err != nil {
			fmt.Println(err)
		}
		line := string(b)
		if line == "(quit)" {
			break
		}
		broadcast_code(line)
	}
}

// Main
func main() {

	runtime.GOMAXPROCS(runtime.NumCPU())
	c := flag.Bool("c", false, "client test code")
	flag.Parse()

	if *c == true {
		do_client()

	} else {
		web.StartWebSocket()
	}
}
