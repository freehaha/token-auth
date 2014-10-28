package main

import (
	"fmt"
	"github.com/freehaha/token-auth"
	"github.com/freehaha/token-auth/memory"
	"github.com/gorilla/mux"
	"net/http"
)

func main() {
	memStore := memstore.New("salty")

	r := mux.NewRouter()
	r.HandleFunc("/login/{id}", func(w http.ResponseWriter, req *http.Request) {
		vars := mux.Vars(req)
		t := memStore.NewToken(vars["id"])
		fmt.Fprintf(w, "hi %s, your token is %s", vars["id"], t)
	})

	rRestrict := mux.NewRouter()
	s := rRestrict.PathPrefix("/restricted/").Subrouter()
	s.HandleFunc("/area", func(w http.ResponseWriter, req *http.Request) {
		token := tauth.Get(req)
		fmt.Fprintf(w, "hi %s", token.Claims("id"))
	})
	tokenAuth := tauth.NewTokenAuth(rRestrict, nil, memStore, nil)
	r.PathPrefix("/restricted").Handler(tokenAuth)

	fmt.Println("listening at :3000")
	http.ListenAndServe(":3000", r)
}
