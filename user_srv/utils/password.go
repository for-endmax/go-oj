package utils

import (
	"crypto/md5"
	"github.com/anaskhan96/go-password-encoder"
	"strings"
)

func EncodePassword(inputPassword string) (salt string, encodedPwd string) {
	options := &password.Options{10, 10000, 50, md5.New}
	salt, encodedPwd = password.Encode(inputPassword, options)
	return
}

func VerifyPassword(inputPassword string, saltandEncodedPwd string) bool {
	data := strings.Split(saltandEncodedPwd, ":")
	options := &password.Options{10, 10000, 50, md5.New}
	return password.Verify(inputPassword, data[0], data[1], options)
}
