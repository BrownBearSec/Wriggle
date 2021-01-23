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
var newURLsfound []string
var foundSubDomains []string
var newSubDomains []string
var maxTimeout int

func help() {
	fmt.Println("Displaying the help page")
	fmt.Println("Usage: wriggle -w <FILE>")
	fmt.Println("-w <FILE> : Specificy a list of domains in scope, one per line. Note: do not right '*.domain.com' just write 'domain.com' ")
	fmt.Println("-b <FILE> : Specificy a list of domains not in scope, one per line ")
	fmt.Println("-t <number> : Set the max timeout (in seconds) for connecting to a URL, default 20 seconds")
	fmt.Println("-s <FILE> : Specifiy the name of the subdomain output file, default is 'subDomainsOf' + time of scan")
	fmt.Println("-u <FILE> : Specifiy the name of the URL output file, default is 'URLsOf' + time of scan")
	fmt.Println("-h : Display this help page")
	os.Exit(3)
}

func getHREFfromURL(url string) {

	client := http.Client{
		Timeout: time.Duration(maxTimeout) * time.Second,
	}

	if strings.Index(url, "http") == 0 {
		resp, err := client.Get(url)
		if err != nil {
			fmt.Println("GetXYZ")
			if strings.Contains(err.Error(), "Client.Timeout") {
				fmt.Println("The get request has timed out, either increase max timeout or check if the site is up")
				return
			}
			if strings.Contains(err.Error(), "connection reset by") {
				fmt.Println("The connection was reset by peer")
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
				if !inArray(foundURLs, foundURL) && !inArray(newURLsfound, foundURL) {
					newURLsfound = append(newURLsfound, foundURL)
				}
			}
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

func inArray(arr []string, toFind string) bool {
	var answer bool = false
	for i := 0; i < len(arr); i++ {
		if arr[i] == toFind {
			answer = true
		}
	}
	return answer
}

func extractSubDomains() {
	for i := 0; i < len(newURLsfound); i++ {
		urlSectionArray := strings.Split(newURLsfound[i], "/")
		for j := 0; j < len(urlSectionArray); j++ {
			for k := 0; k < len(whitelist); k++ {
				if strings.Contains(urlSectionArray[j], whitelist[k]) && !inArray(foundSubDomains, urlSectionArray[j]) && !inArray(newSubDomains, urlSectionArray[j]) && !strings.Contains(urlSectionArray[j], "mailto:") && !strings.Contains(urlSectionArray[j], "=") {
					newSubDomains = append(newSubDomains, urlSectionArray[j])
				}
			}
		}

	}
}

func writeToFile(fileName string, listOfStrings []string) {
	f, err := os.OpenFile(fileName,
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println(err)
	}
	defer f.Close()
	for i := 0; i < len(listOfStrings); i++ {
		if _, err := f.WriteString(listOfStrings[i] + "\n"); err != nil {
			fmt.Println(err)
		}
	}
}

func main() {

	startTimetime := time.Now()
	startTime := startTimetime.String()
	defaultSubdomainName := "subDomainsOf" + startTime
	defaultURLName := "URLsOf" + startTime

	wantHelp := flag.Bool("h", false, "display help page")
	whitelistFile := flag.String("w", "", "the whitelist file for domains")
	blacklistFile := flag.String("b", "", "the blacklist file for domains")
	maxTimeoutOption := flag.String("t", "20", "max timeout for connection timeouts")
	subDomainOutFile := flag.String("s", defaultSubdomainName, "name of subdomain found file")
	URLOutFile := flag.String("u", defaultURLName, "name of URLs found out file")
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
		var url string = fmt.Sprintf("http://%s/", whitelist[i])
		foundURLs = append(foundURLs, url)
	}

	i := 0
	for i < len(foundURLs) {
		fmt.Println("Now processing : ", foundURLs[i])
		getHREFfromURL(foundURLs[i])
		extractSubDomains()

		a := append(foundURLs, newURLsfound...)
		foundURLs = a
		b := append(foundSubDomains, newSubDomains...)
		foundSubDomains = b
		fmt.Println("number of new URLs found : ", len(newURLsfound))
		fmt.Println("number of new subdomains discovered : ", len(newSubDomains))
		//This is the point where you write the newly discovered stuff to a file later
		writeToFile(*subDomainOutFile, newSubDomains)
		writeToFile(*URLOutFile, newURLsfound)
		newURLsfound = nil
		newSubDomains = nil
		//fmt.Println(foundURLs)
		i++
	}

}
