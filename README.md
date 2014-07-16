#token-auth

Golang http middleware to implement token-based authentications

#Usage

wrapping a Handler to enforce token verification with gorilla/mux
```go

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
tokenAuth := tauth.NewTokenAuth(rRestrict, nil, memStore)
r.PathPrefix("/restricted").Handler(tokenAuth)

fmt.Println("listening at :3000")
http.ListenAndServe(":3000", r)

```

or just wrap individual HandleFunc

```go

mux := http.NewServeMux()
memStore := memstore.New("salty")
tokenAuth := tauth.NewTokenAuth(nil, nil, memStore)

mux.HandleFunc("/login", func(w http.ResponseWriter, req *http.Request) {
	t := memStore.NewToken("User1")
	fmt.Fprintf(w, "hi User1, your token is %s", t)
})

mux.HandleFunc("/restricted", tokenAuth.HandleFunc(func(w http.ResponseWriter, req *http.Request) {
	token := tauth.Get(req)
	fmt.Fprintf(w, "hi %s", token.Claims("id").(string))
}))

fmt.Println("listening at :3000")
http.ListenAndServe(":3000", mux)

```

Complete examples are in the example/ folder

#TokenStore

This library comes with a very simple built-in memory storage for your tokens.
You can use your own token store by implementing tauth.TokenStore interface:

```go
type TokenStore interface {
	NewToken(id interface{}) Token
	CheckToken(token string) (Token, error)
}
```

#JWT (JSON Web Tokens)

An token store implementation using (go-jwt)[https://github.com/dgrijalva/jwt-go] is also available.
See example/jwt for example usage.
