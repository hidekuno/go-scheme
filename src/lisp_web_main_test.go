/*
   Go lang 4th study program.
   This is test program for mini scheme subset.

   hidekuno@gmail.com
*/
package main

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"strings"
	"testing"
)

func doRequest(client *http.Client, method string, uri string, data ...string) (int, string) {
	var req *http.Request
	if len(data) == 0 {
		req, _ = http.NewRequest(method, "http://localhost:9000"+uri, nil)
	} else {
		req, _ = http.NewRequest(method, "http://localhost:9000"+uri, bytes.NewBuffer([]byte(data[0])))
	}
	if strings.ToUpper(method) == "POST" {
		req.Header.Set("Content-Type", "application/json")
	}

	res, err := client.Do(req)
	if err != nil {
		return -1, "test failed"
	}
	byteArray, _ := ioutil.ReadAll(res.Body)

	res.Body.Close()
	return res.StatusCode, strings.TrimRight(string(byteArray), "\n")
}
func TestWebOperation(t *testing.T) {
	result := ""
	status := 0

	jar, _ := cookiejar.New(nil)
	client := &http.Client{Jar: jar}

	status, result = doRequest(client, "POST", "/login", `{"user": "admin", "password": "hogehoge"}`)
	if status != http.StatusOK || result != "OK" {
		t.Fatal("failed test: /login")
	}

	status, result = doRequest(client, "GET", "/lisp", "(+ 1 2)")
	if status != http.StatusOK || result != "3" {
		t.Fatal("failed test: (+ 1 2)")
	}

	status, result = doRequest(client, "GET", "/lisp", "(define a 100)")
	if status != http.StatusOK || result != "a" {
		t.Fatal("failed test: define")
	}

	status, result = doRequest(client, "GET", "/lisp", "(+ a 200)")
	if status != http.StatusOK || result != "300" {
		t.Fatal("failed test: (+ a 200)")
	}

	status, result = doRequest(client, "GET", "/logout")
	if status != http.StatusOK || result != "OK" {
		t.Fatal("failed test: /logout")
	}
}
func TestWebOperationErrorCase(t *testing.T) {
	result := ""
	status := 0

	jar, _ := cookiejar.New(nil)
	client := &http.Client{Jar: jar}

	status, result = doRequest(client, "GET", "/login", `{"user": "admin", "password": "damedame"}`)
	if status != http.StatusBadRequest || result != "Bad Request" {
		t.Fatal("failed test: /login error case")
	}

	status, result = doRequest(client, "POST", "/login", `{"user": "admin", "password": "damedame"`)
	if status != http.StatusInternalServerError || result != "Json Parse Error" {
		t.Fatal("failed test: /login error case")
	}

	status, result = doRequest(client, "POST", "/login", `{"user": "admin", "password": "damedame"}`)
	if status != http.StatusUnauthorized || result != "Unauthorized Error" {
		t.Fatal("failed test: /login error case")
	}

	status, result = doRequest(client, "GET", "/lisp", "(+ 1 2)")
	if status != http.StatusForbidden || result != "Forbidden" {
		t.Fatal("failed test: /lisp error case")
	}

	status, result = doRequest(client, "GET", "/logout")
	if status != http.StatusUnauthorized || result != "Unauthorized Error" {
		t.Fatal("failed test: /logout error case")
	}
}
