package main

import (
	// "fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"regexp"
	"strings"

	"github.com/gocolly/colly"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	ignoreRobots bool
	scrapeCmd    = &cobra.Command{
		Use:     "scrape",
		Aliases: []string{"s"},
		Short:   "GCA Wordpress website scraper.",
		Args:    cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			u, err := url.ParseRequestURI(args[0])
			if err != nil {
				log.Fatal("URL not valid: ", err)
			}

			err = os.Mkdir("site/"+u.Host, 0755)
			if err != nil {
				log.Warn("site/"+u.Host+" could not be created: ", err)
			}
			runScrape(u)
		},
	}
	matcher  = regexp.MustCompile("(\\.\\.\\/)?(http(s?):\\/\\/)?([A-Za-z0-9-_]+)?(\\.)?([A-Za-z0-9-_]+)?(\\.)?([A-Za-z0-9-_]+\\/?)+[A-Za-z0-9-_]*\\.(mp3|ogg|ogv|m4v|mp4|webm|ico|png|jpg|jpeg|gif|eot|woff|ttf|svg)")
	elements = []string{"a", "acronym", "article", "audio", "b", "body", "br", "button", "canvas", "caption", "center", "code", "details", "div", "em", "fieldset", "font", "footer", "form", "h1", "h2", "h3", "h4", "h5", "h6", "head", "hr", "i", "iframe", "img", "input", "label", "legend", "li", "link", "menu", "menuitem", "meter", "nav", "ol", "p", "pre", "progress", "section", "select", "small", "source", "span", "strike", "strong", "table", "tbody", "td", "textarea", "tfoot", "th", "thead", "title", "tr", "tt", "ul", "video"}
)

func init() {
	scrapeCmd.PersistentFlags().BoolVarP(&ignoreRobots, "ignore-robots", "i", false, "Ignore restrictions set by a host's robots.txt file.")
}

func runScrape(site *url.URL) {
	var c *colly.Collector

	if ignoreRobots {
		c = colly.NewCollector(
			colly.AllowedDomains(site.Host),
			colly.IgnoreRobotsTxt(),
		)
	} else {
		c = colly.NewCollector(
			colly.AllowedDomains(site.Host),
		)
	}

	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		c.Visit(e.Request.AbsoluteURL(link))
	})
	c.OnHTML("img[src]", func(e *colly.HTMLElement) {
		link := e.Attr("src")
		c.Visit(e.Request.AbsoluteURL(link))
	})
	c.OnHTML("div[data-avia-tooltip]", func(e *colly.HTMLElement) {
		link := e.Attr("data-avia-tooltip")
		c.Visit(e.Request.AbsoluteURL(link))
	})
	c.OnHTML("link[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		c.Visit(e.Request.AbsoluteURL(link))
	})
	c.OnHTML("style", func(e *colly.HTMLElement) {
		matchCSSLinks(e.Request, e.Text)
	})
	c.OnHTML("script[src]", func(e *colly.HTMLElement) {
		link := e.Attr("src")
		c.Visit(e.Request.AbsoluteURL(link))
	})

	for _, v := range elements {
		c.OnHTML(v+"[style]", func(e *colly.HTMLElement) {
			matchCSSLinks(e.Request, e.Attr("style"))
		})
	}

	c.OnResponse(func(r *colly.Response) {
		if strings.Contains(r.Headers.Get("Content-Type"), "text/css") {
			matchCSSLinks(r.Request, string(r.Body))
		}
		if len(r.Request.URL.Path) != 0 {
			p := path.Clean(r.Request.URL.Path[1:])
			if p != "." {
				if !strings.Contains(p, ".") {
					p = p + ".html"
				}

				err := createFilePath(p, site.Host)
				if err != nil {
					log.Warn(p+" could not be created: ", err)
				}
				err = r.Save("site/" + site.Host + "/" + p)
				if err != nil {
					log.Warn("Error saving a scraped file:", err)
				}
			} else {
				err := r.Save("site/" + site.Host + "/index.html")
				if err != nil {
					log.Warn("Error saving a scraped file:", err)
				}
			}
		}
	})

	c.Visit(site.String())
}

func createFilePath(filePath string, site string) error {
	buildFilePath := "site/" + site
	slicedPath := strings.Split(filePath, "/")
	directories := slicedPath[0 : len(slicedPath)-1]

	// fmt.Println(filePath, site)

	for _, s := range directories {
		if len(s) == 0 {
			continue
		}
		err := os.Mkdir(buildFilePath+"/"+s, 0777)
		if err != nil && os.IsExist(err) {
			buildFilePath += "/" + s
			continue
		} else if err != nil {
			return err
		}
		buildFilePath += "/" + s
	}

	return nil
}

func matchCSSLinks(r *colly.Request, stylesheet string) {
	m := matcher.FindAllString(stylesheet, -1)

	for _, v := range m {
		var link string
		p := strings.Split(r.URL.Path, "/")
		if strings.Contains(v, "..") {
			link = strings.Join(p[:len(p)-(1+strings.Count(v, ".."))], "/") + "/" + strings.Replace(v, "../", "", -1)
		} else if strings.Index(v, "/") == 0 {
			link = strings.Join(p[:len(p)-1], "/") + "/" + strings.Replace(v, "../", "", -1)
		} else {
			link = "/" + v
		}

		err := createFilePath(link, r.URL.Host)
		if err != nil {
			log.Warn("Error creating filepath:", err)
		}
		err = downloadResource("site/"+r.URL.Host+link, r.URL.Scheme+"://"+r.URL.Host+link)
		if err != nil {
			log.Warn("Error downloading resource: ", err)
		}
	}
}

func downloadResource(filepath string, url string) (err error) {
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	return nil
}
