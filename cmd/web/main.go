package main

import (
	"database/sql"
	"flag"
	"html/template"
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"snippetbox.alexedwards.net/internal/models"
)

type application struct {
	// 数据操纵model
	snippets *models.SnippetModel
	// 页面模版缓存，避免每次请求要获取页面
	templateCache map[string]*template.Template
	errorLog      *log.Logger
	infoLog       *log.Logger
}

func main() {
	addr := flag.String("addr", ":4000", "HTTP 监听地址")
	dsn := flag.String("dsn", "web:pass@tcp(127.0.0.1:13306)/snippetbox?parseTime=true", "MySQL DSN")
	flag.Parse()

	// 初始化日志
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	db, err := openDB(*dsn)
	if err != nil {
		errorLog.Fatal(err)
	}
	defer db.Close()
	templateCache, err := newTemplateCache()
	if err != nil {
		errorLog.Fatal(err)
	}
	app := &application{
		snippets:      &models.SnippetModel{DB: db},
		templateCache: templateCache,
		infoLog:       infoLog,
		errorLog:      errorLog,
	}

	srv := http.Server{
		Addr:     *addr,
		Handler:  app.routes(),
		ErrorLog: errorLog,
	}

	infoLog.Printf("服务启动于: %s", *addr)
	errorLog.Fatal(srv.ListenAndServe())
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
