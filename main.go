package main

import (
	"code.sk.org/insomnia_test/dal/kv"
	"net/http"
)

func main() {
	kv.InitRedisClient()

	http.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		Ping(w, r)
	})

	http.HandleFunc("/api/v1/lint_rule", func(w http.ResponseWriter, r *http.Request) {
		HandleLintRule(w, r)
	})

	http.ListenAndServe(":8080", nil)
}
