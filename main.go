package main

import (
	"net/http"

	"github.com/gy-kim/golang-daily-practice/myapp"
)

func main() {
	// Code here

	http.ListenAndServe(":3000", myapp.NewHttpHandler())
}
