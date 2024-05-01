package socks

import (
	"fmt"
	"net/http"

	"github.com/caddyserver/certmagic"
)

func RunHTTPChallengeServer(httpAddr string) {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(w, "Hello...")
	})
	magic := certmagic.NewDefault()
	myACME := certmagic.NewACMEIssuer(magic, certmagic.DefaultACME)
	http.ListenAndServe(httpAddr, myACME.HTTPChallengeHandler(mux))
}