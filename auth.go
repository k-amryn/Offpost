package main

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/url"
	"strings"

	"github.com/dghubble/oauth1"
	"github.com/dghubble/oauth1/twitter"
)

//go:embed keys/twitter.txt
var twitterKeys string

func (instances *allInstances) connectTwitter(instIndex int) {
	keys := strings.Split(string(twitterKeys), "***")

	auth := oauth1.Config{
		ConsumerKey:    keys[0],
		ConsumerSecret: keys[1],
		CallbackURL:    "http://localhost:14859/auth",
		Endpoint:       twitter.AuthorizeEndpoint,
	}

	requestToken, requestSecret, err := auth.RequestToken()
	if err != nil {
		fmt.Println(err)
	}

	authorizationURL, _ := auth.AuthorizationURL(requestToken)
	// handle err
	open(authorizationURL.String())

	instances.mu.Lock()
	instances.authComm = make(chan string)
	instances.mu.Unlock()

	instances.authComm <- "Twitter " + instances.c[instIndex].Name

	info, _ := url.Parse(<-instances.authComm)

	// values is the response from the offpost /auth landing page
	values := info.Query()

	_, denied := values["error"]
	if denied {
		fmt.Println(instances.c[instIndex].Name + ": Twitter connection denied.\n")
		close(instances.authComm)
		return
	}

	accessToken, accessSecret, err := auth.AccessToken(values["oauth_token"][0], requestSecret, values["oauth_verifier"][0])
	if err != nil {
		log.Fatal(err)
	}

	instances.c[instIndex].Platforms["twitter"] = accessToken + "***" + accessSecret + "***" + getTwitterUsername(accessToken, accessSecret)
	instances.saveSettings(false, instances.c)

	wsSend <- ""

	close(instances.authComm)
}

func getTwitterUsername(access, secret string) string {
	keys := strings.Split(string(twitterKeys), "***")

	authConf := oauth1.NewConfig(keys[0], keys[1])
	token := oauth1.NewToken(access, secret)

	client := authConf.Client(oauth1.NoContext, token)

	resp, err := client.Get("https://api.twitter.com/2/users/me")
	if err != nil {
		log.Panic(err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	var uBody2 map[string]map[string]string
	json.Unmarshal(body, &uBody2)

	return uBody2["data"]["username"]
}

func (instances *allInstances) refreshUsernames() {
	for _, e := range instances.c {
		for platform := range e.Platforms {
			switch platform {
			case "twitter":
				if e.Platforms["twitter"] == "no-config" {
					break
				}
				keys := strings.Split(e.Platforms["twitter"], "***")
				username := getTwitterUsername(keys[0], keys[1])
				e.Platforms["twitter"] = keys[0] + "***" + keys[1] + "***" + username
			}
		}
	}
}
