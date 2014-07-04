package main

import (
	"fmt"
	"github.com/freehaha/token-auth"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	memStore := tauth.NewMemoryTokenStore("salty")
	tokenAuth := tauth.NewTokenAuth(nil, nil, memStore)

	mux.HandleFunc("/login", func(w http.ResponseWriter, req *http.Request) {
		t := memStore.NewToken("User1")
		fmt.Fprintf(w, "hi User1, your token is %s", t.Token)
	})

	mux.HandleFunc("/restricted", tokenAuth.HandleFunc(func(w http.ResponseWriter, req *http.Request) {
		token := tauth.Get(req)
		fmt.Fprintf(w, "hi %s", token.Id)
	}))

	fmt.Println("listening at :3000")
	http.ListenAndServe(":3000", mux)
}
