# Wriggle

## A go based web crawler and asset discoverer. Currently a WIP

*A project made for me to learn golang*

### If you appreciate what I do please help

https://commerce.coinbase.com/checkout/1f3b4d8f-afb9-47f4-8a18-1e6f4502ce82

### Usage:

```
wriggle -w <FILE>
-w <FILE> : Specificy a list of domains in scope, one per line. Note: do not right '*.domain.com' just write 'domain.com' 
-b <FILE> : Specificy a list of domains not in scope, one per line 
-t <number> : Set the max timeout (in seconds) for connecting to a URL, default 20 seconds
-h : Display this help page
```

### Example:

`wriggle -w domains.txt -b outOfScope.txt`

- where `domains.txt` is a file with a new domain on each line
- where `outOfScope.txt` is a file of any subdomains out of scope, Eg `admin.domain.com`. Also one per line 

### Features:

- Get values from Href tags from html
- deals with relative paths
- Checks for scope 
- Recursion
- Output options
- Subdomain identification


### Features to add:

- Threading
- opposite of verbose mode
- program timer
- timeout counter + url recording
- option not to output to file
- Colours / nice formatting
- JS file identification
- Easy install
- Make logo (least priority)

### Contact me:

If you need to contact me for any reason, bugs, feedback, anything:

email: contact@josephwitten.com

### Obvious legal points

- only use this tool against authorised targets
- Check the scope for the use of automated tools
- use ethically and legally
- etc