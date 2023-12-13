package web

import (
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
	ts, err := template.ParseFiles("./ui/html/login.html")

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
		usernameInput := r.FormValue("username")
		passwordInput := r.FormValue("password")
		if usernameInput != "" && passwordInput != "" {
			if usernameInput == username && passwordInput == password {
				http.Redirect(w, r, "/", 200)
			}
		}
		// Process login (you'll need to add your own logic here)
	}

}

func HandleRegister(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/user/register" {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}
	ts, err := template.ParseFiles("./ui/html/home.html")
	if err != nil {
		log.Printf("EXECUTING TMPL ERROR %e", err)
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}
	switch r.Method {
	case "GET":
		err = ts.Execute(w, nil)
		if err != nil {
			log.Printf("EXECUTING TMPL ERROR %e", err)
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
		if err.Error() == "database is locked" {
			// Handle the database is locked error
			log.Printf("Database is locked, retrying: %v", err)
			// Implement a retry mechanism or return a specific error
		} else {
			log.Printf("Error adding user: %v", err)
		}
	}

	return err

}

func Home(w http.ResponseWriter, r *http.Request) {
	_, err := w.Write([]byte(`This is a home page`))
	if err != nil {
		return
	}
}
