package main

import (
	"flag"
	"fmt"
	"os"
)

func parseCommandLineArguments(arguments []string) ScrapingOptions {
	CommandLine := flag.CommandLine
	scrapingOptions := ScrapingOptions{}

	CommandLine.StringVar(&scrapingOptions.outputDirectory, "OutputDirectory", "./", "The output directory to store the images.")
	CommandLine.IntVar(&scrapingOptions.imageLimit, "ImageLimit", 50, "The output directory to store the images.")
	CommandLine.BoolVar(&scrapingOptions.frontPage, "FrontPage", false, "The output directory to store the images.")

	CommandLine.Parse(arguments)

	return scrapingOptions
}

func main() {
	scrapingOptions := parseCommandLineArguments(os.Args)
	fmt.Println(scrapingOptions)
}
