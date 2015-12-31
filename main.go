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
	"os"
) 
var avatars Avatar = TryAvatars{
	UseFileSystemAvatar,
	UseAuthAvatar}

type templateHandler struct {
	once     sync.Once
	filename string
	templ    *template.Template
}

func (t *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	t.once.Do(func() {
		t.templ = template.Must(template.ParseFiles(filepath.Join(os.Getenv("TEMPLATE_PATH"), t.filename)))
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
	gomniauth.SetSecurityKey(os.Getenv("SECURITY_KEY"))
	gomniauth.WithProviders(
		facebook.New(os.Getenv("FACEBOOK_PROVIDER_KEY"), os.Getenv("FACEBOOK_PROVIDER_SECRET_KEY"), "http://"+os.Getenv("HOST_CALLBACK")+":8080/auth/callback/facebook"),
		github.New(os.Getenv("GITHUB_PROVIDER_KEY"), os.Getenv("GITHUB_PROVIDER_SECRET_KEY"), "http://"+os.Getenv("HOST_CALLBACK")+":8080/auth/callback/github"),
		google.New(os.Getenv("GOOGLE_PROVIDER_KEY"), os.Getenv("GOOGLE_PROVIDER_SECRET_KEY"), "http://"+os.Getenv("HOST_CALLBACK")+":8080/auth/callback/google"),
		)
	
	r := newRoom()
	http.Handle("/", MustAuth(&templateHandler{filename: "chat.html"}))
	http.Handle("/chat", MustAuth(&templateHandler{filename: "chat.html"}))
	http.Handle("/login", &templateHandler{filename: "login.html"})
	http.HandleFunc("/auth/", loginHandler)
	http.Handle("/room", r)
	http.HandleFunc("/logout", func(w http.ResponseWriter, r *http.Request) {
					http.SetCookie(w, &http.Cookie{
						Name: "auth",
						Value: "",
						Path: "/",
						MaxAge: -1,
					})
					w.Header()["Location"] = []string{"/string"}
					w.WriteHeader(http.StatusTemporaryRedirect)
	})
	http.Handle("/upload", &templateHandler{filename: "upload.html"})
	http.HandleFunc("/uploader", uploadHandler)
	http.Handle("/avatars/", http.StripPrefix("/avatars/", http.FileServer(http.Dir("./avatars"))))
	go r.run()
	log.Println("Starting the web server on", *addr)
	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}

}