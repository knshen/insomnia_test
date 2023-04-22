package main

import (
	"net/http"
)

func main() {
	initService()

	http.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		Ping(w, r)
	})

	http.ListenAndServe(":8080", nil)
}

func initService() {

}
