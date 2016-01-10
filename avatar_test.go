package main

import (
	"flag"
	gomniauthtest "github.com/stretchr/gomniauth/test"
	"io/ioutil"
	"os"
	"path"
	"testing"
)

func TestAuthAvatar(t *testing.T) {

	var authAvatar AuthAvatar
	testUser := &gomniauthtest.TestUser{}
	testUser.On("AvatarURL").Return("", ErrNoAvatarURL)
	testChatUser := &chatUser{User: testUser}
	url, err := authAvatar.GetAvatarURL(testChatUser)
	if err != ErrNoAvatarURL {
		t.Error("AuthAvatar.GetAvatarURL should return ErrNoAvatarURL when no value present")
	}

	testUrl := "http://url-to-gravatar/"
	testUser = &gomniauthtest.TestUser{}
	testChatUser.User = testUser
	testUser.On("AvatarURL").Return(testUrl, nil)
	url, err = authAvatar.GetAvatarURL(testChatUser)
	if err != nil {
		t.Error("AuthAvatar.GetAvatarURL should return no error when value present")
	} else {
		if url != testUrl {
			t.Error("AuthAvatar.GetAvatarURL should return correct URL")
		}
	}
}

func TestFileSystemAvatar(t *testing.T) {

	filename := path.Join(*AvatarPath, "abc.jpg")
	ioutil.WriteFile(filename, []byte{}, 0777)
	defer func() { os.Remove(filename) }()

	var fileSystemAvatar FileSystemAvatar
	user := &chatUser{uniqueID: "abc"}
	url, err := fileSystemAvatar.GetAvatarURL(user)

	if err != nil {
		t.Error("FileSystemAvatar.GetAvatarURL should not return an error")
	}
	if url != *AvatarPath+"abc.jpg" {
		t.Errorf("FileSystemAvatar.GetAvatarURL wrongly returned %s and should have been %sabc.jpg", url, *AvatarPath)
	}
}

func TestMain(m *testing.M) {
	flag.Parse()
	AvatarPath = flag.String("avatarPath", "avatars/", "The path to the folder for the avatar images  This is relative to the location from which \"gochat\" is executed.  Can be absolute.")

	os.Exit(m.Run())
}
