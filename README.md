#token-auth

Golang http middleware to implement token-based authentications

#Usage

wrapping a Handler to enforce token verification with gorilla/mux
```go

memStore := tauth.NewMemoryTokenStore("salty")

r := mux.NewRouter()
/* requests to /login are not authenticated */
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

/* authenticate all this sub path */
tokenAuth := tauth.NewTokenAuth(rRestrict, nil, memStore)
r.PathPrefix("/restricted").Handler(tokenAuth)

http.ListenAndServe(":3000", r)

```

or just wrap individual HandleFunc

```go

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

http.ListenAndServe(":3000", mux)

```

Complete examples are in the example/ folder

#TokenStore

This library comes with a very simple built-in memory storage for your tokens.
You can use your own token store by implementing tauth.TokenStore interface:

```go
type TokenStore interface {
	NewToken(id string) *Token
	CheckToken(token string) (*Token, error)
}
```

#TODO
* expire old tokens if new one is generated for the user
