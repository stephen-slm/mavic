package main

import (
	"flag"
	"fmt"
)

func generateScrapingOptions(arguments []string) ScrapingOptions {
	scrapingOptions := options{}
	scrapingOptions.parse(arguments)
	return scrapingOptions
}

func main() {
	scrapingOptions := generateScrapingOptions(flag.Args())
	fmt.Println(scrapingOptions)
}
