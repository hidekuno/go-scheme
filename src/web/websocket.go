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

func StartWebSocket() {

	s := &websocket.Server{Handler: socket}
	http.HandleFunc("/socket", s.ServeHTTP)
	http.HandleFunc("/message", message)
	http.HandleFunc("/index", func(w http.ResponseWriter, r *http.Request) {
		tpl := template.Must(template.ParseFiles("index.tpl"))
		m := map[string]string{
			"Date": time.Now().Format("2006-01-02"),
		}
		tpl.Execute(w, m)
	})
	http.Handle("/wasm/", http.StripPrefix("/wasm/", http.FileServer(http.Dir("./wasm"))))

	port := "9000"
	log.Println("Listening on port:", port)
	http.ListenAndServe(":"+port, nil)
	//log.Fatal(http.ListenAndServe(*listen, http.FileServer(http.Dir(*dir))))
}
