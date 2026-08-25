package main

import (
	"log"
	"net/http"
)

type application struct{}

func main() {
	app := &application{}
	log.Print("服务启动与 http://localhost:4000")
	log.Fatal(http.ListenAndServe(":4000", app.routes()))
}
