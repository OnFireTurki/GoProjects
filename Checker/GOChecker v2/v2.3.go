// @5bub
package main

import (
	"bufio"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"
	
	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttp/fasthttpproxy"
)

var (
  count int = 0
  errorC int = 0
  APICount = 0
  running bool = false
  currentAPI api
  config jsonFile
  sessionIDList []string
  idCount int = 0
  userList []string
  userLine int = 0
  proxyList []fasthttp.DialFunc
  proxyLine int = 0
  rs int = 0
)

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
	Session  string
	Response string
}

// save : save the claim info in txt file
func (c *ClaimedUser) save() {
	text := fmt.Sprintf("Username: %s\nSession ID: %s\nResponse: %s\n", c.Username, c.Session, c.Response)
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
func rsclu(wg *sync.WaitGroup) {
	defer wg.Done()
	for running {
		before := count
		time.Sleep(1 * time.Second)
		rs = count - before
		fmt.Printf("\r[%v] [%v] [%v]     ", count, rs, errorC)
	}
}

// remove : same as templete in C++
func remove[T comparable](arr []T, element T) []T {
	if len(arr) == 0 {
		return nil
	}
	for i, j := range arr {
		if j == element {
			// Found the element now we use it's index to reCreate a list without it
			return append(arr[:i], arr[i+1:]...)
		}
	}
	return arr
}

// move : move user function ...
func move(user *string, Client *fasthttp.Client, bodyCon *string) bool {
	req := fasthttp.AcquireRequest()
	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseRequest(req)   // <- release the variable before the function return
	defer fasthttp.ReleaseResponse(resp) // <- release the variable before the function return
	req.SetRequestURI("https://i.instagram.com/api/v1/accounts/set_username/")
	req.SetHost("i.instagram.com")
	req.Header.SetMethod("POST")
	req.Header.SetUserAgent("Instagram 124.0.0.17.473 Android (28/9; 280dpi; 720x1382; samsung; SM-A105FN; a10; exynos7885; en_US; 192992565)")
	req.Header.Add("Accept-Language", "en-US")
	req.Header.Add("Connection", "close")
	if len(sessionIDList) == 0{
		fmt.Println("Time's up Bro")
		running = false
		return false
	}
	ses := &sessionIDList[idCount]
	idCount++
	if idCount >= len(sessionIDList){
		idCount = 0
	}
	req.Header.Add("Cookie", "sessionid=" + *ses) // sessionid=46765013975%3AFejzuCb9Wo7gdD%3A0
	//req.Header.Add("IG-U-DS-USER-ID", StringWithCharset(15))
	req.Header.SetContentType("application/x-www-form-urlencoded; charset=UTF-8")
	req.SetBodyString("username=" + *user)
	//Client.Dial = proxyList[proxyLine] <           <<<<<<<
	
	proxyLine++
	if proxyLine >= len(proxyList)-1 {
		proxyLine = 0
	}
	if err := Client.Do(req, resp); err != nil {
		return move(user, Client, bodyCon)
	}
	*bodyCon = string(resp.Body())
	if (strings.Contains(*bodyCon , "status\": \"ok") || strings.Contains(*bodyCon , *user)) {
		c := ClaimedUser{*user, *ses, *bodyCon }
		c.save()
    sessionIDList = remove(sessionIDList, *ses)
		fmt.Println("200 - removing claimed session @" + *user)
		if sessionIDList == nil{
			fmt.Println("Time's up Bro")
			running = false
			return true
		}
		return true
	} else if strings.Contains(*bodyCon , "username\":") {
		return false
	} else if strings.Contains(*bodyCon , "login_required") {
		sessionIDList = remove(sessionIDList, *ses)
		fmt.Println("403 - Dead session @" + *user)
		if sessionIDList == nil{
			fmt.Println("Time's up Bro")
			running = false
			return false
		}
		return move(user, Client, bodyCon)
	} else {
		return reg(user, Client, bodyCon)
	}
}

// reg : reg user function ...
func reg(user *string, Client *fasthttp.Client, bodyCon *string) bool {
	req := fasthttp.AcquireRequest()
	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseRequest(req)   // <- release the variable before the function return
	defer fasthttp.ReleaseResponse(resp) // <- release the variable before the function return
	req.SetRequestURI(currentAPI.URL)
	req.SetHost("i.instagram.com")
	req.Header.SetMethod("POST")
	req.Header.SetUserAgent("Instagram 187.0.0.32.120 Android (28/9; 420dpi; 1080x1794; Google/google; AOSP on IA Emulator; generic_x86_arm; ranchu; en_US; 95414347)")
	req.Header.Add("Accept-Language", "en-US")
	req.Header.Add("IG-U-DS-USER-ID", StringWithCharset(15))
	req.Header.SetContentType("application/x-www-form-urlencoded; charset=UTF-8")
	req.SetBodyString(fmt.Sprintf(currentAPI.Payload, *user))
	Client.Dial = proxyList[proxyLine]
	proxyLine++
	if proxyLine >= len(proxyList)-1 {
		proxyLine = 0
		///go APIChanger()
	}
	if err := Client.Do(req, resp); err != nil {
    err = nil
		errorC++
		return reg(user, Client, bodyCon)
	}
	*bodyCon  = string(resp.Body())
	//fmt.Println(body)
	if strings.Contains(*bodyCon , "{\"account_created\": false, \"errors\": {\"phone_number\": [\"This field is required.\"], \"device_id\": [\"This field is required.\"], \"__all__\": [\"You need an email or confirmed phone number.\"]}, \"existing_user\": false, \"status\": \"ok\", \"error_type\": \"required, required, no_contact_point_found\"}") {
		return move(user, Client, bodyCon)
	} else if (strings.Contains(*bodyCon , "taken") || strings.Contains(*bodyCon , "held") || strings.Contains(*bodyCon, "exists")) {
		return false
	} else if strings.Contains(*bodyCon , "wait"){
		errorC++
		return reg(user, Client, bodyCon)
	}
	return false
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

//Test
func TestGC() {
	for running {
		time.Sleep(3 * time.Minute)
		runtime.GC()
		time.Sleep(3 * time.Second)
		fmt.Println("GC Done")
	}
}

// workSub : function to check users
func workSub() {
	// release the variables so the gcg can free the memory {GC}
	Client := &fasthttp.Client{}
	bodyCon := ""
	user := &userList[userLine]
	for running {
		user = &userList[userLine]
		userLine++
		if userLine >= len(userList)-1 {
			userLine = 0
			go APIChanger()
		}
		if !reg(user, Client, &bodyCon) {
			count++
		} else {
			fmt.Println(*user, " - Checked as Available")
		}

		//time.Sleep(1 * time.Second)
	}
}

// ThreadsRun : running threads and walking the list
func ThreadsRun() {
	threads := config.Threads
	for i := 0; i < threads; i++ {
		go workSub()
		// time.Sleep(3 * time.Millisecond)
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
	proxyList = make([]fasthttp.DialFunc, len(Proxies))
	for i := 0; i < len(Proxies); i++ {
		proxyList[i] = fasthttpproxy.FasthttpHTTPDialer(Proxies[i])
	}
}

// Change the current API
func APIChanger() {
	APICount++
	if APICount >= len(config.Apis) {
		APICount = 0
	}
	currentAPI = config.Apis[APICount]
}

// main function
func main() {
	runtime.GC()
	config = jsonFile{}
	userList, _ = readLines("u.txt")
	sessionIDList, _ = readLines("s.txt")
	SetUpProxyList()
	file, _ := ioutil.ReadFile("config.json")
	_ = json.Unmarshal([]byte(file), &config)
	currentAPI = config.Apis[APICount]
	fmt.Printf("[U:%v|P:%v]\n", len(userList), len(proxyList))
	fmt.Println(len(config.Apis), "api/s")
	running = true
	rs = 0
	count = 0
	var wg sync.WaitGroup
	wg.Add(1)
	runtime.GC()
	time.Sleep(3 * time.Second)
	go TestGC()
	go ThreadsRun()
	go rsclu(&wg)
	wg.Wait()
}
