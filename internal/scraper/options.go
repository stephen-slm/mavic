package scraper

type Options struct {
	//  The directory in which we will be downloading all the images into, based on the folder name
	//  of the given sub-reddit.
	OutputDirectory string
	// The total number of images to download max per sub-reddit before we continue to the next one.
	ImageLimit int
	// If set to true, the tool will scrape the front page of reddit for the current most
	// active sub-reddits and then scrape all the links directly from them sub-reddits.
	FrontPage bool
	// if the images are being downloaded directly into the root folder and nothing else.
	RootFolderOnly bool
	// You can change this to adjust on what kind of images you get from Reddits filtering
	// options (hot, new, rising, controversial, top), hot is the default by reddit while
	// also the default in the tool.
	PageType string
	// What subreddits are going to be scrapped for downloading of sad images. If front page is
	// parsed as true then the front page will be pushed onto the sub reddit listings.
	Subreddits []string
	// This is the max number of images that can be downloaded together at anyone time. This is
	// a limit that has to be set because we could end up hitting rate limiting hits on sites
	// that are allowing us to download the images directly.
	MaxConcurrentDownloads int
	// If the loading progress bar should be displayed or not. Simply used for headless progressing
	// or testing that helps with minimising the amount of output that is generated to the console.
	DisplayLoading bool
}
