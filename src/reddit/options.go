package reddit

import (
	"flag"
)

type Options struct {
	//  The directory in which we will be downloading all the images into, based on the folder name
	//  of the given sub-reddit.
	outputDirectory string

	// The total number of images to download max per sub-reddit before we continue to the next one.
	imageLimit int

	// If set to true, the tool will scrape the front page of reddit for the current most active
	// sub-reddits and then scrape all the links directly from them sub-reddits.
	frontPage bool

	// You can change this to adjust on what kind of images you get from Reddits filtering options
	// (hot, new, rising, controversial, top), hot is the default by reddit while also the default
	// in the tool.
	pageType string

	// What subreddits are going to be scrapped for downloading of sad images.
	subreddits []string
}

// parse takes a list of command line arguments and binds them to the scrapingOptions, using default
// values if the arguments are not given.
func (s *Options) Parse(arguments []string) {
	commandLine := flag.CommandLine

	commandLine.StringVar(&s.outputDirectory, "OutputDirectory", "./", "The output directory to store the images.")
	commandLine.IntVar(&s.imageLimit, "ImageLimit", 50, "The output directory to store the images.")
	commandLine.BoolVar(&s.frontPage, "FrontPage", false, "The output directory to store the images.")

	_ = commandLine.Parse(arguments)
}
