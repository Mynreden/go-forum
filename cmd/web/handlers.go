package web

import (
	"fmt"
	"forum/pkg/db"
	"forum/pkg/models"
	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
	_ "golang.org/x/crypto/bcrypt"
	"html/template"
	"log"
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
			log.Printf("ERROR ExEC")
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
	} else {
		DBInstance := db.GetSingleDBInstance()
		err := DBInstance.OpenDB("./mydatabase.db")
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		rows, err := DBInstance.Db.Query("select username, password from users")
		if err != nil {
			log.Printf("handlers.go 36 line")
		}
		var username string
		var password string
		for rows.Next() {
			err = rows.Scan(&username, &password)
			if err != nil {
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			}
		}

		if s := r.FormValue("username"); s != "" {
			if s == username {
				fmt.Println(true)
			}
		}
		// Process login (you'll need to add your own logic here)
	}

}

func HandleRegister(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}
	ts, err := template.ParseFiles("./ui/html/home.html")
	switch r.Method {
	case "GET":
		err = ts.Execute(w, nil)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
			return
		}
	case "POST":
		passwd := r.FormValue("password")
		email := r.FormValue("email")
		name := r.FormValue("username")
		hashedPaswd, err := bcrypt.GenerateFromPassword([]byte(passwd), 16)
		if err != nil {
			log.Printf("HASH ERROR :%e", err)
		}
		user := &models.User{HashedPw: hashedPaswd, Name: name, Email: email}
		err = addUser(user, db.GetSingleDBInstance())
		if err != nil {
			log.Printf("ADD USER ERROR :%e", err)
		}
	}

}

func addUser(user *models.User, db *db.DB) error {
	query := `insert into users (Username, Password, Email) values ($1, $2, $3)`
	err := db.OpenDB("./mydatabase.db")
	if err != nil {
		log.Printf("ERROR OPENING DB")
	}
	_, err = db.Db.Exec(query, user.Name, user.HashedPw, user.Email)
	if err != nil {
		return err
	}
	return nil
}
