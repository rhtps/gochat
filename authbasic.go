package main

import (
	"crypto/md5"
	"fmt"
	auth "github.com/abbot/go-http-auth"
	"github.com/stretchr/objx"
	"io"
	"log"
	"net/http"
	"strings"
)

func Secret(user, realm string) string {
	if len(user) >= 4 {
		// password is "hello"
		return "$1$dlPL2MqE$oQmn16q49SqdmhenQuNgs1"
	}
	return ""
}

func handleAuthBasic(w http.ResponseWriter, r *auth.AuthenticatedRequest) {

	user := NewBasicUser(r.Username, "12345")
	chatUser := &chatUser{User: user}

	m := md5.New()
	io.WriteString(m, strings.ToLower(user.Name()))
	chatUser.uniqueID = fmt.Sprintf("%x", m.Sum(nil))
	avatarURL, err := avatars.GetAvatarURL(chatUser)

	if err != nil {
		log.Fatalln("Error when trying to GetAvatarURL", "-", err)
	}
	authCookieValue := objx.New(map[string]interface{}{
		"userid":     chatUser.uniqueID,
		"name":       user.Name(),
		"avatar_url": avatarURL,
	}).MustBase64()
	http.SetCookie(w, &http.Cookie{
		Name:  "auth",
		Value: authCookieValue,
		Path:  "/"})

	w.Header()["Location"] = []string{"/chat"}
	w.WriteHeader(http.StatusTemporaryRedirect)

}
