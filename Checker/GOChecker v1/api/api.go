package api

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"regexp"
	"strings"
	"time"
)

var APIURL string = "https://i.instagram.com/api/v1/"
var Cookies string
var Client *http.Client = &http.Client{}
var Email string

// new class for send and reseve requests
type ClassReq struct {
	Username string
	Passowrd string
	UUID     string
}

// Login to account and get cookies
func (cr *ClassReq) Login() bool {
	rand.Seed(time.Now().UnixNano())
	request, err := http.NewRequest("POST", APIURL+"accounts/login/", bytes.NewBufferString(fmt.Sprintf("username=%s&password=%s&device_id=%s&phone_id=%s&_csrftoken=missing&login_attempt_count=0", cr.Username, cr.Passowrd, cr.UUID, cr.UUID)))
	request.Host = "i.instagram.com"
	request.Header.Add("User-agent", "Instagram 10.3.2 Android (18/4.3; 320dpi; 720x1280; Xiaomi; HM 1SW; armani; qcom; en_US)")
	request.Header.Add("Accept-Language", "en-US,en;q=0.5")
	request.Header.Add("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	if err != nil {
		return false
	}
	response, err := Client.Do(request)
	if err != nil {
		return false
	}
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return false
	}
	if strings.Contains(string(body), "logged_in_user") {
		cookie := response.Cookies()
		cookie2 := []string{}
		for _, Cook := range cookie {
			cookie2 = append(cookie2, Cook.Name+"="+Cook.Value)
		}
		Cookies = strings.Join(cookie2, "; ")
		// Functoin to Get account Email basic scope Group {}
		func(s *string) bool {
			request, err := http.NewRequest("GET", APIURL+"accounts/current_user/?edit=true", nil)
			request.Host = "i.instagram.com"
			request.Header.Add("User-agent", "Instagram 10.3.2 Android (18/4.3; 320dpi; 720x1280; Xiaomi; HM 1SW; armani; qcom; en_US)")
			request.Header.Add("Accept-Language", "en-US,en;q=0.5")
			request.Header.Add("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
			request.Header.Add("Cookie", Cookies)
			if err != nil {
				return false
			}
			response, err := Client.Do(request)
			if err != nil {
				return false
			}
			body, err := ioutil.ReadAll(response.Body)

			if err != nil {
				return false
			}
			regex, _ := regexp.Compile("email\": \"(.*?)\", ")
			*s = regex.FindStringSubmatch(string(body))[1]
			return true
		}(&Email)
		if Email == "" {
			return false
		} else {
			return true
		}
	} else {
		return false
	}
}

// Changeuser to change account username simple and fast
func (cr *ClassReq) Changeuser(username *string) bool {
	fmt.Println(*username)
	request, err := http.NewRequest("POST", APIURL+"accounts/edit_profile/", bytes.NewBufferString(fmt.Sprintf("username=%s&email=%s", *username, Email)))
	request.Host = "i.instagram.com"
	request.Header.Add("User-agent", "Instagram 10.3.2 Android (18/4.3; 320dpi; 720x1280; Xiaomi; HM 1SW; armani; qcom; en_US)")
	request.Header.Add("Accept-Language", "en-US,en;q=0.5")
	request.Header.Add("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	request.Header.Add("Cookie", Cookies)
	if err != nil {
		return false
	}
	response, err := Client.Do(request)
	if err != nil {
		return false
	}
	if response.StatusCode == 200 {
		return true
	}
	return false
}

// randomInt  just return a random number
func randomInt(min, max int) int {
	return min + rand.Intn(max-min)
}

// SendDiscord Function to send message to the channel {new - claimed} useing webhook
func (cr *ClassReq) SendDiscord(username string, attempt int) {
	// REMOVED
}
