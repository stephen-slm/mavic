package main

import (
	"flag"

	"github.com/tehstun/mavic/src/reddit"
)

// generateScrapingOptions generates some scraping options based on the input
// arguments which would commonly be the command line arguments.
func generateScrapingOptions(arguments []string) reddit.Options {
	options := reddit.Options{}
	options.Parse(arguments)
	return options
}

func main() {
	scrapingOptions := generateScrapingOptions(flag.Args())
	scraper := reddit.NewScraper(scrapingOptions)
	scraper.ProcessSubreddits()
}
