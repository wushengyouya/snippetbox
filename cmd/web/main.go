package main

import (
	"database/sql"
	"flag"
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

type application struct {
	errorLog *log.Logger
	infoLog  *log.Logger
}

func main() {
	addr := flag.String("addr", ":4000", "HTTP 监听地址")
	dsn := flag.String("dsn", "web:pass@tcp(127.0.0.1:13306)/snippetbox?parseTime=true", "MySQL DSN")
	flag.Parse()

	// 初始化日志
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLOg := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	db, err := openDB(*dsn)
	if err != nil {
		errorLOg.Fatal(err)
	}
	defer db.Close()

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

// 打开数据库连接
func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		db.Close()
		return nil, err
	}
	return db, nil
}
