package main

import (
	"flag"
	"log"
	"net/http"
	"os"
)

type application struct {
	errorLog *log.Logger
	infoLog  *log.Logger
}

func main() {
	addr := flag.String("addr", ":4000", "HTTP 监听地址")
	flag.Parse()

	// 初始化日志
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLOg := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	app := &application{
		infoLog:  infoLog,
		errorLog: errorLOg,
	}

	srv := http.Server{
		Addr:     *addr,
		Handler:  app.routes(),
		ErrorLog: errorLOg,
	}

	infoLog.Printf("服务启动于: %s", *addr)
	errorLOg.Fatal(srv.ListenAndServe())
}
