/*
   Go lang 3rd study program.
   This is prototype program mini scheme subset.

   hidekuno@gmail.com
*/
package web

import (
	"bytes"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base32"
	"encoding/gob"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"scheme"
	"strconv"
	"strings"

	"github.com/gorilla/sessions"
)

var (
	rootEnvTbl map[string]*scheme.SimpleEnv
	store      *sessions.CookieStore
)

const (
	loginCookieName = "login-authentication"
	sessionVarName  = "user-info"
	authUser        = "admin"
	authPassword    = "4c716d4cf211c7b7d2f3233c941771ad0507ea5bacf93b492766aa41ae9f720d"
)

type UserInfo struct {
	Name          string
	Authenticated bool
	UseCount      int
}

func doLisp(w http.ResponseWriter, r *http.Request) {

	// check authenticated
	session, _ := store.Get(r, loginCookieName)
	userInfo, ok := session.Values[sessionVarName].(*UserInfo)
	if !ok || !userInfo.Authenticated {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	// main proc
	bufbody := new(bytes.Buffer)
	bufbody.ReadFrom(r.Body)
	e, err := scheme.DoCoreLogic(bufbody.String(), rootEnvTbl[userInfo.Name])

	userInfo.UseCount++
	session.Values[sessionVarName] = userInfo
	session.Save(r, w)

	if err != nil {
		fmt.Fprintln(w, err.Error())
	} else {
		fmt.Fprintln(w, e.String())
	}
	log.Print(r.URL, " ", userInfo.UseCount)
}

func login(w http.ResponseWriter, r *http.Request) {

	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	length, err := strconv.Atoi(r.Header.Get("Content-Length"))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	body := make([]byte, length)
	length, err = r.Body.Read(body)
	if err != nil && err != io.EOF {
		http.Error(w, "Server Error", http.StatusInternalServerError)
		return
	}

	var loginInfo map[string]interface{}
	err = json.Unmarshal(body[:length], &loginInfo)
	if err != nil {
		http.Error(w, "Json Parse Error", http.StatusInternalServerError)
		return
	}

	p, _ := loginInfo["password"].(string)
	password := sha256.Sum256([]byte(p))
	if (loginInfo["user"] == authUser) && (hex.EncodeToString(password[:]) == authPassword) {
		session, _ := store.Get(r, loginCookieName)

		u, _ := loginInfo["user"].(string)
		session.Values[sessionVarName] = &UserInfo{u, true, 0}
		rootEnvTbl[u] = scheme.NewSimpleEnv(nil, nil)

		if err := session.Save(r, w); err != nil {
			log.Println(err)
		}
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "OK")
	} else {
		http.Error(w, "Unauthorized Error", http.StatusUnauthorized)
	}
}

func logout(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, loginCookieName)

	userInfo, ok := session.Values[sessionVarName].(*UserInfo)
	if !ok {
		http.Error(w, "Unauthorized Error", http.StatusUnauthorized)
		return
	}
	userInfo.Authenticated = false
	rootEnvTbl[userInfo.Name] = nil

	session.Values[sessionVarName] = userInfo
	session.Save(r, w)
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "OK")
}
func sessionInit() {
	// 乱数生成
	b := make([]byte, 48)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		panic(err)
	}
	store = sessions.NewCookieStore([]byte(strings.TrimRight(base32.StdEncoding.EncodeToString(b), "=")))
}

// Start Service
func StartApiService() {

	scheme.BuildFunc()
	rootEnvTbl = map[string]*scheme.SimpleEnv{}
	gob.Register(&UserInfo{})

	sessionInit()

	http.HandleFunc("/lisp", doLisp)
	http.HandleFunc("/login", login)
	http.HandleFunc("/logout", logout)

	if err := http.ListenAndServe(":9000", nil); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

// Start Wasm
func StartWebAseembly() {
	listen := flag.String("listen", ":9000", "listen address")
	dir := flag.String("dir", "./wasm", "directory to serve")
	flag.Parse()

	log.Printf("listening on %q...", *listen)
	log.Fatal(http.ListenAndServe(*listen, http.FileServer(http.Dir(*dir))))
}
