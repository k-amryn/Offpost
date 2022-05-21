package main

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/ChimeraCoder/anaconda"
)

func (instance *instance) makePost() {

	f, err := os.OpenFile("./userdata/"+instance.Name+"_queue.txt", os.O_RDWR, 0666)
	if err != nil {
		fmt.Println(err)
	}
	defer f.Close()

	data, err := ioutil.ReadAll(f)
	if err != nil {
		fmt.Println("error:", err)
	}

	newlineIndex := strings.Index(string(data), "\n")
	firstLine := string(data[:newlineIndex])
	theQueue := strings.Split(firstLine, "***")

	data = data[newlineIndex+1:]
	ioutil.WriteFile("./userdata/"+instance.Name+"_queue.txt", data, 0666)

	instance.appendTxtFile([][]string{{firstLine}}, "posted")

	// process caption
	caption := strings.ReplaceAll(instance.Caption, "%{filename}", getBaseName(theQueue[0]))

	fmt.Print(instance.Name + " is posting: [")
	for i := range theQueue {
		if i == len(theQueue)-1 {
			fmt.Print(filepath.Base(theQueue[i]) + "]\n")
		} else {
			fmt.Print(filepath.Base(theQueue[i]) + "; ")
		}

		if filepath.Ext(theQueue[i]) == ".txt" {
			f, _ := os.Open(theQueue[i])
			reader := bufio.NewReader(f)
			content, _ := ioutil.ReadAll(reader)
			caption += "\n" + string(content)
		}
	}

	postLinks := make(map[string]string)

	if len(instance.Platforms) != 0 {
		for key := range instance.Platforms {
			if instance.Platforms[key] == "no-config" {
				fmt.Println(instance.Name + ": no " + key + " config")
				continue
			}
			switch key {
			case "twitter":
				postLinks["twitter"] = postTweet(theQueue, "", instance.Platforms["twitter"])
				// postLinks["twitter"] = "https://tweetlink/"
			case "facebook":
				// postLinks["facebook"] = postFacebook(theQueue, "", instance.Platforms["facebook"])
				postLinks["facebook"] = "https://facebooklink/"
			default:
				fmt.Println(instance.Name + ": " + key + " is an invalid platform")
			}
		}
	} else {
		fmt.Println(instance.Name + ": no post platforms connected")
	}

	for key := range postLinks {
		fmt.Print(postLinks[key] + "\n")
	}
	instance.ItemsInQueue = instance.countQueueItems()
	fmt.Println(instance.Name, "-", instance.countQueueItems(), "items in queue")
	fmt.Print("\n")
}

func postTweet(filepaths []string, caption string, accessKeys string) string {
	if accessKeys == "no" {
		return "no"
	}

	accessKeyList := strings.Split(accessKeys, "***")
	consumerKeyList := strings.Split(twitterKeys, "***")
	api := anaconda.NewTwitterApiWithCredentials(
		accessKeyList[0], accessKeyList[1], consumerKeyList[0], consumerKeyList[1])

	media := url.Values{}
	mediaList := ""
	howManyImg := 0
	for _, path := range filepaths {
		f, _ := os.Open(path)
		reader := bufio.NewReader(f)
		content, _ := ioutil.ReadAll(reader)

		if filepath.Ext(path) != ".txt" {
			howManyImg++
			encoded := base64.StdEncoding.EncodeToString(content)

			response, err := api.UploadMedia(encoded)

			if err != nil {
				fmt.Printf("errtype: %T\nerr value: %v", err, err)
				return "error"
			}

			if howManyImg == 1 {
				mediaList = response.MediaIDString
			}
			if howManyImg != 1 {
				mediaList += "," + response.MediaIDString
			}
		}
	}
	media.Add("media_ids", mediaList)

	//debuginfo
	// fmt.Println(media)
	// fmt.Printf("response type: %T\nresponse value: %v\n", response.MediaIDString, response.MediaIDString)
	// fmt.Printf("err type: %T\nerr value: %v\n\n", err, err)

	fmt.Println(caption)
	response2, err2 := api.PostTweet(caption, media)

	// this prints the formatted response from Twitter when making a Tweet
	// sinfo := reflect.ValueOf(response2)
	// response2type := sinfo.Type()
	// for i := 0; i < sinfo.NumField); i++ {
	// 	fmt.Printf("%s:\t%v\n", response2type.Field(i).Name, sinfo.Field(i).Interface())
	// }

	if err2 != nil {
		fmt.Printf("errtype: %T\nerr value: %v", err2, err2)
		return "error"
	}

	return "https://twitter.com/" + response2.User.ScreenName + "/status/" + response2.IdStr
}

func postFacebook(filepaths []string, caption string, accessKeys string) string {
	if accessKeys == "no" {
		return "no"
	}

	// keyslist[0] = Page ID, keyslist [1] = Access Token
	keysList := strings.Split(accessKeys, "***")

	// this will upload an image to facebook, and receive a fbid of the image
	client := &http.Client{}

	mediaIDs := ""
	howManyTxt := 0
	howManyImg := 0
	for _, path := range filepaths {
		f, _ := os.Open(path)

		if filepath.Ext(path) == ".txt" {
			content, _ := ioutil.ReadFile(path)
			howManyTxt++
			if howManyTxt == 1 {
				caption += string(content)
			} else {
				caption += "\n" + string(content)
			}
		} else {
			howManyImg++
			body := &bytes.Buffer{}
			writer := multipart.NewWriter(body)

			// example arguments to CreateFormFile are "jpg" and "john.jpg"
			part, err := writer.CreateFormFile(filepath.Ext(path)[1:], filepath.Base(f.Name()))
			if err != nil {
				log.Fatal(err)
			}

			io.Copy(part, f)
			f.Close()
			writer.Close()
			request, err := http.NewRequest("POST",
				"https://graph.facebook.com/"+keysList[0]+"/photos"+
					"?published=false"+
					"&access_token="+keysList[1],
				body)
			if err != nil {
				log.Fatal(err)
			}

			request.Header.Add("Content-Type", writer.FormDataContentType())

			resp, err := client.Do(request)
			if err != nil {
				log.Fatal(err)
			}
			defer resp.Body.Close()

			data, _ := ioutil.ReadAll(resp.Body)
			unmar := make(map[string]string)
			_ = json.Unmarshal(data, &unmar)

			// debug info
			// fmt.Printf("\nBody: %v\n", unmar["id"])

			if howManyImg == 1 {
				mediaIDs = `{"media_fbid":"` + unmar["id"] + `"}`
			}
			if howManyImg != 1 {
				mediaIDs += "," + `{"media_fbid":"` + unmar["id"] + `"}`
			}
		}
	}
	// debug info
	// fmt.Println(mediaIDs)

	// ---------------------------

	// this will post one or more images with their fbids

	body := &bytes.Buffer{}

	caption = strings.ReplaceAll(caption, " ", "%20")
	caption = strings.ReplaceAll(caption, "\n", "%0A")
	resp, err := http.Post("https://graph.facebook.com/"+keysList[0]+"/feed"+
		"?message="+caption+
		"&attached_media=["+mediaIDs+"]"+
		"&access_token="+keysList[1],
		"text", body)
	if err != nil {
		fmt.Println("error, ", err)
	}

	data, _ := ioutil.ReadAll(resp.Body)

	// debug info
	// fmt.Printf("\nBody: %v\n", string(data))

	respMap := make(map[string]string)
	_ = json.Unmarshal(data, &respMap)
	postIDs := strings.Split(respMap["id"], "_")

	postLink := "https://www.facebook.com/permalink.php?story_fbid=" + postIDs[1] + "&id=" + postIDs[0]
	return postLink
}
