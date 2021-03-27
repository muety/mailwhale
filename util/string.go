package util

import (
	"crypto/md5"
	"encoding/binary"
	"io"
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
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	return randomString(n, rnd)
}

func RandomStringSeeded(n int, seed string) string {
	md5h := md5.New()
	io.WriteString(md5h, seed)
	var seedInt uint64 = binary.BigEndian.Uint64(md5h.Sum(nil))
	rnd := rand.New(rand.NewSource(int64(seedInt)))
	return randomString(n, rnd)
}

func randomString(n int, rnd *rand.Rand) string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	s := make([]rune, n)
	for i := range s {
		s[i] = letters[rnd.Intn(len(letters))]
	}
	return string(s)
}

func IsEmail(s string) bool {
	return mailRegex.Match([]byte(s))
}
