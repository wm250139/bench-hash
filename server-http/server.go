package main

import (
	"bench-hash/hasher"
	"fmt"
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(resp http.ResponseWriter, req *http.Request) {
		q := req.URL.Query().Get("q")
		hash := hasher.String(q)
		_, _ = fmt.Fprintln(resp, hash)
	})

	fmt.Println("Starting server listening on http://0.0.0.0:3030")
	if err := http.ListenAndServe("0.0.0.0:3030", mux); err != nil {
		log.Fatal(err)
	}
}
