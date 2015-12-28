package main

import (
	"net/http"
	"strings"
	"log"
	"fmt"
	"github.com/stretchr/gomniauth"
	"github.com/stretchr/objx" 
)

type authHandler struct {
	next http.Handler
}

func (h *authHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if _, err := r.Cookie("auth"); err == http.ErrNoCookie {
		//not authenticated
		w.Header().Set("Location", "/login")
		w.WriteHeader(http.StatusTemporaryRedirect)
	} else if err != nil {
		panic(err.Error())
	} else {
		h.next.ServeHTTP(w, r)
	}
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	segs := strings.Split(r.URL.Path, "/")
	action := segs[2]
	provider := segs[3]
	
	log.Println("Action", action, "Provider", provider)
	switch action {
		case "login":
			provider, err := gomniauth.Provider(provider)
			if err != nil {
				log.Fatalln("Error when trying to get provider", provider, "-", err)
			}
			loginUrl, err := provider.GetBeginAuthURL(nil, nil)
			if err != nil {
				log.Fatalln("Error when trying to GetBeginAuthUrl for", provider, "-", err)
			}
			w.Header().Set("Location", loginUrl)
			w.WriteHeader(http.StatusTemporaryRedirect)
		case "callback":
			provider, err := gomniauth.Provider(provider)
			if err != nil {
				log.Fatalln("Error when trying to get provider", provider,"-", err)
			}
			creds, err := provider.CompleteAuth(objx.MustFromURLQuery(r.URL.RawQuery))
			if err != nil {
				log.Fatalln("Error when trying to complete auth for", provider,"-",err)
			}
			user, err := provider.GetUser(creds)
			if err !=nil {
				log.Fatalln("Error when trying to get user from",provider,"-",err)
			}	
			authCookieValue := objx.New(map[string]interface{}{
					"name": user.Name(),
					}).MustBase64()
			http.SetCookie(w, &http.Cookie{
					Name: "auth",
					Value: authCookieValue,
					Path: "/"})
			w.Header()["Location"] = []string{"/chat"}
			w.WriteHeader(http.StatusTemporaryRedirect)
		default:
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprintf(w, "Auth action %s is not supported", action)
			
	}
}

func MustAuth(handler http.Handler) http.Handler {
	return &authHandler{next: handler}
}
