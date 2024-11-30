package main

import (
	"net/http"

	"github.com/Jseongwon/backend-golang/myapp"
)

func main() {
	// Code here

	http.ListenAndServe(":3000", myapp.NewHttpHandler())
}
