package main

import (
	"net/http"
)

func main() {
	http.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		Ping(w, r)
	})

	http.ListenAndServe(":8080", nil)
}
