package main

import (
	"fmt"
	"github.com/freehaha/token-auth"
	"github.com/gorilla/mux"
	"net/http"
)

func main() {
	memStore := tauth.NewMemoryTokenStore("salty")

	r := mux.NewRouter()
	r.HandleFunc("/login", func(w http.ResponseWriter, req *http.Request) {
		t := memStore.NewToken("User1")
		fmt.Fprintf(w, "hi User1, your token is %s", t.Token)
	})

	rRestrict := mux.NewRouter()
	s := rRestrict.PathPrefix("/restricted/").Subrouter()
	s.HandleFunc("/area", func(w http.ResponseWriter, req *http.Request) {
		token := tauth.Get(req)
		fmt.Fprintf(w, "hi %s", token.Id)
	})
	tokenAuth := tauth.NewTokenAuth(rRestrict, nil, memStore)
	r.PathPrefix("/restricted").Handler(tokenAuth)

	fmt.Println("listening at :3000")
	http.ListenAndServe(":3000", r)
}
