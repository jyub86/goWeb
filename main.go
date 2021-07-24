package main

import (
	"net/http"

	"github.com/jyub86/goWeb/app"
)

func main() {

	port := "3000"
	m := app.MakeHandler()
	err := http.ListenAndServe(":"+port, m)
	if err != nil {
		panic(err)
	}
}
