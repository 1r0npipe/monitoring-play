package main

import "net/http"

func main() {
	http.HandleFunc("/", mainHandler)
	http.ListenAndServe(":8080", nil)

}

func mainHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello Jeson"))
}