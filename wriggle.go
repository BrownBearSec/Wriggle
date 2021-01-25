package main

import (
	"flag"
	"fmt"
	"io/ioutil"
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
var numOfTimeout int = 0
var timeOuts []string
var oddURLs []string
var newJSfiles []string
var foundJSfiles []string
var found404URLs []string
var new404URLs []string

var reset string = "\033[0m"
var red string = "\033[31m"
var green string = "\033[32m"
var yellow string = "\033[33m"
var blue string = "\033[34m"
var purple string = "\033[35m"
var cyan string = "\033[36m"
var gray string = "\033[37m"
var white string = "\033[97m"

func help() {
	fmt.Println("Displaying the help page")
	fmt.Println("Usage: wriggle -w <FILE>")
	fmt.Println("-w <FILE> : Specificy a list of domains in scope, one per line. Note: do not right '*.domain.com' just write 'domain.com' ")
	fmt.Println("-b <FILE> : Specificy a list of domains not in scope, one per line ")
	fmt.Println("-t <number> : Set the max timeout (in seconds) for connecting to a URL, default 20 seconds")
	fmt.Println("-s <FILE> : Specifiy the name of the subdomain output file, default is 'subDomainsOf' + time of scan")
	fmt.Println("-u <FILE> : Specifiy the name of the URL output file, default is 'URLsOf' + time of scan")
	fmt.Println("-j <FILE> : Specifiy the name of the JS output file, default is 'JSfilesOf' + time of scan")
	fmt.Println("-4 <FILE> : Specifiy the name of the 404 output file, default is '404FilesOf' + time of scan")
	fmt.Println("-v : verbose mode, not advisiable unless you love spam")
	fmt.Println("-h : Display this help page")
	os.Exit(3)
}

func getHREFfromURL(url string) {

	client := http.Client{
		Timeout: time.Duration(maxTimeout) * time.Second,
	}

	if strings.Index(url, "http") == 0 && !(url[len(url)-3:] == ".js") && !(url[len(url)-4:] == ".css") {

		resp, err := client.Get(url)
		if err != nil {
			//fmt.Println("GetXYZ")
			if strings.Contains(err.Error(), "Client.Timeout") {
				fmt.Println(red + "[Warning]" + reset + " " + white + "The get request has timed out, either increase max timeout or check if the site is up : " + url + reset)
				numOfTimeout++
				timeOuts = append(timeOuts, url)
				return
			}
			if strings.Contains(err.Error(), "connection reset by") {
				fmt.Println(red + "[Warning]" + reset + " " + white + "The connection was reset by peer : " + url + reset)
				return
			}
			fmt.Println(red + "[Warning]" + reset + " " + white + err.Error() + " : " + url + reset)
			oddURLs = append(oddURLs, url)
			return
		}

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("readingXYZ")
			fmt.Println(err)
			oddURLs = append(oddURLs, url)
			return
		}

		defer resp.Body.Close()

		htmlStrig := string(body)

		if resp.StatusCode == 404 {
			//fmt.Println(red+"[Warning]"+reset+" "+white+"Status code 404 on url : ", url)
			new404URLs = append(new404URLs, url)
			return
		}

		htmlArray := strings.Split(htmlStrig, "href=\"")

		for i := 0; i < len(htmlArray); i++ {
			htmlPortion := string(htmlArray[i])
			splitBySpeechmarkArray := strings.Split(htmlPortion, "\"")
			foundURL := splitBySpeechmarkArray[0]

			if len(foundURL) > 1 {
				if foundURL[0] == '/' && foundURL[1] != '/' {
					//fmt.Println("found url : ", foundURL)
					//fmt.Println("url : ", url)
					if url[len(url)-1] == '/' && foundURL[0] == '/' {
						foundURL = foundURL[1:]
					}

					foundURL = url + foundURL
					//fmt.Println("after : ", foundURL)
				}
			}

			isInScope := checkScope(foundURL)
			if isInScope {
				if !inArray(foundURLs, foundURL) && !inArray(newURLsfound, foundURL) {
					//fmt.Println("found this : ", foundURL)
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
		if strings.Index(urlToCheck, "=") > -1 {
			if strings.Index(urlToCheck, "=") < strings.Index(urlToCheck, whitelist[i]) {
				whitelisted = false
			}
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
				//fmt.Println("before : ", urlSectionArray[j])
				if strings.Contains(urlSectionArray[j], whitelist[k]) && !inArray(foundSubDomains, urlSectionArray[j]) && !inArray(newSubDomains, urlSectionArray[j]) && !strings.Contains(urlSectionArray[j], "mailto:") && !strings.Contains(urlSectionArray[j], "=") && !strings.Contains(urlSectionArray[j], "@") {
					//fmt.Println("test : ", urlSectionArray[j])
					//fmt.Println("whitelist? : ", whitelist[k])
					//fmt.Println(newURLsfound[i])
					newSubDomains = append(newSubDomains, urlSectionArray[j])
				}
			}
		}

	}
}

func extractJSfiles() {
	for i := 0; i < len(newURLsfound); i++ {
		url := newURLsfound[i]
		if len(url) > 5 {
			if url[len(url)-3:] == ".js" && !inArray(newJSfiles, url) && !inArray(foundJSfiles, url) {
				if strings.Index(url, "http") != 0 {
					url = "https://" + url
				}
				newJSfiles = append(newJSfiles, url)
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

func printSubDomains(newSub []string) {
	for i := 0; i < len(newSub); i++ {
		fmt.Println(blue + "[Sub domain]" + reset + " " + white + newSub[i])
	}
}

func printJsFiles(newJS []string) {
	for i := 0; i < len(newJS); i++ {
		fmt.Println(yellow + "[JS file]" + reset + " " + white + newJS[i])
	}
}

func main() {

	startTimetime := time.Now()
	startTime := startTimetime.String()
	defaultSubdomainName := "subDomainsOf" + startTime
	defaultURLName := "URLsOf" + startTime
	defaultJSName := "JSfilesOf" + startTime
	default404Name := "404FilesOf" + startTime

	wantHelp := flag.Bool("h", false, "display help page")
	whitelistFile := flag.String("w", "", "the whitelist file for domains")
	blacklistFile := flag.String("b", "", "the blacklist file for domains")
	maxTimeoutOption := flag.String("t", "20", "max timeout for connection timeouts")
	subDomainOutFile := flag.String("s", defaultSubdomainName, "name of subdomain found file")
	URLOutFile := flag.String("u", defaultURLName, "name of URLs found out file")
	JSOutFile := flag.String("j", defaultJSName, "name of JS found out file")
	verbose := flag.Bool("v", false, "verbose output?")
	fourOhFourFile := flag.String("4", default404Name, "file name of 404 output")
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
		fmt.Println(cyan+"[Info]"+reset+" "+white+"whitelist selected : ", *whitelistFile+reset)
	}

	if *blacklistFile == "" {
		fmt.Println(red + "[Warning]" + reset + " " + white + ": no blacklisted domains/subdomains, continuing" + reset)
	} else {
		fmt.Println(cyan+"[Info]"+reset+white+" blacklist selected : ", *blacklistFile)
	}

	maxTimeout, _ = strconv.Atoi(*maxTimeoutOption)

	contentWhitelist, err := ioutil.ReadFile(*whitelistFile)
	if err != nil {
		fmt.Println(err)
	}
	whitelist = strings.Split(string(contentWhitelist), "\n")
	whitelist = whitelist[:len(whitelist)-1]

	if inArray(whitelist, " ") || inArray(whitelist, "") {
		fmt.Println(red + "[ERROR]" + reset + " " + white + ": There is an empty line in your whitelist file (maybe the last line) that would cause scope to accept anything, fix and rerun")
		os.Exit(3)
	}

	if len(*blacklistFile) > 0 {
		contentBlacklist, err := ioutil.ReadFile(*blacklistFile)
		if err != nil {
			fmt.Println(err)
		}
		blacklist = strings.Split(string(contentBlacklist), "\n")
		blacklist = blacklist[:len(blacklist)-1]

	}

	for i := 0; i < len(whitelist); i++ {
		var url string = fmt.Sprintf("https://%s/", whitelist[i])
		foundURLs = append(foundURLs, url)
	}

	i := 0
	for i < len(foundURLs) {
		if *verbose {
			fmt.Println("Now processing : ", foundURLs[i])
		}

		if i%1000 == 0 {
			fmt.Println(green + "[Progress]" + reset + " " + white + strconv.Itoa(i) + "/" + strconv.Itoa(len(foundURLs)))
		}

		getHREFfromURL(foundURLs[i])
		extractSubDomains()
		extractJSfiles()

		a := append(foundURLs, newURLsfound...)
		foundURLs = a
		b := append(foundSubDomains, newSubDomains...)
		foundSubDomains = b
		c := append(foundJSfiles, newJSfiles...)
		foundJSfiles = c
		d := append(found404URLs, new404URLs...)
		found404URLs = d
		if *verbose {
			fmt.Println("number of new URLs found : ", len(newURLsfound))
			fmt.Println("number of new subdomains discovered : ", len(newSubDomains))
			fmt.Println("number of js files found : ", len(newJSfiles))
		}
		printSubDomains(newSubDomains)
		printJsFiles(newJSfiles)
		//This is the point where you write the newly discovered stuff to a file later
		writeToFile(*subDomainOutFile, newSubDomains)
		writeToFile(*URLOutFile, newURLsfound)
		writeToFile(*JSOutFile, newJSfiles)
		writeToFile(*fourOhFourFile, new404URLs)
		newURLsfound = nil
		newSubDomains = nil
		newJSfiles = nil
		new404URLs = nil
		//fmt.Println(foundURLs)
		i++
	}

	fmt.Println(cyan + "[Info]" + reset + " " + white + "Final report")
	fmt.Println(cyan+"[Info]"+reset+" "+white+"Number of subdomains found : ", len(foundSubDomains))
	fmt.Println(cyan+"[Info]"+reset+" "+white+"Number of URLs found : ", len(foundURLs))
	fmt.Println(cyan+"[Info]"+reset+" "+white+"Number of 404 status codes found : ", len(found404URLs))
	if numOfTimeout > 0 {
		timeoutFileName := "URLsThatTimedOut" + startTime
		writeToFile(timeoutFileName, timeOuts)
		fmt.Println(cyan+"[Info] "+reset+" "+white+"Number of links timed out : ", numOfTimeout)
		fmt.Println("(Written to file", timeoutFileName, ")")
	}
	if len(oddURLs) > 0 {
		oddURLsFileName := "OddURLs" + startTime
		writeToFile(oddURLsFileName, oddURLs)
		fmt.Println(cyan + "[Info]" + reset + " " + white + strconv.Itoa(len(oddURLs)) + " number of odd / strange error inducing URLs found, manually reveue is recommended")
		fmt.Println(cyan+"[Info]"+reset+" "+white+"(Written to file", oddURLsFileName, ")")
	}
	elapsed := time.Since(startTimetime).String()
	fmt.Println(cyan + "[Info]" + reset + " " + white + "Scan took " + elapsed + " seconds")
}
