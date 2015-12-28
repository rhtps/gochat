package main

import (
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"sync"
	"flag"
	"github.com/stretchr/gomniauth"
	"github.com/stretchr/gomniauth/providers/facebook"
	"github.com/stretchr/gomniauth/providers/github"
	"github.com/stretchr/gomniauth/providers/google"
	"github.com/stretchr/objx"
) 

type templateHandler struct {
	once     sync.Once
	filename string
	templ    *template.Template
}

func (t *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	t.once.Do(func() {
		t.templ = template.Must(template.ParseFiles(filepath.Join("templates", t.filename)))
	})
	data := map[string]interface{}{
		"Host": r.Host,
	}
	
	if authCookie, err := r.Cookie("auth"); err == nil {
		data["UserData"] = objx.MustFromBase64(authCookie.Value)
	}
	t.templ.Execute(w, data)

}

func main() {
	var addr = flag.String("addr", ":8080", "The addr of the application.")
	flag.Parse()
	
	//set up gomniauth
	gomniauth.SetSecurityKey("e19c1283c925b3206685ff522acfe3e6")
	gomniauth.WithProviders(
		facebook.New("key", "secret", "http://localhost:8080/auth/callback/facebook"),
		github.New("bae0d1c04aa6c419d169", "1724f952f0ea09131468188197dd5558680ce14d", "http://localhost:8080/auth/callback/github"),
		google.New("key", "secret", "http://localhost:8080/auth/callback/google"),
		)
	
	r := newRoom()
	http.Handle("/", MustAuth(&templateHandler{filename: "chat.html"}))
	http.Handle("/chat", MustAuth(&templateHandler{filename: "chat.html"}))
	http.Handle("/login", &templateHandler{filename: "login.html"})
	http.HandleFunc("/auth/", loginHandler)
	http.Handle("/room", r)
	go r.run()
	log.Println("Starting the web server on", *addr)
	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}

}
