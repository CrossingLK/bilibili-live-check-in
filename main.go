package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

const URL = "https://api.live.bilibili.com/msg/send"

var (
	cookie = os.Getenv("COOKIE")
	csrf   = os.Getenv("CSRF")
	roomID = os.Getenv("ROOM_ID")
)

func main() {

	log.SetFlags(log.Lshortfile | log.LstdFlags)

	if cookie == "" || roomID == "" || csrf == "" {
		log.Print("no configuration was read, please check the configuration")
		return
	}

	send(cookie, roomID, csrf)
}

func send(cookie, roomID, csrf string) {

	values := url.Values{
		"bubble":     {"0"},
		"msg":        {"打卡"},
		"color":      {"5566168"},
		"mode":       {"1"},
		"fontsize":   {"25"},
		"rnd":        {strconv.FormatInt(time.Now().Unix(), 10)},
		"roomid":     {roomID},
		"csrf":       {csrf},
		"csrf_token": {csrf},
	}

	request, err := http.NewRequest(http.MethodPost, URL, strings.NewReader(values.Encode()))
	if err != nil {
		log.Print(err)
		return
	}

	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	request.Header.Set("cookie", cookie)

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		log.Print(err)
		return
	}

	defer response.Body.Close()
	responseBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Print(err)
		return
	}

	if response.StatusCode != http.StatusOK {
		log.Print(string(responseBytes))
		return
	}

	data := make(map[string]interface{})
	err = json.Unmarshal(responseBytes, &data)
	if err != nil {
		log.Print(err)
		return
	}

	if data["code"].(float64) != 0 {
		log.Print(string(responseBytes))
		return
	}

	log.Print("message sent successfully")
}
