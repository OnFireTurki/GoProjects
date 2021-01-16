// @5bub
package main

import (
	"GOChecker/api"
	"bufio"
	"bytes"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"
)

var count int = 0
var running bool = false
var userarray []string
var uuid string
var config jsonFile
var userList []string
var userLine int = 0
var proxyList []string
var proxyLine int = 0
var apis *api.ClassReq
var rs int = 0
var xcsr string

// account : struct to collect username and password from accounts array
type account struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// jsonFile : struct to get json file ready
type jsonFile struct {
	Accounts []account `json:"accounts"`
	Threads  int       `json:"threads"`
}

// rsclu : clu the r p s and print it with attempt
func rsclu() {
	for running {
		before := count
		time.Sleep(1 * time.Second)
		rs = count - before
		fmt.Printf("\r[%v] [%v]", count, rs)
	}
}

// checkusername : Check user function ...
func checkusername(user *string) bool {
	request, err := http.NewRequest("POST", "https://i.instagram.com/api/v1/users/check_username/", bytes.NewBufferString(fmt.Sprintf("username=%s&_uuid=%s&_csrftoken=%s", *user, apis.UUID, xcsr)))
	request.Host = "i.instagram.com"
	request.Header.Add("User-agent", "Instagram 10.3.2 Android (18/4.3; 320dpi; 720x1280; Xiaomi; HM 1SW; armani; qcom; en_US)")
	request.Header.Add("Accept-Language", "en-US,en;q=0.5")
	request.Header.Add("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	if err != nil {
		return checkusername(user)
	}
	proxy, _ := url.Parse("http://" + proxyList[proxyLine])
	client := &http.Client{Transport: &http.Transport{Proxy: http.ProxyURL(proxy)}, Timeout: 4 * time.Second}
	proxyLine++
	if proxyLine >= len(proxyList)-1 {
		proxyLine = 0
	}
	response, err := client.Do(request)
	if err != nil {
		return checkusername(user)
	}
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return checkusername(user)
	}
	return strings.Contains(string(body), "available\": true,")
}

// readLines : read text file line by line and return as string array
func readLines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

// workSub : function to check users
func workSub(thread int) {
	for running {
		if len(userarray[thread]) > 0 {
			count++

			if checkusername(&userarray[thread]) && apis.Changeuser(&userarray[thread]) {
				apis.SendDiscord(userarray[thread], count)
				running = false

			} else {
				count++
				userarray[thread] = ""
			}
			time.Sleep(5 * time.Millisecond)
		}

	}
}

// ThreadsRun : running threads and walking the list
func ThreadsRun(wg *sync.WaitGroup) {
	defer wg.Done()
	threads := config.Threads
	userarray = make([]string, threads)
	for i := 0; i < threads; i++ {
		userarray[i] = ""
		go workSub(i)
		time.Sleep(3 * time.Millisecond)
	}
	for running {
		for j := 0; j < threads; j++ {
			if len(userarray[j]) <= 0 {
				if userLine >= len(userList)-1 {
					userLine = 0
				}
				userarray[j] = userList[userLine]
				userLine++
			}
		}
		time.Sleep(50 * time.Millisecond)
	}
}

// tokenget : get mid and csrf from instagram website
func tokenget() {
	r, _ := http.Get("https://www.instagram.com/")
	for _, cook := range r.Cookies() {
		if cook.Name == "csrftoken" {
			xcsr = cook.Value
		}

	}
}

// main function
func main() {
	runtime.GC()
	config = jsonFile{}
	b := make([]byte, 16)
	userList, _ = readLines("list.txt")
	proxyList, _ = readLines("proxy.txt")
	_, _ = rand.Read(b)
	uuid = fmt.Sprintf("%x-%x-%x-%x-%x",
		b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
	file, _ := ioutil.ReadFile("config.json")
	_ = json.Unmarshal([]byte(file), &config)
	tokenget()
	fmt.Println(len(config.Accounts), "session/s")
	for i := 0; i < len(config.Accounts); i++ {
		apis = &api.ClassReq{
			Username: config.Accounts[i].Username, Passowrd: config.Accounts[i].Password, UUID: uuid}
		if apis.Login() {
			running = true
			rs = 0
			count = 0
			var wg sync.WaitGroup
			wg.Add(1)
			runtime.GC()
			time.Sleep(3 * time.Second)
			go rsclu()
			go ThreadsRun(&wg)
			wg.Wait()
		} else {
			fmt.Println("Can't login", config.Accounts[i].Username)
		}
	}

}
