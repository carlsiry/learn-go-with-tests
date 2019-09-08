package main

import (
	"net/http"
	"os"

	"github.com/carlsiry/learn-go-with-tests/di"
)

func GreetHandler(w http.ResponseWriter, r *http.Request) {
	di.Greet(w, r.RequestURI)
}

func main() {
	di.Greet(os.Stdout, "terminal")
	http.ListenAndServe(":8080", http.HandlerFunc(GreetHandler)).Error()
}
