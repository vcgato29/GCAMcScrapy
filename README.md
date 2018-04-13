# GCA McScrapy
### Genesis
To maintain and enhance GCA’s secure operating environment, a secure Website nearly immune to compromise was established in September 2016 and the DMARC micro-site was established in October 2016. Both Websites were created using WordPress in order to allow content creators the ability to modify text on the sites with ease.
 
While WordPress is a popular blogging platform, by its nature, it is prone to potential compromise. WordPress dynamically composes web pages using PHP and JavaScript and thus carries with it a high risk for bugs and security vulnerabilities that serve as a vector for compromise. Because the ability to create and update content on the sites by multiple parties is a necessity, the decision was made to secure the sites by scraping all of the dynamic content into static sites.

### How McScrapy Works
The foundation is a scraping tool that attempts to scrape every piece of a website to be as functional as possible as a static clone, removing potential security issues of third-party services or unnecessary requests. HTML pages are scanned for URLs embedded in element attributes thoroughly including but not limited to: element href attributes, img element src attributes, contents of style tags, and inline styles. If any CSS files are found, their contents are also scanned for potential resources. As they’re found, resources (JS, CSS, images, PDFs, etc.) are downloaded from those URLs and saved relative to the path portion of the URL mimicking the original structure of the website being scraped.

McScrapy includes the ability to debug a scrape, cache files as they’re saved, ignore a website’s robots.txt restrictions, specify a maximum recursive depth to scan HTML pages, and scrape using a specific user agent. These features can be used in any combination, for example, to reduce scan times, acquire more resources, and scrape mobile sites.

A preview function is available to test the functionality of the scraped website to check for completion and general asset availability once the scrape has finished. Using a generic file server, the preview function is able to host a now static clone of the original website including dynamic routing of HTML pages.


## Build
To build the application simply run `make` in the root directory. Alternatively run:

`go build -o bin/mcscrapy github.com/GlobalCyberAlliance.org/McScrapy/cmd/mcscrapy`

## Scrape
To scrape a website run `mcscrapy scrape [domain]`

#### Flags
`-c` `--cache` Specify where requests are caches as files.

`-d` `--debug` Output of debug logs.

`-i` `--ignore-robots` Ignore restrictions set by a host's robots.txt file.

`-m` `--max-depth` Set the max depth of recursion of visited URLs. Leave blank to allow all.

`-o` `--output-dir` Output scraped websites to a specific directory.

`-u` `--user-agent` Set the user agent used by the scraper.

`-v` `--verbose` Verbose output of logs.

## Preview
Preview a scraped website with a built in web server. Paths ending in `/` or without a file extension default to serving the path as an HTML file.

`mcscrapy preview [path_to_site_directory]`

#### Flags
`-a` `--address` Set the address of the preview. Default: `127.0.0.1`

`-p` `--port` Set the port of the preview. Default: `8000`


## License
This repository is licensed under the Apache License version 2.0.

Some of the project's dependencies may be under different licenses.
