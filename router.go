package main

import (
	"net/http"
	"time"
)

// API router

func Ping(w http.ResponseWriter, req *http.Request) {
	w.Write([]byte(time.Now().Format("2006-01-02 15:04:05")))
}
