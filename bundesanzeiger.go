package main

import (
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/Sirupsen/logrus"
)

const (
	domain    = "https://www.bundesanzeiger.de"
	startPage = "https://www.bundesanzeiger.de/ebanzwww/wexsservlet?page.navid=to_nlp_start"
)

// Crawl content from bundesanzeiger.de
func CrawlBundesanzeiger() {
	logrus.Infoln("You Choose to crawl bundesanzeiger.de, please wait...")
	urls := findAllPages()
	logrus.Infoln("get all entry count", len(urls))

	// get cookies we need
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	resp, err := client.Get(startPage)
	if err != nil {
		logrus.Errorln(err)
	}
	defer resp.Body.Close()

	f, err := os.Create(outputFile)
	if err != nil {
		logrus.Fatalln("can not create file", err)
	}
	defer f.Close()
	for _, url := range urls {
		var (
			doc *goquery.Document
			err error
		)
		url = domain + url
		req, _ := http.NewRequest("GET", url, nil)
		resp, _ := client.Do(req)
		if doc, err = goquery.NewDocumentFromResponse(resp); err != nil {
			logrus.Errorln(err)
		}
		link, ok := doc.Find("div.nlp_historylist_download_csv_right_top a.doc").Attr("href")
		if !ok {
			logrus.Errorln("can not find csv link, please download manually", url)
			continue
		}
		link = domain + link
		logrus.Debugln("csv link is", link)
		r := get(*client, link)
		r = strings.Replace(r, `"Positionsinhaber","Emittent","ISIN","Position","Datum"`, "", 1)
		r = strings.Trim(r, "\r\n \t")
		logrus.Infoln("fetched data successfully: ", url)

		if _, err := f.WriteString(r); err != nil {
			logrus.Errorln("can not write to file:", err)
		}
		if err := f.Sync(); err != nil {
			logrus.Errorln(err)
		}
	}

	logrus.Infoln("\nfinish fetching all data, have a nice day")
}

func findAllPages() []string {
	var (
		page int
	)
	ret := make(map[string]int)
	for url := startPage; url != ""; {
		var (
			doc *goquery.Document
			err error
		)
		logrus.Infoln("fetch page url:", url)
		if doc, err = goquery.NewDocument(url); err != nil {
			logrus.Fatal(err)
		}

		if relativeUrl, ok := doc.Find("li.next_page a").Attr("href"); ok {
			url = domain + relativeUrl
		} else {
			break
		}
		doc.Find("a.intern").Each(func(i int, s *goquery.Selection) {
			if href, ok := s.Attr("href"); ok {
				ret[href] = 1
			}
		})
		page++
	}

	retArr := make([]string, len(ret))
	i := 0
	for href := range ret {
		retArr[i] = href
		i++
	}
	return retArr
}

func get(client http.Client, csvLink string) string {
	var (
		resp *http.Response
		err  error
	)
	req, err := http.NewRequest("GET", csvLink, nil)
	if resp, err = client.Do(req); err != nil {
		logrus.Errorln(err)
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logrus.Errorln(err)
	}
	return string(data)
}
