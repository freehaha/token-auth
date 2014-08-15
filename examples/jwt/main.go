package main

import (
	"fmt"
	"github.com/freehaha/token-auth"
	"github.com/freehaha/token-auth/jwt"
	"net/http"
	"time"
)

func main() {
	mux := http.NewServeMux()
	jwtstore := jwtstore.New("my-secret-key")
	tokenAuth := tauth.NewTokenAuth(nil, nil, jwtstore)

	mux.HandleFunc("/login", func(w http.ResponseWriter, req *http.Request) {
		t := jwtstore.NewToken("")
		/* JwtToken implements the ClaimSetter interface */
		/* default expiration time is 30 min, you can set the "exp" claim to modify it */
		t.SetClaim("id", "user1").
			SetClaim("exp", time.Now().Add(time.Minute).Unix())

		fmt.Fprintf(w, "hi User1, your token is %s", t)
	})

	mux.HandleFunc("/restricted", tokenAuth.HandleFunc(func(w http.ResponseWriter, req *http.Request) {
		token := tauth.Get(req)
		fmt.Fprintf(w, "hi %s", token.Claims("id").(string))
	}))

	fmt.Println("listening at :3000")
	http.ListenAndServe(":3000", mux)
}
