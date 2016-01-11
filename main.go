//The entry point for the gochat program
package main

import (
	"flag"
	auth "github.com/abbot/go-http-auth"
	"github.com/stretchr/gomniauth"
	"github.com/stretchr/gomniauth/providers/facebook"
	"github.com/stretchr/gomniauth/providers/github"
	"github.com/stretchr/gomniauth/providers/google"
	"github.com/stretchr/objx"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"
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
var HtpasswdPath *string

//Primary handler
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

/*Main entry point.  Flags, Handlers, and authentication providers configured here.
*
* Basic usage: gochat -host=0.0.0.0:8080
*
* Help: gochat --help
*
 */
func main() {
	var host = flag.String("host", ":8080", "The host address of the application.")
	var callBackHost = flag.String("callBackHost", "http://localhost:8080", "The host address of the application.")
	templatePath = flag.String("templatePath", "templates/", "The path to the HTML templates.  This is relative to the location from which \"gochat\" is executed.  Can be absolute.")
	AvatarPath = flag.String("avatarPath", "avatars/", "The path to the folder for the avatar images.  This is relative to the location from which \"gochat\" is executed.  Can be absolute.")
	HtpasswdPath = flag.String("htpasswdPath", os.Getenv("CHAT_PASSWORD_FILE"), "The path to the htpasswd file for basic auth.  This is relative to the location from which \"gochat\" is executed.  Can be absolute.  By default, this is set to the CHAT_PASSWORD_FILE environment variable.")
	var omniSecurityKey = flag.String("securityKey", "OCc4D4FLADtymUgqI4ircKQ8OJhsWMuEasbPCuJH9KpNayjjeoe3U7hVVKfgwtRG", "The OAuth security key.")
	var facebookProviderKey = flag.String("facebookProviderKey", os.Getenv("KEY_FACEBOOK"), "The FaceBook OAuth provider key.")
	var facebookProviderSecretKey = flag.String("facebookProviderSecretKey", os.Getenv("SECRET_KEY_FACEBOOK"), "The FaceBook OAuth provider secret key.")
	var githubProviderKey = flag.String("githubProviderKey", os.Getenv("KEY_GITHUB"), "The GitHub OAuth provider key.")
	var githubProviderSecretKey = flag.String("githubProviderSecretKey", os.Getenv("SECRET_KEY_GITHUB"), "The GitHub OAuth provider secret key.")
	var googleProviderKey = flag.String("googleProviderKey", os.Getenv("KEY_GOOGLE"), "The Google OAuth provider key.")
	var googleProviderSecretKey = flag.String("googleProviderSecretKey", os.Getenv("SECRET_KEY_GOOGLE"), "The Google OAuth provider secret key.")
	flag.Parse()

	//set up gomniauth
	gomniauth.SetSecurityKey(*omniSecurityKey)
	gomniauth.WithProviders(
		facebook.New(*facebookProviderKey, *facebookProviderSecretKey, *callBackHost+"/auth/callback/facebook"),
		github.New(*githubProviderKey, *githubProviderSecretKey, *callBackHost+"/auth/callback/github"),
		google.New(*googleProviderKey, *googleProviderSecretKey, *callBackHost+"/auth/callback/google"),
	)

	r := newRoom()
	http.Handle("/", MustAuth(&templateHandler{filename: "chat.html"}))
	http.Handle("/chat", MustAuth(&templateHandler{filename: "chat.html"}))
	http.Handle("/login", &templateHandler{filename: "login.html"})
	http.HandleFunc("/auth/", loginHandler)
	http.Handle("/room", r)
	http.HandleFunc("/logout", func(w http.ResponseWriter, r *http.Request) {
		UseOmniAuth = false
		http.SetCookie(w, &http.Cookie{
			Name:   "auth",
			Value:  "",
			Path:   "/",
			MaxAge: -1,
		})
		w.Header()["Location"] = []string{"/logoutpage"}
		w.WriteHeader(http.StatusTemporaryRedirect)
	})
	http.Handle("/logoutpage", &templateHandlerBasicAuth{filename: "logoutpage.html"})
	http.Handle("/upload", &templateHandler{filename: "upload.html"})
	http.Handle("/uploadhtpasswd", &templateHandler{filename: "upload-passwd.html"})
	http.HandleFunc("/uploader", uploadHandler)
	http.HandleFunc("/uploaderpasswd", uploadHtpasswdHandler)
	http.Handle("/avatars/", http.StripPrefix("/avatars/", http.FileServer(http.Dir(*AvatarPath))))

	if len(*HtpasswdPath) > 0 {
		secret := auth.HtpasswdFileProvider(*HtpasswdPath)
		authenticator := auth.NewBasicAuthenticator("gochat", secret)
		http.HandleFunc("/authbasic", authenticator.Wrap(handleAuthBasic))
	} else {
		authenticator := auth.NewBasicAuthenticator("gochat", Secret)
		http.HandleFunc("/authbasic", authenticator.Wrap(handleAuthBasic))
	}

	go r.run()
	log.Println("Starting the web server on", *host)
	if err := http.ListenAndServe(*host, nil); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}

}
