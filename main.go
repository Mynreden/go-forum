package main

import (
	"forum/cmd/web"
	"net/http"
)

//var db *sql.DB

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
