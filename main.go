package main

import (
	"database/sql"
	"forum/cmd/web"
	"forum/pkg/models"
	_ "github.com/mattn/go-sqlite3"
	"html/template"
	"log"
	"net/http"
)

type application struct {
	errorLog *log.Logger
	infoLog  *log.Logger
	users    interface {
		Insert(string, string, string) error
		Authenticate(string, string) (int, error)
		Get(int) (*models.User, error)
	}

	templateCache map[string]*template.Template
}

func main() {

	//addr := flag.String("addr", ":4000", "Сетевой адрес веб-сервера")
	//dsn := flag.String("dsn", "web:pass@/snippetbox?parseTime=true", "Название MySQL источника данных")
	//flag.Parse()
	//
	//infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	//errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	//db, err := openDB(*dsn)
	//if err != nil {
	//	errorLog.Fatal(err)
	//}
	//defer db.Close()

	mux := http.NewServeMux()
	mux.HandleFunc("/", web.HandleLogin)
	http.ListenAndServe(":8000", mux)
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dsn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
