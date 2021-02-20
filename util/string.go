package util

import (
	"math/rand"
	"regexp"
	"time"
)

const (
	MailPattern = "[a-zA-Z0-9_.+-]+@[a-zA-Z0-9-]+\\.[a-zA-Z0-9-.]+"
)

var (
	mailRegex *regexp.Regexp
)

func init() {
	mailRegex = regexp.MustCompile(MailPattern)
}

func RandomString(n int) string {
	gen := rand.New(rand.NewSource(time.Now().UnixNano()))
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	s := make([]rune, n)
	for i := range s {
		s[i] = letters[gen.Intn(len(letters))]
	}
	return string(s)
}

func IsEmail(s string) bool {
	return mailRegex.Match([]byte(s))
}
