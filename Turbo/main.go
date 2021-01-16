package main

import (
	"bufio"
	"bytes"
	"crypto/rand"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const (
	InfoColor    = "\033[1;34m%s\033[0m"
	NoticeColor  = "\033[1;36m%s\033[0m"
	WarningColor = "\033[1;33m%s\033[0m"
	ErrorColor   = "\033[1;31m%s\033[0m"
	SuccessColor = "\033[1;32m%s\033[0m"
)

var sendid string
var email string
var uuid string
var cookie []*http.Cookie
var count int = 0
var rs int = 0

func reqpers() {
	for {
		before := count
		time.Sleep(time.Second * 1)
		rs = count - before
	}

}
func login(user string, pass string) bool {
	Postdata := fmt.Sprintf("username=%s&password=%s&_uuid=%s&device_id=%s&form_reg=false&login_attempt_count=0&_csrftoken=missing", user, pass, uuid, uuid)

	var responser string
	client := http.Client{}
	request, err := http.NewRequest("POST", "https://i.instagram.com/api/v1/accounts/login/", bytes.NewBufferString(Postdata))
	request.Host = "i.instagram.com"
	request.Header.Add("User-Agent", "Instagram 22.0.0.15.68 Android (24/5.0; 515dpi; 1440x2416; huawei/google; Nexus 6P; angler; angler; en_US)")
	request.Header.Add("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	request.Header.Add("Accept-Language", "en;q=0.9")
	if err != nil {
		fmt.Println("1 err")
		responser = fmt.Sprintf("%s", err)
	}

	resp, err := client.Do(request)
	if err != nil {
		fmt.Println("2 err")
		responser = fmt.Sprintf("%s", err)
		return false

	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("3 err")
		responser = fmt.Sprintf("%s", err)
		return false
	}
	responser = string(body)
	if strings.Contains(responser, "logged_in_user") {
		cookie = resp.Cookies()
		return true
	}
	return false
}
func change(user string) bool {
	Postdata := fmt.Sprintf("username=%s&first_name=Go&biography=@5bub&_uuid=%s&_csrftoken=missing&gender=1&email=%s", user, uuid, email)

	var responser string
	client := http.Client{}
	request, err := http.NewRequest("POST", "https://i.instagram.com/api/v1/accounts/edit_profile/", bytes.NewBufferString(Postdata))
	request.Host = "i.instagram.com"
	for _, Cook := range cookie {

		request.AddCookie(&http.Cookie{Name: Cook.Name, Value: Cook.Value})
	}
	request.Header.Add("User-Agent", "Instagram 22.0.0.15.68 Android (24/5.0; 515dpi; 1440x2416; huawei/google; Nexus 6P; angler; angler; en_US)")
	request.Header.Add("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	request.Header.Add("Accept-Language", "en;q=0.9")
	if err != nil {
		fmt.Println("1 err")
		responser = fmt.Sprintf("%s", err)
	}

	resp, err := client.Do(request)
	if err != nil {
		fmt.Println("2 err")
		responser = fmt.Sprintf("%s", err)
		return false

	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("3 err")
		responser = fmt.Sprintf("%s", err)
		return false
	}
	responser = string(body)

	if strings.Contains(responser, "status\": \"ok") {
		fmt.Printf("\n"+SuccessColor+"\n", "Username Changed @"+user)
		//send(user)
		return true
	}
	fmt.Printf("\n"+ErrorColor+"\n", "Username Escaped @"+user)
	return false
}

func send(user string) {
	Postdata := fmt.Sprintf("recipient_users=[[%s]]&thread=0&text=%s&client_context=%s", sendid, "Username Changed @"+user, strings.ToUpper(uuid))

	client := http.Client{}
	request, err := http.NewRequest("POST", "https://i.instagram.com/api/v1/direct_v2/threads/broadcast/text/", bytes.NewBufferString(Postdata))
	request.Host = "i.instagram.com"
	for _, Cook := range cookie {

		request.AddCookie(&http.Cookie{Name: Cook.Name, Value: Cook.Value})
	}
	request.Header.Add("User-Agent", "Instagram 22.0.0.15.68 Android (24/5.0; 515dpi; 1440x2416; huawei/google; Nexus 6P; angler; angler; en_US)")
	request.Header.Add("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	request.Header.Add("Accept-Language", "en;q=0.9")
	if err != nil {
		fmt.Println("1 err")

	}

	_, err = client.Do(request)
	if err != nil {
		fmt.Println("2 err")

	}
}

func getemail() bool {

	client := http.Client{}
	request, err := http.NewRequest("GET", "https://i.instagram.com/api/v1/accounts/current_user/?edit=true", nil)
	request.Host = "i.instagram.com"
	for _, Cook := range cookie {

		request.AddCookie(&http.Cookie{Name: Cook.Name, Value: Cook.Value})
	}
	request.Header.Add("User-Agent", "Instagram 22.0.0.15.68 Android (24/5.0; 515dpi; 1440x2416; huawei/google; Nexus 6P; angler; angler; en_US)")
	request.Header.Add("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	request.Header.Add("Accept-Language", "en;q=0.9")
	if err != nil {
		fmt.Println("1 err")
	}
	resp, err := client.Do(request)
	if err != nil {
		fmt.Println("2 err")
		return false
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("3 err")
		return false
	}
	re, _ := regexp.Compile("email\": \"(.*?)\",")
	email = re.FindStringSubmatch(string(body))[1]
	println(string(body))
	return false

}
func targets(user string) bool {

	client := http.Client{}
	request, err := http.NewRequest("GET", "https://i.instagram.com/api/v1/usertags/"+user+"/feed/username/", nil)
	request.Host = "i.instagram.com"
	for _, Cook := range cookie {

		request.AddCookie(&http.Cookie{Name: Cook.Name, Value: Cook.Value})
	}
	request.Header.Add("User-Agent", "Instagram 22.0.0.15.68 Android (24/5.0; 515dpi; 1440x2416; huawei/google; Nexus 6P; angler; angler; en_US)")
	request.Header.Add("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	request.Header.Add("Accept-Language", "en;q=0.9")
	if err != nil {
		fmt.Println("1 err")

	}

	resp, err := client.Do(request)
	if err != nil {
		fmt.Println("2 err")

		return false

	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("3 err")

		return false
	}

	return strings.Contains(string(body), "must specify valid username")

}
func work(target string) {

	for {
		if targets(target) {
			if change(target) {
				break
			} else {
				os.Exit(1)
			}
		} else {
			count++
			time.Sleep(time.Millisecond * 150)
		}
	}
	os.Exit(0)

}
func getid(user string) string {
	ur, _ := http.Get(fmt.Sprintf("https://www.instagram.com/%s/?__a=1", user))
	re, _ := regexp.Compile("profilePage_(.*?)\",\"")
	body, _ := ioutil.ReadAll(ur.Body)

	return re.FindStringSubmatch(string(body))[1]

}

func input(reader *bufio.Reader) string {
	shit, _ := reader.ReadString('\n')
	shit = strings.Replace(shit, "\r", "", -1)
	shit = strings.Replace(shit, "\n", "", -1)
	return shit
}

func main() {
	fmt.Println("Hello // @5bub")

	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		log.Fatal(err)
	}
	uuid = fmt.Sprintf("%x-%x-%x-%x-%x",
		b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Username : ")
	user := input(reader)
	fmt.Print("Password : ")
	pass := input(reader)
	if login(user, pass) {
		fmt.Printf(SuccessColor+"\n", "Logged in")
		getemail()
		if email != "" {
			fmt.Printf(InfoColor+" %s \n", "email :", email)
		} else {
			fmt.Println("no email found")
			os.Exit(1)
		}
		fmt.Print("Target : ")
		target := input(reader)
		fmt.Printf("DM To  : "+NoticeColor+"\n", "EMPTy")
		//dm,_ := reader.ReadString('\n')
		sendid = "0"
		fmt.Print("Threads : ")
		threads := input(reader)
		thread, _ := strconv.Atoi(threads)
		var th int = thread
		go reqpers()
		for i := 0; i < th; i++ {
			//go work(target)
			go work(target)
			time.Sleep(time.Second)
		}

		for {
			fmt.Printf("\rAttempt: %d | rs: %d | Target: %s", count, rs, target)
			time.Sleep(time.Millisecond * 700)
		}
	} else {
		fmt.Printf(ErrorColor+"\n", "Wrong info / Or / Secure")
	}
}
