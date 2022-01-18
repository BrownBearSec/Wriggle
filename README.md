# Wriggle

## A go based web crawler and asset discoverer. Currently a WIP

*A project made for me to learn golang*

![logo?](other/wriggle.png)

### If you appreciate what I do please help

### Installation

1) `git clone https://github.com/iamSm9l/Wriggle.git`
2) `cd Wriggle`
3) `chmod +x install.sh`
4) `sudo ./install.sh`


### Usage:

```
wriggle -w <FILE>
-w <FILE> : Specificy a list of domains in scope, one per line. Note: do not right '*.domain.com' just write 'domain.com' 
-b <FILE> : Specificy a list of domains not in scope, one per line 
-t <number> : Set the max timeout (in seconds) for connecting to a URL, default 20 seconds
-s <FILE> : Specifiy the name of the subdomain output file, default is 'subDomainsOf' + time of scan
-u <FILE> : Specifiy the name of the URL output file, default is 'URLsOf' + time of scan
-j <FILE> : Specifiy the name of the JS output file, default is 'JSfilesOf' + time of scan
-4 <FILE> : Specifiy the name of the 404 output file, default is '404FilesOf' + time of scan
-v : verbose mode, not advisiable unless you love spam
-h : Display this help page
```

### Example:

`wriggle -w domains.txt -b outOfScope.txt`

- where `domains.txt` is a file with a new domain on each line
- where `outOfScope.txt` is a file of any subdomains out of scope, Eg `admin.domain.com`. Also one per line 

An example of both files can be seen in the examples folder

### Features:

- Get values from Href tags from html
- deals with relative paths
- Checks for scope 
- Recursion
- Output options
- Subdomain identification
- Final report
- Colours / nice formatting
- opposite of verbose mode
- program timer
- JS file identification


### Features to add:

- Threading
- option not to output to file

### Contact me:

If you need to contact me for any reason, bugs, feedback, anything:

email: contact@josephwitten.com

### Obvious legal points

- only use this tool against authorised targets
- Check the scope for the use of automated tools
- use ethically and legally
- etc
