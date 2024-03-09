package handlers

import (
	"forum/internal/handlers/comments"
	"forum/internal/handlers/middlewares"
	"forum/internal/handlers/posts"
	"forum/internal/handlers/user"
	"forum/internal/handlers/web"
	"github.com/gorilla/mux"
	"net/http"
	"path/filepath"
)

func (h *Handler) Routes() http.Handler {
	router := mux.NewRouter()

	authHandler := user.NewAuthHandler(h.service, h.templates)
	commentsHandler := comments.NewCommentsHandler(h.service)
	middleware := middlewares.NewMiddleware(h.service)
	webHandler := web.NewWebHandler(h.service, h.templates)
	postHandler := posts.NewPostsHandler(h.service, h.templates)

	router.PathPrefix("/auth/").Handler(http.StripPrefix("/auth", authHandler.Routes()))
	router.PathPrefix("/post/").Handler(http.StripPrefix("/post", middleware.RequireAuthentication(postHandler.Routes())))
	router.PathPrefix("/comment/").Handler(http.StripPrefix("/comment", middleware.RequireAuthentication(commentsHandler.Routes())))

	fileServer := http.FileServer(http.Dir("./ui/static"))
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static", fileServer))
	router.PathPrefix("/").Handler(http.StripPrefix("", webHandler.Routes()))

	return middleware.Authenticate(router)
}

type neuteredFileSystem struct {
	fs http.FileSystem
}

func (nfs neuteredFileSystem) Open(path string) (http.File, error) {
	f, err := nfs.fs.Open(path)
	if err != nil {
		return nil, err
	}

	s, err := f.Stat()
	if s.IsDir() {
		index := filepath.Join(path, "index.html")
		if _, err := nfs.fs.Open(index); err != nil {
			closeErr := f.Close()
			if closeErr != nil {
				return nil, closeErr
			}

			return nil, err
		}
	}

	return f, nil
}

// func rateLimit(next http.HandlerFunc) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		// here
// 		next.ServeHTTP(w, r)
// 	}
// }

// позже изучить

// func wsHandler(w http.ResponseWriter, r *http.Request) {
// 	// проверяем заголовки
// 	if r.Header.Get("Upgrade") != "websocket" {
// 		return
// 	}
// 	if r.Header.Get("Connection") != "Upgrade" {
// 		return
// 	}
// 	k := r.Header.Get("Sec-Websocket-Key")
// 	if k == "" {
// 		return
// 	}

// 	// вычисляем ответ
// 	sum := k + "258EAFA5-E914-47DA-95CA-C5AB0DC85B11"
// 	hash := sha1.Sum([]byte(sum))
// 	str := base64.StdEncoding.EncodeToString(hash[:])

// 	// Берем под контроль соединение https://pkg.go.dev/net/http#Hijacker
// 	hj, ok := w.(http.Hijacker)
// 	if !ok {
// 		return
// 	}
// 	conn, bufrw, err := hj.Hijack()
// 	if err != nil {
// 		return
// 	}
// 	defer conn.Close()

// 	// формируем ответ
// 	bufrw.WriteString("HTTP/1.1 101 Switching Protocols\r\n")
// 	bufrw.WriteString("Upgrade: websocket\r\n")
// 	bufrw.WriteString("Connection: Upgrade\r\n")
// 	bufrw.WriteString("Sec-Websocket-Accept: " + str + "\r\n\r\n")
// 	bufrw.Flush()

// 	// выводим все, что пришло от клиента
// 	buf := make([]byte, 1024)
// 	for {
// 		n, err := bufrw.Read(buf)
// 		if err != nil {
// 			return
// 		}
// 		fmt.Println(buf[:n])
// 	}
// }
