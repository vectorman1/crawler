package main

import (
	"crawler/crawl"
	"crawler/utils"
	"flag"
	"log"
)

func main() {
	recurseFlag := flag.Bool("recurse", true, "Dig into links.")
	maxFlag := flag.Int("max_heads", 1000, "Max concurrent crawls.")
	depthFlag := flag.Int("depth", 10, "Depth of crawl")

	flag.Parse()

	if len(flag.Args()) == 0 {
		log.Println("Available flags:" +
			"\n\t-recurse=true|false - dig into links found" +
			"\n\t-max_heads=1000 - max concurrent crawls" +
			"\n\t-depth=10 - depth of crawl")

		log.Fatalln("Input urls as args")
	}

	var urls []string
	for _, arg := range flag.Args() {
		urls = append(urls, arg)
	}

	s := utils.NewScannerData()
	crawler := crawl.Crawler{
		Sites:             urls,
		LinksOutput:       make(chan []string),
		UnvisitedWorkList: make(chan string),
		VisitedLinks:      make(map[string]bool),
		Depth:             *depthFlag,
		ScannerData:       *s,
		Recurse:           *recurseFlag,
		MaxCrawls:         *maxFlag,
	}

	log.Printf("Starting crawl with flags: recurse=%v, max_heads=%d, depth=%d", *recurseFlag, *maxFlag, *depthFlag)
	crawler.Run()
}
