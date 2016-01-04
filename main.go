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
	auth "github.com/abbot/go-http-auth"
) 
var avatars Avatar = TryAvatars{
	UseFileSystemAvatar,
	UseAuthAvatar}

type templateHandler struct {
	once     sync.Once
	filename string
	templ    *template.Template
}

var templatePath *string
var AvatarPath *string


func (t *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t.once.Do(func() {
		t.templ = template.Must(template.ParseFiles(filepath.Join(*templatePath, t.filename)))
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
	var host = flag.String("host", ":8080", "The host address of the application.")
	var callBackHost = flag.String("callBackHost", "localhost:8080", "The host address of the application.")
	templatePath = flag.String("templatePath", "templates/", "The path to the HTML templates.  This is relative to the location from which \"gochat\" is executed.  Can be absolute.")
	AvatarPath = flag.String("templatePath", "./avatars", "The path to the folder for the avatar images  This is relative to the location from which \"gochat\" is executed.  Can be absolute.")
	var omniSecurityKey = flag.String("securityKey", "12345", "The OAuth security key.")
	var facebookProviderKey = flag.String("facebookProviderKey", "12345", "The FaceBook OAuth provider key.")
	var facebookProviderSecretKey = flag.String("facebookProviderSecretKey", "12345", "The FaceBook OAuth provider secret key.")
	var githubProviderKey = flag.String("githubProviderKey", "12345", "The GitHub OAuth provider key.")
        var githubProviderSecretKey = flag.String("githubProviderSecretKey", "12345", "The GitHub OAuth provider secret key.")
	var googleProviderKey = flag.String("googleProviderKey", "12345", "The Google OAuth provider key.")
        var googleProviderSecretKey = flag.String("googleProviderSecretKey", "12345", "The Google OAuth provider secret key.")
	flag.Parse()
	
	//set up gomniauth
	gomniauth.SetSecurityKey(*omniSecurityKey)
	gomniauth.WithProviders(
		facebook.New(*facebookProviderKey, *facebookProviderSecretKey, "http://"+*callBackHost+"/auth/callback/facebook"),
		github.New(*githubProviderKey, *githubProviderSecretKey, "http://"+*callBackHost+"/auth/callback/github"),
		google.New(*googleProviderKey, *googleProviderSecretKey, "http://"+*callBackHost+"/auth/callback/google"),
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
	http.Handle("/avatars/", http.StripPrefix("/avatars/", http.FileServer(http.Dir(*AvatarPath))))
	
	authenticator := auth.NewBasicAuthenticator("example.com", Secret)
	http.HandleFunc("/authbasic", authenticator.Wrap(handleAuthBasic))
	
	go r.run()
	log.Println("Starting the web server on", *host)
	if err := http.ListenAndServe(*host, nil); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}

}
