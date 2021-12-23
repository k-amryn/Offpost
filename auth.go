package main

import (
	crand "crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
)

func pkce() (string, string) {
	var vb [64]byte
	io.ReadFull(crand.Reader, vb[:])
	verifier := base64.RawURLEncoding.EncodeToString(vb[:])
	cb := sha256.Sum256([]byte(verifier))
	challenge := base64.RawURLEncoding.EncodeToString(cb[:])
	return challenge, verifier
}

func (instances *allInstances) connectTwitter(i int) {
	twURL := "https://twitter.com/i/oauth2/authorize"
	twURL += "?response_type=code"
	twURL += "&client_id=RWJhQ1NGNGVNTEFYRGd1UUhYaXk6MTpjaQ"
	twURL += "&redirect_uri=http://localhost:14859/auth"
	twURL += "&scope=tweet.write%20tweet.read%20users.read%20offline.access"

	_, state := pkce()
	twURL += "&state=" + state

	ch, ver := pkce()
	twURL += "&code_challenge=" + ch
	twURL += "&code_challenge_method=s256"
	open(twURL)

	instances.mu.Lock()
	instances.authComm = make(chan string)
	instances.mu.Unlock()

	instances.authComm <- "Twitter " + instances.c[i].Name
	info, _ := url.Parse(<-instances.authComm)
	code := info.Query()["code"][0]

	resp, err := http.PostForm("https://api.twitter.com/2/oauth2/token?", url.Values{
		"code":          {code},
		"grant_type":    {"authorization_code"},
		"client_id":     {"RWJhQ1NGNGVNTEFYRGd1UUhYaXk6MTpjaQ"},
		"redirect_uri":  {"http://localhost:14859/auth"},
		"code_verifier": {ver},
	})
	if err != nil {
		log.Panic("Error reaching Twitter for verification", err)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Panic(err)
	}

	var uBody map[string]string
	json.Unmarshal(body, &uBody)

	access := uBody["access_token"]
	refresh := uBody["refresh_token"]

	instances.c[i].Platforms["twitter"] = access + "***" + refresh
	instances.saveSettings(false, instances.c)

	fmt.Println(instances.c[i].Name + ": Connected to twitter.\n")
	close(instances.authComm)
}
