package main
import (
	"errors"
	"io/ioutil"
	"path"
	)

var ErrNoAvatarURL = errors.New("chat: Unable to get an avatar URL.")

type AuthAvatar struct{}
type FileSystemAvatar struct{}

var UseAuthAvatar AuthAvatar
var UseFileSystemAvatar FileSystemAvatar

func (_ AuthAvatar) GetAvatarURL(u ChatUser) (string, error) {
	
	
	url := u.AvatarURL()
	if len(url) > 0 {
		return url, nil
	}
	
	return *AvatarPath+"default.jpg", nil
}

func (_ FileSystemAvatar) GetAvatarURL(u ChatUser) (string, error) {
	if files, err := ioutil.ReadDir(*AvatarPath); err == nil {
		for _, file := range files {
			if file.IsDir() {
				continue
			}
			if match, _ := path.Match(u.UniqueID()+"*", file.Name()); match {
				return *AvatarPath + file.Name(), nil
			}
		}
	}
	return *AvatarPath+"default.jpg", ErrNoAvatarURL
}

type Avatar interface {
	GetAvatarURL(u ChatUser) (string, error)

}

type TryAvatars []Avatar

func (a TryAvatars) GetAvatarURL(u ChatUser) (string, error) {
	for _, avatar := range a {
		if url, err := avatar.GetAvatarURL(u); err == nil {
			return url, nil
		}
	}
	return "", ErrNoAvatarURL
}
