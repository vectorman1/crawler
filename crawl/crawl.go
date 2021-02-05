package crawl

import (
	"crawler/utils"
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"
)

type Crawler struct {
	mux               sync.Mutex
	Sites             []string
	LinksOutput       chan []string
	UnvisitedWorkList chan string
	VisitedLinks      map[string]bool
	Fingerprints      []Fingerprint
	Depth             int
	ScannerData       utils.ScannerData
	Recurse           bool
	MaxCrawls         int
}

func (c *Crawler) crawl(client *http.Client, url string) (*http.Response, error) {
	log.Println("getting", url)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", "golang_crawler/1.0")
	resp, err := client.Do(req)

	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("website responded with %d", resp.StatusCode)
	}

	return resp, nil
}

func (c *Crawler) Run() {
	transport := http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
		IdleConnTimeout: 5 * time.Second,
	}
	client := http.Client{
		Timeout:   5 * time.Second,
		Transport: &transport,
	}

	go func() { c.LinksOutput <- c.Sites }()
	var result []Fingerprint
	totalFingerprints := 0
	for i := 0; i < c.MaxCrawls; i++ {
		go func() {
			for link := range c.UnvisitedWorkList {
				totalFingerprints++
				fingerprint, err := NewFingerprint(&client, link, c)
				if err != nil {
					log.Println("Error generating fingerprint ", err)
					totalFingerprints--
					continue
				}
				c.mux.Lock()
				links, err := fingerprint.UnseenUniqueLinks(c.VisitedLinks)
				c.mux.Unlock()
				if err != nil {
					log.Println("Error getting links ", err)
					totalFingerprints--
					continue
				}
				result = append(result, *fingerprint)

				if c.Recurse {
					go func() {
						c.LinksOutput <- links
					}()
				} else {
					close(c.LinksOutput)
					log.Println("recurse flag is set to false, exiting")
					break
				}
			}
		}()
	}

	for list := range c.LinksOutput {
		if c.Depth == 0 {
			log.Println("Reached max depth")
			log.Println("waiting for", totalFingerprints-len(result), "fingerprints to finish generating")
			for {
				if len(result) == totalFingerprints {
					break
				}
			}
			break
		}
		for _, link := range list {
			if !c.VisitedLinks[link] {
				c.VisitedLinks[link] = true
				c.UnvisitedWorkList <- link
			}
		}
		c.Depth--
	}

	utils.SaveToDiskAsJson(result)
}
