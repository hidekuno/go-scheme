/*
   Go lang 3rd study program.
   This is prototype program mini scheme subset.

   hidekuno@gmail.com
*/
package main

/*
  curl -v -c /tmp/cookie.txt -X POST --header "Content-Type: application/json" --data '{"user": "admin", "password": "hogehoge"}' 'http://localhost:9000/login'
  curl -v -b /tmp/cookie.txt  -c /tmp/cookie.txt 'http://localhost:9000/secret'
  curl -v -b /tmp/cookie.txt  -c /tmp/cookie.txt 'http://localhost:9000/logout'
*/
import (
	"bytes"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base32"
	"encoding/gob"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/gorilla/sessions"
	"io"
	"log"
	"net/http"
	"scheme"
	"strconv"
	"strings"
)

var (
	rootEnv *scheme.SimpleEnv
	store   *sessions.CookieStore
)

const (
	loginCookieName = "login-authentication"
	authenticated   = "authenticated"
	authUser        = "admin"
	authPassword    = "4c716d4cf211c7b7d2f3233c941771ad0507ea5bacf93b492766aa41ae9f720d"
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

type Data struct {
	Count int
	Hoge  int
}

func secret(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, loginCookieName)
	fmt.Println("secret", session)

	// Check if user is authenticated
	auth, ok := session.Values[authenticated].(bool)
	if !ok || !auth {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	// Print secret message
	cnt, _ := session.Values["cnt"].(int)
	session.Values["cnt"] = cnt + 1

	if d, ok := session.Values["data"].(*Data); ok {
		d.Count++
		session.Values["data"] = d
	}
	session.Save(r, w)
	fmt.Fprintln(w, "The cake is a lie! ", cnt, session.Values["cnt"])
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

	var userInfo map[string]interface{}
	err = json.Unmarshal(body[:length], &userInfo)
	if err != nil {
		http.Error(w, "Server Error", http.StatusInternalServerError)
		return
	}

	p, _ := userInfo["password"].(string)
	password := sha256.Sum256([]byte(p))
	if (userInfo["user"] == authUser) && (hex.EncodeToString(password[:]) == authPassword) {
		session, _ := store.Get(r, loginCookieName)
		fmt.Println("login", session)

		session.Values[authenticated] = true
		session.Values["cnt"] = 1
		session.Values["data"] = &Data{0, 10}
		session.Save(r, w)
		w.WriteHeader(http.StatusOK)
	} else {
		http.Error(w, "Unauthorized Error", http.StatusUnauthorized)
	}
}

func logout(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, loginCookieName)
	fmt.Println("logout", session)
	session.Values[authenticated] = false
	session.Save(r, w)
	w.WriteHeader(http.StatusOK)
}
func sessionInit() {
	// 乱数生成
	b := make([]byte, 48)
	_, err := io.ReadFull(rand.Reader, b)
	if err != nil {
		panic(err)
	}
	key := strings.TrimRight(base32.StdEncoding.EncodeToString(b), "=")
	store = sessions.NewCookieStore([]byte(key))
}

// Main
func main() {
	scheme.BuildFunc()
	rootEnv = scheme.NewSimpleEnv(nil, nil)

	gob.Register(Data{})
	sessionInit()

	http.HandleFunc("/lisp", doLisp)
	http.HandleFunc("/secret", secret)
	http.HandleFunc("/login", login)
	http.HandleFunc("/logout", logout)

	if err := http.ListenAndServe(":9000", nil); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
