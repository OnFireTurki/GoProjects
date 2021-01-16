// @5bub
package main

import (
	"bufio"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttp/fasthttpproxy"
)

var count int = 0
var running bool = false
var userarray []string
var uuid string
var currentAPI api
var config jsonFile
var userList []string
var userLine int = 0
var proxyList []fasthttp.DialFunc
var proxyLine int = 0
var rs int = 0
var xcsr string

const charset = "abcdefghijklmnopqrstuvwxyz" +
	"ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

// api : struct to collect url and payload from apis array
type api struct {
	URL     string `json:"url"`
	Payload string `json:"payload"`
}

// jsonFile : struct to get json file ready
type jsonFile struct {
	Apis    []api `json:"apis"`
	Threads int   `json:"threads"`
}

// ClaimedUser : struct to handle the claim and save it
type ClaimedUser struct {
	Username string
	Payload  string
	Response string
}

// save : save the claim info in txt file
func (c *ClaimedUser) save() {
	text := fmt.Sprintf("Username: %s\nPayload: %s\nResponse: %s\n", c.Username, c.Payload, c.Response)
	f, err := os.Create(c.Username + " info.txt")
	if err != nil {
		fmt.Println(text)
	}
	defer f.Close()
	_, err2 := f.WriteString(text)
	if err2 != nil {
		fmt.Println(text)
	}
}

// StringWithCharset : return random string
func StringWithCharset(length int) string {
	bytes := make([]byte, length)
	_, _ = rand.Read(bytes)
	for i, b := range bytes {
		bytes[i] = charset[b%byte(len(charset))]
	}
	return string(bytes)
}

// rsclu : calculate the r p s and print it with the attempts
func rsclu() {
	for running {
		before := count
		time.Sleep(1 * time.Second)
		rs = count - before
		fmt.Printf("\r[%v] [%v]", count, rs)
	}
}

// reg : reg user function ...
func reg(user *string, Client *fasthttp.Client) bool {
	fmt.Println(*user, "now")
	req := fasthttp.AcquireRequest()
	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseRequest(req)   // <- release the variable before the function return
	defer fasthttp.ReleaseResponse(resp) // <- release the variable before the function return
	req.SetRequestURI(currentAPI.URL)
	req.SetHost("i.instagram.com")
	req.Header.SetMethod("POST")
	req.Header.SetUserAgent("Instagram 35.0.0.20.96 Android (28/9; 420dpi; 1080x1794; Google/google; AOSP on IA Emulator; generic_x86_arm; ranchu; en_US; 95414347)")
	req.Header.Add("Accept-Language", "en-US")
	req.Header.Add("IG-U-DS-USER-ID", StringWithCharset(15))
	req.Header.SetContentType("application/x-www-form-urlencoded; charset=UTF-8")
	req.SetBodyString(fmt.Sprintf(currentAPI.Payload, *user, *user, xcsr))
	Client.Dial = proxyList[proxyLine]
	proxyLine++
	if proxyLine >= len(proxyList)-1 {
		proxyLine = 0
	}
	if err := Client.Do(req, resp); err != nil {
		return reg(user, Client)
	}
	body := string(resp.Body())
	//fmt.Println(body)
	if strings.Contains(body, "account_created\": true,") || strings.Contains(body, "/challenge/") {
		c := ClaimedUser{*user, fmt.Sprintf(currentAPI.Payload, *user, *user, xcsr), body}
		c.save()
		return true
	} else if strings.Contains(body, "account_created\": false,") {
		return false
	} else {
		return reg(user, Client)
	}
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
	// release the variables so the gcg can free the memory {GC}
	Client := &fasthttp.Client{}
	for running {
		if len(userarray[thread]) > 0 {
			user := &userarray[thread]
			//  fmt.Println(*user, "now", thread)
			if !reg(user, Client) {
				count++
				userarray[thread] = ""
			} else {
				fmt.Println(*user, "Claimed")
			}
			time.Sleep(1 * time.Millisecond)
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
		time.Sleep(1 * time.Millisecond)
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

// SetUpProxyList : set up the proxies and create Dialer for each proxy
// Proxy Format : username:password@ip:port
// Proxy Format without auth : ip:port
func SetUpProxyList() {
	Proxies, err := readLines("p.txt")
	if err != nil {
		return
	}
	proxyList = make([]fasthttp.DialFunc, len(Proxies)+3)
	for i := 0; i < len(Proxies)-1; i++ {
		proxyList[i] = fasthttpproxy.FasthttpHTTPDialer(Proxies[i])
	}
}

// main function
func main() {
	runtime.GC()
	config = jsonFile{}
	b := make([]byte, 16)
	fmt.Println(StringWithCharset(10))
	userList, _ = readLines("u.txt")
	SetUpProxyList()
	_, _ = rand.Read(b)
	uuid = fmt.Sprintf("%x-%x-%x-%x-%x",
		b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
	file, _ := ioutil.ReadFile("config.json")
	_ = json.Unmarshal([]byte(file), &config)
	tokenget()
	currentAPI = config.Apis[0]
	fmt.Printf("[U:%v|P:%v] - %v\n", len(userList), len(proxyList), xcsr)
	fmt.Println(len(config.Apis), "api/s")
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

}
