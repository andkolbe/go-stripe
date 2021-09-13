package urlsigner

import (
	"fmt"
	"strings"
	"time"

	goalone "github.com/bwmarrin/go-alone"
)

type Signer struct {
	Secret []byte
}

// takes a token, signs it, and hands back the signed string
func (s *Signer) GenerateTokenFromString(data string) string {
	var urlToSign string

	crypt := goalone.New(s.Secret, goalone.Timestamp)
	// if we are trying to send a url that already has url parameters
	if strings.Contains(data, "?") {
		urlToSign = fmt.Sprintf("%s&hash=", data)
	} else {
		urlToSign = fmt.Sprintf("%s?hash=", data)
	}

	tokenBytes := crypt.Sign([]byte(urlToSign))
	token := string(tokenBytes)
	// the token isn't just a token, it is a fully signed url
	return token
}

// use to verify the link that people click on hasn't been changed in any way
func (s *Signer) VerifyToken(token string) bool {
	crypt := goalone.New(s.Secret, goalone.Timestamp)
	_, err := crypt.Unsign([]byte(token))

	if err != nil {
		fmt.Println(err)
		return false
	}
	return true
}

// checks to see if the token is expired
func (s *Signer) Expired(token string, minutesUntilExpire int) bool {
	crypt := goalone.New(s.Secret, goalone.Timestamp)
	ts := crypt.Parse([]byte(token))
	
	return time.Since(ts.Timestamp) > time.Duration(minutesUntilExpire) * time.Minute
}
