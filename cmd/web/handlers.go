package web

import (
	"fmt"
	"forum/pkg/db"
	_ "github.com/mattn/go-sqlite3"
	"html/template"
	"net/http"
)

func HandleLogin(w http.ResponseWriter, r *http.Request) {
	ts, err := template.ParseFiles("./ui/html/home.html")
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}

	if r.Method == "GET" {
		err := ts.Execute(w, nil)
		if err != nil {
			fmt.Println(err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
	} else {
		DBInstance, err := openDB("./mydatabase.db")
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		rows, err := db.Query("select username, password from users")
		var username string
		var password string
		for rows.Next() {
			rows.Scan(&username, &password)
		}

		if s := r.FormValue("username"); s != "" {
			if s == username {
				fmt.Println(true)
			}
		}
		// Process login (you'll need to add your own logic here)
	}

}
