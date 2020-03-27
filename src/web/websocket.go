/*
   Go lang 3rd study program.
   This is prototype program mini scheme subset.

   hidekuno@gmail.com
*/
package web

import (
	"container/list"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	"golang.org/x/net/websocket"
)

const (
	RESOURCEDIR = "wasm"
)

type Event struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

var clients = list.List{}

func socket(conn *websocket.Conn) {

	rand.Seed(time.Now().Unix())
	id := fmt.Sprintf("#%02x%02x%02x", rand.Intn(255), rand.Intn(255), rand.Intn(255))
	logger := log.New(os.Stdout, fmt.Sprintf("[%s]\t", id), 0)

	self := clients.PushBack(conn)

	defer func() {
		conn.Close()
		clients.Remove(self)
		logger.Println("Exited loop")
	}()

	msg := struct {
		Type string
		Text string
	}{}

	ev := &Event{Type: "CONNECT", Text: ";Hello, Web Socket"}
	b, _ := json.Marshal(ev)
	conn.Write(b)
	for {
		if err := websocket.JSON.Receive(conn, &msg); err != nil {
			if err == io.EOF {
				logger.Println("Connection closed:", err)
			} else {
				logger.Println("Unexpected error:", err)
			}
			return
		}
		switch msg.Type {
		case "KEEPALIVE":
		default:
		}
	}
}
func message(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	b, _ := ioutil.ReadAll(r.Body)
	fmt.Println(string(b))
	for e := clients.Front(); e != nil; e = e.Next() {
		e.Value.(*websocket.Conn).Write(b)
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "OK")
}
func index(w http.ResponseWriter, r *http.Request) {
	indexTmpl := `<!doctype html>
<html>
<head>
  <meta charset="utf-8">
  <link href="https://fonts.googleapis.com/css?family=Raleway" rel="stylesheet">
  <link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/skeleton/2.0.4/skeleton.min.css" />
  <title>Go Scheme wasm  {{.Date}}</title>
<style type="text/css">
<!--
.result {
//  position: relative;
  background-color: whitegray;
  font-size: 1.0em;
  font-family: -webkit-body;
}
ul {
  background: #f3fbff;
  border: 2px skyblue;
}
ul li {
  padding 5px 5px;
  color: #0000cd;
  font-size: 0.75em;
  font-family: -webkit-body;
}
-->
</style>
</head>

<body>
  <script src="{{.Resource}}/wasm_exec.js"></script>
  <script src="{{.Resource}}/go.js"></script>
  <center><h2>Mini&nbsp;<span style="color:#FFE900">Scheme&nbsp;</span>Web&nbsp;<span style="color:#28AFB0">Assembly&nbsp;</span>Demo&nbsp;<span style="color:#E53D00">Program</span></h2></center>

  <center>
  <div id="head" class="row">
    <p id="calcResult" style="text-align: center;" class="result">&nbsp;</p>
  </div>
  <div id="body" class="row">
    <textarea id="sExpression" style="font-size: 16px; width: 920px; height:320px;" name="sexpression" cols="128" rows="30"></textarea>
  </div>
  </center>
  <div id="tail" style="float: right">
    <button class="button-primary" id="evalButton" style="font-size: 14px">Eval</button>
  </div>
  <div class="container">
    <ul id="history"></ul>
  </div>
  <script src="{{.Resource}}/websocket.js"></script>
</body>
</html>
`
	tpl := template.Must(template.New("wasm index").Parse(indexTmpl))
	m := map[string]string{
		"Date":     time.Now().Format("2006-01-02"),
		"Resource": RESOURCEDIR,
	}
	if err := tpl.ExecuteTemplate(w, "wasm index", m); err != nil {
		log.Fatal(err)
	}
}
func StartWebSocket() {

	s := &websocket.Server{Handler: socket}
	http.HandleFunc("/socket", s.ServeHTTP)
	http.HandleFunc("/message", message)
	http.HandleFunc("/index", index)
	http.Handle("/"+RESOURCEDIR+"/", http.StripPrefix("/"+RESOURCEDIR+"/", http.FileServer(http.Dir("./"+RESOURCEDIR))))

	port := "9000"
	log.Println("Listening on port:", port)
	http.ListenAndServe(":"+port, nil)
}
