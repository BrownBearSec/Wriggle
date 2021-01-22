package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

var whitelist []string
var blacklist []string
var foundURLs []string
var maxTimeout int

func help() {
	fmt.Println("Displaying the help page")
	fmt.Println("Usage: wriggle -w <FILE>")
	fmt.Println("-w <FILE> : Specificy a list of domains in scope, one per line. Note: do not right '*.domain.com' just write 'domain.com' ")
	fmt.Println("-b <FILE> : Specificy a list of domains not in scope, one per line ")
	fmt.Println("-t <number> : Set the max timeout (in seconds) for connecting to a URL, default 20 seconds")
	fmt.Println("-h : Display this help page")
	os.Exit(3)
}

func getHREFfromURL(url string) {

	client := http.Client{
		Timeout: time.Duration(maxTimeout) * time.Second,
	}

	resp, err := client.Get(url)
	if err != nil {
		fmt.Println("GetXYZ")
		if strings.Contains(err.Error(), "Client.Timeout") {
			fmt.Println("The get request has timed out, either increase max timeout or check if the site is up")
			return
		}
		fmt.Println(err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("readingXYZ")
		panic(err)
	}

	defer resp.Body.Close()

	htmlStrig := string(body)

	if resp.StatusCode == 404 {
		log.Fatal("Status code 404 on url", url)
	}

	htmlArray := strings.Split(htmlStrig, "href=\"")

	fmt.Println("---------------")

	for i := 0; i < len(htmlArray); i++ {
		htmlPortion := string(htmlArray[i])
		splitBySpeechmarkArray := strings.Split(htmlPortion, "\"")
		foundURL := splitBySpeechmarkArray[0]

		if len(foundURL) > 1 {
			if foundURL[0] == '/' && foundURL[1] != '/' {
				url = url + foundURL
			}
		}

		isInScope := checkScope(foundURL)
		if isInScope {
			fmt.Println(foundURL, " is in scope")
		}
	}
}

func checkScope(urlToCheck string) bool {
	var whitelisted bool = false
	for i := 0; i < len(whitelist); i++ {
		if strings.Contains(urlToCheck, whitelist[i]) {
			whitelisted = true
		}
	}

	for i := 0; i < len(blacklist); i++ {
		if strings.Contains(urlToCheck, blacklist[i]) {
			whitelisted = false
		}
	}

	return whitelisted
}

func main() {

	wantHelp := flag.Bool("h", false, "display help page")
	whitelistFile := flag.String("w", "", "the whitelist file for domains")
	blacklistFile := flag.String("b", "", "the blacklist file for domains")
	maxTimeoutOption := flag.String("t", "20", "max timeout for connection timeouts")
	flag.Parse()

	if *wantHelp {
		help()
	}

	if len(os.Args) == 1 {
		fmt.Println("Usage : wriggle -w <WhitelistFile> ")
		fmt.Println("For more options do -h")
		os.Exit(3)
	}

	if *whitelistFile == "" {
		fmt.Println("You have not supplied a whitelist file of domains")
		os.Exit(3)
	} else {
		fmt.Println("You have selected", *whitelistFile, "as the list of domains to whitelist")
	}

	if *blacklistFile == "" {
		fmt.Println("WARNING: no blacklisted domains/subdomains, continuing")
	} else {
		fmt.Println("You have selected", *blacklistFile, "as the file of domains to blacklist")
	}

	maxTimeout, _ = strconv.Atoi(*maxTimeoutOption)

	contentWhitelist, err := ioutil.ReadFile(*whitelistFile)
	if err != nil {
		fmt.Println(err)
	}
	whitelist = strings.Split(string(contentWhitelist), "\n")
	whitelist = whitelist[:len(whitelist)-1]

	if len(*blacklistFile) > 0 {
		contentBlacklist, err := ioutil.ReadFile(*blacklistFile)
		if err != nil {
			fmt.Println(err)
		}
		blacklist = strings.Split(string(contentBlacklist), "\n")
		blacklist = blacklist[:len(blacklist)-1]

	}

	for i := 0; i < len(whitelist); i++ {
		fmt.Println(whitelist[i])

		var url string = fmt.Sprintf("http://%s/", whitelist[i])
		fmt.Println(url)
		getHREFfromURL(url)

	}

}
